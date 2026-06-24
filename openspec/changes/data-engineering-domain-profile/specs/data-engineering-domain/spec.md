# Delta for Data Engineering Domain Profile

## ADDED Requirements

### Requirement: Hybrid Domain Detection

The system MUST detect a data-engineering project by scanning for `template.yaml` and `glue-jobs/*.py`. The system MUST present the result as a preflight hint for confirmation or override. The confirmed domain MUST persist in `openspec/config.yaml` as `domain: data-engineering`.

#### Scenario: Auto-detect with confirmation

- GIVEN a master project containing `template.yaml` and `glue-jobs/ETL_foo.py`
- WHEN `sdd-init` runs preflight
- THEN it proposes `domain: data-engineering`
- AND it writes the confirmed value to `config.yaml`

#### Scenario: Override false-positive detection

- GIVEN a Python repo with Glue-like files
- WHEN the user rejects the detected hint
- THEN `sdd-init` leaves `domain` unset

### Requirement: Domain Config Schema

`openspec/config.yaml` MUST support `domain`, `repos`, and `aws_profiles` fields. The `repos` mapping MUST declare `infra` and `carga` paths. The `aws_profiles` mapping MUST declare `prd`, `dev`, and `usuario`.

#### Scenario: Validate repo paths

- GIVEN `config.yaml` declares `repos.infra` and `repos.carga`
- WHEN `sdd-init` re-validates paths
- THEN it warns if a declared path does not exist

#### Scenario: Resolve AWS profile

- GIVEN `aws_profiles.dev: aws-tcl-ope-set-cloud-895593169121`
- WHEN a skill runs an AWS CLI command for dev
- THEN it resolves the logical name `dev` to the configured CLI profile name

### Requirement: Verify Mode Branch

`sdd-verify` MUST branch to the data-engineering verify path when `domain: data-engineering` is set. The app-dev verify path MUST remain unchanged when no domain is set.

#### Scenario: Data-engineering verify branch

- GIVEN `domain: data-engineering` and a delta spec with sidecars
- WHEN `sdd-verify` runs
- THEN it executes Camino B (SAM deploy both repos + Athena dev-vs-prd comparison)

#### Scenario: App-dev unchanged

- GIVEN a project without `domain` set
- WHEN `sdd-verify` runs
- THEN it executes `go test ./...` and `go vet ./...`

### Requirement: Multi-Repo Coordination and Git Flow

`sdd-tasks` MUST emit a `repo:` prefix on every task using values `infra`, `carga`, or `both`. The system MUST distinguish the master project's GitHub flow from company repos' Bitbucket flow (`feature/ → develop → release/`). The `branch-pr` skill MUST apply only to the master project.

#### Scenario: Repo prefix and flow annotation

- GIVEN a task updates `glue-jobs/ETL_foo.py`
- WHEN `sdd-tasks` writes the task
- THEN it prefixes the task with `[carga]`
- AND it notes the Bitbucket flow and that `branch-pr` does not apply

#### Scenario: Cross-repo task prefix

- GIVEN a task creates a Glue table and its ETL
- WHEN `sdd-tasks` writes the task
- THEN it prefixes the task with `[both]`

### Requirement: AWS Profile Scrubbing

Data-engineering skills MUST reference AWS profiles by logical name (`prd`, `dev`, `usuario`) and resolve them through `aws_profiles` config. Profile names and account IDs MUST be scrubbed from public logs, verify reports, and PR descriptions.

#### Scenario: Hidden profile in report

- GIVEN `aws_profiles.prd: AWSReadFullDat-874970050509`
- WHEN `sdd-verify` writes its report
- THEN it does not echo the profile name

### Requirement: ETL Header Protocol and Authorship

Every generated or modified `.py` ETL file MUST contain a header with Glue job name, `desarrollado por`, `fecha creación`, `modificado por`, `fecha modificación`, and `descripción`. For modifications, `desarrollado por` and `fecha creación` MUST be preserved, while `modificado por` and `fecha modificación` MUST update to the current human author and date. AI-generated or modified code MUST attribute authorship to the directing human and MUST NOT include AI names, `Co-Authored-By`, or generated-by attribution.

#### Scenario: New ETL header

- GIVEN `sdd-apply` generates a new Glue job
- WHEN it writes `glue-jobs/ETL_foo.py`
- THEN the header contains all six fields with the human's name
- AND `desarrollado por` equals `modificado por`

#### Scenario: Modify existing ETL with no AI attribution

- GIVEN `sdd-apply` modifies an existing job
- WHEN it writes the file and prepares a commit
- THEN it preserves `desarrollado por` and `fecha creación`, updates `modificado por` and `fecha modificación`
- AND the commit contains no AI attribution
