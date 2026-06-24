# Exploration: Domain-Aware Model Profiles

> **Status**: working exploration — produced by `sdd-explore` for change
> `domain-aware-model-profiles`. Connects the existing model-profile
> system (`model.Profile` + `components/sdd/profiles.go`) to the
> domain config (`sddconfig.Config.Domain`) introduced by the
> already-merged `data-engineering-domain-profile` change.

## Current State

Two systems exist today and do not talk to each other:

1. **Model profiles** (`internal/model/types.go:193`,
   `internal/components/sdd/profiles.go:1-764`) — a named bundle of
   model assignments (orchestrator + per-phase sub-agent). Lives in
   `opencode.json` as suffixed agent keys
   (`sdd-orchestrator-{name}`, `sdd-{phase}-{name}`). Domain-blind.
2. **Domain config** (`internal/sddconfig/config.go:16-21`,
   already merged in PR #2 / `data-engineering-domain-profile`) —
   parses `openspec/config.yaml` into `Config{Domain, Repos,
   AWSProfiles, Verify}`. `Domain` is `""` for app-dev and
   `"data-engineering"` when set. Read by every SDD phase to branch
   skill content, but NOT read by the model-profile system.

The gap: a profile called `cheap` assigns Haiku to `sdd-explore` and
Sonnet to `sdd-spec` REGARDLESS of whether the project is
app-dev or data-engineering. Per `docs/sdd-profiles-domain-model.md`
and the data-engineering-domain-profile exploration, the recommended
tier for `sdd-explore` in data-engineering is Sonnet (not Haiku)
because data profiling needs judgment about encoding, dates, and
types. Today the user has no way to express that delta.

### What's already solved (don't re-design)

- **Profile struct** (`#1344` decision): add `Domain string` to
  `model.Profile`. Backward-compat: empty = app-dev (today's
  behavior). Empty keeps the field zero-value and every existing
  test continues to pass.
- **Profile overlay generation** (`profiles.go:263-406`): the
  `GenerateProfileOverlay` function builds 11 suffixed agent
  entries. The structure is already domain-friendly — the orchestrator
  prompt's model-assignments table is rendered from
  `renderProfileModelAssignmentsSection` (`profiles.go:649-702`)
  which iterates `profile.PhaseAssignments` and writes a row per
  phase. Adding domain does not change overlay JSON structure; it
  changes the *values* in the table and in each `model:` field.
- **Profile detection** (`profiles.go:156-226`): `DetectProfiles`
  reads `sdd-orchestrator-{name}` keys from `opencode.json`. The
  domain is not encoded in `opencode.json` today; it is purely
  per-profile metadata in Go.
- **Domain config** (already shipped): `sddconfig.LoadConfig(root)`
  is the single reader. `Config.Domain` is the authoritative
  source. TUI/CLI can call it cheaply.
- **CLI profile flow** (`internal/cli/sync.go:142-276`): the
  `--profile name:provider/model` and `--profile-phase
  name:phase:provider/model` flags already parse into
  `[]model.Profile`. New fields plug in here.

### Affected Areas

- `internal/model/types.go` (lines 190-197) — add `Domain string`
  to the `Profile` struct. Add a small comment block citing the
  backward-compat rule (empty = app-dev).
- `internal/components/sdd/profiles.go` (lines 263-406, 649-702) —
  new helper `DefaultModelsForDomain(domain string) (orchestrator,
  map[string]ModelAssignment)` consumed by the TUI pre-fill path
  AND by the `DetectProfiles` round-trip; orchestrator prompt's
  model-assignments table needs a `Domain` column OR a separate
  per-domain sub-section.
- `internal/components/sdd/profiles_test.go` and
  `profiles_lifecycle_test.go` — every existing `makeHaikuProfile()`
  helper must be untouched (empty Domain = app-dev = today's
  behavior). Add new tests for the data-engineering path and the
  "domain-aware prefill" path.
- `internal/cli/sync.go` (lines 142-276, `parseProfileFlags`) — add
  `--profile-domain name:domain` and/or accept a new
  `name:domain:provider/model` form. Backward-compat: existing
  `name:provider/model` form MUST keep working unchanged.
- `internal/tui/screens/profile_create.go` (lines 20-186) — step 0
  (name input) gains a domain step OR a domain detection line that
  reads `openspec/config.yaml` via `sddconfig.LoadConfig`; step 1
  (model picker) pre-fills with `DefaultModelsForDomain(detected)`.
- `internal/tui/screens/model_picker.go` — picker's "Set all phases"
  row needs to be aware of the profile's domain so its label can
  show the domain-appropriate default tier (cheap/mid/premium
  legend). Purely cosmetic.
- `internal/components/sdd/inject.go:513-536` — the loop that
  calls `GenerateProfileOverlay` for each `Selection.Profiles`
  entry is unchanged; the overlay generation itself absorbs the
  domain.
- `openspec/specs/sdd-profiles/spec.md` — the spec governing
  `Profile` CRUD needs an "ADDED Domain field" delta requirement +
  "ADDED Default models per domain" scenario. The spec is
  authoritative: a domain-aware `Profile` round-trip
  (create → serialize → detect → re-edit) MUST preserve the
  domain.

---

## The 6 Architectural Questions

For each question I surface **concrete options** with the trade-off
called out. The recommendation is in the next section.

### Q1. Where does the domain live — Profile struct, profile name, or both?

How is domain encoded in the persisted profile?

| Option | Description | Pros | Cons |
|--------|-------------|------|------|
| **A. Profile struct field only** | Add `Domain string` to `model.Profile`. The name stays as the user typed it (`cheap`, `premium`). | Name is a clean slug; domain is metadata; one profile can be edited across domains; `opencode.json` overlay is unchanged | Two profiles named `cheap` cannot coexist (one app-dev, one data-eng) unless we forbid or auto-suffix |
| **B. Profile name suffix** | The user must name the profile `cheap-app-dev` or `cheap-data-eng`. Domain is derivable from the name. | No struct change; existing detection works; one source of truth (the name) | Ugly names; CLI/TUI must enforce the convention; no way to migrate a `cheap` profile to data-eng without rename |
| **C. Both — struct field is authoritative, name is free** | `Profile.Domain` is the source of truth. The name is whatever the user types. `opencode.json` overlay does NOT encode the domain (the domain only lives in Go state). | Cleanest separation; one source of truth; TUI/CLI/Go code all read from `Profile.Domain`; tests are simple | Two profiles with the same name silently overwrite each other; need a UX guard for this |

**Recommended**: **C**. The domain is metadata, not a name component.
Encoding it in the name would force users to rename profiles they
have already configured. The struct field is the simplest, most
testable, most Go-idiomatic choice. The UX guard (Q6 below) handles
the name-collision edge case.

### Q2. How are per-domain default model recommendations resolved?

When the TUI opens the create flow and `openspec/config.yaml` says
`domain: data-engineering`, what models does it pre-fill?

| Option | Description | Pros | Cons |
|--------|-------------|------|------|
| **A. Hard-coded Go function** | `DefaultModelsForDomain("data-engineering")` returns a fixed `map[string]ModelAssignment` from the table in `docs/sdd-profiles-domain-model.md` | Single source of truth; pure function = trivial unit tests; no new files | Hard to evolve; user can't override defaults per project |
| **B. Reads from `sddconfig` config** | A new `openspec/config.yaml` block like `model_profiles: { cheap: { orchestrator: ..., phases: { ... } } }`; user can pre-declare per-domain defaults | Maximum flexibility; power-user friendly | YAML schema growth; conflation of "domain config" and "model profile config"; the `data-engineering-domain-profile` change is fresh — adding more YAML now invites premature scope |
| **C. Hybrid — Go defaults + TUI override saved per profile** | Go has the canonical defaults per domain; TUI shows them as pre-selected; user can change any row; the resulting `Profile.PhaseAssignments` is what gets persisted | Defaults are deterministic, override is local; tests cover the defaults; UX is clear (you see what the system recommends) | The "recommendation" is Go-hardcoded; if the user wants to start from a different default for a new profile, they still need to walk through the picker |

**Recommended**: **C**. Defaults come from a pure Go function
(`DefaultModelsForDomain`) sourced from the
`docs/sdd-profiles-domain-model.md` table. The TUI pre-fills
those defaults, the user can override any row, and the resulting
`Profile.PhaseAssignments` is what is persisted. The defaults
table is the SAME source of truth used by the orchestrator prompt
rendering, so the two systems cannot drift.

Per the doc, the data-engineering defaults are:

| Phase | App-Dev default | Data-Eng default |
|-------|-----------------|------------------|
| orchestrator | Sonnet (per current TUI) | Sonnet |
| sdd-init | Haiku | Haiku (still cheap) |
| sdd-explore | Haiku | **Sonnet** (data profiling = judgment) |
| sdd-propose | Sonnet | Sonnet |
| sdd-spec | Sonnet | **Opus** (schema + DAG = precision) |
| sdd-design | Opus | **Opus** (insertion-point cascade) |
| sdd-tasks | Sonnet | Sonnet |
| sdd-apply | Sonnet | **Sonnet** (Spark SQL translation) |
| sdd-verify | Sonnet | **Sonnet** (Athena EXCEPT interpretation) |
| sdd-archive | Haiku | Haiku (mechanical) |
| sdd-onboard | Sonnet | Sonnet |

The bolded rows are the meaningful deltas. `sdd-explore` and
`sdd-spec` are the two that move the needle (Haiku→Sonnet and
Sonnet→Opus respectively).

### Q3. Does the opencode.json overlay encode the domain?

When the profile is written to `opencode.json`, does the overlay
JSON carry the domain field anywhere?

| Option | Description | Pros | Cons |
|--------|-------------|------|------|
| **A. Domain is Go-only state** | `Profile.Domain` is persisted via the existing JSON path (`Selection.Profiles` → inject → overlay → opencode.json). The overlay JSON has no domain marker. | The overlay schema is unchanged; backward-compat for hand-edits | On read-back, `DetectProfiles` cannot recover the domain from `opencode.json`; the user must remember the domain for each profile |
| **B. Embed in the orchestrator agent's `description` field** | `description: "SDD Orchestrator ({name} / domain=data-engineering profile)"` | Human-readable; survives hand-edits | String parsing is fragile; `DetectProfiles` would need regex |
| **C. New sidecar field on the orchestrator entry** | Add `"domain": "data-engineering"` to the `sdd-orchestrator-{name}` JSON object. | Structured; testable; survives hand-edits | Schema growth; new field; could be hand-stripped |

**Recommended**: **A** for v1. The domain lives in Go
(`Selection.Profiles[*].Domain` → `Profile.Domain` in the loop in
`inject.go:513-536`). The overlay JSON is unchanged. The
trade-off: if a user hand-edits `opencode.json` and removes a
profile, then re-syncs, the domain is lost. We accept this v1
because:
- Profiles are managed artifacts, not hand-edited config.
- The next sync that creates the profile from CLI/TUI will set
  the domain explicitly.
- A future "v2" can promote the domain to a JSON sidecar without
  breaking v1 profiles (additive).

This is a strict-scope decision: the change is "add a struct field
and pre-fill logic"; the JSON schema is a follow-up if needed.

### Q4. CLI flag shape

How does the user pass a domain via the CLI?

| Option | Description | Pros | Cons |
|--------|-------------|------|------|
| **A. New `--profile-domain name:domain`** | Separate flag, like `--profile-phase`. Mirrors the existing pattern. | Familiar; explicit; can combine with `--profile` for orchestrator and `--profile-phase` for overrides | Three flags now (profile, profile-phase, profile-domain) |
| **B. Extend `--profile` to `name:domain:provider/model`** | Repurpose the colon-separated form to include domain in the middle. | One flag for everything | Breaks every existing `--profile` invocation; the parser at `sync.go:188-212` is well-tested |
| **C. Auto-detect from `sddconfig` if flag is absent** | Domain is optional in CLI. If user runs in a data-eng project and types `--profile cheap:haiku`, the resulting Profile gets `Domain = "data-engineering"` (read from `openspec/config.yaml`). | Zero CLI surface growth; the user just types the same flags they already use; the system is "smart" | Magic; the user might be surprised by the domain; headless sync without `openspec/config.yaml` defaults to app-dev |

**Recommended**: **C (auto-detect from sddconfig) + A (explicit
flag for override)**. The CLI mirrors the TUI's behavior: the
domain comes from `sddconfig.LoadConfig` unless the user
explicitly overrides it via `--profile-domain`. The default
path is zero-friction; the explicit path exists for CI and
headless sync where the user wants to be explicit.

The CLI parser needs:
- New `rawProfileDomains []string` in `SyncFlags`.
- New `parseProfileDomainFlag("name:domain")` mirroring
  `parseProfileFlag`.
- `parseProfileFlags` augments the resulting `Profile.Domain`
  from the auto-detected config OR the explicit override.
- Backward-compat: if `rawProfileDomains` is empty AND
  `sddconfig.LoadConfig` returns `Domain == ""`, the field is
  empty (today's behavior).

### Q5. TUI flow change

How does the TUI step 0 → step 1 transition work?

| Option | Description | Pros | Cons |
|--------|-------------|------|------|
| **A. New step 0.5 — domain selector** | After name, TUI shows "Domain: app-dev (auto-detected) / data-engineering" with two options; user confirms or picks. | Explicit; mirrors the data-eng preflight question; user knows | One more screen in a flow already at 3 steps |
| **B. Auto-detect + show in step 0 footer** | Step 0 detects the domain from sddconfig, displays "Detected domain: data-engineering — pre-filling Sonnet for sdd-explore" in the subtext. User can change via a "Override domain" keybind. | Compact; non-blocking; default path is zero-friction | The user might miss the footer; the override path needs a keybind |
| **C. Domain is a property of the project, not the profile** | The TUI does not ask. The domain is read from sddconfig at every TUI render and applied automatically. The `Profile.Domain` is set to whatever the project says at creation time. | Simplest; no new UI; matches "domain = project, profile = model tier" | The user might want to create an app-dev profile in a data-eng project (rare but possible for cross-cutting work) |

**Recommended**: **C as the default, B as the override path**.
The TUI does not add a domain step. At step 1, the model
picker pre-fills with `DefaultModelsForDomain(sddconfig.Domain)`,
and the orchestrator prompt section in step 2 (summary) shows
"Domain: data-engineering (pre-filled)" so the user sees what
was applied. If the user wants to override, they edit the
individual phase rows in the model picker (today's flow).
This is the smallest TUI delta and the most "it just works"
UX.

### Q6. Naming collision guard

Two profiles with the same name but different domains — what
happens?

| Option | Description | Pros | Cons |
|--------|-------------|------|------|
| **A. Allow + require user to delete one** | Same name + different domain = both persist; `opencode.json` overlay writes the LAST ONE and clobbers the first | None — the user loses data | Clobbering is silent data loss; the orchestrator's `sdd-orchestrator-{name}` key is one slot |
| **B. Disallow — domain forces a unique name** | The TUI suggests `cheap-data-eng` when the user tries to create `cheap` in a data-eng project that already has a `cheap` profile. | No clobber; clear UX | The TUI is opinionated; a power user might want to live with the collision |
| **C. Allow + warn, but the user owns the decision** | TUI shows "Profile 'cheap' already exists. Domain: app-dev. Create a separate data-eng profile? (yes/no)". If user picks yes, name must differ. | Explicit; user owns | One more confirmation in the create flow |

**Recommended**: **C (allow + warn + require disambiguation on
confirm)**. The TUI surfaces the collision in step 2 (summary) or
earlier, and the user must either:
- Pick a different name for the new profile
- Or overwrite the existing one (if the new domain matches
  the existing profile's domain)

The detection happens in the TUI: at step 1, if any existing
profile has the same `Name` AND a different `Domain`, show a
warning. At step 2, refuse to confirm without explicit
disambiguation. This is purely TUI logic; the Go profile
system remains dumb (same name = overwrite, last writer wins)
and the test surface is unchanged.

The CLI mirrors this: if `--profile cheap:haiku` is passed and
`cheap` already exists with a different domain, return a
clear error and require the user to either delete the
existing profile or rename the new one.

---

## Approaches — Summary Matrix

| Q | Recommended | Main risk |
|---|-------------|-----------|
| 1. Where domain lives | C — struct field, not name | Name-collision UX edge case |
| 2. Default models | C — Go defaults + TUI prefill | Defaults drift from doc |
| 3. JSON encoding | A — Go-only state | Domain lost on hand-edit |
| 4. CLI flag | C + A — auto-detect + explicit override | Magic surprise factor |
| 5. TUI flow | C + B — no new step, footer + override | User misses the footer |
| 6. Naming collision | C — warn + disambiguate | UX friction |

### Why these recommendations form a coherent set

- **Q1=C + Q3=A** together mean: domain is Go-only state. The
  Profile struct carries it; the JSON overlay does not. This is
  the smallest possible change to the overlay pipeline.
- **Q2=C + Q5=C** together mean: the TUI pre-fills with Go
  defaults, the user can override any row, and the
  resulting `Profile.PhaseAssignments` is what persists. The
  prefill is a UX nicety, not a contract.
- **Q4=C + Q5=C** together mean: both the CLI and the TUI read
  domain from `sddconfig` by default. Same code path, same
  source of truth, no new surface area for the user to learn.
- **Q6=C** ties the whole thing together: when the user
  attempts to create a profile that would clobber an existing
  one with a different domain, they get a clear warning. This
  is the only place where the domain field is "user-visible" —
  in the disambiguation confirmation.

---

## Recommendation

Adopt all 6 recommendations as a single coherent change,
implemented in 3 small PRs (well under the 400-line budget each):

**PR 1 — Profile struct + Go defaults** (≈150 lines)
- Add `Domain string` to `model.Profile` with backward-compat
  comment.
- New pure function `DefaultModelsForDomain(domain string)
  (orchestrator ModelAssignment, phases map[string]ModelAssignment)`
  in `internal/model/profiles_default.go` (new file). Table-driven
  tests in `_test.go`.
- Existing `Profile` construction sites unchanged (zero-value
  Domain = app-dev). All existing tests pass without modification.

**PR 2 — TUI prefill + collision guard** (≈200 lines)
- New step 0.5 helper that reads `sddconfig.LoadConfig(cwd)` and
  shows the detected domain in the step 0 footer.
- Step 1 model picker pre-fills with
  `DefaultModelsForDomain(detectedDomain)`.
- Step 2 summary shows "Domain: data-engineering (pre-filled,
  can override in model picker)".
- Collision check at step 2: if a profile with the same Name
  and different Domain exists, show a confirmation prompt.
- Tests: unit tests for the prefill + collision logic; snapshot
  test for the rendered step 2.

**PR 3 — CLI + spec** (≈150 lines)
- New `--profile-domain name:domain` flag in `SyncFlags`.
- `parseProfileFlags` augments `Profile.Domain` from
  `sddconfig.LoadConfig` (or the explicit override).
- `openspec/specs/sdd-profiles/spec.md` gets a delta: ADDED
  "Profile Domain" requirement + "Default models per domain"
  scenario.

Total: ~500 lines across 3 PRs, each under 200. Strict TDD on
all Go functions; spec scenarios are validated by the change's
verify phase, not by `go test`.

### Why this is small

The data-engineering-domain-profile change was 800+ lines and
touched 7+ skills + a new `sddconfig` Go core. THIS change is
~500 lines and touches:
- 1 struct field
- 1 new pure function
- 1 new CLI flag
- 1 TUI helper
- 1 spec delta

The reason it's small: the heavy lifting (sddconfig parsing,
domain detection, profile overlay generation, agent naming,
permission scoping) is already done. This change is the
**bridge** between two systems that already exist.

---

## Risks

- **Defaults drift from doc**: if
  `DefaultModelsForDomain("data-engineering")` returns a different
  model than `docs/sdd-profiles-domain-model.md` says, the
  orchestrator's model-assignments table lies. *Mitigation*: the
  function is pure, the doc is the spec, and the test asserts
  each (phase, model) pair against the doc table. CI fails if
  they diverge.
- **Backwards compat in CLI**: if `--profile
  cheap:anthropic/claude-haiku` is parsed today and starts
  producing a different Profile.Domain based on the cwd, the
  user might be surprised. *Mitigation*: only the **profile's**
  domain changes; the orchestrator model and phase
  assignments are still parsed from the existing flag forms.
  The user sees the new domain in `gentle-ai sync --json` or
  the TUI confirmation step.
- **Naming collision in `opencode.json`**: two profiles named
  `cheap` with different domains cannot both persist in
  `opencode.json` — there's only one `sdd-orchestrator-cheap`
  slot. *Mitigation*: the TUI's collision guard catches this
  at creation time; the CLI returns a clear error.
- **`sddconfig` not initialized in headless mode**: if the user
  runs `gentle-ai sync --profile cheap:haiku` in a directory
  with no `openspec/config.yaml`, the auto-detect path returns
  empty Domain = app-dev. *Mitigation*: this is the documented
  behavior; users in data-eng projects will have
  `openspec/config.yaml` (it's required by the data-eng domain
  profile). Headless CI without config = app-dev path, which is
  the safe default.
- **Profile edit flow leaks old domain**: if a user created a
  profile in app-dev, then changed the project to data-eng, the
  edit flow shows the OLD domain. *Mitigation*: at edit time,
  the TUI re-reads sddconfig and offers "Update to current
  domain?" as an explicit prompt. This is a small follow-up if
  the user wants it; v1 keeps the existing domain and the
  user can manually update via the model picker.
- **Spec scenario for round-trip**: the existing
  `sdd-profiles` spec does not test the round-trip
  (create → serialize → re-detect). The new "Profile Domain"
  requirement MUST add a round-trip scenario because the
  domain is not in the JSON. *Mitigation*: the scenario tests
  the Go state, not the JSON. If we promote the domain to
  the JSON in v2, the spec gets a new scenario then.

---

## Ready for Proposal

**Yes — with 1 user-confirmation gate before `sdd-propose` runs:**

The orchestrator should ask the user to confirm the 3-PR
delivery strategy. The trade-off is: one larger PR (~500 lines
total) vs. three small PRs (~150-200 each). The user owns the
risk-vs-velocity call. My recommendation is the 3-PR slice
because:
- Each PR is under the 400-line review budget.
- PR 1 (struct + defaults) is independent and reviewable on
  its own.
- PR 2 (TUI) can land after PR 1; doesn't touch CLI.
- PR 3 (CLI + spec) is the documentation + test surface;
  can land after PR 1+2 are stable.

After the user confirms, the orchestrator runs `sdd-propose`
with change name `domain-aware-model-profiles`. The proposal
should:
- State the 3-PR delivery and explicit per-PR line budget.
- Reference this exploration as the source of decisions.
- List the 6 user-confirmed options in the proposal body.
- Recommend `delivery_strategy: chained` with
  `chain-strategy: stacked-to-feature/aws-dataengineer`
  (PR #2's branch is the landing target).

The proposal artifact should be written to
`openspec/changes/domain-aware-model-profiles/proposal.md`.
