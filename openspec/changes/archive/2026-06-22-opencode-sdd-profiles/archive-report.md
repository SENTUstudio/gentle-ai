# Archive Report: opencode-sdd-profiles

**Change**: `opencode-sdd-profiles`
**Archived**: 2026-06-22
**Mode**: hybrid (OpenSpec + Engram)
**Status**: PASS (after re-verify)

---

## What Was Implemented

SDD profile CRUD system for OpenCode: users can create, list, edit, and delete named model profiles from the TUI or CLI. Each profile generates 1 orchestrator + 10 sub-agents with suffixed keys (e.g., `sdd-apply-cheap`). Sub-agent prompts are extracted to shared files at `~/.config/opencode/prompts/sdd/` via `{file:...}` references. Profile detection reads from `opencode.json` (single source of truth). CLI gains `--profile` and `--profile-phase` flags for headless profile management.

## Verification Status

**PASS** — re-verification after 3 REAL_ISSUE fixes (commit `30c95eb`).

### REAL_ISSUES Found and Fixed

| # | Issue | Resolution | Commit |
|---|-------|-----------|--------|
| 1 | **R-PROF-31**: Missing sync-time model warning | Added `validateProfileModelAssignments()` + `SyncResult.Warnings` pipeline + `RenderSyncReport` warnings section | `592ff6c` |
| 2 | **ScreenProfileCreate cache guard (task 6.2)**: Showed "Continue with defaults" instead of Back-only when model cache missing | Fixed `renderPhaseList` ForProfile branch, `ProfileCreateOptionCount`, and `confirmProfileCreate` to show Back-only | `592ff6c` |
| 3 | **tasks.md checkbox discrepancy**: Overchecked items from initial implementation | Reconciled checkboxes with actual implementation (teatest → model.Update, E2E descoped, snapshot → string-assertion) | `592ff6c` |

### Descope

- **Task 6.1 (E2E test)**: Formally descoped. Coverage provided by unit tests (direct `model.Update()`) and `TestRunSyncWithProfilesIntegration` sync integration test.

## Spec Deltas Merged

| Domain | Action | Details |
|--------|--------|---------|
| `sdd-profiles` | Created (new domain) | 350 lines, 12 requirements covering profile CRUD, naming, agent generation, TUI screens, CLI flags |
| `sdd-profile-sync` | Created (new domain) | 132 lines, 8 requirements covering sync-time detection, prompt maintenance, model preservation, warnings, idempotency |
| `gga` | Updated (merged) | 2 ADDED requirements: Welcome Screen profile option + Sync `--profile` flag (61 lines appended) |

**Overlap check**: `2026-06-10-level-neutral-persona-parity` touches persona behavior files. This change is about SDD profile agent management — distinct spec domains, no overlap.

## Archive Contents

- `proposal.md` ✅
- `specs/` ✅ (3 domains: gga, sdd-profiles, sdd-profile-sync)
- `design.md` ✅
- `tasks.md` ✅ (37/38 tasks complete, 1 descoped)
- `verify-report.md` ✅ (PASS after re-verify)
- `apply-progress.md` ✅ (3 REAL_ISSUES fix batch documented)

## Files Changed (Implementation)

| File | Action |
|------|--------|
| `internal/model/types.go` | Modified — `Profile` struct |
| `internal/model/selection.go` | Modified — `Profiles []Profile` in `SyncOverrides` |
| `internal/components/sdd/profiles.go` | Created — `DetectProfiles`, `GenerateProfileOverlay`, `RemoveProfileAgents`, `ValidateProfileName`, `ProfileAgentKeys` |
| `internal/components/sdd/prompts.go` | Created — `WriteSharedPromptFiles`, `SharedPromptDir` |
| `internal/components/sdd/inject.go` | Modified — profile iteration, `{file:...}` refs |
| `internal/components/sdd/read_assignments.go` | Modified — `DetectProfiles` wrapper |
| `internal/tui/screens/profiles.go` | Created — profile list screen |
| `internal/tui/screens/profile_create.go` | Created — 4-step create/edit flow |
| `internal/tui/screens/profile_delete.go` | Created — delete confirmation screen |
| `internal/tui/screens/model_picker.go` | Modified — ForProfile cache-missing branch |
| `internal/tui/model.go` | Modified — screen constants, state fields, key handling |
| `internal/tui/router.go` | Modified — profile screen routes |
| `internal/tui/screens/welcome.go` | Modified — "OpenCode SDD Profiles" option |
| `internal/cli/sync.go` | Modified — `--profile` flag, `validateProfileModelAssignments`, backup scope |
| `internal/assets/opencode/sdd-overlay-multi.json` | Modified — `{file:...}` references |

## SDD Cycle Complete

The change has been fully planned, implemented, verified, and archived.
Ready for the next change.
