# Exploration: SDD State Review After `main` Reset to Upstream Mirror

> **Working topic only** — not a formal SDD change. Produced by `sdd-explore` to surface
> options after the `main` branch was force-pushed to mirror upstream
> `Gentleman-Programming/gentle-ai` exactly at SHA `510751e`.

## Current State

`main` is an exact mirror of upstream at `510751e` (v1.41.0+). The 7 custom
data-engineering skills (Toyota Chile ETL) live on `feature/kiro-cli` and
`feature/minimax-dataengineer`; the `toyota-chile-etl-template/` is untracked
local scratch. The SDD init was just re-run (`openspec/config.yaml` modified,
skill registry rebuilt).

`openspec/changes/` contains 9 active changes and 8 archived changes.
**Every active change has its implementation already on `main`** (verified via
`git log` and `internal/` package presence). The work is shipped; the SDD
artifacts are not closed.

### Active changes (9)

| # | Change | Files | Tasks `[x]/[ ]` | Verify | Apply-progress | Implementation on `main`? | Status |
|---|--------|-------|------------------|--------|----------------|----------------------------|--------|
| 1 | `agent-builder` | proposal/design/spec/tasks | 0 / 29 | ❌ | ❌ | ✅ `internal/agentbuilder/` (v1.18.0, PR #223) | **Ready** (stale checkboxes) |
| 2 | `antigravity-2.0-compatibility` | exploration/proposal/design/specs/tasks | 15 / 0 | ❌ | ❌ | ✅ `internal/agents/antigravity/` (v1.30.7) | **Ready** |
| 3 | `declarative-picker-navigation` | proposal/design/spec/tasks | 35 / 0 | ❌ | ❌ | ✅ `internal/tui/model.go` (v1.41.0) | **Ready** |
| 4 | `fix-persona-artifact-language-contract` | proposal/design/spec/tasks/apply-progress | 82 / 0 | ❌ | ✅ | ✅ `internal/assets/*/sdd-orchestrator.md`, `internal/components/persona/`, `internal/model/types.go:gentleman-neutral-artifacts` (v1.34.0+v1.37.1) | **Ready** (apply-progress proves completion) |
| 5 | `fix-skill-duplication-and-migrate` | proposal/design/spec/tasks/verify-report | — | ✅ PASS | ❌ | ✅ (PR #458 — LLM-first refactor + picker flags) | **Ready** (verify-report says "Next: sdd-archive") |
| 6 | `hermes-agent-support` | proposal/design/specs/tasks | 17 / 33 | ❌ | ❌ | ✅ `internal/agents/hermes/`, `internal/components/filemerge/yaml.go` (PR #774) | **Partial** (stale checkboxes) |
| 7 | `opencode-sdd-profiles` | proposal/design/tasks/specs/verify-report | 38 / 0 | ✅ PASS w/ warnings | ❌ | ✅ `internal/components/sdd/profiles.go`, `prompts.go` | **Ready w/ 3 minor warnings** |
| 8 | `qwen-code-integration` | proposal/design/spec/tasks/verify-report | 0 / 30 | ✅ PASS | ❌ | ✅ `internal/agents/qwen/` (v1.19.0, PR #263) | **Ready** (verify-report PASS) |
| 9 | `update-experience` | proposal/design/tasks/specs | — | ❌ | ❌ | ✅ Slice 1 + PRs #875/#876/#877/#879/#880 (v1.41.0) | **Ready** (multi-slice, all merged) |

### Archived changes (8)

| Date | Change | Has archive-report.md? | Notes |
|------|--------|------------------------|-------|
| 2026-03-25 | `gga-powershell-support` | — | Standard pre-v2 archive |
| 2026-05-05 | `contextual-skill-loading` | ✅ | Most complete |
| 2026-05-14 | `trae-agent-support` | — | Sparse: only proposal/spec/verify-report. No design/tasks/archive-report. Pre-existing oddity. |
| 2026-05-18 | `all-agent-orchestrator-parity` | ✅ | Standard |
| 2026-05-18 | `claude-opencode-orchestrator-parity` | ✅ | Standard |
| 2026-06-09 | `bind-chained-pr-skill-to-orchestrators` | ✅ | Standard |
| 2026-06-10 | `level-neutral-persona-parity` | ✅ | Related to active `fix-persona-artifact-language-contract` (different change, different scope) |
| 2026-06-14 | `organic-agent-trigger-rules` | — | Sparse: proposal/design/specs/tasks. No verify-report/archive-report. Archived via `b6872c6`. Pre-existing oddity. |

The sparse archives (`trae-agent-support`, `organic-agent-trigger-rules`) were
archived before the current archive-report convention was tightened. They are
correct audit trails, just less verbose than the v2.0 ones.

### Data-engineering staleness

None. The 7 Toyota Chile ETL skills (`data-engineer-*`) live on
`feature/kiro-cli` and `feature/minimax-dataengineer`. They are not on `main`
and are not represented under `openspec/changes/`. The
`toyota-chile-etl-template/` directory is untracked scratch. **No SDD
staleness from the data-engineering work** — the custom fork work that was on
`main` before the reset left no SDD artifacts behind.

## Affected Areas

- `openspec/changes/{9 change folders}/` — all 9 active changes are open
- `openspec/config.yaml` — modified (uncommitted, post-init); mtime fresh
- `openspec/specs/` — may need delta merges for 9 changes (handled by `sdd-archive`)
- `openspec/changes/archive/` — 8 existing, 9 expected to land here
- Working tree — 1 modified file (`openspec/config.yaml`) + 1 untracked dir
  (`toyota-chile-etl-template/`) — both unrelated to the SDD cleanup

## Approaches

### A. Archive all 9 active changes in one batch

**Idea**: Run `sdd-archive` on every change, with a mixed strategy:
- For 3 with `verify-report` (fix-skill-duplication-and-migrate, qwen-code-integration, opencode-sdd-profiles): direct archive.
- For 1 with `apply-progress` (fix-persona-artifact-language-contract): archive with stale-checkbox reconciliation proof.
- For 5 without either (agent-builder, antigravity-2.0-compatibility, declarative-picker-navigation, hermes-agent-support, update-experience): run `sdd-verify` first to produce verify-reports, then archive.

- **Pros**: Maximally clean. Single recommended action for the user. Honors SDD cycle completion. Brings the active change count to 0.
- **Cons**: High token cost (9 `sdd-verify` runs + 9 `sdd-archive` runs). Some `sdd-verify` runs will be 1-line "verified by git log" reports. The sdd-archive skill may block on stale checkboxes for 3 changes (agent-builder, qwen-code-integration, hermes-agent-support) — these need explicit "intentional partial archive" or stale-checkbox reconciliation override.
- **Effort**: High (1 orchestrator session + 14 sub-agent runs)

### B. Archive only the 3 changes with verify-reports; leave the rest as-is

**Idea**: Minimum-risk path. Archive fix-skill-duplication-and-migrate, qwen-code-integration, opencode-sdd-profiles. The other 6 stay open and the user can decide later.

- **Pros**: Lowest risk. No overrides needed. No verify work.
- **Cons**: Leaves 6 stale changes in `openspec/changes/`. The SDD state remains inconsistent with reality.
- **Effort**: Low (3 sdd-archive runs)

### C. Archive in 2 waves: clean ones first, then handle the partials

**Idea**: Wave 1 = 5 fully-clean (antigravity-2.0-compatibility, declarative-picker-navigation, fix-persona-artifact-language-contract, fix-skill-duplication-and-migrate, qwen-code-integration). Wave 2 = 4 with caveats (agent-builder, hermes-agent-support, opencode-sdd-profiles, update-experience) where each gets an explicit user-acknowledged partial archive / stale-checkbox reconciliation.

- **Pros**: Each wave is reviewable. Clean changes are not blocked by dirty ones.
- **Cons**: Still high token cost. Two user interactions instead of one.
- **Effort**: Medium-high (1 orchestrator session + 9 sub-agent runs + 1 user decision)

## Recommendation

**Approach A** — archive all 9 in one batched cleanup. Reasons:

1. **All 9 are objectively done.** Implementation is on `main` for every one of them. The SDD artifacts are an audit-trail gap, not a work gap. The cost of leaving 6 of them open indefinitely (cognitive load, accidental edits, ambiguous status) is higher than the cost of one cleanup session.
2. **The 3 verify-report changes are unambiguous.** `fix-skill-duplication-and-migrate` and `qwen-code-integration` are pure PASS. `opencode-sdd-profiles` has 3 minor warnings (R-PROF-31 missing warning log, ScreenProfileCreate missing-cache guard descope, tasks.md checkboxes) — all are descope candidates, not blockers.
3. **The 6 without verify-reports can be verified inexpensively.** Since implementation is provably on `main` via `git log`, each `sdd-verify` is a 1-pass "static evidence only" report — no test runs needed. That's the cheap path.
4. **The stale-checkboxes are documentation debt, not blockers.** The sdd-archive skill explicitly allows stale-checkbox reconciliation "backed by apply-progress/verify-report proof" — for the 4 changes without either, the `git log` IS the proof. The user must explicitly approve each reconciliation, but that is one decision per change, not per task.

### Recommended execution order

1. **Pre-flight**: Confirm user wants the batch cleanup. Surface the 3 warnings on `opencode-sdd-profiles` (R-PROF-31, ScreenProfileCreate, tasks.md checkboxes) and ask for a descope decision.
2. **Run `sdd-verify` for the 6 changes without verify-reports** (parallel sub-agents):
   - `agent-builder` — static evidence (v1.18.0)
   - `antigravity-2.0-compatibility` — static evidence (v1.30.7)
   - `declarative-picker-navigation` — static evidence (v1.41.0, tests pass)
   - `hermes-agent-support` — static evidence (PR #774, tests pass)
   - `update-experience` — static evidence (5 PRs, tests pass)
   - Plus a top-level verify covering all 6 if a single batch is preferred
3. **Run `sdd-archive` for all 9** (sequential or parallel — each is independent):
   - For 3 with verify-report: standard archive flow
   - For 1 with apply-progress (`fix-persona-artifact-language-contract`): archive with stale-checkbox reconciliation
   - For 5 with new verify-reports: standard archive flow
4. **Post-flight**: Confirm `openspec/changes/` is empty. The 8 prior archives are untouched.

### Delivery strategy

`single-pr-default` is the orchestrator's preflight choice. Recommended: keep it
— this is a cleanup, not new work, and the change folders move but no source
code changes.

## Risks

- **Stale checkbox blocking on 4 changes** (agent-builder, qwen-code-integration, hermes-agent-support, plus the 33 unchecked hermes tasks, plus the 29 unchecked agent-builder tasks). The sdd-archive skill blocks on unchecked tasks by default. Mitigation: each sdd-archive run must carry an explicit orchestrator-approved stale-checkbox reconciliation, with the new verify-report proving each unchecked task was actually done.
- **`opencode-sdd-profiles` has 3 non-critical warnings**: R-PROF-31 missing sync-time model warning log, ScreenProfileCreate missing-cache guard (task 6.2), and tasks.md checkboxes. None block usage. Mitigation: ask user to descope each warning explicitly in the archive report (the sdd-archive skill allows non-critical partial archive with explicit user intent).
- **Workspace contamination**: The working tree has uncommitted `openspec/config.yaml` and untracked `toyota-chile-etl-template/`. If the cleanup runs `go test` or modifies config, it may interact with these. Mitigation: stash or commit them before sdd-verify/sdd-archive work; or restrict action context to `openspec/` only.
- **Two related persona changes** (active `fix-persona-artifact-language-contract` vs. archived `2026-06-10-level-neutral-persona-parity`): they touch different scope (artifact language contract vs. level-neutral L1/L2 behavior parity) but both modify the same files. When archiving the active one, the spec delta should not duplicate what's already merged from the archived one. Mitigation: read the archived change's archive-report before merging deltas.
- **`declarative-picker-navigation` design.md is 17KB and tasks.md is 20KB** — these are large. Reading them via sdd-verify / sdd-archive will consume sub-agent context. Mitigation: budget accordingly; the verify can rely on `git log` proof rather than re-reading the design.
- **Token cost** for running 9 sdd-verify + 9 sdd-archive sub-agents in one session is significant. Mitigation: run in parallel where independent, and use static-evidence verify reports to keep them short.

## Ready for Proposal

**No** — this is an exploration, not a proposed change. The orchestrator should
present the user with the recommendation and ask one question:

> "I found 9 SDD changes in `openspec/changes/` that are all implemented on
> `main` but never archived. Want me to archive them all in a cleanup batch
> (Approach A), or pick a subset?"

After user confirmation, the next step is `sdd-verify` (parallel batch for the
6 without reports) followed by `sdd-archive` (sequential, 9 changes). There is
no `sdd-propose` step because this is cleanup of existing changes, not a new
change.

If the user prefers Approach C, the orchestrator can split the work into two
waves and run a partial batch first.
