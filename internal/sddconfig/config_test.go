package sddconfig

import (
	"os"
	"path/filepath"
	"testing"
)

// writeConfig writes the given content to <root>/openspec/config.yaml, creating
// intermediate directories as needed. Returns the root for chaining.
func writeConfig(t *testing.T, root, content string) {
	t.Helper()
	dir := filepath.Join(root, "openspec")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll(%s): %v", dir, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(config.yaml): %v", err)
	}
}

func TestParseConfigAllFields(t *testing.T) {
	text := `schema: spec-driven
domain: data-engineering
repos:
  infra: ./repositorios/infra-datmos
  carga: ./repositorios/carga-datmos
aws_profiles:
  prd: AWSReadFullDat-874970050509
  dev: aws-tcl-ope-set-cloud-895593169121
  usuario: mi-usuario-profile
verify:
  skip_deploy: false
`
	cfg, err := parseConfig(text)
	if err != nil {
		t.Fatalf("parseConfig() error = %v", err)
	}
	if cfg.Domain != "data-engineering" {
		t.Fatalf("Domain = %q, want data-engineering", cfg.Domain)
	}
	if cfg.Repos.Infra != "./repositorios/infra-datmos" {
		t.Fatalf("Repos.Infra = %q", cfg.Repos.Infra)
	}
	if cfg.Repos.Carga != "./repositorios/carga-datmos" {
		t.Fatalf("Repos.Carga = %q", cfg.Repos.Carga)
	}
	if got := cfg.AWSProfiles["prd"]; got != "AWSReadFullDat-874970050509" {
		t.Fatalf("AWSProfiles[prd] = %q", got)
	}
	if got := cfg.AWSProfiles["dev"]; got != "aws-tcl-ope-set-cloud-895593169121" {
		t.Fatalf("AWSProfiles[dev] = %q", got)
	}
	if got := cfg.AWSProfiles["usuario"]; got != "mi-usuario-profile" {
		t.Fatalf("AWSProfiles[usuario] = %q", got)
	}
	if cfg.Verify.SkipDeploy != false {
		t.Fatalf("Verify.SkipDeploy = true, want false")
	}
}

func TestParseConfigSkipDeployTrue(t *testing.T) {
	// Triangulation: different value exercises the boolean path.
	cfg, err := parseConfig("verify:\n  skip_deploy: true\n")
	if err != nil {
		t.Fatalf("parseConfig() error = %v", err)
	}
	if !cfg.Verify.SkipDeploy {
		t.Fatal("Verify.SkipDeploy = false, want true")
	}
}

func TestParseConfigMissingFieldsAreZero(t *testing.T) {
	// Backward-compat: only some fields present; the rest stay zero.
	cfg, err := parseConfig("domain: data-engineering\n")
	if err != nil {
		t.Fatalf("parseConfig() error = %v", err)
	}
	if cfg.Domain != "data-engineering" {
		t.Fatalf("Domain = %q", cfg.Domain)
	}
	if cfg.Repos.Infra != "" || cfg.Repos.Carga != "" {
		t.Fatalf("Repos = %+v, want zero", cfg.Repos)
	}
	if cfg.AWSProfiles != nil {
		t.Fatalf("AWSProfiles = %v, want nil", cfg.AWSProfiles)
	}
	if cfg.Verify.SkipDeploy != false {
		t.Fatalf("Verify.SkipDeploy = true, want false")
	}
}

func TestParseConfigSkipsBlockScalarContext(t *testing.T) {
	// The real app-dev config.yaml has a multi-line `context: |` block plus
	// unknown top-level blocks (rules:, testing:). The parser MUST skip those
	// and leave domain-related fields at zero values.
	text := `schema: spec-driven
last_init: "2026-06-21"

context: |
  Tech stack: Go 1.25.10, Bubbletea TUI
  Architecture: cmd/ entrypoints, internal/ packages
  Branch: main

strict_tdd: true

rules:
  proposal:
    - Include rollback plan
  specs:
    - Use Given/When/Then format

testing:
  strict_tdd: true
  runner:
    command: "go test ./..."
`
	cfg, err := parseConfig(text)
	if err != nil {
		t.Fatalf("parseConfig() error = %v", err)
	}
	if cfg.Domain != "" {
		t.Fatalf("Domain = %q, want empty (app-dev)", cfg.Domain)
	}
	if cfg.Repos.Infra != "" || cfg.Repos.Carga != "" {
		t.Fatalf("Repos = %+v, want zero", cfg.Repos)
	}
	if cfg.AWSProfiles != nil {
		t.Fatalf("AWSProfiles = %v, want nil", cfg.AWSProfiles)
	}
	if cfg.Verify.SkipDeploy != false {
		t.Fatalf("Verify.SkipDeploy = true, want false")
	}
}

func TestParseConfigStripsTrailingCommentAndQuotes(t *testing.T) {
	text := "domain: \"data-engineering\" # the domain\n"
	cfg, err := parseConfig(text)
	if err != nil {
		t.Fatalf("parseConfig() error = %v", err)
	}
	if cfg.Domain != "data-engineering" {
		t.Fatalf("Domain = %q, want data-engineering", cfg.Domain)
	}
}

func TestParseConfigEmpty(t *testing.T) {
	cfg, err := parseConfig("")
	if err != nil {
		t.Fatalf("parseConfig() error = %v", err)
	}
	if cfg.Domain != "" || cfg.Repos.Infra != "" || cfg.AWSProfiles != nil {
		t.Fatalf("cfg = %+v, want zero Config", cfg)
	}
}

func TestLoadConfigAllFields(t *testing.T) {
	root := t.TempDir()
	writeConfig(t, root, `domain: data-engineering
repos:
  infra: ./repositorios/infra-datmos
  carga: ./repositorios/carga-datmos
aws_profiles:
  prd: AWSReadFullDat-874970050509
  dev: aws-tcl-ope-set-cloud-895593169121
verify:
  skip_deploy: true
`)
	cfg, err := LoadConfig(root)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if cfg.Domain != "data-engineering" {
		t.Fatalf("Domain = %q", cfg.Domain)
	}
	if cfg.Repos.Infra != "./repositorios/infra-datmos" {
		t.Fatalf("Repos.Infra = %q", cfg.Repos.Infra)
	}
	if cfg.Repos.Carga != "./repositorios/carga-datmos" {
		t.Fatalf("Repos.Carga = %q", cfg.Repos.Carga)
	}
	if got := cfg.AWSProfiles["prd"]; got != "AWSReadFullDat-874970050509" {
		t.Fatalf("AWSProfiles[prd] = %q", got)
	}
	if !cfg.Verify.SkipDeploy {
		t.Fatal("Verify.SkipDeploy = false, want true")
	}
}

func TestLoadConfigMissingFileReturnsZeroAndNoError(t *testing.T) {
	// Backward-compat: a project with no openspec/config.yaml is app-dev.
	root := t.TempDir()
	cfg, err := LoadConfig(root)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v, want nil", err)
	}
	if cfg.Domain != "" || cfg.Repos.Infra != "" || cfg.AWSProfiles != nil {
		t.Fatalf("cfg = %+v, want zero Config", cfg)
	}
}

func TestLoadConfigAppDevConfigReturnsZero(t *testing.T) {
	// The current gentle-ai openspec/config.yaml has no domain fields.
	// Loading it MUST yield a zero Config (app-dev) with no error.
	root := t.TempDir()
	writeConfig(t, root, `schema: spec-driven
strict_tdd: true
context: |
  no domain here
rules:
  apply:
    - Follow existing code patterns
`)
	cfg, err := LoadConfig(root)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if cfg.Domain != "" {
		t.Fatalf("Domain = %q, want empty", cfg.Domain)
	}
	if cfg.AWSProfiles != nil {
		t.Fatalf("AWSProfiles = %v, want nil", cfg.AWSProfiles)
	}
}
