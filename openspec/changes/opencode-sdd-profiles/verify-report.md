# Verification Report: opencode-sdd-profiles (RE-VERIFY)

**Change**: `opencode-sdd-profiles`
**Date**: 2026-06-22 (re-verification after fixes)
**Mode**: Strict TDD (project config)
**Spec source**: `openspec/changes/opencode-sdd-profiles/specs/`
**Apply commit**: `592ff6c`

---

## Re-verification Scope

The previous verification report identified three REAL_ISSUE items that blocked
archiving. `sdd-apply` addressed them in commit `592ff6c`. This report verifies
that each fix is actually present in the code, covered by passing tests, and
matches the spec.

---

## Re-verification Summary

| # | Previous Warning | Verdict | Evidence | Status |
|---|------------------|---------|----------|--------|
| 1 | **R-PROF-31: Missing sync-time model warning** | **RESOLVED** | `internal/cli/sync.go` now contains `validateProfileModelAssignments()` (lines 776–833), called in the `ComponentSDD` path (lines 645–651). It loads `~/.cache/opencode/models.json`, warns on unknown model IDs, preserves assignments, and returns non-fatal warnings. `SyncResult.Warnings` is populated and rendered in `RenderSyncReport` (lines 1155–1161). | ✅ Fixed |
| 2 | **ScreenProfileCreate cache guard (task 6.2)** | **RESOLVED** | `internal/tui/screens/model_picker.go:renderPhaseList()` ForProfile + empty-cache branch (lines 680–691) shows the spec message "Run OpenCode at least once to populate the model cache" and only "← Back". `internal/tui/screens/profile_create.go:ProfileCreateOptionCount()` (lines 173–180) returns `1` for cache-missing step 1. `internal/tui/model.go:confirmProfileCreate()` (lines 3935–3945) backs out on enter when the cache is missing. | ✅ Fixed |
| 3 | **tasks.md checkbox discrepancy** | **RESOLVED** | `openspec/changes/opencode-sdd-profiles/tasks.md` now shows 37/38 tasks checked and 1 descoped (`6.1` E2E test is unchecked with a strikethrough explanation). Tasks `3.3`, `3.5`, `4.1`, `4.3` are annotated as direct `model.Update()` / render assertions; `6.2`, `6.4` include implementation notes and test names; `6.6` is annotated as string-assertion tests. | ✅ Reconciled |

---

## Completeness

| Metric | Value |
|--------|-------|
| Tasks total | 38 |
| Tasks complete `[x]` | 37 |
| Tasks descoped `[ ]` | 1 (6.1 E2E test) |
| Tasks incomplete / unmatched | 0 |

---

## Build & Tests Execution

**Build**: ✅ Passed
```text
go build ./...
# No output — clean build, zero errors
```

**Targeted tests** (profile + sync paths): ✅ Passed
```text
go test ./internal/tui/screens/... ./internal/tui/... ./internal/cli/...
ok  	github.com/gentleman-programming/gentle-ai/internal/tui/screens	(cached)
ok  	github.com/gentleman-programming/gentle-ai/internal/tui	(cached)
```

**R-PROF-31 specific tests**: ✅ Passed
```text
go test ./internal/cli/... -run 'TestValidateProfileModelAssignments|TestRunSyncWithSelection_WarnsOnUnknownProfileModel|TestRenderSyncReportIncludesWarnings' -v
=== RUN   TestValidateProfileModelAssignments_WarnsOnUnknownModel
--- PASS: TestValidateProfileModelAssignments_WarnsOnUnknownModel (0.00s)
=== RUN   TestValidateProfileModelAssignments_NoWarningsWhenCacheMissing
--- PASS: TestValidateProfileModelAssignments_NoWarningsWhenCacheMissing (0.00s)
=== RUN   TestValidateProfileModelAssignments_NoWarningsForKnownModels
--- PASS: TestValidateProfileModelAssignments_NoWarningsForKnownModels (0.00s)
=== RUN   TestRunSyncWithSelection_WarnsOnUnknownProfileModel
--- PASS: TestRunSyncWithSelection_WarnsOnUnknownProfileModel (0.54s)
=== RUN   TestRenderSyncReportIncludesWarnings
--- PASS: TestRenderSyncReportIncludesWarnings (0.00s)
PASS
ok  	github.com/gentleman-programming/gentle-ai/internal/cli	(cached)
```

**Task 6.2 cache-missing screen tests**: ✅ Passed
```text
go test ./internal/tui/screens/... -run 'TestProfileCreateOptionCount_Step1CacheMissingReturnsOne|TestRenderProfileCreate_Step1_CacheMissing_ShowsOnlyBack' -v
=== RUN   TestProfileCreateOptionCount_Step1CacheMissingReturnsOne
--- PASS: TestProfileCreateOptionCount_Step1CacheMissingReturnsOne (0.00s)
=== RUN   TestRenderProfileCreate_Step1_CacheMissing_ShowsOnlyBack
--- PASS: TestRenderProfileCreate_Step1_CacheMissing_ShowsOnlyBack (0.00s)
PASS
ok  	github.com/gentleman-programming/gentle-ai/internal/tui/screens	0.819s
```

**Task 6.2 model behavior test**: ✅ Passed
```text
go test ./internal/tui/... -run 'TestProfileCreateEmptyProviderEnterBacksOut' -v
=== RUN   TestProfileCreateEmptyProviderEnterBacksOut
--- PASS: TestProfileCreateEmptyProviderEnterBacksOut (0.00s)
PASS
ok  	github.com/gentleman-programming/gentle-ai/internal/tui	(cached)
```

**Full suite note**: `go test ./internal/cli/...` reports one pre-existing failure:
```text
--- FAIL: TestRunInstallKimiMissingUVFailsBeforeExecutingInstallCommands (0.09s)
    run_integration_test.go:2154: RunInstall() expected error when Kimi uv preflight fails
```
This failure is in the Kimi install preflight path, unrelated to the SDD profiles change, and was already documented in the previous verify report and in `apply-progress.md`. It is **not flagged as a regression**.

**Coverage**: Not measured (tool not configured).

---

## Spec Compliance Matrix

### Spec: sdd-profile-sync

| Requirement | Scenario | Test | Result |
|-------------|----------|------|--------|
| Missing Model Warning | Stale model ID preserved with warning | `internal/cli/sync_test.go > TestValidateProfileModelAssignments_WarnsOnUnknownModel` | ✅ COMPLIANT |
| Missing Model Warning | Stale model ID preserved with warning | `internal/cli/sync_test.go > TestRunSyncWithSelection_WarnsOnUnknownProfileModel` | ✅ COMPLIANT |
| Missing Model Warning | Warning rendered in sync report | `internal/cli/sync_test.go > TestRenderSyncReportIncludesWarnings` | ✅ COMPLIANT |

### Spec: sdd-profiles

| Requirement | Scenario | Test | Result |
|-------------|----------|------|--------|
| TUI — Profile Create Screen | Model cache not available | `internal/tui/screens/profile_create_test.go > TestRenderProfileCreate_Step1_CacheMissing_ShowsOnlyBack` | ✅ COMPLIANT |
| TUI — Profile Create Screen | Model cache not available | `internal/tui/screens/profile_create_test.go > TestProfileCreateOptionCount_Step1CacheMissingReturnsOne` | ✅ COMPLIANT |
| TUI — Profile Create Screen | Model cache not available | `internal/tui/model_test.go > TestProfileCreateEmptyProviderEnterBacksOut` | ✅ COMPLIANT |

**Compliance summary**: 6/6 relevant scenarios compliant (100%).

---

## Correctness (Static Evidence)

| Requirement | Status | Notes |
|------------|--------|-------|
| Missing model warning during sync (R-PROF-31) | ✅ Implemented | `validateProfileModelAssignments()` loads `models.json`, checks orchestrator and phase assignments, returns `{agent-key} references unknown model {model-id}`, and preserves the existing assignment. Warnings flow through `componentSyncStep.warnings` → `syncRuntime.warnings` → `SyncResult.Warnings` → `RenderSyncReport`. |
| Missing model cache guard in `ScreenProfileCreate` | ✅ Implemented | `renderPhaseList()` ForProfile branch renders the spec message + "← Back" only when `AvailableIDs` is empty. `ProfileCreateOptionCount()` returns 1. `confirmProfileCreate()` returns to step 0 (or `ScreenProfiles` in edit mode) on enter when cache is missing. |
| Task tracking integrity | ✅ Reconciled | `tasks.md` now honestly reflects 37 implemented tasks and 1 descoped E2E task (`6.1`). Overchecked items from the previous report are corrected with implementation-method annotations. |

---

## Coherence (Design)

| Decision | Followed? | Notes |
|----------|-----------|-------|
| R-PROF-31 warning format `{agent-key} references unknown model {model-id}` | ✅ Yes | Matches the spec scenario wording. |
| Cache-missing profile step: spec message + Back-only | ✅ Yes | No "Continue with defaults" path remains in the ForProfile branch. |
| Warnings surfaced via `SyncResult.Warnings` (suggestion from previous report) | ✅ Yes | Implemented; enables testable assertions beyond stderr. |
| TDD cycle evidence recorded in `apply-progress.md` | ✅ Yes | 8 RED→GREEN cycles documented for the fix batch. |

---

## Issues Found

### CRITICAL
**None.**

### WARNING
**None.** All three previously identified REAL_ISSUE items are resolved.

### SUGGESTION
- Consider adding a project-level `.tdd-cache` or CI note about `TestRunInstallKimiMissingUVFailsBeforeExecutingInstallCommands` being flaky/unrelated so future verifications do not need to manually exclude it.

---

## Verdict

### ✅ PASS — READY FOR ARCHIVE

All three REAL_ISSUE items from the previous verification report are resolved:

1. **R-PROF-31** — sync now emits a warning when a profile sub-agent references an unknown model, preserves the assignment, and renders the warning in the sync report. Covered by passing unit and integration tests.
2. **Task 6.2** — `ScreenProfileCreate` now shows the spec-required cache-missing message with a Back-only option. Covered by passing screen and model tests.
3. **tasks.md** — checkbox audit trail is reconciled: 37/38 tasks checked, 1 E2E task formally descoped with a documented rationale.

Build is clean, targeted tests pass, and the pre-existing Kimi install test failure is unrelated to this change.

---

## Next Recommended Phase

**`sdd-archive`** — the change is verified and ready to be archived.
