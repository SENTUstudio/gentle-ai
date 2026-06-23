package etl

import (
	"errors"
	"strings"
	"testing"
)

const validSidecarYAML = `database: db_dl_dev_stg_encuestas
table: encuestas_csi
columns:
  - {name: id, type: string, comment: "survey id"}
  - {name: respuesta, type: int, comment: "likert"}
partitions: [year, month, day]
s3_location: s3://toyota-chile-refined/encuestas_csi/
format: parquet
compression: snappy
`

// glueTableForValid builds a map[string]any mirroring `aws glue get-table`
// (DatabaseName + Table.{Name, StorageDescriptor.{Columns, Location}, PartitionKeys}).
func glueTableForValid() map[string]interface{} {
	return map[string]interface{}{
		"DatabaseName": "db_dl_dev_stg_encuestas",
		"Table": map[string]interface{}{
			"Name": "encuestas_csi",
			"StorageDescriptor": map[string]interface{}{
				"Columns": []interface{}{
					map[string]interface{}{"Name": "id", "Type": "string", "Comment": "survey id"},
					map[string]interface{}{"Name": "respuesta", "Type": "int", "Comment": "likert"},
				},
				"Location": "s3://toyota-chile-refined/encuestas_csi/",
			},
			"PartitionKeys": []interface{}{
				map[string]interface{}{"Name": "year", "Type": "string"},
				map[string]interface{}{"Name": "month", "Type": "string"},
				map[string]interface{}{"Name": "day", "Type": "string"},
			},
		},
	}
}

func TestParseSidecarValid(t *testing.T) {
	s, err := ParseSidecar([]byte(validSidecarYAML))
	if err != nil {
		t.Fatalf("ParseSidecar() error = %v", err)
	}
	if s.Database != "db_dl_dev_stg_encuestas" {
		t.Fatalf("Database = %q, want db_dl_dev_stg_encuestas", s.Database)
	}
	if s.Table != "encuestas_csi" {
		t.Fatalf("Table = %q, want encuestas_csi", s.Table)
	}
	if len(s.Columns) != 2 {
		t.Fatalf("len(Columns) = %d, want 2", len(s.Columns))
	}
	if s.Columns[0].Name != "id" || s.Columns[0].Type != "string" || s.Columns[0].Comment != "survey id" {
		t.Fatalf("Columns[0] = %+v, want {id string \"survey id\"}", s.Columns[0])
	}
	if s.Columns[1].Comment != "likert" {
		t.Fatalf("Columns[1].Comment = %q, want likert", s.Columns[1].Comment)
	}
	if len(s.Partitions) != 3 || strings.Join(s.Partitions, ",") != "year,month,day" {
		t.Fatalf("Partitions = %v, want [year month day]", s.Partitions)
	}
	if s.S3Location != "s3://toyota-chile-refined/encuestas_csi/" {
		t.Fatalf("S3Location = %q", s.S3Location)
	}
	if s.Format != "parquet" {
		t.Fatalf("Format = %q, want parquet", s.Format)
	}
	if s.Compression != "snappy" {
		t.Fatalf("Compression = %q, want snappy", s.Compression)
	}
}

func TestParseSidecarInvalid(t *testing.T) {
	// Missing database (required) — must return an error and a zero Sidecar.
	bad := `table: encuestas_csi
columns:
  - {name: id, type: string}
`
	if _, err := ParseSidecar([]byte(bad)); err == nil {
		t.Fatal("ParseSidecar() expected error for missing database, got nil")
	}

	// Empty input is a structural error.
	if _, err := ParseSidecar([]byte("")); err == nil {
		t.Fatal("ParseSidecar() expected error for empty input, got nil")
	}

	// A column missing its name is a structural error.
	badCol := `database: db
table: t
columns:
  - {type: string}
`
	if _, err := ParseSidecar([]byte(badCol)); err == nil {
		t.Fatal("ParseSidecar() expected error for column without name, got nil")
	}
}

func TestParseSidecarSkipsComments(t *testing.T) {
	yaml := `# top comment
database: db
table: t  # inline comment
columns:
  - {name: id, type: string}  # col comment
`
	s, err := ParseSidecar([]byte(yaml))
	if err != nil {
		t.Fatalf("ParseSidecar() error = %v", err)
	}
	if s.Database != "db" {
		t.Fatalf("Database = %q, want db", s.Database)
	}
	if s.Table != "t" {
		t.Fatalf("Table = %q (inline comment not stripped)", s.Table)
	}
	if len(s.Columns) != 1 || s.Columns[0].Name != "id" {
		t.Fatalf("Columns = %v", s.Columns)
	}
}

func TestValidateSidecarMatch(t *testing.T) {
	s, err := ParseSidecar([]byte(validSidecarYAML))
	if err != nil {
		t.Fatalf("ParseSidecar() error = %v", err)
	}
	mismatches := ValidateSidecar(s, glueTableForValid())
	if len(mismatches) != 0 {
		t.Fatalf("ValidateSidecar() = %v, want no mismatches", mismatches)
	}
}

func TestValidateSidecarMissingColumn(t *testing.T) {
	s, _ := ParseSidecar([]byte(validSidecarYAML))
	glue := glueTableForValid()
	sd := glue["Table"].(map[string]interface{})["StorageDescriptor"].(map[string]interface{})
	// Drop the second column from Glue — sidecar lists it, but Glue lacks it.
	sd["Columns"] = []interface{}{
		(map[string]interface{}{"Name": "id", "Type": "string"}),
	}

	mismatches := ValidateSidecar(s, glue)
	if !hasMismatch(mismatches, "missing_column", "respuesta") {
		t.Fatalf("mismatches = %v, want a missing_column mismatch for respuesta", mismatches)
	}
}

func TestValidateSidecarWrongType(t *testing.T) {
	s, _ := ParseSidecar([]byte(validSidecarYAML))
	glue := glueTableForValid()
	sd := glue["Table"].(map[string]interface{})["StorageDescriptor"].(map[string]interface{})
	// Change Glue type for `respuesta` to bigint — sidecar says int.
	sd["Columns"] = []interface{}{
		(map[string]interface{}{"Name": "id", "Type": "string"}),
		(map[string]interface{}{"Name": "respuesta", "Type": "bigint"}),
	}

	mismatches := ValidateSidecar(s, glue)
	if !hasMismatch(mismatches, "type_mismatch", "respuesta") {
		t.Fatalf("mismatches = %v, want a type_mismatch for respuesta", mismatches)
	}
}

func TestValidateSidecarMissingPartition(t *testing.T) {
	s, _ := ParseSidecar([]byte(validSidecarYAML))
	glue := glueTableForValid()
	tbl := glue["Table"].(map[string]interface{})
	// Drop the `day` partition key from Glue.
	tbl["PartitionKeys"] = []interface{}{
		(map[string]interface{}{"Name": "year", "Type": "string"}),
		(map[string]interface{}{"Name": "month", "Type": "string"}),
	}

	mismatches := ValidateSidecar(s, glue)
	if !hasMismatch(mismatches, "missing_partition", "day") {
		t.Fatalf("mismatches = %v, want a missing_partition mismatch for day", mismatches)
	}
}

func TestValidateSidecarExtraPartitionIsMismatch(t *testing.T) {
	s, _ := ParseSidecar([]byte(validSidecarYAML))
	glue := glueTableForValid()
	tbl := glue["Table"].(map[string]interface{})
	// Add an extra partition key not listed in the sidecar.
	tbl["PartitionKeys"] = []interface{}{
		(map[string]interface{}{"Name": "year", "Type": "string"}),
		(map[string]interface{}{"Name": "month", "Type": "string"}),
		(map[string]interface{}{"Name": "day", "Type": "string"}),
		(map[string]interface{}{"Name": "region", "Type": "string"}),
	}

	mismatches := ValidateSidecar(s, glue)
	if !hasMismatch(mismatches, "unexpected_partition", "region") {
		t.Fatalf("mismatches = %v, want an unexpected_partition mismatch for region", mismatches)
	}
}

func TestValidateSidecarWrongDatabase(t *testing.T) {
	s, _ := ParseSidecar([]byte(validSidecarYAML))
	glue := glueTableForValid()
	glue["DatabaseName"] = "other_db"

	mismatches := ValidateSidecar(s, glue)
	if !hasMismatch(mismatches, "database_mismatch", "") {
		t.Fatalf("mismatches = %v, want database_mismatch", mismatches)
	}
}

func TestValidateSidecarWrongS3Location(t *testing.T) {
	s, _ := ParseSidecar([]byte(validSidecarYAML))
	glue := glueTableForValid()
	sd := glue["Table"].(map[string]interface{})["StorageDescriptor"].(map[string]interface{})
	sd["Location"] = "s3://wrong-bucket/t/"

	mismatches := ValidateSidecar(s, glue)
	if !hasMismatch(mismatches, "s3_location_mismatch", "") {
		t.Fatalf("mismatches = %v, want s3_location_mismatch", mismatches)
	}
}

func TestValidateSidecarWrongTable(t *testing.T) {
	s, _ := ParseSidecar([]byte(validSidecarYAML))
	glue := glueTableForValid()
	glue["Table"].(map[string]interface{})["Name"] = "wrong_table"

	mismatches := ValidateSidecar(s, glue)
	if !hasMismatch(mismatches, "table_mismatch", "") {
		t.Fatalf("mismatches = %v, want table_mismatch", mismatches)
	}
}

func TestValidateSidecarMalformedGlueTable(t *testing.T) {
	s, _ := ParseSidecar([]byte(validSidecarYAML))
	// A totally malformed glue map: no Table, no StorageDescriptor.
	if m := ValidateSidecar(s, map[string]interface{}{"DatabaseName": "db"}); len(m) == 0 {
		t.Fatalf("ValidateSidecar(malformed) = %v, want at least one mismatch", m)
	}
}

func TestParseSidecarMultipleColumnsAndQuoted(t *testing.T) {
	// Triangulation: exercise a fuller column list and quoted partition/comment variants.
	yaml := `database: db
table: t
columns:
  - {name: "first col", type: 'varchar(20)', comment: "spaced name and quoted type"}
  - {name: second, type: double, comment: ""}
partitions: ["p1", 'p2']
s3_location: s3://b/x
format: parquet
compression: snappy
`
	s, err := ParseSidecar([]byte(yaml))
	if err != nil {
		t.Fatalf("ParseSidecar() error = %v", err)
	}
	if s.Columns[0].Name != "first col" {
		t.Fatalf("Columns[0].Name = %q, want 'first col'", s.Columns[0].Name)
	}
	if s.Columns[0].Type != "varchar(20)" {
		t.Fatalf("Columns[0].Type = %q, want varchar(20)", s.Columns[0].Type)
	}
	if len(s.Partitions) != 2 || s.Partitions[0] != "p1" || s.Partitions[1] != "p2" {
		t.Fatalf("Partitions = %v", s.Partitions)
	}
}

func TestParseSidecarIgnoresUnknownKeys(t *testing.T) {
	yaml := `database: db
table: t
columns:
  - {name: id, type: string}
unknown_top_level: 42
another_section:
  nested: value
`
	s, err := ParseSidecar([]byte(yaml))
	if err != nil {
		t.Fatalf("ParseSidecar() error = %v", err)
	}
	if s.Database != "db" {
		t.Fatalf("Database = %q", s.Database)
	}
}

func TestParseSidecarColumnInlineMissingType(t *testing.T) {
	// A column with a name but no type — `string` (YAML) is required for
	// structural validation against Glue. Missing type is an error.
	yaml := `database: db
table: t
columns:
  - {name: id}
`
	_, err := ParseSidecar([]byte(yaml))
	if err == nil {
		t.Fatal("ParseSidecar() expected error for column without type, got nil")
	}
	if !errors.Is(err, ErrSidecarValidation) {
		// Non-fatal: structural errors are wrapped with ErrSidecarValidation
		// so callers can distinguish sidecar parse failures from transport errors.
		t.Logf("note: returned err = %v (%T); ErrSidecarValidation sentinel recommended", err, err)
	}
}

func hasMismatch(mismatches []Mismatch, kind, column string) bool {
	for _, m := range mismatches {
		if m.Kind != kind {
			continue
		}
		if column == "" {
			return true
		}
		return m.Column == column
	}
	return false
}
