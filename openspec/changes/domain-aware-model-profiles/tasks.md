# Tasks: Domain-Aware Model Profiles

Scope per user: **collision guard in `DetectProfiles` consumer + TUI prefill**.
The CLI `--profile-domain` flag is deferred to a follow-up change. The proposal
also covers orchestrator prompt table rendering, which the design marks as
cosmetic and out of remaining scope.

Each task lists the deliverable, files touched, and tests required. Tasks are
ordered for incremental reviewability: pure helper → CLI integration → TUI
state → TUI render → regression tests.

---

## Task 1 — Domain-consistency helper in `sdd` package

> Status: ✅ Complete [x]

**Deliverable**: `EnsureProfileDomainConsistency(detected, explicit []model.Profile) error`
returns a formatted error listing every conflicting `(name, detected-domain,
explicit-domain)` triple. Detected and explicit profiles with the same name
AND same effective domain (treating `""` and `"app-dev"` as equivalent) are
NOT conflicts.

**Files**:
- `internal/components/sdd/profiles.go` — add function below `DetectProfiles`.
- `internal/components/sdd/profiles_test.go` — table-driven tests:
  - empty detected → no error
  - empty explicit → no error
  - same name, same domain (`""` / `""`) → no error
  - same name, same domain (`""` / `"app-dev"`) → no error (treat as equivalent)
  - same name, different domain → error mentioning name + both domains
  - multiple conflicts → error lists all

**Lines est.**: ~70 (40 impl + 30 tests).

---

## Task 2 — Sync CLI invokes collision guard

> Status: ✅ Complete [x]

**Deliverable**: After `sdd.DetectProfiles(settingsPath)` populates `profiles`
in `RunSyncWithSelection` (around `internal/cli/sync.go:628`), call the new
helper with `(profiles, s.selection.Profiles)`. On error, wrap with
`fmt.Errorf("sync sdd profile domain collision: %w", err)` and return from
`componentSyncStep`.

**Files**:
- `internal/cli/sync.go` — single call site, no signature changes.
- `internal/cli/sync_test.go` — `TestRunSyncDomainCollisionError`:
  1. Write temp `opencode.json` containing `sdd-orchestrator-cheap` agent
     (no domain in JSON, so `DetectProfiles` returns `Domain=""`).
  2. Build a `model.Selection` with `Profiles: []model.Profile{
        {Name: "cheap", Domain: "data-engineering", OrchestratorModel: ...}}`.
  3. Call `RunSyncWithSelection(home, sel)`; assert error wraps
     `EnsureProfileDomainConsistency` and message contains `"domain"` and
     `"cheap"`.
- `internal/components/sdd/profiles_test.go` — already covered by Task 1.

**Lines est.**: ~60 (5 impl + 55 test).

---

## Task 3 — TUI prefill at step 0→1 transition

> Status: ✅ Complete [x]

**Deliverable**: When the user advances from name input (step 0) to the model
picker (step 1) in the create flow (NOT edit), the TUI model:
1. Calls `sddconfig.LoadConfig(osGetwdFn())`.
2. If `cfg.Domain != ""` AND `len(m.Selection.ModelAssignments) == 0` (idempotency
   guard for re-entry), sets
   `m.Selection.ModelAssignments = model.DefaultModelsForDomain(cfg.Domain)`.
3. Sets `m.ProfileDraft.Domain = cfg.Domain` (always, even on re-entry, so the
   final profile carries the stamp).
4. Edit mode (step 0) leaves domain detection untouched — `ProfileDraft.Domain`
   is already populated from the loaded profile.

**Files**:
- `internal/tui/model.go` — modify the step 0→1 branch in
  `confirmProfileCreate()` (around line 3861). The existing `ModelPicker`
  initialization block stays; add domain detection before it.
- `internal/tui/model_test.go` — `TestProfileCreatePrefillsFromDomain`:
  1. Construct `Model` with `ProfileCreateStep = 0`, empty
     `Selection.ModelAssignments`, `ProfileDraft.Name = "cheap"`.
  2. Stub `osGetwdFn` to a temp dir containing `openspec/config.yaml` with
     `domain: data-engineering`.
  3. Call `confirmProfileCreate()`.
  4. Assert `Selection.ModelAssignments["sdd-explore"].ModelID == "claude-sonnet-4-20250514"`.
  5. Assert `ProfileDraft.Domain == "data-engineering"`.
  - `TestProfileCreatePrefillSkippedWhenAssignmentsPresent`:
    1. Set `Selection.ModelAssignments["sdd-apply"] = ...` (any non-empty).
    2. Run same path; assert `Selection.ModelAssignments` unchanged.
  - `TestProfileCreatePrefillSkippedWhenDomainEmpty`:
    1. No `openspec/config.yaml` in temp dir.
    2. Run; assert `Selection.ModelAssignments` stays empty.
  - `TestProfileCreateEditModeDoesNotDetectDomain`:
    1. `ProfileEditMode = true`, step 0.
    2. Assert `ProfileDraft.Domain` not overwritten (caller already populated it
       from the loaded profile in the screen-transition handler).

**Lines est.**: ~140 (30 impl + 110 tests).

---

## Task 4 — Render domain indicator in step 1 footer

> Status: ✅ Complete [x]

**Deliverable**: `screens.RenderProfileCreate` step 1 appends a single
`SubtextStyle` line — `Domain: <draft.Domain>` — between the existing
`Assign Models` heading and the picker, but only when `draft.Domain != ""`.
Step 0 and step 2 are unchanged. No changes to `ProfileCreateOptionCount`.

**Files**:
- `internal/tui/screens/profile_create.go` — extend
  `renderProfileModelStep` to read `draft.Domain` and render the indicator
  when non-empty. Caller already passes the full `draft` today (line 35),
  so no signature change is needed; the helper must accept it. Cleanest
  refactor: pass `draft model.Profile` as a new parameter to
  `renderProfileModelStep` (currently takes only `profileName`).
- `internal/tui/screens/profile_create_test.go` — two new tests:
  - `TestRenderProfileCreate_Step1_ShowsDomainFooter`: `draft.Domain = "data-engineering"`,
    assert output contains `"data-engineering"`.
  - `TestRenderProfileCreate_Step1_OmitsDomainFooterWhenEmpty`: `draft.Domain = ""`,
    assert output does NOT contain `"Domain:"`.

**Lines est.**: ~40 (10 impl + 30 tests).

---

## Task 5 — Wire `ProfileDraft.Domain` to render call

> Status: ✅ Complete [x] (no-op verification — field plumbs through the existing `draft` parameter at `model.go:979`)

**Deliverable**: The TUI `Model.View()` pass-through to
`screens.RenderProfileCreate` already includes `m.ProfileDraft` (line 979),
so the new field is automatically available to the render. **No code change
required** — this task is the verification step.

**Files**:
- None modified. Re-run `make build` and `go test ./internal/tui/...` to
  confirm the field plumbs through cleanly. Document the dependency in the
  PR description: "Task 5 is a no-op verification; the new
  `ProfileDraft.Domain` is read by the renderer in Task 4 via the existing
  `draft` parameter."

**Lines est.**: 0 (verification only).

---

## Out of Scope (Deferred)

The proposal lists these as in-scope but they are NOT part of the remaining
work the user requested:

- CLI `--profile-domain name:domain` flag
- Auto-detection in `parseProfileFlags` from `sddconfig`
- Orchestrator prompt model table domain column (`renderProfileModelAssignmentsSection`)
- Domain-aware tier legend in the model picker (cosmetic)
- Encoding `Domain` in `opencode.json` overlay (explicitly rejected by Q3=A)
- YAML-declared per-domain defaults (explicitly rejected by Q2=B)

These are tracked in the proposal's "Out of Scope" section and will be
picked up in a follow-up change if the user requests them.

---

## Reviewer-Facing Notes

- **PR size**: Tasks 1–4 land as one PR (~310 new lines incl. tests). Under
  the 800-line budget.
- **No JSON schema change**: `opencode.json` is untouched. `Profile.Domain`
  lives only in Go state and is lost on a hand-edit + sync re-detect
  (accepted v1 trade-off, documented in the proposal's Risks table).
- **No migration**: zero-value `Domain=""` is app-dev; every existing
  profile, every existing test, every existing user behavior is preserved.
- **One pure helper + two call sites**: the collision guard is a single
  30-line function in the `sdd` package, called from exactly two places
  (sync + TUI). Easy to audit, easy to remove.
