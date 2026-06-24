package sddconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Confidence thresholds for domain auto-detection.
//
// Both markers (template.yaml AND glue-jobs/*.py) -> strong signal (0.8).
// A single marker -> weak hint (0.5) that the caller (sdd-init) presents for
// confirmation/override. Zero markers -> app-dev (no domain).
const (
	confidenceBoth   = 0.8
	confidenceSingle = 0.5
	domainDataEng    = "data-engineering"
)

// DetectDomain scans root for data-engineering markers and returns the detected
// domain, a confidence score, and human-readable evidence.
//
// Markers (scanned recursively up to 4 levels deep):
//   - template.yaml anywhere under root (a SAM/CloudFormation Glue template)
//   - at least one *.py file under a glue-jobs/ directory anywhere under root
//
// Both markers present -> ("data-engineering", 0.8). One marker ->
// ("data-engineering", 0.5) as a hint the caller confirms. None -> ("", 0).
//
// Recursive scan supports the master-project pattern where company repos
// (infra + carga) are cloned under repositorios/ and their markers live
// inside subdirectories like repositorios/carga-datos-xxx/glue-jobs/*.py.
func DetectDomain(root string) (domain string, confidence float64, evidence []string, err error) {
	if _, err = os.Stat(root); err != nil {
		return "", 0, nil, fmt.Errorf("detect domain: stat root: %w", err)
	}

	hasTemplate, templatePath := findMarkerRecursive(root, "template.yaml", 4)
	hasGlue, gluePath := findGlueJobsRecursive(root, 4)

	switch {
	case hasTemplate && hasGlue:
		evidence = append(evidence, fmt.Sprintf("found template.yaml at %s", templatePath),
			fmt.Sprintf("found glue-jobs/*.py at %s", gluePath))
		return domainDataEng, confidenceBoth, evidence, nil
	case hasTemplate:
		evidence = append(evidence, fmt.Sprintf("found template.yaml at %s", templatePath))
		return domainDataEng, confidenceSingle, evidence, nil
	case hasGlue:
		evidence = append(evidence, fmt.Sprintf("found glue-jobs/*.py at %s", gluePath))
		return domainDataEng, confidenceSingle, evidence, nil
	default:
		return "", 0, nil, nil
	}
}

// markerExists reports whether path exists (file or directory).
func markerExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// findMarkerRecursive searches for a file named name anywhere under root,
// up to maxDepth levels deep. Returns (true, relativePath) if found.
// Skips .git, .venv, node_modules, __pycache__ directories.
func findMarkerRecursive(root, name string, maxDepth int) (bool, string) {
	var found bool
	var foundPath string
	filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || found {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		depth := strings.Count(rel, string(filepath.Separator))
		if depth > maxDepth {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		base := filepath.Base(path)
		if base == ".git" || base == ".venv" || base == "node_modules" || base == "__pycache__" || base == ".pytest_cache" {
			if d.IsDir() {
				return filepath.SkipDir
			}
		}
		if !d.IsDir() && base == name {
			found = true
			foundPath = rel
			return filepath.SkipAll
		}
		return nil
	})
	return found, foundPath
}

// findGlueJobsRecursive searches for a glue-jobs/ directory containing .py
// files anywhere under root, up to maxDepth levels deep.
// Returns (true, relativePath) if found.
func findGlueJobsRecursive(root string, maxDepth int) (bool, string) {
	var found bool
	var foundPath string
	filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || found {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		depth := strings.Count(rel, string(filepath.Separator))
		if depth > maxDepth {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		base := filepath.Base(path)
		if base == ".git" || base == ".venv" || base == "node_modules" || base == "__pycache__" || base == ".pytest_cache" {
			if d.IsDir() {
				return filepath.SkipDir
			}
		}
		if d.IsDir() && base == "glue-jobs" {
			if dirHasPython(path) {
				found = true
				foundPath = rel
				return filepath.SkipAll
			}
		}
		return nil
	})
	return found, foundPath
}

// dirHasPython reports whether dir contains a .py file at any depth. A missing
// or unreadable directory yields false, so this also guards the glue-jobs
// marker when the directory does not exist.
func dirHasPython(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if entry.IsDir() {
			if dirHasPython(filepath.Join(dir, entry.Name())) {
				return true
			}
			continue
		}
		if strings.HasSuffix(entry.Name(), ".py") {
			return true
		}
	}
	return false
}

// DetectionReport is the JSON shape of `gentle-ai sdd-config --detect`.
type DetectionReport struct {
	Domain     string   `json:"domain"`
	Confidence float64  `json:"confidence"`
	Evidence   []string `json:"evidence"`
}

// RenderDetectionMarkdown renders a human-readable detection summary for the
// default (non-JSON) sdd-config --detect output.
func RenderDetectionMarkdown(report DetectionReport) string {
	lines := []string{
		"## SDD Config: Detect",
		"",
		fmt.Sprintf("domain: %s", presentOrUnset(report.Domain)),
		fmt.Sprintf("confidence: %s", strconv.FormatFloat(report.Confidence, 'f', 1, 64)),
	}
	if len(report.Evidence) > 0 {
		lines = append(lines, "evidence:")
		for _, e := range report.Evidence {
			lines = append(lines, "- "+e)
		}
	} else {
		lines = append(lines, "evidence: (no markers found)")
	}
	return strings.Join(lines, "\n")
}
