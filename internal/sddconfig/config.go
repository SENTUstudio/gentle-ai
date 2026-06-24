package sddconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Config is the parsed view of the domain-relevant subset of openspec/config.yaml.
// Every field is optional: a project that does not set `domain` (the app-dev
// default) produces a zero-value Config, and every downstream branch behaves
// exactly as it does today.
type Config struct {
	Domain      string            `json:"domain"`
	Repos       Repos             `json:"repos"`
	AWSProfiles map[string]string `json:"awsProfiles"`
	Verify      VerifyOpts        `json:"verify"`
}

// Repos holds the logical repository paths declared under `repos:` in config.yaml.
type Repos struct {
	Infra string `json:"infra"`
	Carga string `json:"carga"`
}

// VerifyOpts holds verify-phase tunables declared under `verify:` in config.yaml.
type VerifyOpts struct {
	SkipDeploy bool `json:"skipDeploy"`
}

// LoadConfig reads <root>/openspec/config.yaml and returns the parsed Config.
// A missing config file is not an error: it yields a zero Config (app-dev),
// preserving backward compatibility with projects that never set `domain`.
func LoadConfig(root string) (Config, error) {
	data, err := os.ReadFile(filepath.Join(root, "openspec", "config.yaml"))
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, nil
		}
		return Config{}, err
	}
	return parseConfig(string(data))
}

// parseConfig is a pure, hand-rolled reader for the limited config.yaml subset
// this package owns (domain, repos, aws_profiles, verify). The project
// deliberately avoids gopkg.in/yaml.v3, so this scanner is indentation-aware:
// it skips block scalars (`key: |`) and unknown top-level blocks, leaving every
// unrelated key (schema, context, rules, testing, ...) untouched.
func parseConfig(text string) (Config, error) {
	var cfg Config
	lines := strings.Split(text, "\n")
	for i := 0; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if leadingSpaces(lines[i]) != 0 {
			continue // stray indented line outside a managed block; skip defensively
		}
		key, value := splitKeyValue(trimmed)
		switch key {
		case "domain":
			cfg.Domain = unquote(value)
		case "repos":
			i = parseReposBlock(lines, i, &cfg)
		case "aws_profiles":
			i = parseProfilesBlock(lines, i, &cfg)
		case "verify":
			i = parseVerifyBlock(lines, i, &cfg)
		default:
			if isBlockScalarMarker(value) {
				i = skipIndentedUntilTopLevel(lines, i)
			} else if value == "" {
				i = skipIndentedUntilTopLevel(lines, i)
			}
		}
	}
	return cfg, nil
}

// parseReposBlock consumes the indented children following the `repos:` header
// at lines[start]. Returns the index of the last line it consumed.
func parseReposBlock(lines []string, start int, cfg *Config) int {
	last := start
	for j := start + 1; j < len(lines); j++ {
		trimmed := strings.TrimSpace(lines[j])
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			last = j
			continue
		}
		if leadingSpaces(lines[j]) == 0 {
			break
		}
		key, value := splitKeyValue(trimmed)
		switch key {
		case "infra":
			cfg.Repos.Infra = unquote(value)
		case "carga":
			cfg.Repos.Carga = unquote(value)
		}
		last = j
	}
	return last
}

// parseProfilesBlock consumes the indented children following the
// `aws_profiles:` header, populating cfg.AWSProfiles (logical name -> CLI
// profile name).
func parseProfilesBlock(lines []string, start int, cfg *Config) int {
	if cfg.AWSProfiles == nil {
		cfg.AWSProfiles = map[string]string{}
	}
	last := start
	for j := start + 1; j < len(lines); j++ {
		trimmed := strings.TrimSpace(lines[j])
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			last = j
			continue
		}
		if leadingSpaces(lines[j]) == 0 {
			break
		}
		key, value := splitKeyValue(trimmed)
		if key != "" {
			cfg.AWSProfiles[key] = unquote(value)
		}
		last = j
	}
	return last
}

// parseVerifyBlock consumes the indented children following the `verify:`
// header.
func parseVerifyBlock(lines []string, start int, cfg *Config) int {
	last := start
	for j := start + 1; j < len(lines); j++ {
		trimmed := strings.TrimSpace(lines[j])
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			last = j
			continue
		}
		if leadingSpaces(lines[j]) == 0 {
			break
		}
		key, value := splitKeyValue(trimmed)
		if key == "skip_deploy" {
			if b, err := strconv.ParseBool(unquote(value)); err == nil {
				cfg.Verify.SkipDeploy = b
			}
		}
		last = j
	}
	return last
}

// skipIndentedUntilTopLevel skips block-scalar content and unknown block-map
// children (any indented or blank lines) until the next zero-indent line.
func skipIndentedUntilTopLevel(lines []string, start int) int {
	last := start
	for j := start + 1; j < len(lines); j++ {
		trimmed := strings.TrimSpace(lines[j])
		if trimmed == "" {
			last = j
			continue
		}
		if leadingSpaces(lines[j]) == 0 {
			break
		}
		last = j
	}
	return last
}

func isBlockScalarMarker(value string) bool {
	return strings.HasPrefix(value, "|") || strings.HasPrefix(value, ">")
}

// splitKeyValue splits a trimmed "key: value" (or "key:") line on the first
// colon. The config subset this package owns never contains colons in keys or
// values, so a first-colon split is safe here.
func splitKeyValue(trimmed string) (key, value string) {
	idx := strings.IndexByte(trimmed, ':')
	if idx < 0 {
		return "", ""
	}
	return strings.TrimSpace(trimmed[:idx]), strings.TrimSpace(trimmed[idx+1:])
}

// unquote strips YAML quoting and trailing comments from a scalar value.
func unquote(value string) string {
	v := strings.TrimSpace(value)
	if len(v) > 0 && (v[0] == '"' || v[0] == '\'') {
		q := v[0]
		body := v[1:]
		if idx := strings.IndexByte(body, q); idx >= 0 {
			return body[:idx]
		}
		return body
	}
	if idx := strings.Index(v, " #"); idx >= 0 {
		v = v[:idx]
	}
	return strings.TrimSpace(v)
}

func leadingSpaces(s string) int {
	n := 0
	for _, r := range s {
		if r == ' ' {
			n++
			continue
		}
		break
	}
	return n
}

// CommandArgs holds the parsed flags for `gentle-ai sdd-config`.
type CommandArgs struct {
	CWD           string
	JSON          bool
	Detect        bool
	ValidateRepos bool
}

// ParseCommandArgs parses sdd-config CLI flags. It mirrors sddstatus.ParseCommandArgs:
// --cwd requires a value (rejecting a following flag), and any unexpected token
// is an error (sdd-config takes no positional argument).
func ParseCommandArgs(args []string) (CommandArgs, error) {
	var parsed CommandArgs
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--json":
			parsed.JSON = true
		case "--detect":
			parsed.Detect = true
		case "--validate-repos":
			parsed.ValidateRepos = true
		case "--cwd":
			if i+1 >= len(args) || strings.HasPrefix(args[i+1], "-") {
				return CommandArgs{}, fmt.Errorf("--cwd requires a value")
			}
			parsed.CWD = args[i+1]
			i++
		default:
			if strings.HasPrefix(args[i], "-") {
				return CommandArgs{}, fmt.Errorf("unknown sdd-config argument %q", args[i])
			}
			return CommandArgs{}, fmt.Errorf("unexpected sdd-config argument %q", args[i])
		}
	}
	return parsed, nil
}

// RenderConfigMarkdown renders a human-readable summary of cfg for the default
// (non-JSON) sdd-config output.
func RenderConfigMarkdown(cfg Config) string {
	lines := []string{
		"## SDD Config",
		"",
		fmt.Sprintf("domain: %s", presentOrUnset(cfg.Domain)),
		fmt.Sprintf("repos.infra: %s", presentOrUnset(cfg.Repos.Infra)),
		fmt.Sprintf("repos.carga: %s", presentOrUnset(cfg.Repos.Carga)),
		fmt.Sprintf("aws_profiles: %s", profileLogicalSummary(cfg.AWSProfiles)),
		fmt.Sprintf("verify.skip_deploy: %v", cfg.Verify.SkipDeploy),
	}
	return strings.Join(lines, "\n")
}

func presentOrUnset(value string) string {
	if value == "" {
		return "(unset — app-dev)"
	}
	return value
}

// profileLogicalSummary lists the configured logical profile names (sorted) or
// reports none. It deliberately never prints the real CLI profile names.
func profileLogicalSummary(profiles map[string]string) string {
	if len(profiles) == 0 {
		return "(none)"
	}
	logicals := make([]string, 0, len(profiles))
	for logical := range profiles {
		logicals = append(logicals, logical)
	}
	sort.Strings(logicals)
	return strings.Join(logicals, ", ")
}
