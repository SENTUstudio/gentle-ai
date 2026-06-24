---
name: sdd-design
description: "Create the SDD technical design and architecture approach. Trigger: orchestrator launches design for a change."
disable-model-invocation: true
user-invocable: false
license: MIT
metadata:
  author: gentleman-programming
  version: "2.0"
  delegate_only: true
---

> **ORCHESTRATOR GATE**: If you loaded this skill via the `skill()` tool, you are
> the ORCHESTRATOR — STOP. Do NOT execute these instructions inline. Delegate to
> the dedicated `sdd-design` sub-agent using your platform's delegation primitive
> (e.g., `task(...)`, sub-agent invocation, etc.). This skill is for EXECUTORS
> only.

## Executor Override

If you ARE the `sdd-design` sub-agent (NOT the orchestrator), the gate above does NOT apply to you. Continue with the phase work below. Do NOT delegate. Do NOT call the Skill tool. You are the executor — execute.


## Language Domain Contract

Generated technical artifacts default to English. Do not inherit the user's conversational language or the active persona's regional voice for SDD artifacts unless the user explicitly requests that artifact language or the project convention requires it.

If Spanish technical artifacts are explicitly requested, use neutral/professional Spanish unless the user explicitly asks for a regional variant.

Public/contextual comments follow the target context language by default. Explicit user language or tone overrides win; Spanish comments default to neutral/professional Spanish unless the user or target context clearly calls for regional tone.

## Purpose

You are a sub-agent responsible for TECHNICAL DESIGN. You take the proposal and specs, then produce a `design.md` that captures HOW the change will be implemented — architecture decisions, data flow, file changes, and technical rationale.

## What You Receive

From the orchestrator:
- Change name
- Artifact store mode (`engram | openspec | hybrid | none`)

## Execution and Persistence Contract

> Follow **Section B** (retrieval) and **Section C** (persistence) from `skills/_shared/sdd-phase-common.md`.

- **engram**: Read `sdd/{change-name}/proposal` (required) and `sdd/{change-name}/spec` (optional — may not exist if running in parallel with sdd-spec). Save as `sdd/{change-name}/design`.
- **openspec**: Read and follow `skills/_shared/openspec-convention.md`.
- **hybrid**: Follow BOTH conventions — persist to Engram AND write `design.md` to filesystem. Retrieve dependencies from Engram (primary) with filesystem fallback.
- **none**: Return result only. Never create or modify project files.

## Data-Engineering Domain Branch

> **GATE**: Run `gentle-ai sdd-config --json` once at the start of this phase.
> If `domain == "data-engineering"`, the design document you write MUST follow
> the templates in this section (in addition to the standard design format
> below). If `domain` is absent or anything else, this section is a NO-OP —
> produce the standard design.md exactly as before. Never infer the domain.

When the gate is active, the design must internalize ETL-specific artifacts:

- A **DAG-of-transformations**: nodes = source tables + intermediate temp
  views + target table; edges = data dependencies. This is the single most
  load-bearing diagram for a Glue job — reviewers reason about correctness
  from it, not from prose.
- A **pattern-aware** design choice: the design declares which of the four
  patterns (`incremental`, `multi-step`, `legacy-wrangler`, `glue-studio`)
  detected by `internal/etl.DetectPattern` this change follows, and uses the
  matching scaffold. Mixed patterns MUST be called out explicitly.
- For MODIFICATIONS, an **insertion-point analysis**: WHERE the new logic
  inserts into the existing DAG, and the cascade impact (which downstream
  stages change shape, which remain byte-identical).

### DAG-of-transformations template (data-engineering)

Include this section under "Technical Approach" (replacing or augmenting the
generic data-flow ASCII diagram):

```markdown
### DAG (transformations)

Nodes:
- [S1] source db_dl_dev_stg_encuestas.encuestas_csi  (sidecar: glue-tables/...)
- [V1] tmp_survey_clean
- [V2] tmp_survey_scores
- [V3] tmp_rollup_month
- [T1] target db_dl_dev_ref.encuestas_csi_ref         (sidecar: glue-tables/...)

Edges (data dependencies):
- S1 -> V1
- S1 -> V2
- V1, V2 -> V3
- V3 -> T1

Idempotency window: watermark column `<col>`; re-running over the same window
produces zero new target rows.
```

If the change is a modification rather than a full rebuild, append a
"## Modification Insertion Point" section (template below) and mark the DAG
with `[+] new` / `[~] extended shape` / `[=] byte-identical regression baseline`
on the affected nodes.

### Pattern-aware design choice (data-engineering)

Determine the pattern by inspecting the existing Glue job source (the markers
`internal/etl.DetectPattern` consumes); declare it explicitly:

```markdown
### Pattern

- Classification: `incremental` (or `multi-step` | `legacy-wrangler` | `glue-studio`)
- Confidence: <from DetectPattern, e.g. 0.85>
- Override: none | user-confirmed <pattern> (overrides the heuristic)
- Chosen scaffold (Camino A): <watermark-test / multi-step-DAG-test / Athena-Spark parity / legacy wrangler smoke>
```

For mixed or ambiguous patterns (confidence < 0.85 OR two matched patterns),
the design MUST list the secondary indication and the override choice — never
silently pick one. See `internal/etl.DetectPattern`.

### Modification insertion-point template (data-engineering — MODIFIED only)

```markdown
### Insertion Point

- Insert AFTER node `<V?>` (the last stage whose output remains byte-identical)
- Insert BEFORE node `<V?>` (the first stage whose output shape changes)
- New nodes to add: `[<V?>] <name>` carrying the new logic, declared in the DAG above

### Cascade Impact

| Stage | Pre-change rows | Post-change rows | Reason |
|-------|------------------|------------------|--------|
| <upstream unchanged>   | R | R | byte-identical regression baseline |
| <new stage>             | — | R' | new column(s) / window |
| <downstream extended>   | R | R'' | depends on new stage |
| <target>                | R | R'' | final materialized impact |

### Regression Strategy

- Byte-identical stages: asserted by **Camino A** EXCEPT against a frozen
  pre-change dev run on the same source window.
- Extended stages & target: asserted by the new scenarios in the spec, and by
  **Camino B** dev-vs-prd EXCEPT (`internal/etl.BuildExceptSQL`).
- Sidecar drift detection: `internal/etl.ValidateSidecar` against
  `aws glue get-table` covers schema; row-level drift covered by EXCEPT.
```

### Scenario coverage for ETL

Design-level technical questions to answer (in addition to the standard
design.md sections) when `domain == data-engineering`:

- Is the watermark source reliable (state store vs derived)? What is the
  recovery procedure if the watermark is lost?
- Are the temp views re-usable across jobs (`HasTempViews`-style), or are
  they private to this job? (Affects isolation design.)
- What is the EXCEPT scope for Camino B parity (full table, partition slice,
  watermark window)?
- Which `aws_profiles` logical profile (`prd`/`dev`/`usuario`) is consumed at
  each pipeline stage? (Leaks here are catched by `internal/sddconfig.ScrubProfiles`.)

## What to Do

### Step 1: Load Skills
Follow **Section A** from `skills/_shared/sdd-phase-common.md`.

### Step 1a (data-engineering only): Resolve Domain

If you entered the **Data-Engineering Domain Branch** above (the gate passed),
resolve the domain ONCE here and cache it for the rest of the phase:

```bash
gentle-ai sdd-config --json
```

Parse `.domain`. If `data-engineering`, the DAG + pattern + insertion-point
templates above are MANDATORY. If absent, skip them and proceed with the
app-dev design flow below — the rest of this skill is identical to its prior
behavior.

### Step 2: Read the Codebase

Before designing, read the actual code that will be affected:
- Entry points and module structure
- Existing patterns and conventions
- Dependencies and interfaces
- Test infrastructure (if any)

### Step 3: Write design.md

**IF mode is `openspec` or `hybrid`:** Create the design document:

```
openspec/changes/{change-name}/
├── proposal.md
├── specs/
└── design.md              ← You create this
```

**IF mode is `engram` or `none`:** Do NOT create any `openspec/` directories or files. Compose the design content in memory — you will persist it in Step 4.

#### Design Document Format

```markdown
# Design: {Change Title}

## Technical Approach

{Concise description of the overall technical strategy.
How does this map to the proposal's approach? Reference specs.}

## Architecture Decisions

### Decision: {Decision Title}

**Choice**: {What we chose}
**Alternatives considered**: {What we rejected}
**Rationale**: {Why this choice over alternatives}

### Decision: {Decision Title}

**Choice**: {What we chose}
**Alternatives considered**: {What we rejected}
**Rationale**: {Why this choice over alternatives}

## Data Flow

{Describe how data moves through the system for this change.
Use ASCII diagrams when helpful.}

    Component A ──→ Component B ──→ Component C
         │                              │
         └──────── Store ───────────────┘

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `path/to/new-file.ext` | Create | {What this file does} |
| `path/to/existing.ext` | Modify | {What changes and why} |
| `path/to/old-file.ext` | Delete | {Why it's being removed} |

## Interfaces / Contracts

{Define any new interfaces, API contracts, type definitions, or data structures.
Use code blocks with the project's language.}

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | {What} | {How} |
| Integration | {What} | {How} |
| E2E | {What} | {How} |

## Migration / Rollout

{If this change requires data migration, feature flags, or phased rollout, describe the plan.
If not applicable, state "No migration required."}

## Open Questions

- [ ] {Any unresolved technical question}
- [ ] {Any decision that needs team input}
```

### Step 4: Persist Artifact

**This step is MANDATORY — do NOT skip it.**

Follow **Section C** from `skills/_shared/sdd-phase-common.md`.
- artifact: `design`
- topic_key: `sdd/{change-name}/design`
- type: `architecture`

### Step 5: Return Summary

Return to the orchestrator:

```markdown
## Design Created

**Change**: {change-name}
**Location**: `openspec/changes/{change-name}/design.md` (openspec/hybrid) | Engram `sdd/{change-name}/design` (engram) | inline (none)

### Summary
- **Approach**: {one-line technical approach}
- **Key Decisions**: {N decisions documented}
- **Files Affected**: {N new, M modified, K deleted}
- **Testing Strategy**: {unit/integration/e2e coverage planned}

### Open Questions
{List any unresolved questions, or "None"}

### Next Step
Ready for tasks (sdd-tasks).
```

## Rules

- ALWAYS read the actual codebase before designing — never guess
- Every decision MUST have a rationale (the "why")
- Include concrete file paths, not abstract descriptions
- Use the project's ACTUAL patterns and conventions, not generic best practices
- If you find the codebase uses a pattern different from what you'd recommend, note it but FOLLOW the existing pattern unless the change specifically addresses it
- Keep ASCII diagrams simple — clarity over beauty
- Apply any `rules.design` from `openspec/config.yaml`
- If you have open questions that BLOCK the design, say so clearly — don't guess
- **Size budget**: Design artifact MUST be under 800 words. Architecture decisions as tables (option | tradeoff | decision). Code snippets only for non-obvious patterns.
- Return envelope per **Section D** from `skills/_shared/sdd-phase-common.md`.
