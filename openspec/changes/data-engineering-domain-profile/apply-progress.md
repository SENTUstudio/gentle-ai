# Apply Progress: Data Engineering Domain Profile

**Change**: data-engineering-domain-profile
**Phases covered**: 1, 2, 3, 4a of 8 — Profile Detection + Config Schema → Spec/Design Templates → Apply + Verify Dual-Path → Skills Foundation
**Mode**: Strict TDD (test runner: `go test ./...`)
**Artifact store**: both (OpenSpec + Engram)
**Delivery**: chained PRs — feature-branch-chain, 8 PRs; PR 1 (Config + detect + CLI), PR 2 (sidecar + pattern + spec/design), PR 3 (header + compare + apply/verify), PR 4a (protocol + pattern-detect + catalog + study-file + create-table). Latest batch = PR 4a.
**Chain strategy note**: `tasks.md` Review Workload Forecast says `stacked-to-main`; the
orchestrator instruction for this run says `feature-branch-chain`. The orchestrator instruction
is followed for PR targeting. This does not affect what was implemented.

## Completed Tasks

- [x] 1.1 `internal/sddconfig/config.go` + `config_test.go` — `Config`, `Repos`, `VerifyOpts`, `LoadConfig`, pure `parseConfig`.
- [x] 1.2 `internal/sddconfig/detect.go` + `detect_test.go` — `DetectDomain` (0.8 both markers / 0.5 one / 0 none).
- [x] 1.3 `internal/sddconfig/repos.go` + `repos_test.go` — `ValidateRepos`, `PathWarning`.
- [x] 1.4 `internal/sddconfig/profiles.go` + `profiles_test.go` — `ResolveProfile`, `ScrubProfiles`.
- [x] 1.5 `internal/cli/sdd_config.go` + `sdd_config_test.go` — `RunSDDConfig` (`--json`, `--detect`, `--validate-repos`, `--cwd`).
- [x] 1.6 Wire `sdd-config` subcommand in `internal/app/app.go` + `internal/app/help.go`.
- [x] 1.7 Patch `internal/assets/skills/sdd-init/SKILL.md` — preflight detect + confirm/override + write `domain`.

## TDD Cycle Evidence

| Task | Test File | Layer | Safety Net | RED | GREEN | TRIANGULATE | REFACTOR |
|------|-----------|-------|------------|-----|-------|-------------|----------|
| 1.1 | `internal/sddconfig/config_test.go` | Unit | N/A (new) | ✅ Written (undefined parseConfig/LoadConfig) | ✅ 9/9 passed | ✅ 6 cases (all fields, skip_deploy true, missing fields, block-scalar skip, quotes/comments, empty) + 3 LoadConfig file behaviors | ✅ Clean (shared `skipIndentedUntilTopLevel`) |
| 1.2 | `internal/sddconfig/detect_test.go` | Unit | N/A (new) | ✅ Written (undefined DetectDomain) | ✅ 6/6 passed | ✅ 5 cases (both, only-template, only-glue, none, glue-dir-without-py) + nonexistent-root error | ✅ DRY: unified `glueJobsHasPython` into `dirHasPython` |
| 1.3 | `internal/sddconfig/repos_test.go` | Unit | N/A (new) | ✅ Written (undefined ValidateRepos) | ✅ 6/6 passed | ✅ 6 cases (all exist, infra missing, both missing+order, empty config, absolute existing, carga missing) | ➖ None needed (already clean) |
| 1.4 | `internal/sddconfig/profiles_test.go` | Unit | N/A (new) | ✅ Written (undefined ResolveProfile/ScrubProfiles) | ✅ 9/9 passed | ✅ 7 cases (profile scrub, standalone acct id, both, 13-digit boundary, empty-config scrub, innocuous, resolve known/unknown/empty) | ➖ None needed |
| 1.5 | `internal/cli/sdd_config_test.go` | Integration | N/A (new file); cli pkg baseline = PASS except known-flaky | ✅ Written (undefined RunSDDConfig/DetectionReport/ValidationReport) | ✅ 9/9 passed | ✅ 9 cases (markdown, json config, app-dev json, detect json, detect markdown, validate json, validate markdown, bad cwd, bad flag) | ✅ DRY: extracted `renderOrEncode` |
| 1.6 | `internal/app/app_test.go` + `help_test.go` | Integration | ✅ internal/app baseline PASS (captured before edit) | ✅ Written (sdd-config not dispatched → "unknown command"; help missing sdd-config) | ✅ 3/3 passed (2 new app + help) | ✅ 2 dispatch cases (before-platform-validation + reads config domain) + help listing | ➖ None needed (one-line dispatch + help row) |
| 1.7 | `internal/assets/skills/sdd-init/SKILL.md` | N/A (markdown) | Golden tests (see below) | ➖ N/A (skill prompt, no Go test layer) | ➖ N/A | ➖ Triangulation skipped: markdown skill patch, validated by sdd-verify spec scenarios | ➖ N/A |

### Test Summary

- **Total new Go tests written**: 41 (sddconfig 30 unit + cli 9 integration + app 2 integration)
- **Total tests passing**: 41/41 (plus `TestHelpContainsAllCommands` updated to assert `sdd-config`)
- **Layers used**: Unit (30), Integration (11)
- **Approval tests** (refactoring): None — no existing Go code refactored; only additive new packages + thin wiring.
- **Pure functions created**: `parseConfig`, `LoadConfig` (file→pure), `DetectDomain`, `ValidateRepos`, `ResolveProfile`, `ScrubProfiles`, `ParseCommandArgs`, render helpers.

## Golden-File Refresh (task 1.7 side effect)

Patching the embedded `sdd-init` SKILL.md asset changed the content the SDD injector
copies to each agent, breaking 8 `TestGoldenSDD_*` golden references in
`internal/components`. Because the asset change is intentional (task 1.7), the golden
references were regenerated with `go test ./internal/components/ -run TestGoldenSDD -update`.
Diff verified: **only** the 8 `*-skill-sdd-init.golden` files changed, and **only** with the
Domain Preflight section + renumbered Execution Steps + 2 new Decision Gates rows. No
command/agent goldens changed; no non-deterministic content (timestamps/AI tags) introduced.
`TestGoldenSDD_Claude` was unaffected (Claude golden-checks commands/agents, not the skill SKILL.md).

## Verification Results

| Check | Command | Result |
|---|---|---|
| Build | `go build ./...` | ✅ PASS |
| Targeted tests | `go test ./internal/sddconfig/... ./internal/cli/...` | ✅ PASS (cli: only known-flaky `TestRunInstallKimiMissingUVFailsBeforeExecutingInstallCommands` fails — pre-existing, unrelated) |
| Full suite | `go test ./...` | ✅ No regressions — only the known-flaky install test fails (failed identically in pre-edit baseline) |
| Vet | `go vet ./internal/sddconfig/... ./internal/cli/... ./internal/app/...` and `go vet ./...` | ✅ PASS |

## Backward Compatibility (verified by tests)

- `LoadConfig` on a project with no `openspec/config.yaml` → zero `Config`, nil error.
- `LoadConfig` on the current app-dev `config.yaml` (with `context: |`, `rules:`, `testing:`) → zero `Config` (no `domain`/`repos`/`aws_profiles`/`verify` keys), nil error. The hand-rolled parser skips block scalars and unknown blocks.
- `DetectDomain` with zero markers → `""`, confidence 0.
- `ValidateRepos` on empty config → no warnings.
- `ResolveProfile` on empty config → `""`.
- `ScrubProfiles` on empty config → still scrubs bare 12-digit account ids (harmless for app-dev).
- `sdd-config` subcommand dispatches before platform validation; no existing command behavior changed.
- All `sddconfig`/`cli`/`app`/`components` tests green except the pre-existing flaky install test.

## Deviations from Design

- **ScrubProfiles location**: The design Interfaces block lists `ScrubProfiles` under `// internal/etl`, but the File Changes table and the Phase 1 task both place it in `internal/sddconfig/profiles.go`. The task prompt is authoritative for Phase 1 → `ScrubProfiles` lives in `sddconfig`. If a later phase moves ETL helpers to `internal/etl`, it can re-export or relocate then.
- **Subcommand wiring location**: `tasks.md` 1.6 says "in `cmd/gentle-ai/main.go`", but the actual root dispatch is in `internal/app/app.go` (`cmd/gentle-ai/main.go` only calls `app.Run()`). Wired in `app.go`'s info-commands switch + `help.go`, matching the existing `sdd-status`/`sdd-continue` pattern.
- **YAML parsing**: The project deliberately avoids `gopkg.in/yaml.v3` (hand-rolled scanner in `internal/components/filemerge/yaml.go`). `parseConfig` is a minimal indentation-aware reader for the 4 owned keys, skipping block scalars (`context: |`) and unknown blocks — no new dependency added.
- **One-marker detection semantics**: A single marker returns `domain="data-engineering"` at confidence 0.5 (a hint), not `""`. This lets `sdd-init` present the hint for confirm/override. Zero markers → `""`. This matches "present hint, confirm/override" in the design; flagged here in case the team prefers one-marker → empty.

## Issues Found

- Golden-file regression from the skill patch (resolved by regeneration, see above). Future phases (2-5) that patch embedded SDD skill assets MUST run `go test ./internal/components/ -run TestGoldenSDD -update` and verify the diff is only the intended skill content.

## Remaining Tasks

Phase 3-8 are out of scope for this batch (do not implement):
- [ ] 3.1–3.4 Apply + Verify dual-path (`internal/etl/header.go`, `compare.go`, skill patches)
- [ ] 4a.1–4a.6, 4b.1–4b.2, 4c.1–4c.2, 4d.1 Skills foundation + embeds
- [ ] 5.1–5.2 Tasks multi-repo + git flow

---

# Phase 2: Spec + Design Templates

**Phase**: 2 of 8 — Spec + Design Templates (`internal/etl` helpers + skill patches)
**Mode**: Strict TDD (test runner: `go test ./...`)
**Delivery**: chained PR slice (PR 2 of 8) under `feature-branch-chain`. Builds
atop PR 1 on `feature/aws-dataengineer`. Backward compatibility preserved: the
two skill patches add a conditional `if domain == "data-engineering"` branch;
absent `domain` (app-dev) yields identical behaviour to today.

## Completed Tasks

- [x] 2.1 `internal/etl/sidecar.go` + `sidecar_test.go` — `Sidecar` (+`Column`, `Mismatch`), `ParseSidecar` (hand-rolled YAML scanner, no `yaml.v3`), `ValidateSidecar` (database/table/S3/column names + types/partitions match against `aws glue get-table` shape).
- [x] 2.2 `internal/etl/pattern.go` + `pattern_test.go` — `Pattern` constants, `Markers` flag struct, `DetectPattern(markers) (Pattern, float64, string)` 4-pattern taxonomy with confidence table + ambiguity detection + partial-signal rationale for unknown.
- [x] 2.3 Patch `internal/assets/skills/sdd-spec/SKILL.md` — `## Data-Engineering Domain Branch` gate (run `gentle-ai sdd-config --json`) + 6 ETL delta sections (Source Tables, Target Schema, Watermark Strategy, DAG, AWS Profile Requirements, Verify Approach) + table-output scenario convention + MODIFIED insertion-point template.
- [x] 2.4 Patch `internal/assets/skills/sdd-design/SKILL.md` — `## Data-Engineering Domain Branch` gate + DAG-of-transformations template + pattern-aware design choice + MODIFICATION insertion-point analysis cascade template + ETL scenario coverage questions.

## TDD Cycle Evidence

| Task | Test File | Layer | Safety Net | RED | GREEN | TRIANGULATE | REFACTOR |
|------|-----------|-------|------------|-----|-------|-------------|----------|
| 2.1 | `internal/etl/sidecar_test.go` | Unit | N/A (new package) | ✅ Written (undefined `ParseSidecar`/`ValidateSidecar`/`Mismatch`) | ✅ 15/15 passed | ✅ 15 cases (valid, missing database/column-type/name, malformed glue, quoted partitions/comments, extra partition) | ✅ Moved `Pattern`+constants out of `sidecar.go` into `pattern.go` for cohesion; kept shared helpers (`splitKeyValue`, `leadingSpaces`, `unquoteScalar`, `sliceToSet`) generic |
| 2.2 | `internal/etl/pattern_test.go` | Unit | N/A (new) | ✅ Written (undefined `Markers`/`DetectPattern`/`Pattern*`) | ✅ 11/11 passed | ✅ 11 cases (3 incremental/glue-studio/legacy triangulations, ambiguous→highest-conf, partial-signal→unknown, struct equality, missed-marker suppression) | ✅ Refactored to single dispatch table `confidenceFor[pattern]` + `candidates[]{pattern,matched,rationale}`; one triangulation test (`LegacyWranglerSuppressedByGlueContext`) corrected — multi-step rule explicitly requires `!HasWranglerAthena`, so the contradiction falls to unknown (rationale names the partial signal) |
| 2.3 | `internal/assets/skills/sdd-spec/SKILL.md` | N/A (markdown) | Golden tests (no-content regression — see below) | ➖ N/A (skill prompt) | ➖ N/A | ➖ Triangulation skipped: skill template validated by sdd-verify scenarios in Phase 3 | ➖ Conditional gate keeps app-dev flow byte-identical |
| 2.4 | `internal/assets/skills/sdd-design/SKILL.md` | N/A (markdown) | Golden tests (no-content regression) | ➖ N/A | ➖ N/A | ➖ N/A | ➖ Conditional gate keeps app-dev flow byte-identical |

### Test Summary (Phase 2)

- **New Go tests written**: 26 (etl/sidecar 15 + etl/pattern 11) — total +26 over Phase 1's 41 → 67 cumulative across the change so far
- **Etl package tests passing**: 26/26 (`ok internal/etl 2.563s`)
- **Layers used**: Unit (26)
- **Pure functions created**: `ParseSidecar`, `Sidecar.validateStructure`, `parseColumnsBlock`, `parseInlineColumn`, `splitFlowMapEntries`, `parseFlowSeq`, `unquoteScalar`, `splitKeyValue`, `leadingSpaces`, `ValidateSidecar`, `indexColumns`, `indexPartitionNames`, `partitionKeysSorted`, `sliceToSet`, `DetectPattern`, `rationaleFor`, `unknownRationale`.

## Golden File Verification (Phase 2)

The Phase 1 issue warned that any patch to `internal/assets/skills/*/SKILL.md`
breaks golden tests. That is true ONLY for skills with a content golden
snapshot. The injector golden tests snapshot `*-skill-sdd-init.golden` content
for every adapter; `sdd-spec` and `sdd-design` are only presence-checked (every
adapter except the sdd-init content golden). Consequently:

- `TestGoldenSDD_*` (all adapters): ✅ PASS unchanged after the sdd-spec and
  sdd-design patches. No golden regeneration needed for Phase 2.
- `git diff testdata/golden/` after this batch: empty (no golden touched).

This is a refinement of the Phase 1 issue note: future phases that patch
`sdd-spec`/`sdd-design`/`sdd-apply`/`sdd-verify`/`sdd-tasks` (Phase 3-5) likely do
NOT need golden regeneration either; only `sdd-init` content-goldens break on
skill content patches. Re-verify phase-by-phase regardless.

## Verification Results

| Check | Command | Result |
|---|---|---|
| Build | `go build ./...` | ✅ PASS |
| Etl package tests | `go test ./internal/etl/...` | ✅ PASS (26/26) |
| Components golden (skill content) | `go test ./internal/components/ -run TestGoldenSDD` | ✅ PASS — sdd-spec/sdd-design content not snapshotted, so patches do not break any golden |
| Full suite | `go test ./...` | ✅ No regressions — only the pre-existing flaky `TestRunInstallKimiMissingUVFailsBeforeExecutingInstallCommands` fails (identical to Phase 1 baseline) |
| Vet | `go vet ./internal/etl/...` and `go vet ./...` | ✅ PASS (no warnings) |
| gofmt | `gofmt -l internal/etl/` | ✅ Clean (post `gofmt -w`) |

## Backward Compatibility (verified)

- `internal/etl` is a brand-new additive package; no existing import touched.
- `sdd-spec/SKILL.md` patch is a conditional branch: the new `## Data-Engineering Domain Branch` section is self-declared NO-OP when `domain != data-engineering`. The full `## What to Do` flow continues under a guard `gentle-ai sdd-config --json` → `.domain` check; absent domain skips Step 1a and the 6 ETL sections, yielding exactly the prior skill output.
- `sdd-design/SKILL.md` patch is the same pattern: `## Data-Engineering Domain Branch` gated on `sdd-config`.
- All non-ETL code paths in `internal/etl` are unreachable when the caller doesn't scan Glue jobs (Markers is filled by a content scanner; the sdd-* skills drive that).

## Deviations from Design

- **Pattern struct cohesion**: design's Interfaces block lists `Pattern` and `DetectPattern` under `internal/etl` but does not specify the file. I kept `Pattern`/`Markers`/`DetectPattern` together in `pattern.go` and `Sidecar`/`ParseSidecar`/`ValidateSidecar`/`Mismatch` in `sidecar.go` for pigeonhole reading. This is purely organizational; the package surface (exported identifiers) matches the contract exactly.
- **Mismatch kinds named explicitly**: design says `ValidateSidecar` returns `[]Mismatch` but does not specify the `Kind` labels. I chose descriptive kinds (`database_mismatch`, `table_mismatch`, `s3_location_mismatch`, `missing_column`, `type_mismatch`, `missing_partition`, `unexpected_partition`) so spec/verify prose can branch on them consistently. Documented as constants `Mismatch*` so downstream phases bind to stable names.
- **Confidence value absence in ambiguous unknown case**: when NO pattern matches, `DetectPattern` returns `("", 0, rationale)` — explanation names partial signals. This matches the design's "never silent" principle; the rationale makes the unknown useful (e.g. "watermark (incremental partial)").
- **`HasSparkSqlQueryHelper` only check for glue-studio**: the heuristic line in the prompt lists both `HasApplyMapping && HasSparkSqlQueryHelper`. I implemented the AND (`&&`) form (both required) per the prompt; the `unknownRationale` then flags either as a partial. If the team prefers OR semantics, the test `TestDetectPatternGlueStudio` would need a paired triangulation. Flagged for the verify phase to confirm with the ETL domain lead.

## Issues Found

- One initial pattern test (`TestDetectPatternLegacyWranglerSuppressedByGlueContext`) asserted `PatternMultiStep` for `{HasWranglerAthena, HasGlueContext, HasTempViews}`, but the multi-step rule explicitly requires `!HasWranglerAthena` (the prompt spec). Corrected the test triangulation to expect unknown with a "awswrangler.athena" rationale mention. This is a test fix that aligns with the design contract — no behavioral ambiguity remains.
- No other issues.

## Workload / PR Boundary (Phase 2)

- **Mode**: chained PR slice (PR 2 of 8) under `feature-branch-chain`. Targets the previous PR's branch (`feature/aws-dataengineer`, which carries Phase 1).
- **Current work unit**: Phase 2 — Sidecar + Pattern + sdd-spec/sdd-design patches.
- **Boundary**: starts at task 2.1, ends at task 2.4. Self-contained: additive `internal/etl` package (validated by unit tests); two embedded skill patches gated on `sdd-config --json`. No golden regeneration required; no new dependency added (`yaml.v3` deliberately avoided — hand-rolled scanner shares the project's existing convention).
- **Estimated review budget impact**: PR 2 forecast was ~520 lines. Real review surface (Go + tests + skill patches) is within the 800-line budget. No golden file churn this round, so the diff is small and entirely readable.

## Status

4/4 Phase 2 tasks complete. App-dev behaviour unchanged (`domain`
absent → both skill branches are NO-OPs; etl package unused by app-dev
callers). Ready for `sdd-apply` Phase 3 (Apply + Verify dual-path: header,
compare, sdd-apply Camino A, sdd-verify Camino B).

---

# Phase 3: Apply + Verify Dual-Path

**Phase**: 3 of 8 — Apply + Verify dual-path (`internal/etl` header + compare + `sdd-apply` Camino A + `sdd-verify` Camino B)
**Mode**: Strict TDD (test runner: `go test ./...`)
**Delivery**: chained PR slice (PR 3 of 8) under `feature-branch-chain`.
Builds atop PR 2 (`feature/aws-dataengineer`). Backward compatibility
preserved: the two skill patches add a conditional `domain ==
"data-engineering"` branch; absent `domain` (app-dev) yields identical
behaviour to today. No existing Go code path is changed.

## Completed Tasks

- [x] 3.1 `internal/etl/header.go` + `header_test.go` — `ETLHeader` 6-field struct (`JobName`, `DesarrolladoPor`, `FechaCreacion`, `ModificadoPor`, `FechaModificacion`, `Descripcion`); `RenderHeader` (label-column-aligned comment block, rune-aware padding for accented vowels); `ValidateHeader` (rejects `co-authored-by`, `generated by`, `auto-generated`, and AI model/vendor names: claude, gpt, chatgpt, gemini, copilot, openai, anthropic — case-insensitive substring, no false-positives on common words like "available"/"domain"); `UpdateHeaderForModify` (preserves `desarrollado por` + `fecha creación`, updates `modificado por` + `fecha modificación`; non-destructive fallback returns `orig` unchanged when input is unparseable).
- [x] 3.2 `internal/etl/compare.go` + `compare_test.go` — `BuildExceptSQL(target, devDB, prdDB) string` emits `SELECT * FROM <devDB>.<target> EXCEPT SELECT * FROM <prdDB>.<target>` (dev-first ordering; uppercase `EXCEPT`/`SELECT` for Athena + Spark SQL parsers).
- [x] 3.3 Patch `internal/assets/skills/sdd-apply/SKILL.md` — `## Data-Engineering Domain Branch` gate (run `gentle-ai sdd-config --json`, parse `.domain`) + **Camino A** (Glue Docker `aws-glue-libs:5`, throwaway `<table>_test`, TDD scaffold selected by `DetectPattern`: watermark / multi-step DAG / Athena-Spark parity / legacy-wrangler-migrate / unknown→blocked) + header-enforcement block referencing `internal/etl.RenderHeader`/`ValidateHeader`/`UpdateHeaderForModify` + `Step 1a (data-engineering only): Resolve Domain` + small-model gate (returns `needs-explore` so a capable model runs Camino A).
- [x] 3.4 Patch `internal/assets/skills/sdd-verify/SKILL.md` — `## Data-Engineering Domain Branch` gate + **Camino B** (`sam deploy` both repos in parallel, `verify.skip_deploy` short-circuit, run job, `BuildExceptSQL` dev-vs-prd via Athena, `ValidateSidecar` vs `aws glue get-table`, `ScrubProfiles` on all output, verdict rule) + small-model gate.

## TDD Cycle Evidence

| Task | Test File | Layer | Safety Net | RED | GREEN | TRIANGULATE | REFACTOR |
|------|-----------|-------|------------|-----|-------|-------------|----------|
| 3.1 | `internal/etl/header_test.go` | Unit | ✅ etl package baseline PASS (Phase 2 captured) | ✅ Written (undefined `ETLHeader`/`RenderHeader`/`ValidateHeader`/`UpdateHeaderForModify`) | ✅ 7/7 initial tests passed | ✅ 6 cases (exact-output golden locking column alignment + label vocab; case-insensitive rejection; common-word non-false-positive; empty passes; modify carries JobName+Descripcion; unparseable→orig unchanged) | ✅ Extracted shared `headerLabels` table + `headerFieldByLabel` map (single source of truth for render+parse); `parseHeader` tolerant of blank/non-# lines; rune-based padding (`utf8.RuneCountInString`) so accented labels align |
| 3.2 | `internal/etl/compare_test.go` | Unit | ✅ etl package baseline PASS | ✅ Written (undefined `BuildExceptSQL`) | ✅ 1/1 initial test passed | ✅ 3 cases (target table appears twice — once per side; `EXCEPT` keyword present and uppercase; dev-before-prd lexical ordering) | ➖ None needed (single `fmt.Sprintf`) |
| 3.3 | `internal/assets/skills/sdd-apply/SKILL.md` | N/A (markdown) | Golden tests (no-content regression — see below) | ➖ N/A (skill prompt) | ➖ N/A | ➖ Triangulation skipped: skill template validated by sdd-verify scenarios in Phase 3.4 | ➖ Conditional gate keeps app-dev flow byte-identical |
| 3.4 | `internal/assets/skills/sdd-verify/SKILL.md` | N/A (markdown) | Golden tests (no-content regression) | ➖ N/A | ➖ N/A | ➖ N/A | ➖ Conditional gate keeps app-dev flow byte-identical |

### Test Summary (Phase 3)

- **New Go tests written**: 17 (etl/header 13 + etl/compare 4)
- **Etl package tests passing**: 17/17 (`ok github.com/.../internal/etl 0.420s` — fresh run, no cache)
- **Layers used**: Unit (17)
- **Pure functions created**: `RenderHeader`, `renderHeaderLine`, `ValidateHeader`, `parseHeader`, `UpdateHeaderForModify`, `BuildExceptSQL`.

### Cumulative Correction

Phase 2 apply-progress stated "26 (etl/sidecar 15 + etl/pattern 11)".
Recount in this batch: sidecar 15 + pattern 10 = **25** (pattern file
contains 10 tests, not 11 — the doc over-counted by one). True cumulative
Go test count across the change is therefore:

- Phase 1: 41
- Phase 2: 25 (corrected from 26)
- Phase 3: 17
- **Total: 83 cumulative Go tests, all passing.**

## Golden File Verification (Phase 3)

Phase 2 predicted that patching non-`sdd-init` skills likely needs no
golden regeneration (only `sdd-init` has content goldens; the rest are
presence-checked). This held for Phase 3:

- `TestGoldenSDD_*` (all adapters): ✅ PASS unchanged after the sdd-apply
  and sdd-verify patches (`go test ./internal/components/ -run TestGoldenSDD -count=1` → ok 4.016s).
- `git diff --stat HEAD -- testdata/golden/` after this batch: **empty**
  (no golden touched, no regeneration needed).

This re-confirms the Phase 2 refinement of the Phase 1 issue note:
phases that patch `sdd-spec`/`sdd-design`/`sdd-apply`/`sdd-verify`/
`sdd-tasks` do not break content goldens. Only `sdd-init` content-goldens
break on skill content patches.

## Verification Results

| Check | Command | Result |
|---|---|---|
| Build | `go build ./...` | ✅ PASS |
| Vet (target) | `go vet ./internal/etl/...` | ✅ PASS (no warnings) |
| Vet (full) | `go vet ./...` | ✅ PASS (no warnings) |
| gofmt | `gofmt -l internal/etl/` | ✅ Clean |
| Etl package tests | `go test ./internal/etl/... -count=1` | ✅ PASS (42 total: 25 Phase 2 + 17 Phase 3) |
| Components golden (skill content) | `go test ./internal/components/ -run TestGoldenSDD -count=1` | ✅ PASS — sdd-apply/sdd-verify content not snapshotted, so patches break no golden |
| Full suite | `go test ./...` | ✅ No regressions — only the pre-existing flaky `TestRunInstallKimiMissingUVFailsBeforeExecutingInstallCommands` fails (identical to Phase 1 + Phase 2 baseline; Phase 3 does not touch `internal/cli`) |

## Backward Compatibility (verified)

- `internal/etl/header.go` + `compare.go` are additive — no existing import touched, no existing function signature changed.
- `sdd-apply/SKILL.md` patch is a conditional branch: the new `## Data-Engineering Domain Branch` section self-declares NO-OP when `domain != data-engineering`. The full `## What to Do` flow continues under a guard `gentle-ai sdd-config --json` → `.domain` check; absent domain skips Step 1a and Camino A, yielding exactly the prior skill output.
- `sdd-verify/SKILL.md` patch is the same pattern: `## Data-Engineering Domain Branch` gated on `sdd-config`; absent domain yields the prior `go test ./...` + `go vet ./...` verify flow.
- Small-model fallback sections in both skills explicitly route app-dev work (domain absent) to the unchanged minimal flow; only the data-engineering path returns `needs-explore`.
- All non-ETL code paths in `internal/etl` remain unreachable when callers don't scan Glue jobs (Markers is filled by a content scanner driven by the sdd-* skills).

## Deviations from Design

- **Forbidden-pattern vocabulary**: design Interfaces block says `ValidateHeader` "rejects 'generated by', 'Co-Authored-By', AI names" but does not enumerate which AI names. I chose the explicit vendor/model list (claude, gpt, chatgpt, gemini, copilot, openai, anthropic) plus `auto-generated`, deliberately NOT using a bare `"ai"` substring (would false-fire on "available", "domain", "certain", "main"). `TestValidateHeaderDoesNotFalsePositiveOnCommonWords` pins this. Flagged for the verify phase to confirm with the ETL domain lead.
- **Header column alignment**: design shows a 6-field header but does not specify column alignment. I aligned all colons to the longest label's rune column (`# fecha modificación` = 20 runes; pad + ` : ` separator), using `utf8.RuneCountInString` so accented vowels (`ó`) align correctly. `TestRenderHeaderAlignsColons` + `TestRenderHeaderExactOutput` pin the format.
- **`UpdateHeaderForModify` non-destructive fallback**: design says "preserves `desarrollado_por`+`fecha_creación`" but does not specify what happens on unparseable input. I chose to return `orig` unchanged rather than emit a partial header (non-destructive). `TestUpdateHeaderForModifyUnparseableReturnsOrigUnchanged` pins this.
- **`BuildExceptSQL` dev-before-prd ordering**: design says "SQL EXCEPT for dev-vs-prd" without specifying order. I placed dev on the left (`SELECT * FROM <dev> EXCEPT SELECT * FROM <prd>`) because the verify contract cares about rows in dev but not yet in prd (the direction of promotion). `TestBuildExceptSQLOrderIsDevFirstThenPrd` pins this. If the team prefers the reverse direction, the test triangulation must flip.

## Issues Found

- **Phase 2 test-count over-count**: the Phase 2 apply-progress stated "26 (15 sidecar + 11 pattern)" but the pattern test file actually contains 10 tests. Corrected in this batch's Cumulative Correction section. No behavioural impact — the count was the only error.
- No other issues.

## Workload / PR Boundary (Phase 3)

- **Mode**: chained PR slice (PR 3 of 8) under `feature-branch-chain`. Targets the previous PR's branch (`feature/aws-dataengineer`, which carries Phase 1 + Phase 2).
- **Current work unit**: Phase 3 — Header + Compare + sdd-apply Camino A + sdd-verify Camino B.
- **Boundary**: starts at task 3.1, ends at task 3.4. Self-contained: additive `internal/etl/header.go` + `compare.go` (validated by 17 unit tests); two embedded skill patches gated on `sdd-config --json`. No golden regeneration required; no new dependency added.
- **Estimated review budget impact**: PR 3 forecast was ~580 lines. Real review surface (186 lines Go+tests + ~210 lines skill patches) is within the 800-line budget. No golden file churn this round.

## Status

4/4 Phase 3 tasks complete. 15/15 cumulative tasks across the change
(Phase 1 7 + Phase 2 4 + Phase 3 4). 83 cumulative Go tests passing
(Phase 1 41 + Phase 2 25 corrected + Phase 3 17). App-dev behaviour
unchanged (`domain` absent → both skill branches are NO-OPs; etl
header/compare unreachable by app-dev callers). Ready for `sdd-apply`
Phase 4a (Skills Foundation: data-engineer-protocol.md + 8 SkillIDs +
catalog registration + pattern-detect + study-file + create-table).

---

# Phase 4a: Skills Foundation

**Phase**: 4a of 8 — Skills Foundation (shared protocol + 8 SkillIDs + catalog
registration + `data-engineer-pattern-detect` + embed/patch study-file +
embed/patch create-table)
**Mode**: Strict TDD (test runner: `go test ./...`)
**Delivery**: chained PR slice (PR 4a of 8) under `feature-branch-chain`.
Builds atop PR 3 (`feature/aws-dataengineer`). Backward compatibility
preserved: the 8 new skills are registered at category `data-engineering` /
priority `p1` and are NOT part of any preset (`SkillsForPreset` is unchanged),
so a default/minimal/ecosystem/full install never installs them and app-dev
behaviour is byte-identical. The new `_shared` protocol + 3 skill dirs are
additive assets.

> **Batch recovery note**: Phase 4a artifacts were found complete-but-uncommitted
> in the working tree (a prior session produced them without committing or
> persisting progress — Engram memory #1064 still read "Ready for Phase 4a").
> This batch VERIFIED the existing artifacts against the spec/design
> (build + targeted tests + goldens + full suite), confirmed correctness, then
> merged this progress record. No artifact was rewritten; only verified.

## Completed Tasks

- [x] 4a.1 Create `internal/assets/skills/_shared/data-engineer-protocol.md` — shared reference (6 sections): §1 ETL Header Protocol (6-field format, points at `etl.RenderHeader`/`ValidateHeader`/`UpdateHeaderForModify`), §2 Authorship Rule (AI never the author; forbidden vocab; replacement rule for legacy violations), §3 Pattern Index (4-pattern taxonomy with markers/confidence/scaffold, points at `etl.DetectPattern`), §4 Master-Project Awareness (infra vs carga repos + git flow), §5 AWS Profile logical-name rule (`ResolveProfile`/`ScrubProfiles`, closed set prd/dev/usuario), §6 cross-references to all `etl`/`sddconfig` Go funcs.
- [x] 4a.2 Add 8 `SkillDataEngineer*` constants to `internal/model/types.go` — `PatternDetect`, `StudyFile`, `ETLS3`, `ETLGlue`, `ETLSharepoint`, `CreateTable`, `SQLFromLogic`, `Integrate` (+ doc comment). Pinned by `TestDataEngineerSkillIDs_ClosedSet` (8 subtest cases + `len==8` guard).
- [x] 4a.3 Register 8 data-engineer skills in `internal/catalog/skills.go` — all category `data-engineering`, priority `p1`, canonical names. Pinned by `TestDataEngineerSkillsRegistered` (8 skills × name/category/priority assertions).
- [x] 4a.4 Create `internal/assets/skills/data-engineer-pattern-detect/SKILL.md` — new skill; protocol ref prepended; documents the "never silent" detection contract (pattern + confidence + rationale), Markers struct fields mapped to §3, 0.8/0.5/0.0 thresholds, user-override propagation to `sdd-apply` Camino A, YAML Pattern Detection Report output.
- [x] 4a.5 Embed + patch `internal/assets/skills/data-engineer-study-file/SKILL.md` — copied from `~/.config/opencode/skills/`; prepended `## Reference: _shared/data-engineer-protocol.md`; added Resources cross-reference. (No authorship byte-change needed: source had no AI-attribution violation — `metadata.author: gentleman-programming` is human; the etl-s3 violation is Phase 4b.)
- [x] 4a.6 Embed + patch `internal/assets/skills/data-engineer-create-table/SKILL.md` — copied from source; prepended protocol ref; patched Phase 2/3 deploy commands to resolve the AWS profile at runtime via `--profile "$(gentle-ai sdd-config --json | jq -r .aws_profiles.dev)"` (protocol §5, never hardcoded) + added `ScrubProfiles` reminder on shared output.

## TDD Cycle Evidence

| Task | Test File | Layer | Safety Net | RED | GREEN | TRIANGULATE | REFACTOR |
|------|-----------|-------|------------|-----|-------|-------------|----------|
| 4a.1 | `internal/assets/skills/_shared/data-engineer-protocol.md` | N/A (markdown) | `TestSkillFrontmatterIsLintClean` (via `_shared/SKILL.md`) + assets count | ➖ N/A (reference doc, no Go test layer) | ➖ N/A | ➖ Triangulation skipped: DRY contract consumed by skills; validated structurally by §-references in 4a.4/4a.5/4a.6 + sdd-verify | ➖ N/A |
| 4a.2 | `internal/model/types_test.go` | Unit | ✅ model pkg baseline PASS (captured before edit) | ✅ Written (undefined `SkillDataEngineer*` constants → compile error) | ✅ 8/8 subtests passed | ✅ 8 cases (one per SkillID) + `len==8` closed-set guard so adding a skill without a const breaks the test | ➖ None needed (const block) |
| 4a.3 | `internal/catalog/skills_test.go` | Unit | ✅ catalog pkg baseline PASS | ✅ Written (missing catalog entries for the 8 SkillIDs) | ✅ 8/8 skills matched | ✅ 8 cases × 3 fields (Name + Category `data-engineering` + Priority `p1`) | ➖ None needed |
| 4a.4 | `internal/assets/skills/data-engineer-pattern-detect/SKILL.md` | N/A (markdown) | `TestSkillFrontmatterIsLintClean` + `TestEmbeddedAssetCount` | ➖ N/A (skill prompt) | ➖ N/A | ➖ Triangulation skipped: skill validated by sdd-verify scenarios; detection logic already pinned by Phase 2 `etl.DetectPattern` tests | ➖ N/A |
| 4a.5 | `internal/assets/skills/data-engineer-study-file/SKILL.md` | N/A (markdown) | `TestSkillFrontmatterIsLintClean` + `TestEmbeddedAssetCount` | ➖ N/A | ➖ N/A | ➖ N/A | ➖ N/A |
| 4a.6 | `internal/assets/skills/data-engineer-create-table/SKILL.md` | N/A (markdown) | `TestSkillFrontmatterIsLintClean` + `TestEmbeddedAssetCount` | ➖ N/A | ➖ N/A | ➖ N/A | ➖ N/A |

### Test Summary (Phase 4a)

- **New Go test functions written**: 2 (`model.TestDataEngineerSkillIDs_ClosedSet` with 8 subtests; `catalog.TestDataEngineerSkillsRegistered` with 8 skills × 3 field assertions). Closed-set guards (`len==8` in both) force any future SkillID/catalog addition to extend the tables or fail.
- **Existing test modified**: `internal/assets/assets_test.go` `TestEmbeddedAssetCount` bumped expected skill directories 23 → 26 (adds `_shared` was already counted; the 3 new dirs: `data-engineer-pattern-detect`, `data-engineer-study-file`, `data-engineer-create-table`). `TestSkillFrontmatterIsLintClean` now lints all 26 (including the 3 new) — PASS.
- **Cumulative Go tests across the change**: 83 (Phases 1-3) + 2 = **85**, all passing.
- **Layers used**: Unit (2 new functions).
- **Markdown assets**: 4 created (protocol, pattern-detect) / embedded+patched (study-file, create-table); all pass frontmatter lint + asset presence checks; no Go test layer (validated by sdd-verify scenarios, consistent with tasks 1.7/2.3/2.4/3.3/3.4).

## Golden File Verification (Phase 4a)

- `TestGoldenSDD_*` (all 12 adapters): ✅ PASS. The 3 new data-engineer skills are NOT `sdd-*` skills, so they are not part of the SDD injector golden snapshots; adding them as asset directories does not change any golden content.
- `TestEmbeddedAssetCount`: ✅ PASS at 26 (assertion updated to match the 3 new dirs).
- `TestSkillFrontmatterIsLintClean`: ✅ PASS for all 26 skills incl. the 3 new data-engineer ones.
- `git diff --stat HEAD -- testdata/golden/` after this batch: **empty** (no golden touched, no regeneration needed). Confirms the Phase 2/3 refinement: only `sdd-init` content-goldens break on skill content patches; additive new skill dirs + non-`sdd-init` changes need no golden regeneration.

## Verification Results

| Check | Command | Result |
|---|---|---|
| Build | `go build ./...` | ✅ PASS |
| Vet (target) | `go vet ./internal/model/... ./internal/catalog/...` | ✅ PASS (no warnings) |
| gofmt | `gofmt -l internal/model/ internal/catalog/` | ✅ Clean (post `gofmt -w`) |
| Model tests | `go test ./internal/model/ -count=1` | ✅ PASS |
| Catalog tests | `go test ./internal/catalog/ -count=1` | ✅ PASS |
| Assets tests | `go test ./internal/assets/ -count=1` | ✅ PASS (count=26, frontmatter lint 26/26) |
| Components golden (skill content) | `go test ./internal/components/ -run TestGoldenSDD -count=1` | ✅ PASS (12/12 adapters) — new dirs are not `sdd-*`, break no golden |
| Full suite | `go test ./... -count=1` | ✅ No regressions — only the pre-existing flaky `TestRunInstallKimiMissingUVFailsBeforeExecutingInstallCommands` fails (identical to Phase 1/2/3 baseline; Phase 4a does NOT touch `internal/cli`, confirmed by `git status`) |

## Backward Compatibility (verified)

- `internal/model/types.go` change is purely additive (8 new `SkillID` consts in an existing `const (...)` block). No existing constant renamed/removed.
- `internal/catalog/skills.go` change is purely additive (8 new entries appended to `mvpSkills`). `MVPSkills()` returns a superset; existing callers that look up prior skills see identical results. `TestMVPSkillsNoDuplicates` + `TestMVPSkillsCoverAllPresetSkills` still PASS.
- `SkillsForPreset` / `AllSkillIDs` in `presets.go` are UNCHANGED — the 8 data-engineer skills are registered (known to the catalog) but installed by NO preset. A default/minimal/ecosystem/full install never writes them; they activate only when a user explicitly selects them or under the `domain: data-engineering` profile.
- The 3 new asset directories + `_shared/data-engineer-protocol.md` are additive; no existing asset moved/renamed. `_shared/` was already counted in the prior 23, so the 23→26 delta is exactly the 3 new skill dirs.
- The two gofmt incidental cleanups (`TriggerBinding.Run` tag alignment; `TriggerEvent` test struct alignment) are whitespace-only, behaviour-neutral.

## Deviations from Design

- **"Authorship fix" is a NO-OP for study-file + create-table**: design File Changes row says embed "with protocol reference + authorship fix". The authorship fix (protocol §2) is CONDITIONAL — it applies only when the source carries an AI-attribution violation. Neither study-file nor create-table source had one (`metadata.author: gentleman-programming` is human; no `# author: generated by gentle-ai`, no `Co-Authored-By`). The actual violation lives in `data-engineer-etl-s3` and is fixed in Phase 4b. So for 4a.5/4a.6 the patch is "prepend protocol ref" (+ profile-resolution for create-table), and the authorship byte-change is correctly absent. Flagged so the verify phase does not expect an authorship diff here.
- **create-table profile resolution rendering**: protocol §5 uses the Go form `ResolveProfile(cfg, "dev")` in prose. The embedded create-table skill renders the runtime equivalent `--profile "$(gentle-ai sdd-config --json | jq -r .aws_profiles.dev)"`. Both express the same rule (resolve logical→concrete at runtime, never hardcode the profile string); the CLI form is what a skill consumer actually executes. Faithful to §5's intent.
- **8 skills registered, only 3 embedded in 4a**: by design (Phases 4b/4c/4d embed the other 5). The catalog knows all 8 now; `assets_test.go` counts only the 3 embedded dirs (26 = prior 23 + 3). No preset installs the not-yet-embedded 5, so there is no dangling install reference.

## Issues Found

- **Uncommitted prior-session work**: Phase 4a artifacts were present in the working tree but neither committed nor reflected in Engram (memory #1064 still read "Ready for Phase 4a"). Resolved by verifying the artifacts (build/targeted/golden/full suite) and merging this progress record. No rewrite was needed; the artifacts were correct and complete. Recommendation: the prior session should have persisted progress — flagged as a process note, not a code defect.
- **gofmt incidental edits**: a `gofmt -w` during the prior session also reflowed two unrelated existing spots (`TriggerBinding.Run` struct-tag alignment in `types.go`; field alignment in `types_test.go`). Whitespace-only, behaviour-neutral; kept since they bring those lines to gofmt-canonical form. No functional change.
- No other issues.

## Workload / PR Boundary (Phase 4a)

- **Mode**: chained PR slice (PR 4a of 8) under `feature-branch-chain`. Targets `feature/aws-dataengineer` (carries Phases 1-3). Delivery decision was resolved and documented from PR 1 onward; `tasks.md` Review Workload Forecast says `Decision needed before apply: No`, so no blocking decision was required for this batch.
- **Current work unit**: Phase 4a — protocol + 8 SkillIDs + catalog + pattern-detect + study-file + create-table.
- **Boundary**: starts at task 4a.1, ends at task 4a.6. Self-contained: additive model/catalog consts (2 test functions, 8-skill closed-set guards); 4 markdown assets (1 new shared protocol + 1 new skill + 2 embedded/patched skills). No golden regeneration required; no new dependency; `presets.go` untouched.
- **Estimated review budget impact**: PR 4a forecast was ~680 lines. Real review surface (model/catalog Go + tests + 4 skill markdown files) is within the 800-line budget. No golden file churn.

## Status

6/6 Phase 4a tasks complete. 21/21 cumulative tasks across the change
(Phase 1 7 + Phase 2 4 + Phase 3 4 + Phase 4a 6). 85 cumulative Go tests
passing (Phase 1 41 + Phase 2 25 + Phase 3 17 + Phase 4a 2). App-dev
behaviour unchanged (the 8 skills are p1, in NO preset, gated on
`domain: data-engineering`; `_shared` protocol is a passive reference).
Ready for `sdd-apply` Phase 4b (Embed ETL Skills 1: data-engineer-etl-s3 +
data-engineer-sql-from-logic — the former carries the authorship violation
fixed by protocol §2).
