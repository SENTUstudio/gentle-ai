# Delta for ETL Spec Format

## ADDED Requirements

### Requirement: ETL Delta Sections

An ETL delta spec MUST include `## Source Tables`, `## Target Schema`, `## Watermark Strategy`, `## DAG`, `## AWS Profile Requirements`, and `## Verify Approach` when `domain: data-engineering`. Behavioral requirements MUST use Given/When/Then focused on data comparisons, not code assertions.

#### Scenario: New table spec

- GIVEN a new table ETL with sources and target schema
- WHEN `sdd-spec` writes the delta
- THEN it includes all six ETL sections
- AND scenarios assert table output, not function output

#### Scenario: Modification spec

- GIVEN an additive modification to an existing ETL
- WHEN `sdd-spec` writes the delta
- THEN it documents the insertion point and regression expectations

### Requirement: Sidecar YAML Format

Each target Glue table MUST have a sidecar file at `glue-tables/{db}.{table}.yaml` adjacent to the delta spec. The sidecar MUST declare `database`, `table`, `columns` (name, type, comment), `partitions`, `s3_location`, `format`, and `compression`. `sdd-verify` MUST validate the sidecar against the deployed Glue table via `aws glue get-table`.

#### Scenario: Sidecar schema

- GIVEN a target table `db_dl_dev_stg_encuestas.encuestas_csi`
- WHEN `sdd-spec` writes the sidecar
- THEN it creates `glue-tables/db_dl_dev_stg_encuestas.encuestas_csi.yaml`
- AND the file declares schema and partitions

#### Scenario: Sidecar validation

- GIVEN a sidecar and a deployed Glue table
- WHEN `sdd-verify` runs
- THEN it validates the sidecar columns, partitions, and location
- AND it reports mismatches

### Requirement: Pattern Template Selection

`sdd-spec` MUST select one of four ETL pattern templates based on `data-engineer-pattern-detect` output or user override: incremental watermark, multi-step Spark SQL, legacy wrangler, or Glue Studio visual. The spec MUST name the selected pattern and the confidence score when auto-detected.

#### Scenario: Incremental pattern

- GIVEN a single-source S3 CSV with watermark logic
- WHEN the pattern is detected
- THEN the spec uses the incremental watermark template

#### Scenario: Multi-step pattern

- GIVEN multiple source tables and temp-view stages
- WHEN the pattern is detected
- THEN the spec uses the multi-step Spark SQL template

#### Scenario: User override

- GIVEN `data-engineer-pattern-detect` reports low confidence
- WHEN the user selects a pattern in `sdd-propose`
- THEN the spec uses the user-selected pattern

### Requirement: DAG Representation

The `## DAG` section MUST describe transformations as a directed acyclic graph: nodes (sources, temp views, target) and edges (dependencies). For modifications, it MUST highlight the optimal insertion point and list downstream stages affected by the change.

#### Scenario: New ETL DAG

- GIVEN a multi-step ETL with three SQL stages
- WHEN the DAG section is written
- THEN it lists each stage and its input dependencies

#### Scenario: Modification insertion point

- GIVEN an additive modification requiring a new stage
- WHEN the DAG section is written
- THEN it marks the insertion point and affected downstream stages

### Requirement: Watermark and Whitelist Documentation

Incremental ETL specs MUST document the watermark query, fallback value, file-timestamp extraction rule, and whitelist filter. The documentation MUST explain that idempotency comes from the whitelist (only new files), not from the write mode.

#### Scenario: Watermark spec

- GIVEN an incremental load keyed on `fecha_carga_proveedor`
- WHEN the watermark strategy is written
- THEN it states `MAX(fecha_carga_proveedor)` with `2000-01-01` fallback
- AND it describes the S3 file filter by filename timestamp

### Requirement: Verify Approach Documentation

The `## Verify Approach` section MUST document Camino A (`sdd-apply`: Glue Docker `aws-glue-libs:5`, throwaway table, header/authorship checks) and Camino B (`sdd-verify`: SAM deploy both repos, run job, Athena dev-vs-prd `EXCEPT` comparison).

#### Scenario: Dual-path verify

- GIVEN an ETL change with local and deployed verify steps
- WHEN the verify approach is written
- THEN it lists Camino A commands and Camino B SQL comparisons

#### Scenario: Skip-deploy option

- GIVEN `verify.skip_deploy: true` in config
- WHEN `sdd-verify` runs
- THEN it validates sidecar and SQL logic without SAM deploy
