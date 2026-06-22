package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gentleman-programming/gentle-ai/internal/sddconfig"
)

func writeCLIConfig(t *testing.T, root, content string) {
	t.Helper()
	dir := filepath.Join(root, "openspec")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll(%s): %v", dir, err)
	}
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
}

func writeCLIFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll(%s): %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%s): %v", path, err)
	}
}

const dataEngConfig = `domain: data-engineering
repos:
  infra: ./repositorios/infra-datmos
  carga: ./repositorios/carga-datmos
aws_profiles:
  prd: AWSReadFullDat-874970050509
  dev: aws-tcl-ope-set-cloud-895593169121
verify:
  skip_deploy: false
`

func TestRunSDDConfigPrintsMarkdownForDataEngConfig(t *testing.T) {
	root := t.TempDir()
	writeCLIConfig(t, root, dataEngConfig)

	var stdout bytes.Buffer
	if err := RunSDDConfig([]string{"--cwd", root}, &stdout); err != nil {
		t.Fatalf("RunSDDConfig() error = %v", err)
	}
	out := stdout.String()
	for _, want := range []string{"## SDD Config", "data-engineering"} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q:\n%s", want, out)
		}
	}
}

func TestRunSDDConfigJSONDecodesConfig(t *testing.T) {
	root := t.TempDir()
	writeCLIConfig(t, root, dataEngConfig)

	var stdout bytes.Buffer
	if err := RunSDDConfig([]string{"--cwd", root, "--json"}, &stdout); err != nil {
		t.Fatalf("RunSDDConfig() error = %v", err)
	}
	var cfg sddconfig.Config
	if err := json.Unmarshal(stdout.Bytes(), &cfg); err != nil {
		t.Fatalf("JSON decode error = %v\n%s", err, stdout.String())
	}
	if cfg.Domain != "data-engineering" {
		t.Fatalf("Domain = %q, want data-engineering", cfg.Domain)
	}
	if cfg.Repos.Infra != "./repositorios/infra-datmos" {
		t.Fatalf("Repos.Infra = %q", cfg.Repos.Infra)
	}
	if got := cfg.AWSProfiles["prd"]; got != "AWSReadFullDat-874970050509" {
		t.Fatalf("AWSProfiles[prd] = %q", got)
	}
}

func TestRunSDDConfigAppDevJSONIsEmptyConfig(t *testing.T) {
	// Backward-compat: no config.yaml -> zero Config, no error.
	root := t.TempDir()
	var stdout bytes.Buffer
	if err := RunSDDConfig([]string{"--cwd", root, "--json"}, &stdout); err != nil {
		t.Fatalf("RunSDDConfig() error = %v", err)
	}
	var cfg sddconfig.Config
	if err := json.Unmarshal(stdout.Bytes(), &cfg); err != nil {
		t.Fatalf("JSON decode error = %v\n%s", err, stdout.String())
	}
	if cfg.Domain != "" {
		t.Fatalf("Domain = %q, want empty (app-dev)", cfg.Domain)
	}
}

func TestRunSDDConfigDetectBothMarkersJSON(t *testing.T) {
	root := t.TempDir()
	writeCLIFile(t, filepath.Join(root, "template.yaml"), "AWSTemplateFormatVersion: ...\n")
	writeCLIFile(t, filepath.Join(root, "glue-jobs", "etl.py"), "# glue job\n")

	var stdout bytes.Buffer
	if err := RunSDDConfig([]string{"--cwd", root, "--detect", "--json"}, &stdout); err != nil {
		t.Fatalf("RunSDDConfig() error = %v", err)
	}
	var report sddconfig.DetectionReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("JSON decode error = %v\n%s", err, stdout.String())
	}
	if report.Domain != "data-engineering" {
		t.Fatalf("Domain = %q, want data-engineering", report.Domain)
	}
	if report.Confidence != 0.8 {
		t.Fatalf("Confidence = %v, want 0.8", report.Confidence)
	}
	if len(report.Evidence) != 2 {
		t.Fatalf("Evidence = %v, want 2 entries", report.Evidence)
	}
}

func TestRunSDDConfigDetectMarkdownContainsDomainAndConfidence(t *testing.T) {
	root := t.TempDir()
	writeCLIFile(t, filepath.Join(root, "template.yaml"), "AWSTemplateFormatVersion: ...\n")
	writeCLIFile(t, filepath.Join(root, "glue-jobs", "etl.py"), "# glue job\n")

	var stdout bytes.Buffer
	if err := RunSDDConfig([]string{"--cwd", root, "--detect"}, &stdout); err != nil {
		t.Fatalf("RunSDDConfig() error = %v", err)
	}
	out := stdout.String()
	if !strings.Contains(out, "data-engineering") || !strings.Contains(out, "0.8") {
		t.Fatalf("detect markdown missing domain/confidence:\n%s", out)
	}
}

func TestRunSDDConfigValidateReposJSONReportsMissingPath(t *testing.T) {
	root := t.TempDir()
	writeCLIConfig(t, root, dataEngConfig) // repos point at non-existent paths

	var stdout bytes.Buffer
	if err := RunSDDConfig([]string{"--cwd", root, "--validate-repos", "--json"}, &stdout); err != nil {
		t.Fatalf("RunSDDConfig() error = %v", err)
	}
	var report sddconfig.ValidationReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("JSON decode error = %v\n%s", err, stdout.String())
	}
	if len(report.Warnings) != 2 {
		t.Fatalf("Warnings = %v, want 2 (infra + carga missing)", report.Warnings)
	}
}

func TestRunSDDConfigValidateReposMarkdownMentionsWarning(t *testing.T) {
	root := t.TempDir()
	writeCLIConfig(t, root, dataEngConfig)

	var stdout bytes.Buffer
	if err := RunSDDConfig([]string{"--cwd", root, "--validate-repos"}, &stdout); err != nil {
		t.Fatalf("RunSDDConfig() error = %v", err)
	}
	out := stdout.String()
	if !strings.Contains(out, "does not exist") {
		t.Fatalf("validate markdown missing warning:\n%s", out)
	}
}

func TestRunSDDConfigRejectsNonexistentCWD(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope")
	var stdout bytes.Buffer
	if err := RunSDDConfig([]string{"--cwd", missing}, &stdout); err == nil {
		t.Fatal("RunSDDConfig() expected error for nonexistent cwd")
	}
}

func TestRunSDDConfigRejectsUnknownFlag(t *testing.T) {
	root := t.TempDir()
	var stdout bytes.Buffer
	if err := RunSDDConfig([]string{"--cwd", root, "--bogus"}, &stdout); err == nil {
		t.Fatal("RunSDDConfig() expected error for unknown flag")
	}
}
