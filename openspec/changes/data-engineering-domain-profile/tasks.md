# Tasks: Data Engineering Domain Profile

## Open Questions Resolved

1. **Phase 4 split**: Split into 4a (protocol + pattern-detect + catalog + study-file + create-table), 4b (etl-s3 + sql-from-logic), 4c (etl-glue + integrate), 4d (etl-sharepoint).
2. **sdd-config vs sdd-status**: Keep `sdd-config` separate; do not modify `sdd-status`.
3. **Confidence thresholds**: 0.8 (both markers) / 0.5 (one marker).

## Task Conventions

- All tasks are `[master]` unless prefixed `[infra]`/`[carga]`.
- Go tasks are test-first (`*_test.go` before implementation).
- Skill patches are gated on `domain: data-engineering`.
- Phase 3/5 depend on Phase 1 config; Phase 2 templates use Phase 2 `etl` helpers.

## Review Workload Forecast

| PR | Scope | ~Lines |
|---|---|---|
| 1 | Config + detect + CLI | 730 |
| 2 | Sidecar + pattern + spec/design | 520 |
| 3 | Header + compare + apply/verify | 580 |
| 4a | Protocol + pattern-detect + catalog + 2 skills | 680 |
| 4b | Embed etl-s3 + sql-from-logic | 540 |
| 4c | Embed etl-glue + integrate | 570 |
| 4d | Embed etl-sharepoint | 340 |
| 5 | Tasks prefix + git flow | 190 |

Decision needed before apply: No
Chained PRs recommended: Yes
Chain strategy: stacked-to-main
400-line budget risk: High

## Phase 1: Detection + Config

- [x] 1.1 `internal/sddconfig/config.go` + `_test.go`: `Config`, `Repos`, `AWSProfiles`, `VerifyOpts`, `LoadConfig`.
- [x] 1.2 `internal/sddconfig/detect.go` + `_test.go`: `DetectDomain` (0.8/0.5 confidence).
- [x] 1.3 `internal/sddconfig/repos.go` + `_test.go`: `ValidateRepos`.
- [x] 1.4 `internal/sddconfig/profiles.go` + `_test.go`: `ResolveProfile`, `ScrubProfiles`.
- [x] 1.5 `internal/cli/sdd_config.go` + `_test.go`: `RunSDDConfig` (`--json`, `--detect`, `--validate-repos`).
- [x] 1.6 Wire `sdd-config` subcommand in `cmd/gentle-ai/main.go`.
- [x] 1.7 Patch `internal/assets/skills/sdd-init/SKILL.md`: preflight detect + confirm/override.

## Phase 2: Spec + Design Templates

- [x] 2.1 `internal/etl/sidecar.go` + `_test.go`: `Sidecar`, `ParseSidecar`, `ValidateSidecar`.
- [x] 2.2 `internal/etl/pattern.go` + `_test.go`: `Pattern`, `DetectPattern` (4 patterns + ambiguous).
- [x] 2.3 Patch `internal/assets/skills/sdd-spec/SKILL.md`: ETL delta sections + sidecar template.
- [x] 2.4 Patch `internal/assets/skills/sdd-design/SKILL.md`: DAG + insertion-point template.

## Phase 3: Apply + Verify Dual-Path

- [x] 3.1 `internal/etl/header.go` + `_test.go`: `ETLHeader`, `RenderHeader`, `ValidateHeader`, `UpdateHeaderForModify`.
- [x] 3.2 `internal/etl/compare.go` + `_test.go`: `BuildExceptSQL`.
- [x] 3.3 Patch `internal/assets/skills/sdd-apply/SKILL.md`: Camino A branch.
- [x] 3.4 Patch `internal/assets/skills/sdd-verify/SKILL.md`: Camino B branch.

## Phase 4a: Skills Foundation

- [ ] 4a.1 Create `internal/assets/skills/_shared/data-engineer-protocol.md`.
- [ ] 4a.2 Add 8 `SkillDataEngineer*` constants to `internal/model/types.go`.
- [ ] 4a.3 Register 8 data-engineer skills in `internal/catalog/skills.go`.
- [ ] 4a.4 Create `internal/assets/skills/data-engineer-pattern-detect/SKILL.md`.
- [ ] 4a.5 Embed + patch `internal/assets/skills/data-engineer-study-file/SKILL.md`.
- [ ] 4a.6 Embed + patch `internal/assets/skills/data-engineer-create-table/SKILL.md`.

## Phase 4b: Embed ETL Skills 1

- [ ] 4b.1 Embed + patch `internal/assets/skills/data-engineer-etl-s3/SKILL.md`.
- [ ] 4b.2 Embed + patch `internal/assets/skills/data-engineer-sql-from-logic/SKILL.md`.

## Phase 4c: Embed ETL Skills 2

- [ ] 4c.1 Embed + patch `internal/assets/skills/data-engineer-etl-glue/SKILL.md`.
- [ ] 4c.2 Embed + patch `internal/assets/skills/data-engineer-integrate/SKILL.md`.

## Phase 4d: Embed ETL Skills 3

- [ ] 4d.1 Embed + patch `internal/assets/skills/data-engineer-etl-sharepoint/SKILL.md`.

## Phase 5: Tasks Multi-Repo + Git Flow

- [ ] 5.1 `internal/etl/tasks.go` + `_test.go`: `RepoPrefix`, `GitFlowForRepo`.
- [ ] 5.2 Patch `internal/assets/skills/sdd-tasks/SKILL.md`: emit `repo: infra|carga|both` prefix + git-flow annotation.
