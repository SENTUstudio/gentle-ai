package sddconfig

import (
	"regexp"
	"sort"
	"strings"
)

// accountIDPattern matches a standalone 12-digit AWS account id (bounded by
// word boundaries, so longer digit runs are not touched).
var accountIDPattern = regexp.MustCompile(`\b\d{12}\b`)

// ResolveProfile maps a logical profile name (e.g. "prd", "dev", "usuario") to
// the CLI AWS profile name configured under aws_profiles in config.yaml. An
// unknown logical name or an app-dev config (no profiles) returns "".
func ResolveProfile(cfg Config, logical string) string {
	if cfg.AWSProfiles == nil {
		return ""
	}
	return cfg.AWSProfiles[logical]
}

// ScrubProfiles redacts sensitive AWS identifiers from text destined for logs
// or reports: each configured CLI profile name is replaced with
// `<aws-profile:logical>`, and any remaining standalone 12-digit account id is
// replaced with `<account-id>`.
//
// With an app-dev config (no profiles), only bare account ids are scrubbed.
func ScrubProfiles(text string, cfg Config) string {
	out := text
	if len(cfg.AWSProfiles) > 0 {
		// Replace longest profile names first so a name that is a prefix of
		// another cannot be partially matched.
		logicals := make([]string, 0, len(cfg.AWSProfiles))
		for logical := range cfg.AWSProfiles {
			logicals = append(logicals, logical)
		}
		sort.Slice(logicals, func(i, j int) bool {
			return len(cfg.AWSProfiles[logicals[i]]) > len(cfg.AWSProfiles[logicals[j]])
		})
		for _, logical := range logicals {
			name := cfg.AWSProfiles[logical]
			if name == "" {
				continue
			}
			out = strings.ReplaceAll(out, name, "<aws-profile:"+logical+">")
		}
	}
	out = accountIDPattern.ReplaceAllString(out, "<account-id>")
	return out
}
