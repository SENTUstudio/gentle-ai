# Apply Progress: opencode-sdd-profiles (REAL_ISSUES fix batch)

**Change**: `opencode-sdd-profiles`
**Date**: 2026-06-22
**Mode**: Strict TDD (project config)
**Scope**: Fix the 3 REAL_ISSUES identified in `verify-report.md`

---

## Goal

Resolve the three warnings the verification report classified as REAL_ISSUE
before archiving:

1. **R-PROF-31** — missing sync-time model warning (`internal/cli/sync.go`)
2. **Task 6.2** — cache-missing guard shows "Continue with defaults" instead of
   Back-only (`internal/tui/screens/profile_create.go` + `model_picker.go` +
   `model.go`)
3. **tasks.md** — checkbox reconciliation with actual implementation

---

## TDD Cycle Evidence

| Task | RED (test written first) | GREEN (implementation passes) | REFACTOR |
|------|--------------------------|-------------------------------|----------|
| R-PROF-31 unit | `TestValidateProfileModelAssignments_WarnsOnUnknownModel` — undefined function → build fail | `validateProfileModelAssignments` implemented; test passes | No refactor needed |
| R-PROF-31 cache-missing | `TestValidateProfileModelAssignments_NoWarningsWhenCacheMissing` — build fail | Passes (Stat guard returns nil when file absent) | — |
| R-PROF-31 known-models | `TestValidateProfileModelAssignments_NoWarningsForKnownModels` — build fail | Passes | — |
| R-PROF-31 integration | `TestRunSyncWithSelection_WarnsOnUnknownProfileModel` — `SyncResult.Warnings` undefined → build fail | `Warnings` field + pipeline wiring; test passes | — |
| R-PROF-31 report | `TestRenderSyncReportIncludesWarnings` — `Warnings` field undefined → build fail | `RenderSyncReport` renders warnings section; test passes | — |
| Task 6.2 render | `TestRenderProfileCreate_Step1_CacheMissing_ShowsOnlyBack` — output contained "Continue with defaults" | `renderPhaseList` ForProfile branch shows Back-only; test passes | — |
| Task 6.2 option count | `TestProfileCreateOptionCount_Step1CacheMissingReturnsOne` — returned 2 | `ProfileCreateOptionCount` returns 1 for empty cache; test passes | — |
| Task 6.2 behavior | `TestProfileCreateEmptyProviderEnterBacksOut` (replaced `...ContinuesAndBacksOut`) — advanced to step 2 | `confirmProfileCreate` empty-cache = Back only; test passes | Updated existing test to match spec |

All 8 TDD cycles followed RED → GREEN. No task was implemented without a
failing test first.

---

## Completed Tasks

- [x] R-PROF-31: Added `validateProfileModelAssignments` + `SyncResult.Warnings` + pipeline wiring + report rendering
- [x] Task 6.2: Cache-missing profile create step shows Back-only (spec message + single option)
- [x] tasks.md: Reconciled checkboxes with actual implementation (teatest → model.Update, E2E descoped, snapshot → string-assertion)

---

## Files Changed

| File | Action | What Was Done |
|------|--------|---------------|
| `internal/cli/sync.go` | Modified | Added `Warnings []string` to `SyncResult`; `warnings` to `syncRuntime`/`componentSyncStep`; `validateProfileModelAssignments` + `modelExistsInCache` functions; validation call in ComponentSDD case; `result.Warnings` in `RunSyncWithSelection`; warnings section in `RenderSyncReport`; `opencode` import |
| `internal/cli/sync_test.go` | Modified | Added 5 tests: `TestValidateProfileModelAssignments_WarnsOnUnknownModel`, `_NoWarningsWhenCacheMissing`, `_NoWarningsForKnownModels`, `TestRunSyncWithSelection_WarnsOnUnknownProfileModel`, `TestRenderSyncReportIncludesWarnings`; `writeTestModelCache` helper |
| `internal/tui/screens/profile_create.go` | Modified | `ProfileCreateOptionCount` returns 1 for empty cache (Back-only) |
| `internal/tui/screens/profile_create_test.go` | Modified | Added `TestRenderProfileCreate_Step1_CacheMissing_ShowsOnlyBack`; replaced `...IncludesContinueAndBack` with `...CacheMissingReturnsOne` (expects 1) |
| `internal/tui/screens/model_picker.go` | Modified | `renderPhaseList` ForProfile + empty-cache branch: spec message + Back-only (removed "Continue with defaults") |
| `internal/tui/model.go` | Modified | `confirmProfileCreate` empty-cache case: Back-only (removed Continue cursor-0 path) |
| `internal/tui/model_test.go` | Modified | Replaced `TestProfileCreateEmptyProviderEnterContinuesAndBacksOut` with `TestProfileCreateEmptyProviderEnterBacksOut` (Back-only behavior) |
| `openspec/changes/opencode-sdd-profiles/tasks.md` | Modified | 3.3/3.5/4.1/4.3: teatest → model.Update()/render; 6.1: unchecked + descoped; 6.2/6.4: kept checked with implementation notes; 6.6: snapshot → string-assertion |
| `openspec/changes/opencode-sdd-profiles/apply-progress.md` | Created | This file |

---

## Deviations from Design

None — implementation matches the spec requirements:
- R-PROF-31 warning format: `{agent-key} references unknown model {model-id}` (matches spec scenario)
- Task 6.2: spec message "Run OpenCode at least once to populate the model cache" + Back-only (matches spec scenario "Model cache not available")

---

## Issues Found

None. The pre-existing `TestRunInstallKimiMissingUVFailsBeforeExecutingInstallCommands` failure in `internal/cli` is unrelated to this change (Kimi install preflight path) and was already noted in the verify report.

---

## Build & Test Results

```
go build ./...                        → PASS (no output)
go test ./internal/tui/screens/...     → PASS
go test ./internal/tui/...             → PASS
go test ./internal/cli/...             → 1 pre-existing unrelated failure (Kimi install)
```

All profile/sync-related tests pass.

---

## Remaining Tasks

- [ ] 6.1 E2E test — formally descoped (coverage via unit + integration tests)

### Status

37/38 tasks complete (1 descoped). Ready for `sdd-verify`.

---

## Next Recommended Phase

**`sdd-verify`** — re-run verification to confirm the three REAL_ISSUES are
resolved, then `sdd-archive`.
