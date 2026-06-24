# Apply Progress: Domain-Aware Model Profiles

**Change**: domain-aware-model-profiles
**Mode**: Strict TDD (`openspec/config.yaml` → `strict_tdd: true`, Go standard testing)
**Executor**: sdd-apply (single batch, Tasks 1–5)
**Started from**: Pre-completed foundation — `model.Profile.Domain`, `model.DefaultModelsForDomain`, and 5 `model` package tests (already merged in prior work).

---

## Completed Tasks

- [x] **Task 1** — `EnsureProfileDomainConsistency` helper in `sdd` package
- [x] **Task 2** — Sync CLI invokes collision guard
- [x] **Task 3** — TUI prefill at step 0→1 transition
- [x] **Task 4** — Render domain indicator in step 1 footer
- [x] **Task 5** — Wire `ProfileDraft.Domain` to render call (no-op verification)

**Status: 5/5 tasks complete. Ready for verify.**

---

## Files Changed

| File | Action | What Was Done |
|------|--------|---------------|
| `internal/components/sdd/profiles.go` | Modified | Added `EnsureProfileDomainConsistency(detected, explicit)` + private `effectiveProfileDomain` normalizer (treats `""`/`"app-dev"` as equivalent). Placed below `DetectProfiles` per the task. |
| `internal/components/sdd/profiles_test.go` | Modified | Added table-driven `TestEnsureProfileDomainConsistency` (7 subtests: empty detected, empty explicit, same domain `""`/`""`, `""`/`"app-dev"` equivalence, single conflict, multiple conflicts, non-overlapping names). |
| `internal/cli/sync.go` | Modified | Hoisted on-disk profile detection into a `detected` variable; added collision guard `EnsureProfileDomainConsistency(detected, s.selection.Profiles)` that returns `fmt.Errorf("sync sdd profile domain collision: %w", err)` from `componentSyncStep`. Guard runs only when explicit profiles are provided. |
| `internal/cli/sync_test.go` | Modified | Added `sdd` import. Added `TestRunSyncDomainCollisionError` (two-sync: create `cheap` app-dev → re-sync `cheap` data-engineering → wrapped error) and triangulation `TestRunSyncDomainConsistency_AllowsSameEffectiveDomain` (`""` vs `"app-dev"` re-sync succeeds). |
| `internal/tui/model.go` | Modified | Added `sddconfig` import. Added domain-aware prefill in `handleProfileNameInput` create step 0→1 branch: `sddconfig.LoadConfig(osGetwdFn())` → stamp `ProfileDraft.Domain` always, seed `Selection.ModelAssignments` with `DefaultModelsForDomain` only when domain non-empty AND assignments empty (idempotent). |
| `internal/tui/model_test.go` | Modified | Added `writeDomainConfig` helper + 4 tests: prefill from domain, skip when assignments present (idempotency), skip when domain empty, edit mode does not re-detect. |
| `internal/tui/screens/profile_create.go` | Modified | Refactored `renderProfileModelStep` to accept `draft model.Profile` (was `profileName string`); renders `Domain: <draft.Domain>` footer between the Assign Models heading and the picker, only when `draft.Domain != ""`. `ProfileCreateOptionCount` unchanged. |
| `internal/tui/screens/profile_create_test.go` | Modified | Added `TestRenderProfileCreate_Step1_ShowsDomainFooter` and `TestRenderProfileCreate_Step1_OmitsDomainFooterWhenEmpty`. |

**Pre-existing foundation (not in this batch, already done):** `internal/model/types.go` (`Profile.Domain`), `internal/model/selection.go` (`DefaultModelsForDomain`), `internal/model/selection_test.go` (5 tests).

---

## TDD Cycle Evidence

| Task | Test File | Layer | Safety Net | RED | GREEN | TRIANGULATE | REFACTOR |
|------|-----------|-------|------------|-----|-------|-------------|----------|
| 1 | `internal/components/sdd/profiles_test.go` | Unit | ✅ sdd ok (cached) | ✅ `undefined: EnsureProfileDomainConsistency` | ✅ 7/7 subtests pass | ✅ 7 cases (no-conflict ×3, equivalence, single conflict, multi-conflict, non-overlap) | ✅ Extracted `effectiveProfileDomain`, sorted for determinism |
| 2 | `internal/cli/sync_test.go` | Integration | ✅ cli sync tests ok | ✅ run2 returned nil (no guard) | ✅ both tests pass | ✅ 2 cases (collision errors + same-effective-domain succeeds) | ✅ Hoisted detection, single guard call site |
| 3 | `internal/tui/model_test.go` | Unit (TUI state) | ✅ tui ok | ✅ Tests 1&2 fail (no Domain stamp/prefill); 3&4 guard tests pass pre-impl | ✅ 4/4 pass | ✅ 4 cases (prefill, idempotency guard, empty-domain skip, edit-mode no-detect) | ✅ Focused block, reuses `osGetwdFn`/`setOSGetwdForTest` |
| 4 | `internal/tui/screens/profile_create_test.go` | Unit (render) | ✅ screens ok | ✅ Test 1 fails (no footer); Test 2 guard passes pre-impl | ✅ 2/2 pass | ✅ 2 cases (shows when non-empty, omits when empty) | ✅ Signature refactor to `draft model.Profile`, one caller |
| 5 | — (no-op) | — | ✅ tui ok | ➖ N/A (verification) | ✅ build + `./internal/tui/...` green | ➖ N/A | ➖ N/A — field already plumbed via `draft` at `model.go:979` |

### Test Summary
- **Total new tests written**: 15 (7 table subtests + 4 TUI + 2 render + 2 cli)
- **Total new tests passing**: 15
- **Layers used**: Unit (sdd, tui, screens), Integration (cli sync)
- **Approval tests** (refactoring): None — Task 4 was a signature refactor validated by new behavior tests, not approval tests (the renderer had no prior domain-footer behavior to preserve).
- **Pure functions created**: 2 (`EnsureProfileDomainConsistency`, `effectiveProfileDomain`)

---

## Deviations from Design

1. **Task 2 — guard call argument (`detected` vs `profiles`).**
   The task specifies calling `EnsureProfileDomainConsistency(profiles, s.selection.Profiles)` immediately after `DetectProfiles` populates `profiles`. Taken literally, this **cannot produce the error the spec's own test requires**: when explicit profiles are provided (`s.selection.Profiles` non-empty), the existing code skips `DetectProfiles` entirely, so `profiles == s.selection.Profiles` and the helper would compare a list against itself — never a conflict. The design's *intent* (and its test) is to reject an **explicit** profile that reuses the name of an **existing on-disk** profile with a different effective domain. To realize that intent, detection was hoisted into a `detected` variable that always runs (for non-external-strategy syncs), and the guard is called as `EnsureProfileDomainConsistency(detected, s.selection.Profiles)`. The guard is invoked at a single site and only when explicit profiles are present (a profile cannot collide with itself, so the no-explicit re-sync path skips it). Behavior for all existing tests is preserved (their explicit profiles carry `Domain=""`, same effective domain as detected `""`).

2. **Task 3 — implementation site (`handleProfileNameInput`, not `confirmProfileCreate`).**
   The task text says "modify the step 0→1 branch in `confirmProfileCreate()` (around line 3861)". The design anchors on "line ~3861", which is the **create** step 0→1 transition inside `handleProfileNameInput` — not `confirmProfileCreate`. The dispatch (`model.go:850`) routes create-mode step 0 keys to `handleProfileNameInput`; `confirmProfileCreate` step 0 only handles **edit** mode. Implementing in `confirmProfileCreate` would mean the real create flow never prefills. The prefill was therefore placed in `handleProfileNameInput` (matching the design's line anchor and the real create flow), and tests drive it via `m.Update(tea.KeyEnter)` — identical to every existing profile test in `model_test.go`. The spec's "Call `confirmProfileCreate()`" instruction was adapted to the actual dispatch. Edit mode is untouched (it advances via `confirmProfileCreate`, which has no prefill), satisfying "edit mode leaves domain detection untouched".

3. **Task 2 test — two-sync fixture instead of pre-writing `opencode.json`.**
   The task says "Write temp `opencode.json` containing `sdd-orchestrator-cheap`". The robust, proven pattern in this repo (e.g. `TestRunSyncDetectsExistingProfilesOnRegularSync`) is to establish the on-disk profile via a first sync, then trigger the collision with a second. This exercises the real `DetectProfiles → guard → wrapped-error` path without fragile direct file writes that the sync's deep-merge could disturb. The test asserts the exact same facts (error wraps `EnsureProfileDomainConsistency`, contains `"domain"` and `"cheap"`). It also required stubbing `backup.UserHomeDirFn` (not just `osUserHomeDir`) so the pipeline's rollback path validation succeeds and the original collision error propagates instead of being masked by a rollback failure.

---

## Issues Found

- **Pre-existing test failure (NOT introduced by this change):** `TestRunInstallKimiMissingUVFailsBeforeExecutingInstallCommands` in `internal/cli/run_integration_test.go` fails on `main` before any of these changes (Kimi `uv` preflight install logic — unrelated to profiles/sync/domain). Captured in the safety-net baseline; left untouched per strict-TDD safety-net rules. It is the **only** failure in the full `go test ./...` run; every other package is `ok`, including all touched packages and all 15 new tests.

---

## Build / Test Results

- `go build ./...` → **exit 0** (clean)
- `go vet` on changed packages (sdd, cli, tui, tui/screens, model) → **exit 0** (clean)
- `go test ./...` → all packages `ok` except the single pre-existing Kimi failure noted above.
- New tests (consolidated): `TestEnsureProfileDomainConsistency` (7 subtests), `TestRunSyncDomainCollisionError`, `TestRunSyncDomainConsistency_AllowsSameEffectiveDomain`, `TestProfileCreatePrefillsFromDomain`, `TestProfileCreatePrefillSkippedWhenAssignmentsPresent`, `TestProfileCreatePrefillSkippedWhenDomainEmpty`, `TestProfileCreateEditModeDoesNotDetectDomain`, `TestRenderProfileCreate_Step1_ShowsDomainFooter`, `TestRenderProfileCreate_Step1_OmitsDomainFooterWhenEmpty` — **all pass**.

---

## Remaining Tasks

None. All 5 tasks complete. Recommended next step: **sdd-verify** (run `go test ./...` + `go vet ./...` and compare against every spec scenario), then **sdd-archive**.

## Workload / PR Boundary

- Mode: **single PR** (~310 new lines incl. tests, under the 800-line budget per the tasks.md Reviewer-Facing Notes).
- Current work unit: Tasks 1–5 (collision guard + TUI prefill + render).
- Boundary: complete change in one batch; no chained/stacked PR strategy needed.
