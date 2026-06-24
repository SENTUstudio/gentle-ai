package etl

import (
	"fmt"
	"strings"
)

// Pattern labels supported by DetectPattern. Empty Pattern == unknown.
const (
	PatternIncremental    Pattern = "incremental"
	PatternMultiStep      Pattern = "multi-step"
	PatternLegacyWrangler Pattern = "legacy-wrangler"
	PatternGlueStudio     Pattern = "glue-studio"
)

// Pattern is the ETL implementation taxonomy returned by DetectPattern.
// Possible values are the constants above; the zero value "" means unknown.
type Pattern string

// Markers is the set of file-content flags a caller derives from scanning a
// Glue job's `.py` source. DetectPattern consumes a Markers value and never
// touches the filesystem, so the heuristics remain table-testable.
type Markers struct {
	// HasWranglerAthena: source uses awswrangler.athena.* reads/queries
	// (legacy pattern that predates the Glue-Context migration).
	HasWranglerAthena bool
	// HasApplyMapping: source uses ApplyMapping mappings (a Glue Studio node
	// tell). The 4th data-engineer pattern is GUI-generated Glue Studio jobs.
	HasApplyMapping bool
	// HasSparkSqlQueryHelper: source uses the Glue Studio Spark Sql Query Helper.
	HasSparkSqlQueryHelper bool
	// HasTempViews: source creates temporary SQL views chained across stages
	// (the multi-step pattern's signature).
	HasTempViews bool
	// HasWatermark: source uses an incremental watermark column pattern
	// (the incremental pattern's signature).
	HasWatermark bool
	// HasBoto3S3: source uses boto3.client('s3') direct calls, paired with
	// the watermark pattern in the incremental pipeline.
	HasBoto3S3 bool
	// HasGlueContext: source builds a GlueContext. Absence distinguishes the
	// legacy Wrangler pattern from the Glue-Context-based ones.
	HasGlueContext bool
}

// Confidence values per pattern. The 4-pattern taxonomy mirrors the
// observations captured for the Toyota Chile data-engineering domain; see
// availability design notes. A caller MUST surface the confidence alongside
// the pattern (never report the pattern silently).
const (
	confIncremental    = 0.85
	confMultiStep      = 0.80
	confLegacyWrangler = 0.85
	confGlueStudio     = 0.90
)

// confidenceFor maps a Pattern to its base confidence; helper so DetectPattern
// stays a single dispatch table.
var confidenceFor = map[Pattern]float64{
	PatternIncremental:    confIncremental,
	PatternMultiStep:      confMultiStep,
	PatternLegacyWrangler: confLegacyWrangler,
	PatternGlueStudio:     confGlueStudio,
}

// DetectPattern applies the 4-pattern ETL taxonomy over a Markers scan of a
// Glue job's source and returns:
//
//   - the detected Pattern (zero value "" == unknown)
//   - a confidence score (the higher the more specific) — derived from the
//     Markers that uniquely satisfy each pattern's required set
//   - a human-readable rationale explaining the verdict, including an
//     "ambiguous" note when more than one pattern matched and we picked the
//     highest-confidence one
//
// Heuristic rules:
//
//   - incremental:    HasWatermark && HasBoto3S3 && HasGlueContext              (0.85)
//   - multi-step:      HasTempViews && HasGlueContext && !HasWranglerAthena     (0.80)
//   - legacy-wrangler: HasWranglerAthena && !HasGlueContext                     (0.85)
//   - glue-studio:    HasApplyMapping && HasSparkSqlQueryHelper                 (0.90)
//
// When multiple patterns match, the highest-confidence one wins and the
// rationale flags the ambiguity. The four rules are mutually reasoned (each
// is checked against its full required marker set), so a single partial
// marker never produces a false match — unknown wins instead, but the
// rationale names the partial signal so the caller can surface a hint.
func DetectPattern(m Markers) (Pattern, float64, string) {
	candidates := []struct {
		pattern   Pattern
		matched   bool
		rationale string
	}{
		{
			pattern:   PatternIncremental,
			matched:   m.HasWatermark && m.HasBoto3S3 && m.HasGlueContext,
			rationale: "watermark column + boto3 s3 + GlueContext — matches the incremental pipeline signature",
		},
		{
			pattern:   PatternMultiStep,
			matched:   m.HasTempViews && m.HasGlueContext && !m.HasWranglerAthena,
			rationale: "chained temp views + GlueContext without awswrangler.athena — matches the multi-step pattern",
		},
		{
			pattern:   PatternLegacyWrangler,
			matched:   m.HasWranglerAthena && !m.HasGlueContext,
			rationale: "awswrangler.athena usage without a GlueContext — matches the legacy wrangler pattern",
		},
		{
			pattern:   PatternGlueStudio,
			matched:   m.HasApplyMapping && m.HasSparkSqlQueryHelper,
			rationale: "ApplyMapping + SparkSqlQueryHelper — matches a Glue Studio generated job",
		},
	}

	var hits []Pattern
	for _, c := range candidates {
		if c.matched {
			hits = append(hits, c.pattern)
		}
	}

	// Unknown: no full pattern satisfied. Explain which partial signals were
	// present so the caller can decide whether to override.
	if len(hits) == 0 {
		why := unknownRationale(m)
		return "", 0, why
	}

	// Single hit: report it directly. Multiple hits: pick the highest
	// confidence (ties broken by declared order above).
	best := hits[0]
	for _, p := range hits[1:] {
		if confidenceFor[p] > confidenceFor[best] {
			best = p
		}
	}

	why := rationaleFor(best)
	if len(hits) > 1 {
		why = fmt.Sprintf("ambiguous: patterns %v matched; classified as %s (highest confidence). %s", hits, best, why)
	}
	return best, confidenceFor[best], why
}

// rationaleFor returns the base rationale string for a matched pattern.
func rationaleFor(p Pattern) string {
	switch p {
	case PatternIncremental:
		return "watermark column + boto3 s3 + GlueContext — matches the incremental pipeline signature"
	case PatternMultiStep:
		return "chained temp views + GlueContext without awswrangler.athena — matches the multi-step pattern"
	case PatternLegacyWrangler:
		return "awswrangler.athena usage without a GlueContext — matches the legacy wrangler pattern"
	case PatternGlueStudio:
		return "ApplyMapping + SparkSqlQueryHelper — matches a Glue Studio generated job"
	default:
		return ""
	}
}

// unknownRationale describes the partial markers that were observed without
// fully satisfying any pattern, so the caller can surface a useful hint
// instead of a silent unknown.
func unknownRationale(m Markers) string {
	var signals []string
	if m.HasWatermark {
		signals = append(signals, "watermark (incremental partial)")
	}
	if m.HasBoto3S3 {
		signals = append(signals, "boto3 s3 (incremental partial)")
	}
	if m.HasTempViews {
		signals = append(signals, "temp views (multi-step partial)")
	}
	if m.HasWranglerAthena {
		signals = append(signals, "awswrangler.athena (legacy partial — blocked by present GlueContext)")
	}
	if m.HasApplyMapping {
		signals = append(signals, "ApplyMapping (glue-studio partial)")
	}
	if m.HasSparkSqlQueryHelper {
		signals = append(signals, "SparkSqlQueryHelper (glue-studio partial)")
	}
	if len(signals) == 0 {
		return "unknown: no recognizable ETL pattern markers present"
	}
	return fmt.Sprintf("unknown: no pattern fully satisfied; partial signals observed: %s", strings.Join(signals, ", "))
}
