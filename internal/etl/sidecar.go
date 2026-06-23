// Package etl holds the pure-function ETL helpers that back the
// data-engineering domain profile. Each helper is a small, table-testable
// function consumed by the sdd-* skills (via the gentle-ai CLI shell-out) and
// by the apply/verify branch points. The package deliberately avoids
// gopkg.in/yaml.v3, mirroring the hand-rolled scanners in internal/sddconfig
// and internal/components/filemerge so the sidecar parser stays inside the
// project's "no external YAML dependency" convention.
package etl

import (
	"errors"
	"fmt"
	"strings"
)

// ErrSidecarValidation is the sentinel for sidecar parse/validate failures so
// callers can distinguish a sidecar structural problem from a transport error.
var ErrSidecarValidation = errors.New("sidecar: validation error")

// Sidecar is the parsed view of a glue-tables/{db}.{table}.yaml sidecar file.
// The shape mirrors the subset of Glue's `get-table` output that this package
// structurally validates against (table name, columns, partitions, S3
// location). Comments are preserved so the spec/design prose can quote them.
type Sidecar struct {
	Database    string   `json:"database"`
	Table       string   `json:"table"`
	Columns     []Column `json:"columns"`
	Partitions  []string `json:"partitions"`
	S3Location  string   `json:"s3Location"`
	Format      string   `json:"format"`
	Compression string   `json:"compression"`
}

// Column is one column entry of a Sidecar.
type Column struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Comment string `json:"comment"`
}

// Mismatch is one structural drift found by ValidateSidecar between a parsed
// Sidecar and the live Glue table (`aws glue get-table` output). Kind names the
// drift category; Column names the offending column where relevant (empty for
// database/table/s3 mismatches). Expected/Actual carry human-readable context.
type Mismatch struct {
	Kind     string `json:"kind"`
	Column   string `json:"column"`
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
	Message  string `json:"message"`
}

// Mismatch kinds produced by ValidateSidecar.
const (
	MismatchDatabase            = "database_mismatch"
	MismatchTable               = "table_mismatch"
	MismatchS3Location          = "s3_location_mismatch"
	MismatchMissingColumn       = "missing_column"
	MismatchTypeMismatch        = "type_mismatch"
	MismatchMissingPartition    = "missing_partition"
	MismatchUnexpectedPartition = "unexpected_partition"
)

// ParseSidecar parses a serialized glue-table sidecar (YAML) into a Sidecar.
// It rejects structurally invalid sidecars (empty input, missing database,
// columns without name/type) by returning an error wrapping ErrSidecarValidation.
// Comment lines and trailing inline `# ...` comments are stripped, but content
// inside quotes is preserved verbatim.
func ParseSidecar(data []byte) (Sidecar, error) {
	text := strings.TrimSpace(string(data))
	if text == "" {
		return Sidecar{}, fmt.Errorf("%w: empty sidecar", ErrSidecarValidation)
	}
	var s Sidecar
	lines := strings.Split(text, "\n")
	for i := 0; i < len(lines); i++ {
		raw := lines[i]
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if leadingSpaces(raw) != 0 {
			// Stray indented line outside a managed block; skip defensively.
			continue
		}
		key, value := splitKeyValue(trimmed)
		switch key {
		case "database":
			s.Database = unquoteScalar(value)
		case "table":
			s.Table = unquoteScalar(value)
		case "s3_location":
			s.S3Location = unquoteScalar(value)
		case "format":
			s.Format = unquoteScalar(value)
		case "compression":
			s.Compression = unquoteScalar(value)
		case "partitions":
			s.Partitions = parseFlowSeq(unquoteScalar(value))
		case "columns":
			cols, next, err := parseColumnsBlock(lines, i)
			if err != nil {
				return Sidecar{}, err
			}
			s.Columns = cols
			i = next - 1
		default:
			// Unknown top-level key: ignore (forward-compat).
		}
	}
	if err := s.validateStructure(); err != nil {
		return Sidecar{}, err
	}
	return s, nil
}

// validateStructure enforces the minimal structural contract a sidecar must
// hold before it can be validated against a Glue table.
func (s Sidecar) validateStructure() error {
	if s.Database == "" {
		return fmt.Errorf("%w: missing database", ErrSidecarValidation)
	}
	if s.Table == "" {
		return fmt.Errorf("%w: missing table", ErrSidecarValidation)
	}
	if len(s.Columns) == 0 {
		return fmt.Errorf("%w: no columns", ErrSidecarValidation)
	}
	for idx, c := range s.Columns {
		if c.Name == "" {
			return fmt.Errorf("%w: column #%d missing name", ErrSidecarValidation, idx)
		}
		if c.Type == "" {
			return fmt.Errorf("%w: column %q missing type", ErrSidecarValidation, c.Name)
		}
	}
	return nil
}

// parseColumnsBlock consumes the indented children following the `columns:`
// header at lines[start]. Each item is expected to be an inline flow-map
// (e.g. `- {name: id, type: string, comment: "x"}`). Returns the parsed
// columns, the index AFTER the last consumed line, and any structural error.
func parseColumnsBlock(lines []string, start int) ([]Column, int, error) {
	cols := []Column{}
	last := start
	for j := start + 1; j < len(lines); j++ {
		trimmed := strings.TrimSpace(lines[j])
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			last = j
			continue
		}
		if leadingSpaces(lines[j]) == 0 {
			break
		}
		if !strings.HasPrefix(trimmed, "- ") && trimmed != "-" {
			// Indented non-dash line we don't manage — ignore defensively.
			last = j
			continue
		}
		body := strings.TrimSpace(strings.TrimPrefix(trimmed, "-"))
		col, err := parseInlineColumn(body)
		if err != nil {
			return nil, j, err
		}
		cols = append(cols, col)
		last = j
	}
	return cols, last, nil
}

// parseInlineColumn parses an inline flow-map column entry like:
//
//	{name: id, type: string, comment: "survey id"}
//
// Quotes (single or double) are honored so values may contain spaces or
// `#`. A column that lacks a name yields ErrSidecarValidation.
func parseInlineColumn(body string) (Column, error) {
	var c Column
	body = strings.TrimSpace(strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(body), "{")))
	body = strings.TrimSuffix(body, "}")
	entries, err := splitFlowMapEntries(body)
	if err != nil {
		return Column{}, fmt.Errorf("%w: %v", ErrSidecarValidation, err)
	}
	for _, entry := range entries {
		k, v := splitKeyValue(strings.TrimSpace(entry))
		switch k {
		case "name":
			c.Name = unquoteScalar(v)
		case "type":
			c.Type = unquoteScalar(v)
		case "comment":
			c.Comment = unquoteScalar(v)
		}
	}
	if c.Name == "" {
		return Column{}, fmt.Errorf("%w: column missing name", ErrSidecarValidation)
	}
	if c.Type == "" {
		return Column{}, fmt.Errorf("%w: column %q missing type", ErrSidecarValidation, c.Name)
	}
	return c, nil
}

// splitFlowMapEntries splits a flow-map body `a: 1, b: "x, y", c: 'p'` into
// individual entries on top-level commas, honoring quotes.
func splitFlowMapEntries(body string) ([]string, error) {
	var entries []string
	var buf strings.Builder
	inSingle, inDouble := false, false
	for i := 0; i < len(body); i++ {
		ch := body[i]
		switch {
		case ch == '\'' && !inDouble:
			inSingle = !inSingle
			buf.WriteByte(ch)
		case ch == '"' && !inSingle:
			inDouble = !inDouble
			buf.WriteByte(ch)
		case ch == ',' && !inSingle && !inDouble:
			if e := strings.TrimSpace(buf.String()); e != "" {
				entries = append(entries, e)
			}
			buf.Reset()
		default:
			buf.WriteByte(ch)
		}
	}
	if e := strings.TrimSpace(buf.String()); e != "" {
		entries = append(entries, e)
	}
	if inSingle || inDouble {
		return nil, fmt.Errorf("unterminated quote in column entry %q", body)
	}
	return entries, nil
}

// parseFlowSeq parses an inline flow-sequence value `[year, month, day]` or
// `["p1", 'p2']` into []string, honoring quotes. A scalar (no brackets) is
// treated as a single-element sequence so callers can pass either form.
func parseFlowSeq(value string) []string {
	v := strings.TrimSpace(value)
	if v == "" {
		return nil
	}
	v = strings.TrimSpace(strings.TrimPrefix(v, "["))
	v = strings.TrimSuffix(v, "]")
	if v == "" {
		return nil
	}
	var items []string
	var buf strings.Builder
	inSingle, inDouble := false, false
	for i := 0; i < len(v); i++ {
		ch := v[i]
		switch {
		case ch == '\'' && !inDouble:
			inSingle = !inSingle
		case ch == '"' && !inSingle:
			inDouble = !inDouble
		case ch == ',' && !inSingle && !inDouble:
			if it := unquoteScalar(strings.TrimSpace(buf.String())); it != "" {
				items = append(items, it)
			}
			buf.Reset()
		default:
			buf.WriteByte(ch)
		}
	}
	if it := unquoteScalar(strings.TrimSpace(buf.String())); it != "" {
		items = append(items, it)
	}
	return items
}

// unquoteScalar strips YAML quoting and trailing inline `# ...` comments from a
// scalar value, preserving `#` and `:` characters that appear inside quotes.
func unquoteScalar(value string) string {
	v := strings.TrimSpace(value)
	if v == "" {
		return ""
	}
	if v[0] == '"' || v[0] == '\'' {
		q := v[0]
		body := v[1:]
		// Find the matching unescaped close quote.
		for j := 0; j < len(body); j++ {
			if body[j] == q {
				return body[:j]
			}
		}
		return body
	}
	// Bare scalar: strip a trailing ` # ...` comment (space required so we
	// don't mangle URLs like s3://bucket#frag).
	if idx := strings.Index(v, " #"); idx >= 0 {
		v = v[:idx]
	}
	return strings.TrimSpace(v)
}

// splitKeyValue splits a "key: value" (or "key:") line on the first top-level
// colon. Keys/values in this schema never contain colons outside quotes.
func splitKeyValue(trimmed string) (key, value string) {
	inSingle, inDouble := false, false
	for i := 0; i < len(trimmed); i++ {
		ch := trimmed[i]
		switch {
		case ch == '\'' && !inDouble:
			inSingle = !inSingle
		case ch == '"' && !inSingle:
			inDouble = !inDouble
		case ch == ':' && !inSingle && !inDouble:
			return strings.TrimSpace(trimmed[:i]), strings.TrimSpace(trimmed[i+1:])
		}
	}
	return strings.TrimSpace(trimmed), ""
}

func leadingSpaces(s string) int {
	n := 0
	for _, r := range s {
		if r == ' ' {
			n++
			continue
		}
		break
	}
	return n
}

// ValidateSidecar compares a parsed Sidecar against a Glue table represented
// as map[string]interface{} (the JSON shape of `aws glue get-table`). It
// returns the list of Mismatch items found; an empty slice means the sidecar
// structurally matches the live table. A malformed glue map (missing
// Table/StorageDescriptor) yields one database/table/s3 mismatch per missing
// block rather than an error, so callers can render a partial report.
func ValidateSidecar(sidecar Sidecar, glueTable map[string]interface{}) []Mismatch {
	var mismatches []Mismatch

	glueDB, _ := glueTable["DatabaseName"].(string)
	if sidecar.Database != glueDB {
		mismatches = append(mismatches, Mismatch{
			Kind:     MismatchDatabase,
			Expected: sidecar.Database,
			Actual:   glueDB,
			Message:  "sidecar database does not match glue DatabaseName",
		})
	}

	tableMap, ok := glueTable["Table"].(map[string]interface{})
	if !ok {
		mismatches = append(mismatches, Mismatch{
			Kind:    MismatchTable,
			Actual:  "(glue Table block absent)",
			Message: "glue get-table result has no Table object",
		})
		return mismatches
	}

	glueTable_, _ := tableMap["Name"].(string)
	if sidecar.Table != glueTable_ {
		mismatches = append(mismatches, Mismatch{
			Kind:     MismatchTable,
			Expected: sidecar.Table,
			Actual:   glueTable_,
			Message:  "sidecar table name does not match glue Table.Name",
		})
	}

	sd, ok := tableMap["StorageDescriptor"].(map[string]interface{})
	if !ok {
		mismatches = append(mismatches, Mismatch{
			Kind:    MismatchS3Location,
			Actual:  "(StorageDescriptor absent)",
			Message: "glue Table has no StorageDescriptor; cannot validate columns or s3 location",
		})
		return mismatches
	}

	glueLoc, _ := sd["Location"].(string)
	if sidecar.S3Location != "" && sidecar.S3Location != glueLoc {
		mismatches = append(mismatches, Mismatch{
			Kind:     MismatchS3Location,
			Expected: sidecar.S3Location,
			Actual:   glueLoc,
			Message:  "sidecar s3_location does not match StorageDescriptor.Location",
		})
	}

	glueColumns := indexColumns(sd["Columns"])
	for _, sc := range sidecar.Columns {
		gc, ok := glueColumns[sc.Name]
		if !ok {
			mismatches = append(mismatches, Mismatch{
				Kind:     MismatchMissingColumn,
				Column:   sc.Name,
				Expected: sc.Type,
				Message:  fmt.Sprintf("sidecar column %q absent from glue table", sc.Name),
			})
			continue
		}
		glueType, _ := gc["Type"].(string)
		if sc.Type != glueType {
			mismatches = append(mismatches, Mismatch{
				Kind:     MismatchTypeMismatch,
				Column:   sc.Name,
				Expected: sc.Type,
				Actual:   glueType,
				Message:  fmt.Sprintf("sidecar column %q type %q != glue type %q", sc.Name, sc.Type, glueType),
			})
		}
	}

	gluePartitions := indexPartitionNames(tableMap["PartitionKeys"])
	sidecarPartitions := sliceToSet(sidecar.Partitions)
	for _, p := range sidecar.Partitions {
		if !gluePartitions[p] {
			mismatches = append(mismatches, Mismatch{
				Kind:    MismatchMissingPartition,
				Column:  p,
				Message: fmt.Sprintf("sidecar partition %q absent from glue PartitionKeys", p),
			})
		}
	}
	for _, p := range partitionKeysSorted(tableMap["PartitionKeys"]) {
		if !sidecarPartitions[p] {
			mismatches = append(mismatches, Mismatch{
				Kind:    MismatchUnexpectedPartition,
				Column:  p,
				Message: fmt.Sprintf("glue partition %q not listed in sidecar", p),
			})
		}
	}
	return mismatches
}

// indexColumns builds a name->{Name,Type,...} map from a Columns block of the
// Glue StorageDescriptor.
func indexColumns(raw interface{}) map[string]map[string]interface{} {
	out := map[string]map[string]interface{}{}
	cols, ok := raw.([]interface{})
	if !ok {
		return out
	}
	for _, c := range cols {
		m, ok := c.(map[string]interface{})
		if !ok {
			continue
		}
		if name, _ := m["Name"].(string); name != "" {
			out[name] = m
		}
	}
	return out
}

// indexPartitionNames builds a set of partition key names from a Glue
// PartitionKeys block.
func indexPartitionNames(raw interface{}) map[string]bool {
	out := map[string]bool{}
	keys, ok := raw.([]interface{})
	if !ok {
		return out
	}
	for _, k := range keys {
		m, ok := k.(map[string]interface{})
		if !ok {
			continue
		}
		if name, _ := m["Name"].(string); name != "" {
			out[name] = true
		}
	}
	return out
}

// partitionKeysSorted returns the partition key names in Glue's declared
// order. The order is preserved so unexpected_partition mismatches remain
// deterministic across runs.
func partitionKeysSorted(raw interface{}) []string {
	keys, ok := raw.([]interface{})
	if !ok {
		return nil
	}
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		m, ok := k.(map[string]interface{})
		if !ok {
			continue
		}
		if name, _ := m["Name"].(string); name != "" {
			out = append(out, name)
		}
	}
	return out
}

func sliceToSet(s []string) map[string]bool {
	out := make(map[string]bool, len(s))
	for _, v := range s {
		out[v] = true
	}
	return out
}
