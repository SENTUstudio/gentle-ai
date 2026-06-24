<!-- section:model-capable -->
---
name: sdd-verify
description: "Trigger: SDD verification phase, verify change. Execute tests and prove implementation matches specs, design, and tasks."
disable-model-invocation: true
user-invocable: false
license: MIT
metadata:
  author: gentleman-programming
  version: "3.0"
  delegate_only: true
---

> **ORCHESTRATOR GATE**: If you loaded this skill via the `skill()` tool, you are
> the ORCHESTRATOR — STOP. Do NOT execute these instructions inline. Delegate to
> the dedicated `sdd-verify` sub-agent using your platform's delegation primitive
> (e.g., `task(...)`, sub-agent invocation, etc.). This skill is for EXECUTORS
> only.

## Executor Override

If you ARE the `sdd-verify` sub-agent (NOT the orchestrator), the gate above does NOT apply to you. Continue with the phase work below. Do NOT delegate. Do NOT call the Skill tool. You are the executor — execute.

## Language Domain Contract

Generated technical artifacts default to English. Do not inherit the user's conversational language or the active persona's regional voice for SDD artifacts unless the user explicitly requests that artifact language or the project convention requires it.

If Spanish technical artifacts are explicitly requested, use neutral/professional Spanish unless the user explicitly asks for a regional variant.

Public/contextual comments follow the target context language by default. Explicit user language or tone overrides win; Spanish comments default to neutral/professional Spanish unless the user or target context clearly calls for regional tone.

## Activation Contract

Run when the orchestrator launches verification for an SDD change. You are the quality gate: prove completion with source inspection plus real execution evidence.

The orchestrator should provide structured status from `skills/_shared/sdd-status-contract.md`. Use its `schemaName`, `planningHome`, `changeRoot`, `artifactPaths`, `contextFiles`, task progress, dependency states, and `actionContext` before judging artifacts.

## Data-Engineering Domain Branch

> **GATE**: Run `gentle-ai sdd-config --json` once at the start of the phase.
> If `domain == "data-engineering"`, follow THIS section (Camino B) AND the
> rest of the skill as the outer container. If `domain` is absent or anything
> else, the instructions below are a NO-OP — execute the rest of the skill
> exactly as before (app-dev behavior is untouched: `go test ./...` +
> `go vet ./...`). Never infer the domain from the prompt.

When the gate is active, verification follows **Camino B** (deploy + parity)
instead of the app-dev Go test suite. The contract of a Glue job is the table
it materializes, so Camino B proves dev-vs-prd parity and sidecar structural
drift — not function return values. Every report, log, and PR description
produced here MUST run through `internal/sddconfig.ScrubProfiles` so no AWS
profile name or account ID leaks.

### Camino B — deploy + dev-vs-prd parity

1. **Resolve config once.**
   ```bash
   gentle-ai sdd-config --json
   ```
   Read `.domain` (MUST be `data-engineering`), `.repos.{infra,carga}`,
   `.aws_profiles.{prd,dev,usuario}`, and `.verify.skip_deploy`. Profile names
   are resolved by `internal/sddconfig.ResolveProfile` — never echo a raw
   profile name or account ID into any output.

2. **Deploy both repos (parallel), unless `verify.skip_deploy`.** Run
   `sam deploy` on the `infra` and `carga` repos in parallel. If
   `verify.skip_deploy` is `true`, short-circuit to sidecar + EXCEPT-only
   verification (local mode) and record the skip in the report.

3. **Run the Glue job.**
   ```bash
   aws glue start-job-run --job-name <job> --profile $(ResolveProfile dev)
   ```
   Wait for `JobRunState == SUCCEEDED`. Resolve the profile through
   `ResolveProfile`; scrub the command + its captured output before logging.

4. **Dev-vs-prd parity (Athena EXCEPT).** For each target table, build and run
   the parity query via Athena:
   ```sql
   <internal/etl.BuildExceptSQL(target, devDB, prdDB)>
   ```
   The query is `SELECT * FROM <devDB>.<target> EXCEPT SELECT * FROM
   <prdDB>.<target>`. A non-empty result is a parity CRITICAL finding (drift
   between dev and prd). Record the row count, never the profile name.

5. **Sidecar structural validation.** For each `glue-tables/{db}.{table}.yaml`
   sidecar in the change, fetch the live table and run
   `internal/etl.ValidateSidecar(sidecar, glueTable)`:
   ```bash
   aws glue get-table --database <db> --name <table> --profile $(ResolveProfile dev)
   ```
   Each returned `Mismatch` (`database_mismatch`, `table_mismatch`,
   `s3_location_mismatch`, `missing_column`, `type_mismatch`,
   `missing_partition`, `unexpected_partition`) is a CRITICAL finding.

6. **Scrub all output.** Run the final report, every captured command, and any
   PR description through `internal/sddconfig.ScrubProfiles` against the
   resolved config. Profile names and 12-digit account IDs MUST NOT appear in
   any persisted artifact.

7. **Verdict.** Camino B yields `PASS` only when every EXCEPT result is empty
   AND `ValidateSidecar` returns no mismatches AND the job run succeeded. Any
   parity drift or sidecar mismatch is `FAIL`.

## Hard Rules

- Read all available status `contextFiles` before judging implementation. Full spec-driven verification reads proposal, specs, design, and tasks; partial artifact sets degrade as described below.
- Execute relevant tests; static analysis alone is never verification.
- A spec scenario is compliant only when a covering test passed at runtime.
- Compare specs first, design second, task completion third.
- Do not fix issues; report them for the orchestrator/user.
- Persist `verify-report` according to mode: Engram, openspec file, hybrid both, or inline-only for `none`.
- If Strict TDD is active, load `strict-tdd-verify.md` from this skill directory; if inactive, never load it.
- Return the Section D envelope from `../_shared/sdd-phase-common.md`.

## Decision Gates

| Condition | Action |
|---|---|
| Orchestrator says `STRICT TDD MODE IS ACTIVE` | Treat as authoritative. |
| Cached/config `strict_tdd: true` and runner exists | Strict TDD verify; load module. |
| Strict TDD false or no runner | Standard verify; skip TDD checks. |
| `actionContext.mode: workspace-planning` | STOP; full workspace implementation verification is not supported in this slice. |
| Only tasks artifact exists | Verify task completion only; skip spec/design correctness and record skipped checks. |
| Tasks + specs exist | Verify completeness and correctness; skip design coherence and record skipped checks. |
| Proposal/specs/design/tasks exist | Verify all dimensions. |
| Task incomplete | CRITICAL for core task, WARNING for cleanup task. |
| Test command exits non-zero | CRITICAL. |
| Spec scenario has no passing covering test | CRITICAL `UNTESTED` or `FAILING`. |
| Design deviation exists | WARNING unless it breaks a spec. |

## Execution Steps

1. Load relevant skills via shared SDD Section A.
2. Retrieve artifacts via shared Section B for the active persistence mode, or read the concrete `contextFiles` from structured status.
3. Resolve testing/TDD mode from cached capabilities, config, or project files.
   - **(data-engineering only)** Resolve domain with `gentle-ai sdd-config --json`. If `domain == "data-engineering"`, switch to Camino B (deploy + dev-vs-prd parity, see the Data-Engineering Domain Branch above) and skip the app-dev Go-suite steps below. Honor `verify.skip_deploy` for local-only mode.
4. Count completed and incomplete tasks. Any unchecked implementation task is CRITICAL and blocks archive readiness.
5. If specs exist, map each spec requirement/scenario to implementation evidence and tests.
6. If design exists, check design decisions against changed code. If design is missing, skip design coherence and record why.
7. Run test, build/type-check, and coverage commands when available. For full spec verification, preserve gentle-ai's stricter runtime evidence: source inspection alone does not prove spec scenario compliance.
8. Build the behavioral compliance matrix from actual test results when specs/scenarios exist.
9. Persist and return the verification report, including skipped dimensions for missing artifacts.

## Output Contract

Return `## Verification Report` with change, mode, completeness table, build/tests/coverage evidence, spec compliance matrix, correctness table, design coherence table, issues grouped as CRITICAL/WARNING/SUGGESTION, and final verdict `PASS`, `PASS WITH WARNINGS`, or `FAIL`.

## Graceful Artifact Handling

- **Tasks only**: verify objective task completion only. Do not claim spec correctness or design coherence. If all tasks are checked and no runtime evidence is available, verdict may be `PASS WITH WARNINGS` for task completion only.
- **Tasks + specs**: verify task completeness and requirement/scenario correctness. Runtime test evidence is still required for full spec scenario compliance; missing covering tests are CRITICAL for required scenarios unless project config explicitly allows manual verification.
- **Full artifacts**: verify completeness, correctness, and coherence.
- **Unchecked tasks**: always remain CRITICAL, even when other artifacts are missing or warnings-only.

## References

- [references/report-format.md](references/report-format.md) — full report template, compliance statuses, and command evidence fields.
- [strict-tdd-verify.md](strict-tdd-verify.md) — load only when Strict TDD is active.
- `../_shared/sdd-phase-common.md` — skill loading, retrieval, persistence, and return envelope.
<!-- /section:model-capable -->

<!-- section:model-small -->
---
name: sdd-verify
description: "Trigger: SDD verification phase, verify change. Execute tests and prove implementation matches specs, design, and tasks."
disable-model-invocation: true
user-invocable: false
license: MIT
metadata:
  author: gentleman-programming
  version: "3.0"
  delegate_only: true
---

> **ORCHESTRATOR GATE**: If you loaded this skill via the `skill()` tool, you are the ORCHESTRATOR — STOP. Do NOT execute these instructions inline. Do NOT delegate, do NOT call task/delegate, do NOT launch sub-agents. Read this SKILL.md and follow it exactly.


## Language Domain Contract

Generated technical artifacts default to English. Do not inherit the user's conversational language or the active persona's regional voice for SDD artifacts unless the user explicitly requests that artifact language or the project convention requires it.

If Spanish technical artifacts are explicitly requested, use neutral/professional Spanish unless the user explicitly asks for a regional variant.

Public/contextual comments follow the target context language by default. Explicit user language or tone overrides win; Spanish comments default to neutral/professional Spanish unless the user or target context clearly calls for regional tone.

## Purpose

You are a VERIFY sub-agent. Your job: check implemented changes match spec acceptance criteria. Do NOT delegate.

## Hard Rules

- Read spec acceptance criteria only
- Inspect changed files listed in apply-progress (or tasks) — limit to those files
- Use structured status when provided; stop on workspace-planning action context
- Do NOT run tests unless `strict_tdd` is active and test runner is explicitly provided
- Do not fix issues; report them for the orchestrator/user
- Return minimal report

## Data-Engineering Domain Gate (small-model)

Run `gentle-ai sdd-config --json` and parse `.domain`. If `domain ==
"data-engineering"`, Camino B is required (SAM deploy both repos + Athena
dev-vs-prd EXCEPT via `internal/etl.BuildExceptSQL` + sidecar
`ValidateSidecar` against `aws glue get-table` + `ScrubProfiles` on output,
honoring `verify.skip_deploy`). That loop exceeds this small-model read budget
— STOP and return `needs-explore` so a capable model runs Camino B. App-dev
verify (domain absent) continues with the minimal report below.

## Return Minimal Report

```json
{
  "status": "pass|fail|warning",
  "checks": [{"criterion": "text", "result": "pass|fail", "evidence": "one-line"}],
  "next": "ready-for-archive|fixes-required"
}
```
<!-- /section:model-small -->
