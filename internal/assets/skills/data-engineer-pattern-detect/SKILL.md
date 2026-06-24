---
name: data-engineer-pattern-detect
description: "Heuristic scan for ETL patterns (incremental, multi-step, legacy-wrangler, glue-studio). Trigger: ETL pattern detection before scaffold selection."
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## Reference: _shared/data-engineer-protocol.md

This skill follows the shared data-engineer protocol. Read it before
proceeding: header authorship rule (§2), the 4-pattern taxonomy (§3), and the
master-project awareness rule (§4).

## When to Use

- When `sdd-apply` Camino A needs to pick a TDD scaffold before writing the
  first ETL test
- When the user asks "what kind of Glue job is this?" before refactoring
- When `sdd-propose` needs a confidence-backed pattern override for the record
- When onboarding a legacy ETL repo and the first task is "classify every job"

Do NOT use this skill for app-dev work. It is gated on
`gentle-ai sdd-config --json` → `.domain == "data-engineering"`; absent the
domain flag, the detection contract does not apply and this skill is inert.

## Critical Patterns

### Detection is Never Silent

Every detection result carries three components:

1. **Pattern** — one of `incremental`, `multi-step`, `legacy-wrangler`,
   `glue-studio`, or `unknown`
2. **Confidence** — `0.0` to `1.0` (0.8 = both markers fire, 0.5 = one marker
   fires as a hint, 0.0 = no markers)
3. **Rationale** — human-readable explanation naming which markers fired,
   which did not, and why a partial signal did not promote to a full verdict

A bare `(pattern, confidence)` pair without a rationale is a contract
violation. The rationale is what lets the user trust OR override the verdict.

### How to Run the Heuristic Scan

Two equivalent paths — pick whichever the caller already supports:

**Path A — gentle-ai CLI (preferred when available):**

```bash
# Returns JSON: { "pattern": "...", "confidence": 0.8, "rationale": "..." }
gentle-ai sdd-config --detect-pattern --path glue-jobs/<job>.py
```

(If this subcommand is not yet wired, fall back to Path B.)

**Path B — inline marker scan:**

Read the target `.py` file and populate the `Markers` struct consumed by
`internal/etl.DetectPattern`. The struct fields (see `internal/etl/pattern.go`):

| Field                      | Fires when the source contains                                |
|----------------------------|---------------------------------------------------------------|
| `HasFromISO`               | `from_catalog(` or `from_options(` with a dynamic frame       |
| `HasWatermarkColumn`       | A column compared against a partition/snapshot date           |
| `HasGlueContext`           | `GlueContext()` instantiated                                  |
| `HasTempViews`             | `createOrReplaceTempView`                                     |
| `HasWranglerAthena`        | `awswrangler.athena.read_sql` (or `.unload`)                  |
| `HasApplyMapping`          | `ApplyMapping.apply(...)` (Glue Studio transform)             |
| `HasSparkSqlQueryHelper`   | `SparkSqlQueryHelper` (Glue Studio SQL helper)                |

Map fields to patterns per §3 of the shared protocol. When the AND of two
markers defines a pattern (e.g. glue-studio requires BOTH `ApplyMapping` AND
`SparkSqlQueryHelper`), a single-marker match is a partial signal → contributes
to the `unknown` rationale, never silently promotes.

### Confidence Thresholds

| Confidence | Meaning                                                              | Action                          |
|------------|----------------------------------------------------------------------|---------------------------------|
| 0.8        | Both markers of a pattern fired                                      | Apply the matching scaffold     |
| 0.5        | One marker fired — hint, not verdict                                 | Surface hint + ask user         |
| 0.0        | No markers fired                                                     | Surface "unknown" + ask user    |

The thresholds (0.8 / 0.5) are fixed by `sdd-tasks` Open Questions Resolved
#3. Tuning them requires a spec change, not a skill edit.

### User Override

When the user overrides the detected pattern during `sdd-propose`:

1. Record the override in the proposal as `pattern_override: <pattern>` with
   the human rationale.
2. Carry the override forward to `sdd-apply` Camino A.
3. `sdd-apply` uses the override INSTEAD of re-running detection. Detection is
   only run once per change; the override is authoritative from proposal onward.

Never silently drop an override. If `sdd-apply` re-runs detection and gets a
different answer, it surfaces BOTH the override and the fresh detection, then
defers to the override.

## Output

The skill produces a **Pattern Detection Report**:

```yaml
job: <job-name>
file: glue-jobs/<job>.py
pattern: incremental          # incremental | multi-step | legacy-wrangler | glue-studio | unknown
confidence: 0.8
markers_fired:
  - HasFromISO
  - HasWatermarkColumn
markers_missed: []
rationale: >
  Both incremental markers present (from_catalog + watermark column).
  Scaffold: watermark test.
override_allowed: true
```

The report is consumed by `sdd-apply` Camino A to select the TDD scaffold and
by `sdd-propose` to record the verdict (or override) in the proposal.

## Resources

- Shared protocol: `_shared/data-engineer-protocol.md`
- Go detector: `internal/etl.DetectPattern` (`internal/etl/pattern.go`)
- Confidence thresholds: `openspec/changes/data-engineering-domain-profile/tasks.md` Open Questions Resolved #3
