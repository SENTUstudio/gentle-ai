package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gentleman-programming/gentle-ai/internal/sddconfig"
)

// RunSDDConfig is the CLI entry point for `gentle-ai sdd-config [--json]
// [--detect] [--validate-repos] [--cwd <dir>]`.
//
// It mirrors cli.RunSDDStatus: parse flags, resolve the workspace root, run the
// requested sddconfig operation, and render JSON or markdown. The pure logic
// lives in internal/sddconfig so every branch point is unit-testable.
func RunSDDConfig(args []string, stdout io.Writer) error {
	parsed, err := sddconfig.ParseCommandArgs(args)
	if err != nil {
		return err
	}
	root, err := resolveConfigRoot(parsed.CWD)
	if err != nil {
		return err
	}

	switch {
	case parsed.Detect:
		domain, confidence, evidence, derr := sddconfig.DetectDomain(root)
		if derr != nil {
			return fmt.Errorf("detect domain: %w", derr)
		}
		report := sddconfig.DetectionReport{
			Domain:     domain,
			Confidence: confidence,
			Evidence:   evidence,
		}
		return renderOrEncode(stdout, parsed.JSON, sddconfig.RenderDetectionMarkdown(report), report)

	case parsed.ValidateRepos:
		cfg, lerr := sddconfig.LoadConfig(root)
		if lerr != nil {
			return fmt.Errorf("load config: %w", lerr)
		}
		report := sddconfig.ValidationReport{
			Warnings: sddconfig.ValidateRepos(cfg, root),
		}
		return renderOrEncode(stdout, parsed.JSON, sddconfig.RenderValidationMarkdown(report), report)

	default:
		cfg, lerr := sddconfig.LoadConfig(root)
		if lerr != nil {
			return fmt.Errorf("load config: %w", lerr)
		}
		return renderOrEncode(stdout, parsed.JSON, sddconfig.RenderConfigMarkdown(cfg), cfg)
	}
}

// renderOrEncode writes markdown for human output or indented JSON for machine
// output, collapsing the shared tail of the three sdd-config branches.
func renderOrEncode(stdout io.Writer, asJSON bool, markdown string, v any) error {
	if asJSON {
		return encodeSDDConfigJSON(stdout, v)
	}
	_, err := fmt.Fprintln(stdout, markdown)
	return err
}

// resolveConfigRoot turns the --cwd value into an absolute, existing directory.
// An empty --cwd falls back to the process working directory, mirroring
// sddstatus.absOrCWD. A missing or non-directory root is an error.
func resolveConfigRoot(cwd string) (string, error) {
	root := strings.TrimSpace(cwd)
	if root == "" {
		wd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("resolve cwd: %w", err)
		}
		root = wd
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("workspace root is not a directory: %s", abs)
	}
	return abs, nil
}

// encodeSDDConfigJSON writes v as indented JSON, mirroring sddstatus's encoder.
func encodeSDDConfigJSON(stdout io.Writer, v any) error {
	encoder := json.NewEncoder(stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}