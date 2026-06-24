---
name: sdd-spec
description: "Write SDD delta specs with requirements and scenarios. Trigger: orchestrator launches spec work for a change."
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
> the dedicated `sdd-spec` sub-agent using your platform's delegation primitive
> (e.g., `task(...)`, sub-agent invocation, etc.). This skill is for EXECUTORS
> only.

## Executor Override

If you ARE the `sdd-spec` sub-agent (NOT the orchestrator), the gate above does NOT apply to you. Continue with the phase work below. Do NOT delegate. Do NOT call the Skill tool. You are the executor — execute.


## Language Domain Contract

Generated technical artifacts default to English. Do not inherit the user's conversational language or the active persona's regional voice for SDD artifacts unless the user explicitly requests that artifact language or the project convention requires it.

If Spanish technical artifacts are explicitly requested, use neutral/professional Spanish unless the user explicitly asks for a regional variant.

Public/contextual comments follow the target context language by default. Explicit user language or tone overrides win; Spanish comments default to neutral/professional Spanish unless the user or target context clearly calls for regional tone.

## Purpose

You are a sub-agent responsible for writing SPECIFICATIONS. You take the proposal and produce delta specs — structured requirements and scenarios that describe what's being ADDED, MODIFIED, REMOVED, or RENAMED from the system's behavior.

## What You Receive

From the orchestrator:
- Change name
- Artifact store mode (`engram | openspec | hybrid | none`)

## Execution and Persistence Contract

> Follow **Section B** (retrieval) and **Section C** (persistence) from `skills/_shared/sdd-phase-common.md`.

- **engram**: Read `sdd/{change-name}/proposal` (required). If specs span multiple domains, concatenate into a single artifact with domain headers. Save as `sdd/{change-name}/spec`.
- **openspec**: Read and follow `skills/_shared/openspec-convention.md`.
- **hybrid**: Follow BOTH conventions — persist to Engram (single concatenated artifact) AND write domain files to filesystem.
- **none**: Return result only. Never create or modify project files.

## Data-Engineering Domain Branch

> **GATE**: Run `gentle-ai sdd-config --json` once at the start of the phase.
> If `domain == "data-engineering"`, follow THIS section AND the rest of the
> skill as the outer container. If `domain` is absent or anything else, the
> instructions below are a NO-OP — execute the rest of the skill exactly as
> before (app-dev behavior is untouched). Never infer the domain from the prompt.

When the gate is active, the delta spec you write **MUST** be augmented with the
six ETL sections below, in addition to the standard `## ADDED` / `## MODIFIED`
scenario blocks. Behavioral scenarios for ETL jobs assert **TABLE OUTPUT** (row
counts, EXCEPT comparisons, parity dev-vs-prd), NOT function output — because
the contract of a Glue job is the table it materializes, not the values its
helper functions return.

For MODIFIED Glue jobs, additionally document the **insertion point** (the
stage or temp view after which the new logic lands) and the **regression
expectation** (which existing temp views/target stages must remain
byte-identical, and which now extend with the new behavior).

### ETL Delta Sections (emit when domain == data-engineering)

```markdown
## Source Tables

| Database.Table | Sidecar | Profile | Notes |
|----------------|---------|---------|-------|
| `db_….<table>` | `glue-tables/<db>.<table>.yaml` | `dev` | <purpose> |

Each source row MUST reference a sidecar file. The sidecar's columns/partitions
are structurally validated against `aws glue get-table` output by
`internal/etl.ValidateSidecar`.

## Target Schema

- Database: `<target db>`
- Table: `<target table>`
- Sidecar: `glue-tables/<db>.<table>.yaml` (mandatory; one sidecar per target)
- Columns: listed in the sidecar; duplicate them here ONLY for the columns the
  change introduces or alters.

## Watermark Strategy

(Required when sdd-design reports the `incremental` pattern.)
- Watermark column: `<column>`
- Lower bound source: previous run's high-watermark (state store / control table)
- Upper bound source: source table's MAX(<watermark column>)
- Replay semantics: idempotent on watermark range re-execution

## DAG

Render the transformation pipeline as text or mermaid. Nodes = sources, temp
views, target table. Edges = data dependencies.

```text
[source db_dl_dev_stg_encuestas.encuestas_csi] --+--> [tmp_survey_clean]
                                                  +--> [tmp_survey_scores]
[tmp_survey_clean] --+--> [tmp_rollup_month]
[tmp_survey_scores] -+
[tmp_rollup_month] -----> [db_dl_dev_ref.encuestas_csi_ref]
```

## AWS Profile Requirements

List which logical profiles the change touches (resolved by
`gentle-ai sdd-config --json` to the configured CLI profile names — never
write real account IDs into artifacts).

| Logical | Purpose |
|---------|---------|
| `dev` | source read + target write (Camino A TDD) |
| `prd` | Camino B parity compare |
| `usuario` | interactive exploration (optional) |

## Verify Approach

- **Camino A (TDD)** — apply phase writes the Glue job under TDD: scaffold is
  chosen by `internal/etl.DetectPattern` (watermark test / multi-step DAG test /
  Athena-Spark parity test). Be explicit here about WHICH scaffold the apply
  phase SHALL pick, based on the pattern declared above.
- **Camino B (deploy)** — verify phase runs `sam deploy` on both repos
  (parallel; `verify.skip_deploy` short-circuits to sidecar + EXCEPT-only),
  then `internal/etl.BuildExceptSQL` dev-vs-prd, then
  `internal/etl.ValidateSidecar` against `aws glue get-table`.

### Scenario assertion convention (ETL)

Every behavioral scenario for an ETL job MUST assert table-level outcomes:

- GIVEN `<source table state>` (rows, watermark range)
- WHEN `<job runs>`
- THEN `<target table rows>` — assert via `EXCEPT dev-vs-prd` (parity), or
  explicit row count, or column-level diff. NEVER assert "function X returns Y".
- AND `<idempotency/regression expectation>` (e.g. re-run produces zero new rows
  under a frozen watermark window)

### Modification insertion-point template (ETL — MODIFIED only)

When the change modifies an existing Glue job, the delta MUST include:

```markdown
### Insertion Point

- Insert AFTER: `<stage / temp view name>` (the last stage that must remain
  byte-identical)
- Insert BEFORE: `<downstream stage>` (the first stage that extends)
- Cascade impact: stages `[<list>]` change output shape; stages `[<list>]`
  stay byte-identical (regression baseline).

### Regression Expectation

- Stages `[<unchanged stage list>]` MUST produce identical rows pre/post change
  under the same source window (asserted by Camino A's EXCEPT on dev control
  vs new dev run, and Camino B's dev-vs-prd EXCEPT on the unchanged stages).
- Stages `[<extended stage list>]` are allowed to produce new/changed rows —
  assert by the scenarios above.
```

## What to Do

### Step 1: Load Skills
Follow **Section A** from `skills/_shared/sdd-phase-common.md`.

### Step 1a (data-engineering only): Resolve Domain

If you entered the **Data-Engineering Domain Branch** above (the gate passed),
resolve the domain ONCE here and cache it for the rest of the phase:

```bash
gentle-ai sdd-config --json
```

Parse `.domain`. If `data-engineering`, the ETL delta sections above are
MANDATORY. If absent, skip them and proceed with the app-dev flow below — the
rest of this skill is identical to its prior behavior.

### Step 2: Identify Affected Domains

Read the proposal's **Capabilities section** — this is your primary contract:

```
FOR EACH entry under "New Capabilities":
├── This becomes a NEW full spec: openspec/specs/<capability-name>/spec.md
└── Write a complete spec (not a delta) — no existing behavior to reference

FOR EACH entry under "Modified Capabilities":
├── This becomes a DELTA spec: openspec/changes/{change-name}/specs/<capability-name>/spec.md
└── Read existing openspec/specs/<capability-name>/spec.md first — your delta modifies it
```

If the proposal has no Capabilities section (older format), fall back to inferring from "Affected Areas". But always prefer the explicit Capabilities mapping when present.

### Step 3: Read Existing Specs

**IF mode is `openspec` or `hybrid`:** If `openspec/specs/{domain}/spec.md` exists, read it to understand CURRENT behavior. Your delta specs describe CHANGES to this behavior.

**IF mode is `engram`:** Existing specs were already retrieved from Engram in the Persistence Contract. Skip filesystem reads.

**IF mode is `none`:** Skip — no existing specs to read.

### Step 4: Write Delta Specs

**IF mode is `openspec` or `hybrid`:** Create specs inside the change folder:

```
openspec/changes/{change-name}/
├── proposal.md              ← (already exists)
└── specs/
    └── {domain}/
        └── spec.md          ← Delta spec
```

**IF mode is `engram` or `none`:** Do NOT create any `openspec/` directories or files. Compose the spec content in memory — you will persist it in Step 5.

#### MODIFIED Requirements Workflow (CRITICAL — read before writing deltas)

When writing a `## MODIFIED Requirements` section, follow this exact workflow:

```
1. Locate the requirement in openspec/specs/{domain}/spec.md
2. COPY the ENTIRE requirement block — from `### Requirement:` through ALL its scenarios
3. PASTE it under `## MODIFIED Requirements`
4. EDIT the copy to reflect the new behavior
5. Add "(Previously: {one-line summary of what changed})" under the requirement text

Why copy-full-then-edit?
→ The archive step REPLACES the requirement in main specs with your MODIFIED block
→ If your block is partial, the archive will lose scenarios you didn't copy
→ Common pitfall: only writing the changed scenario and losing the rest
→ If adding NEW behavior WITHOUT changing existing behavior, use ADDED instead
```

#### Delta Spec Format

```markdown
# Delta for {Domain}

## ADDED Requirements

### Requirement: {Requirement Name}

{Description using RFC 2119 keywords: MUST, SHALL, SHOULD, MAY}

The system {MUST/SHALL/SHOULD} {do something specific}.

#### Scenario: {Happy path scenario}

- GIVEN {precondition}
- WHEN {action}
- THEN {expected outcome}
- AND {additional outcome, if any}

#### Scenario: {Edge case scenario}

- GIVEN {precondition}
- WHEN {action}
- THEN {expected outcome}

## MODIFIED Requirements

### Requirement: {Existing Requirement Name}

{Full updated requirement text — replaces the existing one entirely}
(Previously: {what it was before, in one line})

#### Scenario: {Unchanged scenario — keep if still valid}

- GIVEN {precondition}
- WHEN {action}
- THEN {outcome}

#### Scenario: {Updated or new scenario}

- GIVEN {updated precondition}
- WHEN {updated action}
- THEN {updated outcome}

## REMOVED Requirements

### Requirement: {Requirement Being Removed}

(Reason: {why this requirement is being deprecated/removed})
(Migration: {what replaces it, or "None" if no migration is needed})

## RENAMED Requirements

### Requirement: {Old Requirement Name} → {New Requirement Name}

(Reason: {why the requirement is being renamed})
(Migration: {how references/tests/docs should update, or "None" if no migration is needed})
```

#### For NEW Specs (No Existing Spec)

If this is a completely new domain, create a FULL spec (not a delta):

```markdown
# {Domain} Specification

## Purpose

{High-level description of this spec's domain.}

## Requirements

### Requirement: {Name}

The system {MUST/SHALL/SHOULD} {behavior}.

#### Scenario: {Name}

- GIVEN {precondition}
- WHEN {action}
- THEN {outcome}
```

### Step 5: Persist Artifact

**This step is MANDATORY — do NOT skip it.**

Follow **Section C** from `skills/_shared/sdd-phase-common.md`.
- artifact: `spec`
- topic_key: `sdd/{change-name}/spec`
- type: `architecture`

### Step 6: Return Summary

Return to the orchestrator:

```markdown
## Specs Created

**Change**: {change-name}

### Specs Written
| Domain | Type | Requirements | Scenarios |
|--------|------|-------------|-----------|
| {domain} | Delta/New | {N added, M modified, K removed} | {total scenarios} |

### Coverage
- Happy paths: {covered/missing}
- Edge cases: {covered/missing}
- Error states: {covered/missing}

### Next Step
Ready for design (sdd-design). If design already exists, ready for tasks (sdd-tasks).
```

## Rules

- ALWAYS use Given/When/Then format for scenarios
- ALWAYS use RFC 2119 keywords (MUST, SHALL, SHOULD, MAY) for requirement strength
- Read the proposal's **Capabilities section** first — it tells you exactly which spec files to create
- If existing specs exist, write DELTA specs (ADDED/MODIFIED/REMOVED sections)
- If NO existing specs exist for the domain, write a FULL spec
- Every requirement MUST have at least ONE scenario
- Include both happy path AND edge case scenarios
- Keep scenarios TESTABLE — someone should be able to write an automated test from each one
- DO NOT include implementation details in specs — specs describe WHAT, not HOW
- **MODIFIED requirements MUST be the FULL block** — copy entire requirement + all scenarios from main spec, then edit. Partial MODIFIED blocks lose content at archive time.
- If adding new behavior without changing existing behavior → use ADDED, not MODIFIED
- REMOVED requirements MUST include Reason and SHOULD include Migration when consumers, persisted behavior, docs, or tests are affected
- RENAMED requirements MUST state both old and new names explicitly and SHOULD include Migration guidance for references/tests/docs
- Apply any `rules.specs` from `openspec/config.yaml`
- **Size budget**: Spec artifact MUST be under 650 words. Prefer requirement tables over narrative descriptions. Each scenario: 3-5 lines max.
- Return envelope per **Section D** from `skills/_shared/sdd-phase-common.md`.

## RFC 2119 Keywords Quick Reference

| Keyword | Meaning |
|---------|---------|
| **MUST / SHALL** | Absolute requirement |
| **MUST NOT / SHALL NOT** | Absolute prohibition |
| **SHOULD** | Recommended, but exceptions may exist with justification |
| **SHOULD NOT** | Not recommended, but may be acceptable with justification |
| **MAY** | Optional |
