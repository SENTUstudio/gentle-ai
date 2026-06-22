# Apply Progress: Data Engineering Domain Profile

**Change**: data-engineering-domain-profile
**Phase**: 1 of 8 — Profile Detection + Config Schema
**Mode**: Strict TDD (test runner: `go test ./...`)
**Artifact store**: both (OpenSpec + Engram)
**Delivery**: chained PRs — feature-branch-chain, 8 PRs; this batch = PR 1 (Config + detect + CLI)
**Chain strategy note**: `tasks.md` Review Workload Forecast says `stacked-to-main`; the
orchestrator instruction for this run says `feature-branch-chain`. The orchestrator instruction
is followed for PR targeting. This does not affect what was implemented.

## Completed Tasks

- [x] 1.1 `internal/sddconfig/config.go` + `config_test.go` — `Config`, `Repos`, `VerifyOpts`, `LoadConfig`, pure `parseConfig`.
- [x] 1.2 `internal/sddconfig/detect.go` + `detect_test.go` — `DetectDomain` (0.8 both markers / 0.5 one / 0 none).
- [x] 1.3 `internal/sddconfig/repos.go` + `repos_test.go` — `ValidateRepos`, `PathWarning`.
- [x] 1.4 `internal/sddconfig/profiles.go` + `profiles_test.go` — `ResolveProfile`, `ScrubProfiles`.
- [x] 1.5 `internal/cli/sdd_config.go` + `sdd_config_test.go` — `RunSDDConfig` (`--json`, `--detect`, `--validate-repos`, `--cwd`).
- [x] 1.6 Wire `sdd-config` subcommand in `internal/app/app.go` + `internal/app/help.go`.
- [x] 1.7 Patch `internal/assets/skills/sdd-init/SKILL.md` — preflight detect + confirm/override + write `domain`.

## TDD Cycle Evidence

| Task | Test File | Layer | Safety Net | RED | GREEN | TRIANGULATE | REFACTOR |
|------|-----------|-------|------------|-----|-------|-------------|----------|
| 1.1 | `internal/sddconfig/config_test.go` | Unit | N/A (new) | ✅ Written (undefined parseConfig/LoadConfig) | ✅ 9/9 passed | ✅ 6 cases (all fields, skip_deploy true, missing fields, block-scalar skip, quotes/comments, empty) + 3 LoadConfig file behaviors | ✅ Clean (shared `skipIndentedUntilTopLevel`) |
| 1.2 | `internal/sddconfig/detect_test.go` | Unit | N/A (new) | ✅ Written (undefined DetectDomain) | ✅ 6/6 passed | ✅ 5 cases (both, only-template, only-glue, none, glue-dir-without-py) + nonexistent-root error | ✅ DRY: unified `glueJobsHasPython` into `dirHasPython` |
| 1.3 | `internal/sddconfig/repos_test.go` | Unit | N/A (new) | ✅ Written (undefined ValidateRepos) | ✅ 6/6 passed | ✅ 6 cases (all exist, infra missing, both missing+order, empty config, absolute existing, carga missing) | ➖ None needed (already clean) |
| 1.4 | `internal/sddconfig/profiles_test.go` | Unit | N/A (new) | ✅ Written (undefined ResolveProfile/ScrubProfiles) | ✅ 9/9 passed | ✅ 7 cases (profile scrub, standalone acct id, both, 13-digit boundary, empty-config scrub, innocuous, resolve known/unknown/empty) | ➖ None needed |
| 1.5 | `internal/cli/sdd_config_test.go` | Integration | N/A (new file); cli pkg baseline = PASS except known-flaky | ✅ Written (undefined RunSDDConfig/DetectionReport/ValidationReport) | ✅ 9/9 passed | ✅ 9 cases (markdown, json config, app-dev json, detect json, detect markdown, validate json, validate markdown, bad cwd, bad flag) | ✅ DRY: extracted `renderOrEncode` |
| 1.6 | `internal/app/app_test.go` + `help_test.go` | Integration | ✅ internal/app baseline PASS (captured before edit) | ✅ Written (sdd-config not dispatched → "unknown command"; help missing sdd-config) | ✅ 3/3 passed (2 new app + help) | ✅ 2 dispatch cases (before-platform-validation + reads config domain) + help listing | ➖ None needed (one-line dispatch + help row) |
| 1.7 | `internal/assets/skills/sdd-init/SKILL.md` | N/A (markdown) | Golden tests (see below) | ➖ N/A (skill prompt, no Go test layer) | ➖ N/A | ➖ Triangulation skipped: markdown skill patch, validated by sdd-verify spec scenarios | ➖ N/A |

### Test Summary

- **Total new Go tests written**: 41 (sddconfig 30 unit + cli 9 integration + app 2 integration)
- **Total tests passing**: 41/41 (plus `TestHelpContainsAllCommands` updated to assert `sdd-config`)
- **Layers used**: Unit (30), Integration (11)
- **Approval tests** (refactoring): None — no existing Go code refactored; only additive new packages + thin wiring.
- **Pure functions created**: `parseConfig`, `LoadConfig` (file→pure), `DetectDomain`, `ValidateRepos`, `ResolveProfile`, `ScrubProfiles`, `ParseCommandArgs`, render helpers.

## Golden-File Refresh (task 1.7 side effect)

Patching the embedded `sdd-init` SKILL.md asset changed the content the SDD injector
copies to each agent, breaking 8 `TestGoldenSDD_*` golden references in
`internal/components`. Because the asset change is intentional (task 1.7), the golden
references were regenerated with `go test ./internal/components/ -run TestGoldenSDD -update`.
Diff verified: **only** the 8 `*-skill-sdd-init.golden` files changed, and **only** with the
Domain Preflight section + renumbered Execution Steps + 2 new Decision Gates rows. No
command/agent goldens changed; no non-deterministic content (timestamps/AI tags) introduced.
`TestGoldenSDD_Claude` was unaffected (Claude golden-checks commands/agents, not the skill SKILL.md).

## Verification Results

| Check | Command | Result |
|---|---|---|
| Build | `go build ./...` | ✅ PASS |
| Targeted tests | `go test ./internal/sddconfig/... ./internal/cli/...` | ✅ PASS (cli: only known-flaky `TestRunInstallKimiMissingUVFailsBeforeExecutingInstallCommands` fails — pre-existing, unrelated) |
| Full suite | `go test ./...` | ✅ No regressions — only the known-flaky install test fails (failed identically in pre-edit baseline) |
| Vet | `go vet ./internal/sddconfig/... ./internal/cli/... ./internal/app/...` and `go vet ./...` | ✅ PASS |

## Backward Compatibility (verified by tests)

- `LoadConfig` on a project with no `openspec/config.yaml` → zero `Config`, nil error.
- `LoadConfig` on the current app-dev `config.yaml` (with `context: |`, `rules:`, `testing:`) → zero `Config` (no `domain`/`repos`/`aws_profiles`/`verify` keys), nil error. The hand-rolled parser skips block scalars and unknown blocks.
- `DetectDomain` with zero markers → `""`, confidence 0.
- `ValidateRepos` on empty config → no warnings.
- `ResolveProfile` on empty config → `""`.
- `ScrubProfiles` on empty config → still scrubs bare 12-digit account ids (harmless for app-dev).
- `sdd-config` subcommand dispatches before platform validation; no existing command behavior changed.
- All `sddconfig`/`cli`/`app`/`components` tests green except the pre-existing flaky install test.

## Deviations from Design

- **ScrubProfiles location**: The design Interfaces block lists `ScrubProfiles` under `// internal/etl`, but the File Changes table and the Phase 1 task both place it in `internal/sddconfig/profiles.go`. The task prompt is authoritative for Phase 1 → `ScrubProfiles` lives in `sddconfig`. If a later phase moves ETL helpers to `internal/etl`, it can re-export or relocate then.
- **Subcommand wiring location**: `tasks.md` 1.6 says "in `cmd/gentle-ai/main.go`", but the actual root dispatch is in `internal/app/app.go` (`cmd/gentle-ai/main.go` only calls `app.Run()`). Wired in `app.go`'s info-commands switch + `help.go`, matching the existing `sdd-status`/`sdd-continue` pattern.
- **YAML parsing**: The project deliberately avoids `gopkg.in/yaml.v3` (hand-rolled scanner in `internal/components/filemerge/yaml.go`). `parseConfig` is a minimal indentation-aware reader for the 4 owned keys, skipping block scalars (`context: |`) and unknown blocks — no new dependency added.
- **One-marker detection semantics**: A single marker returns `domain="data-engineering"` at confidence 0.5 (a hint), not `""`. This lets `sdd-init` present the hint for confirm/override. Zero markers → `""`. This matches "present hint, confirm/override" in the design; flagged here in case the team prefers one-marker → empty.

## Issues Found

- Golden-file regression from the skill patch (resolved by regeneration, see above). Future phases (2-5) that patch embedded SDD skill assets MUST run `go test ./internal/components/ -run TestGoldenSDD -update` and verify the diff is only the intended skill content.

## Remaining Tasks

Phase 2-8 are out of scope for this batch (do not implement):
- [ ] 2.1–2.4 Spec + Design templates (`internal/etl/sidecar.go`, `pattern.go`, skill patches)
- [ ] 3.1–3.4 Apply + Verify dual-path (`internal/etl/header.go`, `compare.go`, skill patches)
- [ ] 4a.1–4a.6, 4b.1–4b.2, 4c.1–4c.2, 4d.1 Skills foundation + embeds
- [ ] 5.1–5.2 Tasks multi-repo + git flow

## Workload / PR Boundary

- **Mode**: chained PR slice (PR 1 of 8)
- **Current work unit**: Phase 1 — Detection + Config (`sddconfig` core + `sdd-config` CLI + `sdd-init` preflight)
- **Boundary**: starts at task 1.1, ends at task 1.7. Self-contained: new `internal/sddconfig` package, thin `cli` entry, one-line `app` dispatch, `sdd-init` skill patch + golden refresh. No behavior change when `domain` is absent.
- **Estimated review budget impact**: PR 1 forecast was ~730 lines. Real review surface (Go + skill patch + tests) is well under the 800-line budget. The 8 regenerated golden files add +312/−56 mechanical lines that are auto-generated references (diff verified to contain only the intended skill content) — reviewable as a sanity check, not line-by-line logic.

## Status

7/7 Phase 1 tasks complete. Ready for `sdd-apply` Phase 2 (or `sdd-verify` if the orchestrator prefers to verify this slice first).
