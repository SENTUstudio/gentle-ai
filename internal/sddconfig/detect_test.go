package sddconfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeMarkerFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll(%s): %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%s): %v", path, err)
	}
}

func TestDetectDomainBothMarkersHighConfidence(t *testing.T) {
	root := t.TempDir()
	writeMarkerFile(t, filepath.Join(root, "template.yaml"), "AWSTemplateFormatVersion: ...\n")
	writeMarkerFile(t, filepath.Join(root, "glue-jobs", "etl_encuestas.py"), "# glue job\n")

	domain, confidence, evidence, err := DetectDomain(root)
	if err != nil {
		t.Fatalf("DetectDomain() error = %v", err)
	}
	if domain != "data-engineering" {
		t.Fatalf("domain = %q, want data-engineering", domain)
	}
	if confidence != 0.8 {
		t.Fatalf("confidence = %v, want 0.8", confidence)
	}
	if !containsEvidence(evidence, "template.yaml") || !containsEvidence(evidence, "glue-jobs") {
		t.Fatalf("evidence = %v, want both markers", evidence)
	}
}

func TestDetectDomainOnlyTemplateIsLowConfidenceHint(t *testing.T) {
	root := t.TempDir()
	writeMarkerFile(t, filepath.Join(root, "template.yaml"), "AWSTemplateFormatVersion: ...\n")

	domain, confidence, evidence, err := DetectDomain(root)
	if err != nil {
		t.Fatalf("DetectDomain() error = %v", err)
	}
	if domain != "data-engineering" {
		t.Fatalf("domain = %q, want data-engineering (hint)", domain)
	}
	if confidence != 0.5 {
		t.Fatalf("confidence = %v, want 0.5", confidence)
	}
	if !containsEvidence(evidence, "template.yaml") {
		t.Fatalf("evidence = %v, want template marker", evidence)
	}
	if containsEvidence(evidence, "glue-jobs") {
		t.Fatalf("evidence = %v, did not expect glue marker", evidence)
	}
}

func TestDetectDomainOnlyGlueJobsIsLowConfidenceHint(t *testing.T) {
	// Triangulation: the other single-marker path.
	root := t.TempDir()
	writeMarkerFile(t, filepath.Join(root, "glue-jobs", "job.py"), "# glue job\n")

	domain, confidence, evidence, err := DetectDomain(root)
	if err != nil {
		t.Fatalf("DetectDomain() error = %v", err)
	}
	if domain != "data-engineering" {
		t.Fatalf("domain = %q, want data-engineering (hint)", domain)
	}
	if confidence != 0.5 {
		t.Fatalf("confidence = %v, want 0.5", confidence)
	}
	if !containsEvidence(evidence, "glue-jobs") {
		t.Fatalf("evidence = %v, want glue marker", evidence)
	}
}

func TestDetectDomainNoMarkersIsEmptyAndZeroConfidence(t *testing.T) {
	root := t.TempDir()

	domain, confidence, evidence, err := DetectDomain(root)
	if err != nil {
		t.Fatalf("DetectDomain() error = %v", err)
	}
	if domain != "" {
		t.Fatalf("domain = %q, want empty (app-dev)", domain)
	}
	if confidence != 0 {
		t.Fatalf("confidence = %v, want 0", confidence)
	}
	if len(evidence) != 0 {
		t.Fatalf("evidence = %v, want empty", evidence)
	}
}

func TestDetectDomainGlueJobsDirWithoutPythonIsNotAMarker(t *testing.T) {
	// A glue-jobs/ directory with no .py files does not count as a marker.
	root := t.TempDir()
	writeMarkerFile(t, filepath.Join(root, "glue-jobs", "README.md"), "docs\n")

	domain, confidence, _, err := DetectDomain(root)
	if err != nil {
		t.Fatalf("DetectDomain() error = %v", err)
	}
	if domain != "" {
		t.Fatalf("domain = %q, want empty", domain)
	}
	if confidence != 0 {
		t.Fatalf("confidence = %v, want 0", confidence)
	}
}

func TestDetectDomainNonexistentRootReturnsError(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "does-not-exist")
	_, _, _, err := DetectDomain(missing)
	if err == nil {
		t.Fatal("DetectDomain() expected error for nonexistent root")
	}
}

func TestDetectDomainRecursiveMasterProjectPattern(t *testing.T) {
	// Master project pattern: markers live inside repositorios/carga-.../glue-jobs/
	// and repositorios/infra-.../template.yaml — not at root.
	root := t.TempDir()
	writeMarkerFile(t, filepath.Join(root, "repositorios", "infra-datos-foo", "template.yaml"), "AWSTemplateFormatVersion: ...\n")
	writeMarkerFile(t, filepath.Join(root, "repositorios", "carga-datos-foo", "glue-jobs", "etl_foo.py"), "# glue job\n")

	domain, confidence, evidence, err := DetectDomain(root)
	if err != nil {
		t.Fatalf("DetectDomain() error = %v", err)
	}
	if domain != "data-engineering" {
		t.Fatalf("domain = %q, want data-engineering", domain)
	}
	if confidence != 0.8 {
		t.Fatalf("confidence = %v, want 0.8 (both markers found recursively)", confidence)
	}
	if !containsEvidence(evidence, "template.yaml") || !containsEvidence(evidence, "glue-jobs") {
		t.Fatalf("evidence = %v, want both markers", evidence)
	}
}

func TestDetectDomainRecursiveNestedGlueJobs(t *testing.T) {
	// Encuestas pattern: markers live inside migracion/PAP/carga-.../glue-jobs/
	root := t.TempDir()
	writeMarkerFile(t, filepath.Join(root, "migracion", "PAP", "carga-datos-encuestas", "glue-jobs", "ETL_encuestas.py"), "# glue job\n")
	writeMarkerFile(t, filepath.Join(root, "migracion", "PAP", "infra-datos-encuestas", "template.yaml"), "AWSTemplateFormatVersion: ...\n")

	domain, confidence, _, err := DetectDomain(root)
	if err != nil {
		t.Fatalf("DetectDomain() error = %v", err)
	}
	if domain != "data-engineering" {
		t.Fatalf("domain = %q, want data-engineering", domain)
	}
	if confidence != 0.8 {
		t.Fatalf("confidence = %v, want 0.8", confidence)
	}
}

func containsEvidence(evidence []string, needle string) bool {
	for _, e := range evidence {
		if strings.Contains(e, needle) {
			return true
		}
	}
	return false
}
