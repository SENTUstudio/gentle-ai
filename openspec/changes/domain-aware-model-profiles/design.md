# Design: Domain-Aware Model Profiles

## Technical Approach

This change bridges the existing `sddconfig.Config.Domain` field (from PR #2)
to the profile system by treating domain as **metadata on `model.Profile`**, not
a name suffix or JSON overlay field. Two surfaces consume it:

1. **Sync collision guard** — `DetectProfiles` reads `opencode.json` and returns
   existing profiles with `Domain = ""` (backward-compat). A new helper
   `EnsureProfileDomainConsistency(detected, explicit)` rejects the case where
   the same profile name appears in both lists with different domains.
2. **TUI prefill** — at step 0→1 in the profile create flow, detect the project
   domain via `sddconfig.LoadConfig(workspaceDir)`. If domain is non-empty AND
   `Selection.ModelAssignments` is empty, seed it with
   `model.DefaultModelsForDomain(domain)`. Stamp `ProfileDraft.Domain` so the
   domain is carried into the generated profile.

The CLI `--profile-domain` flag is **out of scope** for this slice per user
direction; the proposal lists it, but the user explicitly scoped remaining work
to collision guard + TUI prefill only. The CLI path will pick up domain
auto-detection in a follow-up.

## Architecture Decisions

| Decision | Choice | Alternative | Rationale |
|----------|--------|-------------|-----------|
| Where domain lives | Field on `model.Profile` | JSON overlay key, agent name suffix | Spec mandates metadata-only; `Domain=""` = app-dev preserves full backward compat. Suffix would explode the name space (`cheap-de`). |
| Collision check site | Dedicated helper called by sync + TUI | Inside `DetectProfiles` itself | `DetectProfiles` only knows the JSON (domain-blind). The conflict is between detected (`""`) and explicit (`"data-engineering"`) — needs both lists. |
| Backward-compat for `DetectProfiles` | Returned profiles get `Domain=""` | Try to infer domain from JSON | JSON has no domain field by design (Q3=A). Inferring would be guessing. |
| TUI prefill trigger | Step 0→1 transition | Screen entry, confirm step | Step transition is where the picker is initialized today; idempotent re-entry must not overwrite user changes. |
| Prefill condition | Empty `Selection.ModelAssignments` AND non-empty domain | Always overwrite | User who manually selected models must not be clobbered. |
| Domain in TUI footer | One-line indicator in step 1 header | Full tier legend table | Per proposal: "cosmetic". One line is enough signal. |
| Tier legend expansion | Defer to follow-up | Add Sonnet/Opus tier chips now | Cosmetic; out of remaining-scope budget. |

## Data Flow

```
ScreenProfileCreate step 0 → step 1 (enter on valid name)
        │
        ├─► osGetwdFn() → workspaceDir
        ├─► sddconfig.LoadConfig(workspaceDir) → cfg
        ├─► if cfg.Domain != "" && len(Selection.ModelAssignments) == 0:
        │       Selection.ModelAssignments = DefaultModelsForDomain(cfg.Domain)
        ├─► ProfileDraft.Domain = cfg.Domain
        │
        ▼
RenderProfileCreate(step=1, draft, assignments, ...) 
   → footer shows "Domain: <cfg.Domain>" (or omitted when "")
        │
        ▼
Continue → confirmProfileCreate()
        │
        ├─► copies assignments into ProfileDraft.{OrchestratorModel, PhaseAssignments}
        ├─► ProfileDraft.Domain already set from step 0→1
        │
        ▼
Confirm → PendingSyncOverrides.Profiles = [draft]
        │
        ▼
Sync → EnsureProfileDomainConsistency(detected, explicit) → error if conflict
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `internal/components/sdd/profiles.go` | Modify | Add `EnsureProfileDomainConsistency(detected, explicit []model.Profile) error` — returns formatted error listing conflicting `(name, detected-domain, explicit-domain)` tuples. |
| `internal/components/sdd/profiles_test.go` | Modify | New tests: no conflict (empty detected), no conflict (same domain), conflict (same name, different domain), multiple conflicts. |
| `internal/cli/sync.go` | Modify | After `sdd.DetectProfiles(settingsPath)` returns and `profiles` is assigned, call `EnsureProfileDomainConsistency(profiles, selection.Profiles)`. On error, return wrapped error from `RunSync`. |
| `internal/cli/sync_test.go` | Modify | Add `TestRunSyncDomainCollisionError` — writes opencode.json with `sdd-orchestrator-cheap`, calls `RunSync` with explicit `Profile{Name: "cheap", Domain: "data-engineering"}`, asserts error contains "domain". |
| `internal/tui/screens/profile_create.go` | Modify | Step 1 render: accept and display `domain string` in footer. `ProfileCreateOptionCount` unchanged. |
| `internal/tui/screens/profile_create_test.go` | Modify | New tests: step 1 shows domain footer when domain != "", omits when domain == "". |
| `internal/tui/model.go` | Modify | At step 0→1 transition (line ~3861), detect domain via `sddconfig.LoadConfig(osGetwdFn())`, prefill `Selection.ModelAssignments` if empty, set `ProfileDraft.Domain`. Pass `ProfileDraft.Domain` to `RenderProfileCreate`. |
| `internal/tui/model_test.go` | Modify | New test: prefill happens when domain set + assignments empty; prefill skipped when assignments non-empty; `ProfileDraft.Domain` stamped. |

## Interfaces / Contracts

```go
// EnsureProfileDomainConsistency returns an error if any profile name appears
// in both detected and explicit lists with different effective domains.
// "Effective domain" treats "" and "app-dev" as equivalent.
//
// The error message lists each conflicting (name, detected, explicit) triple
// so the user can resolve the collision by renaming the new profile.
func EnsureProfileDomainConsistency(detected, explicit []model.Profile) error
```

TUI render signature gains one parameter:

```go
func RenderProfileCreate(
    step int,
    draft model.Profile,        // draft.Domain drives the footer
    nameInput string, namePos int, nameErr string,
    editMode bool,
    assignments map[string]model.ModelAssignment,
    picker ModelPickerState,
    cursor int,
) string
```

No new public types. No JSON schema changes.

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit (sdd) | `EnsureProfileDomainConsistency` — 4 cases | Table-driven: no conflict, same domain, single conflict, multiple conflicts. |
| Unit (model) | `DefaultModelsForDomain` for all phases (already covered by 5 existing tests) | None new. |
| Unit (screens) | `RenderProfileCreate` step 1 with/without domain footer | Snapshot substring assertions. |
| Unit (TUI model) | Prefill on step 0→1, skip when assignments non-empty, `ProfileDraft.Domain` stamping | Construct `Model`, call `confirmProfileCreate()`, assert state. |
| Integration (sync) | CLI: detected `cheap` + explicit `cheap(data-engineering)` → error containing "domain" | Temp `opencode.json` fixture, call `RunSync`, assert error. |

## Migration / Rollout

No migration required. `Domain` is additive and zero-value for all existing
profiles. The collision guard fires only when a user actively tries to create a
profile whose name matches an existing one with a different domain — a path
nobody can hit today because no profile has a non-empty domain.

## Open Questions

None — the proposal locks decisions Q1–Q6 and the user has scoped remaining
work to two surfaces. The CLI `--profile-domain` flag is deferred.
