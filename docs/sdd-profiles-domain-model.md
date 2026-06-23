# SDD Profiles: Domain × Model

← [Back to README](../README.md)

---

## Two Orthogonal Axes

SDD has **two independent profile axes** that combine to determine how each phase executes:

| Axis | Question | Existing? | Config Location |
|------|----------|-----------|-----------------|
| **Model Profile** | Which AI model runs each phase? | ✅ Yes (`sdd-profiles`) | `opencode.json` agent overlays, TUI model picker |
| **Domain Profile** | Which domain rules/skills/templates apply? | ✅ New (`data-engineering-domain`) | `openspec/config.yaml` `domain:` field |

They are independent — a project can be `domain: data-engineering` with a `cheap` model profile, or `domain: app-dev` with a `premium` profile. Both axes are resolved at session start and applied to every phase.

---

## Architecture Diagram

```mermaid
flowchart TB
    subgraph SESSION["Session Start"]
        PREFLIGHT["SDD Session Preflight"]
    end

    subgraph AXIS1["Axis 1 — Model Profile (existing)"]
        direction LR
        MP_CHEAP["cheap<br/>Haiku/Mini → explore, archive<br/>Opus/Sonnet → design, apply"]
        MP_MID["mid<br/>Sonnet → all phases"]
        MP_PREMIUM["premium<br/>Opus → all phases"]
    end

    subgraph AXIS2["Axis 2 — Domain Profile (new)"]
        direction LR
        DP_APP["app-dev<br/>go-testing, branch-pr<br/>go test ./...<br/>software specs"]
        DP_DATA["data-engineering<br/>study-file, etl-s3, create-table<br/>Glue Docker + SAM deploy<br/>ETL specs + sidecar YAML"]
    end

    subgraph CORE["SDD 8-Phase Core (unchanged)"]
        direction LR
        P1[explore] --> P2[propose] --> P3[spec] --> P4[design]
        P4 --> P5[tasks] --> P6[apply] --> P7[verify] --> P8[archive]
    end

    PREFLIGHT --> AXIS1
    PREFLIGHT --> AXIS2
    AXIS1 -->|"which model per phase"| CORE
    AXIS2 -->|"which skills/templates/verify per phase"| CORE

    classDef modelProfile fill:#a78bfa,stroke:#6d28d9,stroke-width:2px,color:#000
    classDef domainProfile fill:#60a5fa,stroke:#1d4ed8,stroke-width:2px,color:#000
    classDef phase fill:#f9a8d4,stroke:#be185d,stroke-width:2px,color:#000
    classDef session fill:#fbbf24,stroke:#b45309,stroke-width:2px,color:#000

    class MP_CHEAP,MP_MID,MP_PREMIUM modelProfile
    class DP_APP,DP_DATA domainProfile
    class P1,P2,P3,P4,P5,P6,P7,P8 phase
    class PREFLIGHT session
```

---

## Domain Profile Comparison

```mermaid
flowchart LR
    subgraph APP_DEV["App-Dev Profile (default)"]
        direction TB
        A_EXPLORE["explore<br/>Read codebase<br/>Skill: sdd-explore<br/>Model: per profile"]
        A_SPEC["spec<br/>Software requirements<br/>Given/When/Then (code)<br/>Skill: sdd-spec"]
        A_DESIGN["design<br/>Architecture patterns<br/>Skill: sdd-design"]
        A_TASKS["tasks<br/>Implementation tasks<br/>Skill: sdd-tasks"]
        A_APPLY["apply<br/>Write Go/code<br/>TDD: go test<br/>Skill: sdd-apply + go-testing"]
        A_VERIFY["verify<br/>go test ./...<br/>go vet, go build<br/>Skill: sdd-verify"]

        A_EXPLORE --> A_SPEC --> A_DESIGN --> A_TASKS --> A_APPLY --> A_VERIFY
    end

    subgraph DATA_ENG["Data-Engineering Profile (new)"]
        direction TB
        D_EXPLORE["explore<br/>Study source data (CSV, tables)<br/>Skill: study-file + sdd-explore<br/>ETL delta sections"]
        D_SPEC["spec<br/>ETL requirements<br/>Source → Transform → Dest<br/>Sidecar YAML: schema, partitions<br/>Watermark strategy"]
        D_DESIGN["design<br/>DAG of transformations<br/>Insertion-point analysis<br/>Pattern-aware (4 patterns)"]
        D_TASKS["tasks<br/>[infra] / [carga] / [both] prefix<br/>Git flow: Bitbucket feature→develop→release"]
        D_APPLY["apply (Camino A)<br/>Glue Docker TDD loop<br/>Header protocol enforcement<br/>Authorship: human name"]
        D_VERIFY["verify (Camino B)<br/>SAM deploy both repos<br/>Athena EXCEPT dev vs prd<br/>Sidecar validation<br/>Profile scrubbing"]
        D_ARCHIVE["archive<br/>Same as app-dev"]

        D_EXPLORE --> D_SPEC --> D_DESIGN --> D_TASKS --> D_APPLY --> D_VERIFY --> D_ARCHIVE
    end

    classDef appPhase fill:#86efac,stroke:#16a34a,stroke-width:2px,color:#000
    classDef dataPhase fill:#fca5a5,stroke:#dc2626,stroke-width:2px,color:#000

    class A_EXPLORE,A_SPEC,A_DESIGN,A_TASKS,A_APPLY,A_VERIFY appPhase
    class D_EXPLORE,D_SPEC,D_DESIGN,D_TASKS,D_APPLY,D_VERIFY,D_ARCHIVE dataPhase
```

---

## Per-Agent Model Assignment

Each SDD phase runs as a sub-agent. The model assigned to each sub-agent can vary by domain — data-engineering phases may benefit from different model tiers than app-dev phases.

```mermaid
flowchart TB
    subgraph MODEL_ASSIGNMENT["Model Assignment per Agent"]
        direction TB

        subgraph APP_MODELS["App-Dev Recommended Models"]
            direction LR
            AM_EXPLORE["sdd-explore<br/>→ cheap (Haiku)<br/>Codebase scan is mechanical"]
            AM_PROPOSE["sdd-propose<br/>→ mid (Sonnet)<br/>Product decisions"]
            AM_SPEC["sdd-spec<br/>→ mid (Sonnet)<br/>Requirement writing"]
            AM_DESIGN["sdd-design<br/>→ premium (Opus)<br/>Architecture reasoning"]
            AM_TASKS["sdd-tasks<br/>→ mid (Sonnet)<br/>Decomposition"]
            AM_APPLY["sdd-apply<br/>→ mid (Sonnet)<br/>Code generation"]
            AM_VERIFY["sdd-verify<br/>→ mid (Sonnet)<br/>Test execution"]
            AM_ARCHIVE["sdd-archive<br/>→ cheap (Haiku)<br/>File moves"]
        end

        subgraph DATA_MODELS["Data-Engineering Recommended Models"]
            direction LR
            DM_EXPLORE["sdd-explore + study-file<br/>→ mid (Sonnet)<br/>Data profiling needs judgment<br/>(encoding, dates, types)"]
            DM_PROPOSE["sdd-propose<br/>→ mid (Sonnet)<br/>ETL scope decisions"]
            DM_SPEC["sdd-spec<br/>→ premium (Opus)<br/>Schema design + DAG planning<br/>Sidecar YAML needs precision"]
            DM_DESIGN["sdd-design<br/>→ premium (Opus)<br/>Insertion-point cascade analysis<br/>4-pattern detection"]
            DM_TASKS["sdd-tasks<br/>→ mid (Sonnet)<br/>Multi-repo task split"]
            DM_APPLY["sdd-apply (Camino A)<br/>→ mid (Sonnet)<br/>Spark SQL + PySpark gen<br/>Header enforcement"]
            DM_VERIFY["sdd-verify (Camino B)<br/>→ mid (Sonnet)<br/>SAM deploy + Athena EXCEPT<br/>Sidecar validation"]
            DM_ARCHIVE["sdd-archive<br/>→ cheap (Haiku)<br/>File moves"]
        end
    end

    classDef appModel fill:#86efac,stroke:#16a34a,stroke-width:1px,color:#000
    classDef dataModel fill:#fca5a5,stroke:#dc2626,stroke-width:1px,color:#000

    class AM_EXPLORE,AM_PROPOSE,AM_SPEC,AM_DESIGN,AM_TASKS,AM_APPLY,AM_VERIFY,AM_ARCHIVE appModel
    class DM_EXPLORE,DM_PROPOSE,DM_SPEC,DM_DESIGN,DM_TASKS,DM_APPLY,DM_VERIFY,DM_ARCHIVE dataModel
```

### Why data-engineering needs different model tiers

| Phase | App-Dev (why cheap is OK) | Data-Eng (why it changes) |
|-------|---------------------------|---------------------------|
| **explore** | Codebase scan = mechanical | Data profiling = judgment (encoding, DD/MM vs MM/DD, types) |
| **spec** | Requirements = straightforward | Schema + DAG + sidecar = precision critical |
| **design** | Architecture patterns = well-known | Insertion-point cascade = complex dependency analysis |
| **apply** | Code gen from tests = clear | Spark SQL translation (Presto→Spark) + header protocol = nuanced |
| **verify** | `go test` = binary pass/fail | Athena EXCEPT + sidecar validation = interpretation needed |

---

## How Model + Domain Combine at Runtime

```mermaid
sequenceDiagram
    participant User
    participant Orchestrator
    participant SDDConfig as gentle-ai sdd-config
    participant SubAgent as Phase Sub-Agent

    User->>Orchestrator: /sdd-explore (or any phase)
    Orchestrator->>SDDConfig: sdd-config --json

    SDDConfig-->>Orchestrator: {domain: "data-engineering", repos: {...}}

    Note over Orchestrator: Resolve BOTH axes:
    Note over Orchestrator: 1. Domain → data-engineering skills/templates
    Note over Orchestrator: 2. Model profile → which model for this phase

    Orchestrator->>SubAgent: Launch with:
    Note right of SubAgent: • Model: per profile (e.g. Sonnet for explore)
    Note right of SubAgent: • Skills: study-file + sdd-explore (data-eng)
    Note right of SubAgent: • Templates: ETL delta sections
    Note right of SubAgent: • Verify approach: Camino A (TDD) or B (deploy)
    Note right of SubAgent: • Header protocol + authorship rule

    SubAgent-->>Orchestrator: Result contract
```

---

## Config: How Both Profiles Are Declared

### Model Profile (existing — in opencode.json)

```json
{
  "agent": {
    "sdd-orchestrator": { "model": "anthropic/claude-sonnet-4-20250514" },
    "sdd-explore": { "model": "anthropic/claude-haiku-4-5-20250315" },
    "sdd-design": { "model": "anthropic/claude-opus-4-20250514" },
    "sdd-apply": { "model": "anthropic/claude-sonnet-4-20250514" }
  }
}
```

### Domain Profile (new — in openspec/config.yaml)

```yaml
domain: data-engineering
repos:
  infra: ./repositorios/infra-datos-trs-posventa
  carga: ./repositorios/carga-datos-trs-posventa
aws_profiles:
  prd: AWSReadFullDat-874970050509
  dev: aws-tcl-ope-set-cloud-895593169121
  usuario: aws-tcl-ope-set-devdat-516363283643
verify:
  skip_deploy: false
```

### Combined resolution

```
Phase: sdd-explore
  Model:  opencode.json → claude-haiku-4-5 (from model profile)
  Domain: config.yaml → data-engineering (from domain profile)
  Result: Haiku runs sdd-explore WITH study-file skill + ETL delta sections
```

---

## See Also

- [SDD Ecosystem](sdd-ecosystem.md) — full ecosystem diagram with skills, MCP, Engram
- [OpenSpec Config](openspec-config.md) — domain profile config fields
- [OpenCode Profiles](opencode-profiles.md) — model profile configuration
- [Skill Registry](skill-registry.md) — how skills are resolved per domain
