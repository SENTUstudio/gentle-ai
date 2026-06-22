# Exploration: Data Engineering Domain Profile for SDD

> **Status**: working exploration — produced by `sdd-explore` for change
> `data-engineering-domain-profile`. Maps the territory before
> `sdd-propose`. Reads 22 domain observations in Engram under topic
> prefix `architecture/data-engineer-domain/*` (the user-captured domain
> knowledge from prior sessions).

## Current State

SDD's 8-phase core (`explore → propose → spec → design → tasks → apply → verify → archive`) is **domain-agnostic** and the orchestrator pattern is solid. The mismatch lives in the **4 layers around the core**, all of which currently assume app-dev:

1. **Skills** — `~/.config/opencode/skills/` has app-dev skills
   (`go-testing`, `branch-pr`, `work-unit-commits`, `chained-pr`,
   `judgment-day`) and 7 `data-engineer-*` skills that exist but
   pre-date SDD, pre-date the 4-pattern taxonomy, pre-date the
   two-repo invariant, and pre-date the master-project pattern. They
   have no awareness of SDD phases.
2. **Verify** — `openspec/config.yaml` mandates `go test ./...` and
   `go vet ./...`. There is no concept of table-comparison, dev-vs-prd,
   or SAM-deploy validation. The verify skill is a code-assertion gate.
3. **Spec templates** — OpenSpec delta specs use generic
   `Requirement / Scenario` (Given/When/Then) format with no
   glue-table, partition, watermark, or transformation-DAG concepts.
4. **Tools** — Context7 MCP for library docs is wired in. There is
   no AWS CLI profile awareness, no Glue Docker integration, no
   SAM-deploy convention.

### What's already solved (don't re-design)

- **Master project pattern** (`#1036`): user wraps company repos
  under `repositorios/` with `openspec/` at the master level. SDD
  artifacts already live in the right place.
- **Two-repo invariant** (`#1031`, `#1035`): every table = infra repo
  + carga repo, always both, no exceptions. This is structural.
- **Verify dual-path resolution** (`#1043`): Camino A (TDD in
  `sdd-apply` with Glue Docker + throwaway table) + Camino B
  (deploy verify in `sdd-verify` with SAM deploy both repos + dev vs
  prd).
- **Authorship rule** (`#1045`): human name carries legal
  responsibility; AI attribution is forbidden (already mirrors the
  gentle-ai persona rule).
- **ETL header protocol** (`#1044`): 6-field mandatory header on
  every `.py` — Glue job name, desarrollado por, fecha creación,
  modificado por, fecha modificación, descripción.
- **AWS profile system** (`#1033`): prd / dev / usuario with account
  IDs `874970050509 / 895593169121 / 516363283643`.
- **Floci decision** (`#1038`): investigated, **descarted**. User
  has dev AWS access; no need for local AWS emulation.
- **Glue Docker** (`#1039`): `public.ecr.aws/glue/aws-glue-libs:5`
  is the chosen local Spark runtime for Camino A.
- **Data-driven integration** (`#1051`): DDT / data-driven dev /
  DataOps / data-contracts are existing-practice formalization, not
  new adoption.

### The 22 captured observations break into 5 themes

| Theme | Observations | What it locks in |
|-------|--------------|------------------|
| **Two-repo invariant** | #1031, #1035, #1036, #1042 | Multi-repo coordination is structural |
| **Study → Infra → Load** | #1032, #1034 | sdd-explore = source study, not codebase exploration |
| **ETL patterns (4)** | #1040, #1048, #1049, #1050 | Pattern detection guides design + apply + verify |
| **Verify paradigm** | #1037, #1041, #1043, #1051 | Data-comparison + dual-path, not code assertions |
| **Company protocol** | #1044, #1045, #1046, #1047 | Authorship + header + additive + insertion-point rules |

### Affected Areas

- `openspec/config.yaml` — needs `domain: data-engineering` field,
  multi-repo declarations, AWS profile mapping, verify command
  override
- `~/.config/opencode/skills/sdd-init` — needs domain detection +
  preflight question
- `~/.config/opencode/skills/sdd-explore` — needs data-engineering
  exploration branch (source study, not codebase exploration)
- `~/.config/opencode/skills/sdd-spec` — needs data-engineering
  spec template (table schema, partitions, watermark, DAG)
- `~/.config/opencode/skills/sdd-design` — needs DAG-of-transformations
  template + insertion-point analysis for modifications
- `~/.config/opencode/skills/sdd-tasks` — needs `repo: infra|carga`
  prefix per task
- `~/.config/opencode/skills/sdd-apply` — needs Camino A branch
  (Glue Docker TDD) + header-protocol enforcement
- `~/.config/opencode/skills/sdd-verify` — needs Camino B branch
  (SAM deploy + Athena compare) + data-driven testing paradigm
- `~/.config/opencode/skills/data-engineer-*` (7 skills) —
  need SDD-phase awareness, 4-pattern coverage, master-project
  awareness, header-protocol + authorship-rule baked in
- New: `~/.config/opencode/skills/data-engineer-pattern-detect`
  (or merge into `integrate`) — heuristic scan of carga repo for
  pattern markers
- `openspec/specs/sdd-profiles/` — existing profile spec is
  app-dev only; needs domain-profile extension

---

## The 6 Architectural Questions

For each question I surface **concrete options** with the trade-off
called out. The recommendation is in the next section.

### Q1. Domain profile detection mechanism

How does SDD know it's in data-engineering mode vs app-dev mode?

| Option | Description | Pros | Cons |
|--------|-------------|------|------|
| **A. Session preflight question** | `sdd-init` asks "domain?" as a new preflight field alongside delivery strategy, persistence mode, review budget | Explicit, predictable, easy to override | One more preflight question; easy to skip |
| **B. Config-file driven** | `openspec/config.yaml` gets `domain: data-engineering`; `sdd-init` writes it; subsequent phases read it | Single source of truth, project-portable, survives sessions | Requires writing before phases; bootstrap question |
| **C. Auto-detection** | Scan repo for Glue jobs, SAM templates, `.py` ETL files; infer domain | Zero user friction | False positives (any Python repo has ETL-like scripts), magic = brittle |
| **D. Hybrid B+C** | Auto-detect as hint (writes `domain: auto-detected-as-X`), preflight confirms | Best UX: user can confirm or correct | More code in `sdd-init` |

**Recommended**: **D (hybrid)**. Auto-detect proposes, user confirms,
config persists. The preflight question becomes "Domain detected as
data-engineering (based on `template.yaml` + `glue-jobs/*.py`
matches). Confirm or override?".

### Q2. Skills adaptation strategy

How do the 7 existing `data-engineer-*` skills get updated, and do
new skills get added?

| Option | Description | Pros | Cons |
|--------|-------------|------|------|
| **A. Patch existing 7 in-place** | Update each `SKILL.md` to know its SDD phase, the 4 patterns, the header protocol, and the multi-repo structure | Smallest surface area; skills stay self-contained | Each skill becomes large; some changes duplicated across skills |
| **B. Profile SKILL umbrella** | One new `data-engineer-profile/SKILL.md` that holds the shared concepts (4 patterns, header protocol, AWS profiles); the 7 skills reference it via `## Reference: data-engineer-profile` | DRY; one place to evolve shared concepts | Cross-skill references are easy to break; loses locality |
| **C. Patch + new orchestrator skill** | Patch the 7 + add `data-engineer-pattern-detect` (heuristic scan) + add `data-engineer-verify` (verify orchestration) | Skills stay focused; orchestration explicit | More skills to maintain |

**Recommended**: **C**. The 4-pattern taxonomy + header protocol +
authorship rule are shared knowledge that EVERY skill needs. Patch
the 7 with a shared "protocol block" at the top (header +
authorship + 4-pattern index), and add `data-engineer-pattern-detect`
as a new skill that the orchestrator can call before
`sdd-design`.

### Q3. Verify paradigm

How does `sdd-verify` know it's verifying a data pipeline, not a
Go module?

| Option | Description | Pros | Cons |
|--------|-------------|------|------|
| **A. Domain-agnostic core, domain skill** | `sdd-verify` stays generic; when `domain: data-engineering`, it delegates to `data-engineer-verify` (new skill) | Core stays simple; each domain owns its verify | More skills; cross-domain consistency harder |
| **B. Mode branches inside `sdd-verify`** | `sdd-verify` reads `domain` from config, branches to data-engineering code path | One skill, one entry point | `sdd-verify` becomes large; app-dev changes risk breaking data-eng path |
| **C. Replace verify for data engineering** | When `domain: data-engineering`, the orchestrator routes verification to `data-engineer-verify` directly, bypassing `sdd-verify` | Cleanest separation; both skills stay focused | Two verify skills means two sets of conventions; risk of drift |

**Recommended**: **B (mode branches)**. `sdd-verify` is already
domain-aware at config level (it reads `testing` section). Adding a
`data_engineering` section to config and branching inside
`sdd-verify` is the smallest change with the least skill sprawl.

### Q4. Spec template

What does an ETL change spec look like in OpenSpec?

| Option | Description | Pros | Cons |
|--------|-------------|------|------|
| **A. Reuse OpenSpec delta, extend with sections** | Keep `ADDED / MODIFIED / REMOVED / RENAMED`; add optional sections like `## Source Tables`, `## Target Schema`, `## Watermark Strategy`, `## DAG` | Backwards compatible; one format | Spec gets long for table-heavy changes |
| **B. New domain-specific spec type** | `openspec/changes/{name}/specs/glue-tables/{db}/{table}/spec.md`; schema-lifecycle is the unit | Matches the deliverable atomic (table) | Doesn't compose with multi-table changes well |
| **C. Hybrid — OpenSpec delta + sidecar YAML** | Delta spec for behavior; `glue-tables/{db}.{table}.yaml` sidecar for schema/partition/S3-path; both reviewed together | Clear separation: prose vs machine-readable | Two files to keep in sync |

**Recommended**: **C (hybrid)**. The behavioral requirements (Given/
When/Then for the data comparison) live in the OpenSpec delta. The
schema/partition/S3-path/compression format live in a YAML sidecar
that matches `glue-tables/*.yaml` in the infra repo. `sdd-verify`
runs the behavioral matrix AND validates the sidecar against the
infra repo's deployed Glue table.

### Q5. Multi-repo coordination

How does SDD know which company repo a task belongs to?

| Option | Description | Pros | Cons |
|--------|-------------|------|------|
| **A. Master project pattern only** | User already does this; SDD just respects the master project layout | No code change | Tasks have no machine-readable repo hint; cross-repo refactors are error-prone |
| **B. Config declares repos** | `openspec/config.yaml` lists `repos: { infra: ./repositorios/infra-..., carga: ./repositorios/carga-... }` | Discoverable; tools can run commands against the right path | Path drift if user moves repos |
| **C. Tasks get `repo:` prefix** | `tasks.md` entries become `- [ ] 1.1 [infra] Update glue-tables/<table>.yaml`; `- [ ] 2.1 [carga] Update glue-jobs/<job>.py` | Machine-readable; cross-repo visibility at the task level | Requires tasks authoring discipline |

**Recommended**: **B + C together**. Config declares repos (used by
verify commands to find the right path); tasks annotate which repo
each task touches. The `sdd-tasks` skill emits the `repo:` prefix
automatically based on which subdomain the task targets (infra-only
tasks get `[infra]`, carga-only get `[carga]`, two-repo tasks get
`[both]`).

### Q6. ETL pattern detection

How does SDD know which of the 4 ETL patterns a given change is?

| Option | Description | Pros | Cons |
|--------|-------------|------|------|
| **A. New skill `data-engineer-pattern-detect`** | Standalone skill: scan carga repo for markers (wr.athena.read_sql_query → legacy; ApplyMapping + sparkSqlQuery → Glue Studio; temp-view count + GlueContext → multi-step; boto3+awswrangler+watermark → incremental); report detected pattern | Decoupled; reusable in any phase | Yet another skill; one-shot scanner |
| **B. Enhancement to `data-engineer-integrate`** | Add pattern-detection logic to the integrator skill | One orchestrator skill owns the pattern knowledge | Integrator becomes a god-skill |
| **C. User declares in proposal** | `sdd-propose` asks "which pattern?" and user picks | Zero code, zero heuristics | Friction; users often don't know their pattern by name |

**Recommended**: **A + C combined**. New skill does auto-detection
(markers + structure scoring) and reports a confidence. `sdd-propose`
uses that signal as the default but lets the user override. For
brand-new ETLs (no existing code), `sdd-propose` asks directly using
the 4-pattern menu.

---

## Approaches — Summary Matrix

For each question, the recommended option and its main risk:

| Q | Recommended | Main risk |
|---|-------------|-----------|
| 1. Detection | D (hybrid auto + preflight confirm) | Auto-detection false positive on hybrid Python repos |
| 2. Skills | C (patch + new pattern-detect) | Skill maintenance overhead |
| 3. Verify | B (mode branch in sdd-verify) | Cross-domain drift in single skill |
| 4. Spec | C (delta + sidecar YAML) | Two files to keep in sync |
| 5. Multi-repo | B+C (config + task prefix) | Repo path drift |
| 6. Pattern detect | A+C (new skill + user override) | Heuristic misclassification |

### Why these recommendations form a coherent set

- **Detection (D)** feeds **verify (B)** and **spec (C)**: once
  `domain: data-engineering` is in config, every phase branches
  without re-asking.
- **Skills (C)** is the enabler for **spec (C)**: the YAML sidecar
  format mirrors what `data-engineer-create-table` already
  generates, so the sidecar format reuses existing knowledge.
- **Multi-repo (B+C)** is the only way **verify (B)** can run
  `sam deploy --config-file infra/samconfig.toml` against the right
  path.
- **Pattern detection (A+C)** is what makes **apply (Camino A)**
  smart: the TDD-loop pattern (incremental watermark test vs
  multi-step DAG test vs Athena-vs-Spark parity test) differs by
  pattern.

---

## Recommendation

Adopt all 6 recommendations as a single coherent domain profile,
implemented in this order:

1. **Phase 1 — Profile detection + config schema**
   - Add `domain`, `repos`, `aws_profiles` fields to
     `openspec/config.yaml`
   - `sdd-init` preflight adds the domain question
   - Auto-detect hints based on `template.yaml` + `glue-jobs/*.py`
     presence
2. **Phase 2 — Spec + design templates**
   - New OpenSpec delta sections: `## Source Tables`,
     `## Target Schema`, `## Watermark Strategy`, `## DAG`,
     `## AWS Profile Requirements`, `## Verify Approach`
   - Sidecar `glue-tables/*.yaml` lives next to delta spec
   - `sdd-spec` and `sdd-design` gain data-engineering templates
3. **Phase 3 — Apply + Verify dual-path**
   - `sdd-apply` gains Camino A branch: Glue Docker scaffold,
     throwaway-table test, header-protocol enforcement,
     authorship-rule
   - `sdd-verify` gains Camino B branch: SAM deploy both repos,
     job run, Athena dev-vs-prd comparison via SQL EXCEPT
4. **Phase 4 — Skills update**
   - Patch all 7 `data-engineer-*` skills with shared protocol
     block (header + authorship + 4-pattern index + master-project
     awareness + AWS profile awareness)
   - Add new `data-engineer-pattern-detect` skill
5. **Phase 5 — Tasks multi-repo**
   - `sdd-tasks` emits `repo: infra | carga | both` prefix
   - `sdd-tasks` reads `repos:` from config and validates path
     exists

This is a 5-phase change that touches every SDD phase. The review
budget is 800 lines; recommend chained PRs (5 PRs, one per
phase) — each PR is a vertical slice that can land independently
behind a feature flag if needed.

---

## Risks

- **Bifurcation drift**: once we split `sdd-verify` into two
  branches, app-dev changes risk breaking data-eng and vice versa.
  *Mitigation*: shared test fixtures in
  `~/.config/opencode/skills/sdd-verify/testdata/` that exercise
  both branches; CI runs both.
- **Heuristic pattern misclassification**: if
  `data-engineer-pattern-detect` picks the wrong pattern, the TDD
  scaffold will be wrong (e.g., watermark test where DAG test is
  needed). *Mitigation*: user override in `sdd-propose`; report
  confidence score; never silently misclassify.
- **YAML sidecar drift**: the sidecar `glue-tables/*.yaml` and the
  delta spec sections could diverge if the user edits one and not
  the other. *Mitigation*: `sdd-verify` validates sidecar against
  deployed Glue table schema via `aws glue get-table`.
- **Repo path fragility**: if the user moves `repositorios/`, the
  config `repos:` block becomes stale and verify commands fail.
  *Mitigation*: `sdd-init` re-validates paths at the start of every
  session; warn if path missing.
- **SAM deploy cost**: each `sdd-verify` run does TWO `sam
  deploy` calls. If `sam deploy` is slow (>5 min each), the verify
  cycle becomes painful. *Mitigation*: parallel SAM deploys via
  background jobs; cache stack outputs; allow `verify:
  skip-deploy` for local-only verify when not at the deploy-check
  gate.
- **Multi-repo PRs**: the user has Bitbucket, not GitHub, for the
  two company repos. PR conventions differ. *Mitigation*: the
  existing `branch-pr` skill is GitHub-flavored; data-eng changes
  touch Bitbucket — verify skill needs to know this. The SDD
  artifact changes live in the master project (this repo); the
  Bitbucket changes are PR'd separately.
- **Authorship rule enforcement**: if a sub-agent forgets the rule
  and inserts AI attribution, the company header is invalid.
  *Mitigation*: code generation in `sdd-apply` uses a shared
  header template with `desarrollado_por` and `modificado_por`
  pulled from session config (not hardcoded to AI model).
- **AWS profile leak in logs**: if `sdd-verify` accidentally
  logs `--profile prd` commands, sensitive ops show up in audit
  trails. *Mitigation*: scrub profile names from log output;
  document which commands are dev-only.

---

## Ready for Proposal

**Yes — but with 4 user-confirmation gates before `sdd-propose`
runs:**

The orchestrator should ask the user to confirm or override each of
these decisions. They are stated as recommendations but each has a
real trade-off the user owns:

1. **Domain detection**: hybrid (auto + preflight confirm)?
   Alternative: pure preflight question, no auto-detect.
2. **Verify approach**: mode-branch inside `sdd-verify`?
   Alternative: separate `data-engineer-verify` skill that
   `sdd-verify` delegates to.
3. **Spec format**: OpenSpec delta + YAML sidecar?
   Alternative: pure OpenSpec delta, schema inside a fenced
   block; or pure YAML sidecar, prose outside.
4. **Pattern detection**: new `data-engineer-pattern-detect`
   skill? Alternative: heuristic lives inside
   `data-engineer-integrate`, no new skill.

After these 4 are locked, the orchestrator runs `sdd-propose` with
the change name `data-engineering-domain-profile`. The proposal
should:

- State the 5-phase delivery (one phase per PR slice)
- Explicitly declare the 800-line budget per PR
- Reference this exploration as the source of decisions
- List the 4 user-confirmed options in the proposal body
- Recommend `delivery_strategy: chained` with `chain-strategy:
  stacked-to-main`

The proposal artifact should be written to
`openspec/changes/data-engineering-domain-profile/proposal.md`.
