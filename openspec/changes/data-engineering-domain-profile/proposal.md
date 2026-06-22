# Proposal: Data Engineering Domain Profile for SDD

## Intent

SDD's 8-phase core (`explore → propose → spec → design → tasks → apply → verify → archive`) is domain-agnostic and stays untouched. The mismatch lives in the **4 layers around the core** — skills, verify, spec templates, tools — all of which currently assume **app-dev** (Go/TUI: `go test ./...`, code-assertion verify, generic Given/When/Then specs). Data engineering needs a **domain profile** that swaps these layers: two-repo (infra + carga) coordination, 4 ETL patterns, data-comparison verify (not code assertions), Glue-table schema specs, AWS profile awareness, and the company Bitbucket git flow. Source of decisions: exploration #1052 + 22 domain observations (#1030–#1055).

## Scope

### In Scope
- `openspec/config.yaml` schema: `domain`, `repos`, `aws_profiles` fields
- Hybrid domain detection (auto-detect + preflight confirm)
- ETL spec format: OpenSpec delta + `glue-tables/{db}.{table}.yaml` sidecar
- `sdd-apply` Camino A (Glue Docker TDD) + `sdd-verify` Camino B (SAM deploy + Athena dev-vs-prd) mode-branch
- Patch 7 `data-engineer-*` skills + new `data-engineer-pattern-detect` skill
- `sdd-tasks` multi-repo `repo: infra|carga|both` prefix + git-flow annotation

### Out of Scope (Non-Goals)
- Does NOT modify the SDD 8-phase core or the orchestrator dispatch
- Does NOT modify the company repos (infra/carga) directly — only the SDD layer that coordinates them
- Does NOT create a new SDD orchestrator skill
- Does NOT adopt Floci local AWS emulation (descarted, #1038 — user has dev AWS access)
- Does NOT touch the model-profiles feature (`openspec/specs/sdd-profiles` governs model assignment, unrelated to domain profiles)

## Capabilities

> Contract for `sdd-spec`. Researched `openspec/specs/` first. Existing `sdd-profiles` = model-assignment profiles (no overlap with domain profiles).

### New Capabilities
- `data-engineering-domain`: Domain profile behavior — hybrid detection (`template.yaml` + `glue-jobs/*.py`), config schema (`domain`/`repos`/`aws_profiles`), phase-branching contract (verify mode-branch when `domain: data-engineering`), multi-repo coordination (`repo:` prefix), dual git-flow awareness (master project = GitHub; company repos = Bitbucket `feature/ → develop → release/`).
- `etl-spec-format`: OpenSpec delta + `glue-tables/{db}.{table}.yaml` sidecar. Behavioral Given/When/Then (data comparison) in delta; schema/partitions/S3-path/compression in sidecar; reviewed together; validated against deployed Glue table via `aws glue get-table`.

### Modified Capabilities
- None. The 8-phase core, orchestrator, and model-profiles spec are untouched.

## Confirmed Decisions

| # | Decision | Resolution |
|---|----------|------------|
| Q1 | Domain detection | **Hybrid** — auto-detect (`template.yaml` + `glue-jobs/*.py`) + preflight confirm; persists `domain: data-engineering` in config |
| Q2 | Skills adaptation | **Patch 7 + new skill** — shared protocol block (header + authorship + 4-pattern index); see Phase 4 |
| Q3 | Verify paradigm | **Mode-branch inside `sdd-verify`** — not a separate skill; branches when `domain: data-engineering` |
| Q4 | Spec template | **OpenSpec delta + sidecar YAML** — behavior in delta, schema in `glue-tables/*.yaml` |
| Q5 | Multi-repo | **Config declares repos + `repo:` task prefix** — `sdd-tasks` emits prefix automatically |
| Q6 | Pattern detection | **New `data-engineer-pattern-detect` skill + user override** — heuristic scan + confidence; 4-pattern menu for new ETLs |
| +  | Company git flow | Company repos (Bitbucket): `feature/ → PR develop → validate → clone release/`; hotfix injects into `release/`. `branch-pr` skill (GitHub) does NOT apply to company repo tasks. Master project (this repo) uses GitHub flow. |

## Approach — 5-Phase Vertical Slices

| Phase | Deliverable | Lands independently |
|-------|-------------|---------------------|
| **1. Detection + config** | `domain`/`repos`/`aws_profiles` in `config.yaml`; `sdd-init` preflight + auto-detect hint | Yes — config-only, no behavior change until phases read it |
| **2. Spec + design templates** | ETL delta sections (`## Source Tables`, `## Target Schema`, `## Watermark Strategy`, `## DAG`, `## AWS Profile Requirements`, `## Verify Approach`); sidecar YAML format; `sdd-spec`/`sdd-design` templates | Yes — templates only |
| **3. Apply + Verify dual-path** | `sdd-apply` Camino A (Glue Docker `aws-glue-libs:5`, throwaway table, header-protocol + authorship enforcement); `sdd-verify` Camino B (SAM deploy both repos, job run, Athena dev-vs-prd via SQL EXCEPT) | Yes — behind `domain` config gate |
| **4. Skills update** | Patch 7 `data-engineer-*` with shared protocol block; add `data-engineer-pattern-detect` | Yes — additive skill changes |
| **5. Tasks multi-repo** | `sdd-tasks` emits `repo: infra|carga|both` prefix; reads `repos:` from config; annotates git flow per task | Yes — depends on Phase 1 config |

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `openspec/config.yaml` | Modified | Add `domain`, `repos`, `aws_profiles`, `data_engineering` verify section |
| `~/.config/opencode/skills/sdd-init` | Modified | Preflight domain question + auto-detect |
| `~/.config/opencode/skills/sdd-spec` | Modified | ETL delta + sidecar template |
| `~/.config/opencode/skills/sdd-design` | Modified | DAG-of-transformations + insertion-point analysis |
| `~/.config/opencode/skills/sdd-apply` | Modified | Camino A branch (Glue Docker TDD + header protocol) |
| `~/.config/opencode/skills/sdd-verify` | Modified | Camino B branch (SAM deploy + Athena compare) |
| `~/.config/opencode/skills/sdd-tasks` | Modified | `repo:` prefix + git-flow annotation |
| `~/.config/opencode/skills/data-engineer-*` (7) | Modified | Shared protocol block (header/authorship/4-pattern/master-project/AWS-profile awareness) |
| `~/.config/opencode/skills/data-engineer-pattern-detect` | New | Heuristic pattern scan + confidence report |
| `openspec/specs/data-engineering-domain/` | New | Domain profile spec |
| `openspec/specs/etl-spec-format/` | New | Delta + sidecar format spec |

## Delivery Strategy

- **5 chained PRs** to master project `main` (GitHub flow), one per phase.
- `chain-strategy: stacked-to-main`; each PR ≤ **800 lines** (review budget).
- Each PR is a vertical slice that can land independently behind the `domain` config gate.
- Company-repo code changes (infra/carga) are PR'd separately via Bitbucket `feature/ → develop → release/` — out of scope for these 5 PRs.

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Verify bifurcation drift (app-dev vs data-eng in one skill) | Med | Shared test fixtures exercising both branches; CI runs both |
| Pattern-detect misclassification → wrong TDD scaffold | Med | Confidence score + user override in `sdd-propose`; never silent |
| YAML sidecar ↔ delta spec drift | Med | `sdd-verify` validates sidecar vs deployed Glue table (`aws glue get-table`) |
| Repo path fragility (moved `repositorios/`) | Low | `sdd-init` re-validates paths each session; warn if missing |
| SAM deploy cost (2 deploys/verify run, >5 min each) | Med | Parallel deploys; cache stack outputs; `verify: skip-deploy` for local-only |
| Authorship rule breach (AI attribution slips in) | Low | Shared header template; `desarrollado_por`/`modificado_por` from session config |
| AWS profile leak in logs | Low | Scrub profile names from log output; document dev-only commands |

## Rollback Plan

- Each phase lands behind the `domain: data-engineering` config gate; app-dev repos are unaffected — revert = unset `domain` field.
- Per-phase revert: `git revert <phase-PR>`; phases are ordered and mostly independent (Phase 5 depends on Phase 1 config only).
- Skills changes are additive (new branches, new skill) — removing the branch restores prior behavior; the new `data-engineer-pattern-detect` skill can be deleted with no downstream dependency.
- Config schema additions are backward-compatible (new optional fields); old `config.yaml` still parses.

## Dependencies

- Glue Docker image `public.ecr.aws/glue/aws-glue-libs:5` for Camino A (already chosen, #1039).
- AWS CLI with `prd`/`dev`/`usuario` profiles (#1033) for Camino B.
- Existing master-project layout (`repositorios/` + `openspec/` at master level, #1036).

## Success Criteria

- [ ] `openspec/config.yaml` with `domain: data-engineering` makes every phase branch without re-asking
- [ ] A data-engineering change produces a delta spec + `glue-tables/*.yaml` sidecar reviewed together
- [ ] `sdd-apply` runs a Camino A TDD loop against Glue Docker; `sdd-verify` runs Camino B SAM-deploy + Athena dev-vs-prd comparison
- [ ] `data-engineer-pattern-detect` reports a pattern + confidence; user can override
- [ ] `sdd-tasks` emits `repo: infra|carga|both` prefixes and annotates git flow per task
- [ ] Each of the 5 PRs ≤ 800 lines and lands on `main` independently
- [ ] App-dev repos (no `domain` field) behave identically to before (no regression)
