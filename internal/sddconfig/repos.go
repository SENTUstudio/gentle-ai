package sddconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// PathWarning describes a configured repository path that failed validation.
type PathWarning struct {
	Label   string `json:"label"`   // "infra" | "carga"
	Path    string `json:"path"`    // resolved path that was checked
	Message string `json:"message"` // human-readable reason
}

// ValidateRepos resolves and existence-checks the infra/carga repository paths
// declared in cfg.Repos (relative to root, or absolute if already absolute).
//
// Unconfigured (empty) paths produce no warning — an app-dev project has no
// repos and validates cleanly. Configured-but-missing paths produce one
// PathWarning each, in stable order: infra, then carga.
func ValidateRepos(cfg Config, root string) []PathWarning {
	var warnings []PathWarning
	if w, ok := checkRepoPath("infra", cfg.Repos.Infra, root); ok {
		warnings = append(warnings, w)
	}
	if w, ok := checkRepoPath("carga", cfg.Repos.Carga, root); ok {
		warnings = append(warnings, w)
	}
	return warnings
}

// checkRepoPath returns a PathWarning (and ok=true) when label's path is
// configured but does not exist. An empty path (unconfigured) returns ok=false.
func checkRepoPath(label, configured, root string) (PathWarning, bool) {
	if configured == "" {
		return PathWarning{}, false
	}
	resolved := resolveRepoPath(configured, root)
	if pathExists(resolved) {
		return PathWarning{}, false
	}
	return PathWarning{
		Label:   label,
		Path:    resolved,
		Message: fmt.Sprintf("%s repo path does not exist: %s", label, resolved),
	}, true
}

// resolveRepoPath turns a configured repo path into an absolute-ish path. An
// absolute configured path is used as-is; a relative path is joined to root.
func resolveRepoPath(configured, root string) string {
	if filepath.IsAbs(configured) {
		return filepath.Clean(configured)
	}
	return filepath.Join(root, configured)
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ValidationReport is the JSON shape of `gentle-ai sdd-config --validate-repos`.
type ValidationReport struct {
	Warnings []PathWarning `json:"warnings"`
}

// RenderValidationMarkdown renders a human-readable repo-validation summary for
// the default (non-JSON) sdd-config --validate-repos output.
func RenderValidationMarkdown(report ValidationReport) string {
	lines := []string{
		"## SDD Config: Validate Repos",
		"",
		fmt.Sprintf("warnings: %d", len(report.Warnings)),
	}
	if len(report.Warnings) == 0 {
		lines = append(lines, "(no warnings — all configured repo paths exist)")
	} else {
		for _, w := range report.Warnings {
			lines = append(lines, fmt.Sprintf("- [%s] %s", w.Label, w.Message))
		}
	}
	return strings.Join(lines, "\n")
}
