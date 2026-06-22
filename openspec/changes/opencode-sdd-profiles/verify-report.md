# Verification Report: opencode-sdd-profiles

**Change**: `opencode-sdd-profiles`
**Date**: 2026-06-21 (warning investigation)
**Mode**: Strict TDD (project config)
**Spec source**: `openspec/changes/opencode-sdd-profiles/specs/`

---

## Investigation Scope

The user explicitly chose to investigate the three warnings recorded in the previous verification report before archiving. This report re-examines each warning, records the evidence found in the current codebase, and gives a verdict for each.

---

## Warning Investigation Summary

| # | Warning | Verdict | Evidence | Recommended Action |
|---|---|---|---|---|
| 1 | **R-PROF-31: Missing warning log when syncing profiles** | **REAL_ISSUE** | `internal/cli/sync.go` (ComponentSDD path, lines ~601-657) resolves profiles but never loads `~/.cache/opencode/models.json` or validates profile model IDs. No warning is emitted. Model preservation works by deep merge, but the spec's MUST-level warning is absent. | Implement the warning (load cache + warn on unknown model) or formally descope R-PROF-31 in the spec and tasks. |
| 2 | **ScreenProfileCreate: Missing cache guard (task 6.2)** | **REAL_ISSUE** | `internal/tui/model.go` enters `ScreenProfileCreate` without checking the model cache. After the name step, `handleProfileNameInput` advances to step 1 even when the cache is missing, and `internal/tui/screens/model_picker.go:renderPhaseList` shows an empty-state message plus **"Continue with defaults"** and **"← Back"**. The spec scenario *"Model cache not available"* requires the message **and only a "Back" option**. | Add a cache-missing guard when entering `ScreenProfileCreate` (Back-only message) or update the spec/task to match the current "Continue with defaults" behavior. |
| 3 | **tasks.md checkbox discrepancy** | **REAL_ISSUE** | `tasks.md` now shows **all 38 tasks checked**, but several checked items do not match the actual implementation: teatest-based tasks (3.3, 3.5, 4.1, 4.3) are render unit tests, task 6.1 (E2E test) does not exist, tasks 6.2 and 6.4 are partially/not implemented as specified, and task 6.6 is covered by string-assertion tests rather than snapshot tests. | Reconcile `tasks.md` with the actual implementation (uncheck or implement the missing/partial items). |

### Detailed Evidence

#### 1. R-PROF-31 — Missing sync-time model warning

- **Spec**: `specs/sdd-profile-sync/spec.md`, Requirement *Missing Model Warning*:
  > "If a profile sub-agent references a model that no longer exists in `~/.cache/opencode/models.json`, sync MUST emit a warning and preserve the existing model assignment."
- **Implementation path**: `internal/cli/sync.go` `componentSyncStep.Run()` for `model.ComponentSDD`.
- **Observation**: The code calls `sdd.DetectProfiles` or uses explicit profiles, then passes them to `sdd.Inject`. It never reads the model cache (`opencode.LoadModels` / `LoadModelsOrEmpty`) and never validates whether an assigned model ID exists. There is no `fmt.Fprintf(os.Stderr, "WARNING: ...")` or equivalent.
- **Conclusion**: The warning is genuinely missing. The behavior is non-fatal and preserves models, but it does not satisfy the spec's MUST.

#### 2. ScreenProfileCreate — Missing cache guard

- **Spec**: `specs/sdd-profiles/spec.md`, Scenario *Model cache not available*:
  > "GIVEN `~/.cache/opencode/models.json` does not exist ... THEN a message 'Run OpenCode at least once to populate the model cache' is shown AND only a 'Back' option is available."
- **Task 6.2**: "Handle missing OpenCode model cache edge case in `ScreenProfileCreate`: if `~/.cache/opencode/models.json` does not exist, show ... message and only offer 'Back'."
- **Implementation paths**:
  - `internal/tui/model.go` lines 1388-1397 (`n` key) and 1722-1732 (`Create new profile` enter) transition to `ScreenProfileCreate` without checking the cache.
  - `internal/tui/model.go` `handleProfileNameInput` (lines 3861-3868) advances to step 1 with an empty `ModelPickerState` when the cache is missing.
  - `internal/tui/screens/model_picker.go:renderPhaseList` (lines 680-695) renders the empty-state message and offers **"Continue with defaults"** plus **"← Back"**.
- **Conclusion**: The user can still create a profile with default assignments when the cache is missing. This contradicts the spec's "Back only" requirement and makes task 6.2 overchecked.

#### 3. tasks.md — Checkbox accuracy

- Current `tasks.md` marks **all 38 tasks `[x]`**.
- Verified mismatches:
  - **3.3** `[x]` "Write teatest test for `ScreenProfiles` ..." — `internal/tui/screens/profiles_test.go` contains render/unit tests, no `teatest` import or message-loop tests.
  - **3.5** `[x]` "Write teatest test for `ScreenProfileCreate` step flow ..." — `internal/tui/screens/profile_create_test.go` contains render/unit tests, no `teatest`.
  - **4.1** `[x]` "Write teatest test for edit flow ..." — same file, no `teatest`.
  - **4.3** `[x]` "Write teatest test for `ScreenProfileDelete` ..." — `internal/tui/screens/profile_delete_test.go` contains render/unit tests, no `teatest`.
  - **6.1** `[x]` "E2E: profile creation, sync, list display, edit, re-sync, delete ..." — no E2E test file exists for this change.
  - **6.2** `[x]` "Handle missing OpenCode model cache edge case ... only offer 'Back'" — implemented as empty-state "Continue with defaults" path, not Back-only.
  - **6.4** `[x]` "Handle sync-time missing model warning (R-PROF-31) ..." — warning not implemented.
  - **6.6** `[x]` "TUI snapshot tests for welcome screen ..." — `internal/tui/screens/welcome_test.go` has string-assertion unit tests, not snapshot/golden tests.
- **Conclusion**: The task list is currently **overchecked**. It claims test coverage and edge-case behavior that are not present in the codebase.

---

## Completeness

| Metric | Value |
|--------|-------|
| Tasks total | 38 |
| Tasks marked complete `[x]` | 38 |
| Tasks with verified implementation | ~30 (core feature + most unit tests) |
| Tasks overchecked / not as specified | 8 (3.3, 3.5, 4.1, 4.3, 6.1, 6.2, 6.4, 6.6) |

> The previous report noted that tasks were under-checked but implemented. The current `tasks.md` has swung the other way: it is now overchecked relative to actual test coverage and edge-case implementation.

---

## Build & Tests Execution

**Build**: ✅ Passed
```text
go build ./...
# No output — clean build, zero errors
```

**Targeted tests** (profile + sync paths): ✅ Passed
```text
go test ./internal/components/sdd/... ./internal/tui/... ./internal/tui/screens/... ./internal/cli/ -run 'Profile|Sync'
ok  	github.com/gentleman-programming/gentle-ai/internal/components/sdd	5.448s
ok  	github.com/gentleman-programming/gentle-ai/internal/tui	0.951s
ok  	github.com/gentleman-programming/gentle-ai/internal/tui/screens	1.277s
ok  	github.com/gentleman-programming/gentle-ai/internal/cli	18.408s
```

**Full suite** (`go test ./...`): ❌ One unrelated failure in `internal/cli`
```text
--- FAIL: TestRunInstallKimiMissingUVFailsBeforeExecutingInstallCommands (0.13s)
    run_integration_test.go:2154: RunInstall() expected error when Kimi uv preflight fails
```
This failure is in the Kimi install preflight path and is not related to the SDD profiles change. The profile/sync code paths pass.

**Coverage**: Not measured (tool not configured).

---

## Spec Compliance Matrix (relevant rows)

### Spec: sdd-profile-sync

| Requirement | Scenario | Test | Result |
|-------------|----------|------|--------|
| Missing Model Warning | Stale model ID preserved with warning | None found | ❌ UNTESTED / MISSING WARNING |

### Spec: sdd-profiles

| Requirement | Scenario | Test | Result |
|-------------|----------|------|--------|
| TUI — Profile Create Screen | Model cache not available | `model_picker.go > renderPhaseList` empty state | ⚠️ PARTIAL (message shown, but "Continue with defaults" offered instead of Back only) |

**Compliance summary**: The core CRUD, agent generation, shared prompts, CLI flags, and sync integration scenarios remain compliant. The two warnings above are the remaining spec gaps.

---

## Correctness (Static — Relevant Updates)

| Requirement | Status | Notes |
|------------|--------|-------|
| Missing model warning during sync (R-PROF-31) | ❌ Missing | No model-cache load or validation in `internal/cli/sync.go` ComponentSDD path. |
| Missing model cache guard in `ScreenProfileCreate` | ⚠️ Partial | Empty-state message exists, but spec requires Back-only option. |
| Task tracking integrity | ❌ Overchecked | `tasks.md` claims teatest/E2E coverage and edge-case behavior not present. |

---

## Strict TDD Compliance (additional observation)

| Check | Result | Details |
|-------|--------|---------|
| TDD Evidence reported | ❌ | `apply-progress.md` does not exist for this change; no TDD Cycle Evidence table available. |
| Existing tests pass | ✅ | All profile/sync-related tests pass (targeted run). |
| Assertion quality | ✅ | No tautologies or ghost loops observed in profile-related tests. |

The missing `apply-progress.md` is an audit-trail gap. It does not block the feature from working, but it means the strict-TDD paper trail cannot be validated for this change.

---

## Issues Found

### CRITICAL
**None.** The feature is functionally complete and the targeted tests pass. The issues below are spec/audit gaps, not runtime failures.

### WARNING (investigated)

1. **R-PROF-31 missing sync-time model warning** — **REAL_ISSUE**. Spec MUST is not implemented. A warning should be emitted when a profile sub-agent references a model missing from the OpenCode cache.
2. **ScreenProfileCreate missing cache guard** — **REAL_ISSUE**. The spec requires a Back-only message when the model cache is missing; the current screen allows "Continue with defaults".
3. **tasks.md overchecked** — **REAL_ISSUE**. Several checked tasks do not match the actual implementation (missing teatest/E2E coverage, partial edge-case handling).

### SUGGESTION
- Add a `Warnings []string` field to `SyncResult` so missing-model warnings can be tested via `RenderSyncReport`, not only observed on stderr.
- If teatest/E2E coverage is intentionally out of scope, record the explicit descope in `tasks.md` and the spec, rather than leaving checked boxes for unimplemented work.

---

## Verdict

### ❌ FAIL — DO NOT ARCHIVE WITHOUT FIX OR FORMAL DESCOPE

The implementation is feature-complete, builds cleanly, and all profile/sync-related tests pass. However, the three investigated warnings are **real issues**, not acceptable descopes:

1. **R-PROF-31** is a spec non-compliance (MUST-level warning missing).
2. **Task 6.2** is a spec non-compliance (cache-missing screen allows Continue instead of Back only).
3. **tasks.md** is currently overchecked, which corrupts the audit trail.

Archiving now would record the change as fully complete when it is not. Either fix the two code gaps and reconcile `tasks.md`, or formally descope R-PROF-31 and the Back-only cache guard in the spec/tasks and update the checkboxes honestly.

---

## Next Recommended Phase

**`sdd-apply`** — implement the R-PROF-31 warning, add the ScreenProfileCreate cache-missing guard, and reconcile `tasks.md`. Then re-run `sdd-verify` before archiving.
