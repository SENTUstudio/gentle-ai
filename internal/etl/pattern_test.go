package etl

import (
	"strings"
	"testing"
)

func TestDetectPatternIncremental(t *testing.T) {
	m := Markers{HasWatermark: true, HasBoto3S3: true, HasGlueContext: true}
	p, conf, why := DetectPattern(m)
	if p != PatternIncremental {
		t.Fatalf("pattern = %q, want increment", p)
	}
	if conf != 0.85 {
		t.Fatalf("confidence = %v, want 0.85", conf)
	}
	if !strings.Contains(why, "watermark") {
		t.Fatalf("rationale = %q, want watermark mention", why)
	}
}

func TestDetectPatternMultiStep(t *testing.T) {
	// TempViews + GlueContext + !WranglerAthena + !ApplyMapping + !Watermark
	m := Markers{HasTempViews: true, HasGlueContext: true}
	p, conf, why := DetectPattern(m)
	if p != PatternMultiStep {
		t.Fatalf("pattern = %q, want multi-step", p)
	}
	if conf != 0.80 {
		t.Fatalf("confidence = %v, want 0.80", conf)
	}
	if !strings.Contains(why, "temp") {
		t.Fatalf("rationale = %q, want temp-views mention", why)
	}
}

func TestDetectPatternLegacyWrangler(t *testing.T) {
	// HasWranglerAthena + !GlueContext. No SparkSql/ApplyMapping/Watermark.
	m := Markers{HasWranglerAthena: true}
	p, conf, why := DetectPattern(m)
	if p != PatternLegacyWrangler {
		t.Fatalf("pattern = %q, want legacy-wrangler", p)
	}
	if conf != 0.85 {
		t.Fatalf("confidence = %v, want 0.85", conf)
	}
	if !strings.Contains(strings.ToLower(why), "wrangler") {
		t.Fatalf("rationale = %q, want wrangler mention", why)
	}
}

func TestDetectPatternGlueStudio(t *testing.T) {
	// ApplyMapping + SparkSqlQueryHelper — independent of GlueContext.
	m := Markers{HasApplyMapping: true, HasSparkSqlQueryHelper: true}
	p, conf, why := DetectPattern(m)
	if p != PatternGlueStudio {
		t.Fatalf("pattern = %q, want glue-studio", p)
	}
	if conf != 0.90 {
		t.Fatalf("confidence = %v, want 0.90", conf)
	}
	if !strings.Contains(strings.ToLower(why), "applymapping") && !strings.Contains(strings.ToLower(why), "spark") {
		t.Fatalf("rationale = %q, want ApplyMapping/Spark mention", why)
	}
}

func TestDetectPatternUnknownIsZeroConfidence(t *testing.T) {
	// No markers at all.
	m := Markers{}
	p, conf, why := DetectPattern(m)
	if p != "" {
		t.Fatalf("pattern = %q, want empty (unknown)", p)
	}
	if conf != 0 {
		t.Fatalf("confidence = %v, want 0", conf)
	}
	if why == "" {
		t.Fatalf("rationale = empty, want a non-empty explanation for unknown")
	}
}

func TestDetectPatternAmbiguousPrefersHighestConfidence(t *testing.T) {
	// Markers lit in both Glue-Studio (0.90) and Incremental (0.85). Studio
	// has the higher confidence, so it must win; rationale must mention
	// ambiguity and the winning pattern.
	m := Markers{
		HasWatermark:           true,
		HasBoto3S3:             true,
		HasGlueContext:         true,
		HasApplyMapping:        true,
		HasSparkSqlQueryHelper: true,
	}
	p, conf, why := DetectPattern(m)
	if p != PatternGlueStudio {
		t.Fatalf("pattern = %q, want glue-studio (highest confidence)", p)
	}
	if conf != 0.90 {
		t.Fatalf("confidence = %v, want 0.90", conf)
	}
	if !strings.Contains(strings.ToLower(why), "ambig") {
		t.Fatalf("rationale = %q, want to call out ambiguity", why)
	}
}

func TestDetectPatternIncrementalMissingBoto3IsMultiStep(t *testing.T) {
	// Triangulation: incremental needs Watermark+Boto3S3+GlueContext. With
	// Boto3S3 absent, Watermark+GlueContext+TempViews falls into multi-step
	// (HasTempViews + GlueContext + !HasWranglerAthena still win at 0.80).
	m := Markers{HasWatermark: true, HasGlueContext: true, HasTempViews: true}
	p, conf, _ := DetectPattern(m)
	if p != PatternMultiStep {
		t.Fatalf("pattern = %q, want multi-step when boto3 absent", p)
	}
	if conf != 0.80 {
		t.Fatalf("confidence = %v, want 0.80", conf)
	}
}

func TestDetectPatternLegacyWranglerSuppressedByGlueContext(t *testing.T) {
	// Triangulation: legacy-wrangler requires !GlueContext. When GlueContext
	// is also lit alongside WranglerAthena, legacy is suppressed. The
	// multi-step rule ALSO requires !HasWranglerAthena, so this constellation
	// (WranglerAthena + GlueContext + TempViews) does not match any rule —
	// unknown wins, and the rationale must mention the contradiction.
	m := Markers{HasWranglerAthena: true, HasGlueContext: true, HasTempViews: true}
	p, conf, why := DetectPattern(m)
	if p != "" {
		t.Fatalf("pattern = %q, want empty (contradiction wrangler+gluecontext)", p)
	}
	if conf != 0 {
		t.Fatalf("confidence = %v, want 0", conf)
	}
	if !strings.Contains(strings.ToLower(why), "awswrangler.athena") {
		t.Fatalf("rationale = %q, want awswrangler.athena mention", why)
	}
}

func TestDetectPatternMarkersStructIsComparable(t *testing.T) {
	// The Markers struct is the contract a caller fills from a content scan.
	// Equality semantics matter because pattern-detect's caller may diff two
	// scans. Confirm all flags participate in a simple equality check.
	a := Markers{HasWatermark: true, HasGlueContext: true}
	b := a
	c := Markers{HasWatermark: true, HasGlueContext: false}
	if a != b {
		t.Fatalf("Markers equality broken: a==b expected (a=%v b=%v)", a, b)
	}
	if a == c {
		t.Fatalf("Markers equality broken: a!=c expected (a=%v c=%v)", a, c)
	}
}

func TestDetectPatternWatermarkOnlyIsUnknown(t *testing.T) {
	// Triangulation: a lone Watermark partials neither increment (needs all
	// 3 markers) nor multi-step. It must surface as unknown rather than
	// silently hint at increment. The rationale should still be useful (say
	// that watermark was the partial signal).
	m := Markers{HasWatermark: true}
	p, conf, why := DetectPattern(m)
	if p != "" {
		t.Fatalf("pattern = %q, want empty (no full marker set)", p)
	}
	if conf != 0 {
		t.Fatalf("confidence = %v, want 0", conf)
	}
	if !strings.Contains(strings.ToLower(why), "watermark") && !strings.Contains(strings.ToLower(why), "unknown") {
		t.Fatalf("rationale = %q, want partial-signal mention", why)
	}
}
