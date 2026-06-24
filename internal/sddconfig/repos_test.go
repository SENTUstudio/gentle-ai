package sddconfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateReposAllExistNoWarnings(t *testing.T) {
	root := t.TempDir()
	infra := filepath.Join(root, "repositorios", "infra-datmos")
	carga := filepath.Join(root, "repositorios", "carga-datmos")
	mkdirs(t, infra, carga)

	cfg := Config{Repos: Repos{
		Infra: "./repositorios/infra-datmos",
		Carga: "./repositorios/carga-datmos",
	}}
	warnings := ValidateRepos(cfg, root)
	if len(warnings) != 0 {
		t.Fatalf("warnings = %v, want none", warnings)
	}
}

func TestValidateReposInfraMissingReturnsWarning(t *testing.T) {
	root := t.TempDir()
	mkdirs(t, filepath.Join(root, "repositorios", "carga-datmos"))

	cfg := Config{Repos: Repos{
		Infra: "./repositorios/infra-datmos",
		Carga: "./repositorios/carga-datmos",
	}}
	warnings := ValidateRepos(cfg, root)
	if len(warnings) != 1 {
		t.Fatalf("warnings = %v, want 1", warnings)
	}
	if warnings[0].Label != "infra" {
		t.Fatalf("Label = %q, want infra", warnings[0].Label)
	}
	if !strings.Contains(warnings[0].Path, "infra-datmos") {
		t.Fatalf("Path = %q, want it to contain infra-datmos", warnings[0].Path)
	}
	if warnings[0].Message == "" {
		t.Fatal("Message empty, want a human-readable reason")
	}
}

func TestValidateReposBothMissingReturnsTwoWarnings(t *testing.T) {
	// Triangulation: both paths missing -> both warned, in stable order (infra, carga).
	root := t.TempDir()
	cfg := Config{Repos: Repos{
		Infra: "./repositorios/infra-datmos",
		Carga: "./repositorios/carga-datmos",
	}}
	warnings := ValidateRepos(cfg, root)
	if len(warnings) != 2 {
		t.Fatalf("warnings = %v, want 2", warnings)
	}
	if warnings[0].Label != "infra" || warnings[1].Label != "carga" {
		t.Fatalf("order = [%s, %s], want [infra, carga]", warnings[0].Label, warnings[1].Label)
	}
}

func TestValidateReposEmptyConfigNoWarnings(t *testing.T) {
	// Backward-compat: app-dev config has no repos -> nothing to validate.
	root := t.TempDir()
	warnings := ValidateRepos(Config{}, root)
	if len(warnings) != 0 {
		t.Fatalf("warnings = %v, want none for empty config", warnings)
	}
}

func TestValidateReposAbsoluteExistingPathNoWarning(t *testing.T) {
	// Triangulation: an absolute configured path that exists must not warn.
	root := t.TempDir()
	abs := t.TempDir()
	mkdirs(t, abs)

	cfg := Config{Repos: Repos{Infra: abs}}
	warnings := ValidateRepos(cfg, root)
	if len(warnings) != 0 {
		t.Fatalf("warnings = %v, want none for existing absolute path", warnings)
	}
}

func TestValidateReposCargaMissingInfraExists(t *testing.T) {
	root := t.TempDir()
	mkdirs(t, filepath.Join(root, "repositorios", "infra-datmos"))

	cfg := Config{Repos: Repos{
		Infra: "./repositorios/infra-datmos",
		Carga: "./repositorios/carga-datmos",
	}}
	warnings := ValidateRepos(cfg, root)
	if len(warnings) != 1 {
		t.Fatalf("warnings = %v, want 1", warnings)
	}
	if warnings[0].Label != "carga" {
		t.Fatalf("Label = %q, want carga", warnings[0].Label)
	}
}

func mkdirs(t *testing.T, dirs ...string) {
	t.Helper()
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatalf("MkdirAll(%s): %v", d, err)
		}
	}
}
