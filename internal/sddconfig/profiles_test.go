package sddconfig

import (
	"strings"
	"testing"
)

func sampleProfileConfig() Config {
	return Config{AWSProfiles: map[string]string{
		"prd":     "AWSReadFullDat-874970050509",
		"dev":     "aws-tcl-ope-set-cloud-895593169121",
		"usuario": "mi-usuario-profile",
	}}
}

func TestResolveProfileKnownLogical(t *testing.T) {
	cfg := sampleProfileConfig()
	if got := ResolveProfile(cfg, "prd"); got != "AWSReadFullDat-874970050509" {
		t.Fatalf("ResolveProfile(prd) = %q", got)
	}
	if got := ResolveProfile(cfg, "dev"); got != "aws-tcl-ope-set-cloud-895593169121" {
		t.Fatalf("ResolveProfile(dev) = %q", got)
	}
}

func TestResolveProfileUnknownLogicalReturnsEmpty(t *testing.T) {
	cfg := sampleProfileConfig()
	if got := ResolveProfile(cfg, "staging"); got != "" {
		t.Fatalf("ResolveProfile(staging) = %q, want empty", got)
	}
}

func TestResolveProfileEmptyConfigReturnsEmpty(t *testing.T) {
	// Backward-compat: app-dev has no profiles.
	if got := ResolveProfile(Config{}, "prd"); got != "" {
		t.Fatalf("ResolveProfile(prd) on empty config = %q, want empty", got)
	}
}

func TestScrubProfilesReplacesProfileNames(t *testing.T) {
	cfg := sampleProfileConfig()
	text := "Running job with profile AWSReadFullDat-874970050509 for dev aws-tcl-ope-set-cloud-895593169121"
	out := ScrubProfiles(text, cfg)

	if strings.Contains(out, "AWSReadFullDat-874970050509") {
		t.Fatalf("scrubbed text still contains prd profile name:\n%s", out)
	}
	if strings.Contains(out, "aws-tcl-ope-set-cloud-895593169121") {
		t.Fatalf("scrubbed text still contains dev profile name:\n%s", out)
	}
	if !strings.Contains(out, "<aws-profile:prd>") || !strings.Contains(out, "<aws-profile:dev>") {
		t.Fatalf("scrubbed text missing placeholders:\n%s", out)
	}
}

func TestScrubProfilesReplacesStandaloneAccountIDs(t *testing.T) {
	cfg := sampleProfileConfig()
	text := "account 874970050509 deployed stack"
	out := ScrubProfiles(text, cfg)
	if strings.Contains(out, "874970050509") {
		t.Fatalf("scrubbed text still contains account id:\n%s", out)
	}
	if !strings.Contains(out, "<account-id>") {
		t.Fatalf("scrubbed text missing <account-id> placeholder:\n%s", out)
	}
}

func TestScrubProfilesReplacesProfileAndStandaloneAccountID(t *testing.T) {
	// Triangulation: profile name (which embeds an account id) and a separate
	// standalone account id both get scrubbed.
	cfg := sampleProfileConfig()
	text := "prd=AWSReadFullDat-874970050509 target account 895593169121"
	out := ScrubProfiles(text, cfg)
	if strings.Contains(out, "AWSReadFullDat") || strings.Contains(out, "895593169121") {
		t.Fatalf("scrubbed text still contains sensitive data:\n%s", out)
	}
	if !strings.Contains(out, "<aws-profile:prd>") || !strings.Contains(out, "<account-id>") {
		t.Fatalf("scrubbed text missing placeholders:\n%s", out)
	}
}

func TestScrubProfilesDoesNotMatchThirteenDigitNumber(t *testing.T) {
	// A 13-digit run is NOT a 12-digit account id and must be left untouched.
	cfg := sampleProfileConfig()
	text := "ref=1234567890123"
	out := ScrubProfiles(text, cfg)
	if strings.Contains(out, "<account-id>") {
		t.Fatalf("13-digit number wrongly scrubbed:\n%s", out)
	}
	if !strings.Contains(out, "1234567890123") {
		t.Fatalf("13-digit number altered:\n%s", out)
	}
}

func TestScrubProfilesEmptyConfigStillScrubsAccountIDs(t *testing.T) {
	// Backward-compat: even app-dev (no profiles) scrubs bare account ids.
	text := "account 874970050509"
	out := ScrubProfiles(text, Config{})
	if strings.Contains(out, "874970050509") {
		t.Fatalf("scrubbed text still contains account id:\n%s", out)
	}
	if !strings.Contains(out, "<account-id>") {
		t.Fatalf("scrubbed text missing <account-id> placeholder:\n%s", out)
	}
}

func TestScrubProfilesLeavesInnocuousTextUntouched(t *testing.T) {
	cfg := sampleProfileConfig()
	text := "the quick brown fox jumps over the lazy dog"
	if out := ScrubProfiles(text, cfg); out != text {
		t.Fatalf("innocuous text altered: %q", out)
	}
}
