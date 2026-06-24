# SDD Ecosystem

← [Back to README](../README.md)

---

## Overview

Spec-Driven Development (SDD) is the workflow layer at the center of Gentle-AI. It does not operate in isolation. SDD orchestrates phases (explore → propose → spec → design → tasks → apply → verify → archive), and each phase leverages three supporting pillars:

| Pillar | Role | What it provides to SDD |
|--------|------|-------------------------|
| **Artifact Store** | Persistence | Where specs, designs, tasks, and reports live — Engram (memory), OpenSpec (files), or both (hybrid) |
| **Skills** | Execution patterns | Curated `SKILL.md` files loaded by sub-agents before each phase — testing, PR creation, doc design, work-unit commits, etc. |
| **MCP Servers** | External capabilities | Context7 (live documentation lookup), Engram MCP (memory access from any agent), and project-specific servers (Notion, Jira, etc.) |

SDD is the conductor. The pillars are the instruments. The orchestrator agent coordinates; sub-agents execute with the right skills, memory, and tools for each phase.

---

## Ecosystem Diagram

```mermaid
flowchart TB
    subgraph ECOSYSTEM["Gentle-AI Ecosystem"]
        subgraph SDD["SDD — Workflow Layer"]
            direction LR
            EXPLORE[explore] --> PROPOSE[propose]
            PROPOSE --> SPEC[spec]
            SPEC --> DESIGN[design]
            DESIGN --> TASKS[tasks]
            TASKS --> APPLY[apply]
            APPLY --> VERIFY[verify]
            VERIFY --> ARCHIVE[archive]
        end

        subgraph ARTIFACTS["Artifact Store — Persistence"]
            direction LR
            ENGRAM["Engram<br/>(memory)"]
            OPENSPEC["OpenSpec<br/>(files)"]
            HYBRID["Hybrid<br/>(both)"]
        end

        subgraph SKILLS["Skills — Execution Patterns"]
            SDD_SKILLS["SDD Phase Skills<br/>sdd-init, sdd-explore,<br/>sdd-apply, sdd-verify, etc."]
            WORKFLOW_SKILLS["Workflow Skills<br/>branch-pr, chained-pr,<br/>work-unit-commits, etc."]
            DOMAIN_SKILLS["Domain Skills<br/>go-testing, cognitive-doc-design,<br/>comment-writer, etc."]
        end

        subgraph MCP["MCP Servers — External Capabilities"]
            CONTEXT7["Context7<br/>(live docs lookup)"]
            ENGRAM_MCP["Engram MCP<br/>(memory access)"]
            PROJECT_MCP["Project MCP<br/>(Notion, Jira, etc.)"]
        end

        subgraph AGENTS["AI Coding Agents (15+)"]
            CLAUDE["Claude Code"]
            OPENCODE["OpenCode"]
            CURSOR["Cursor"]
            OTHERS["12 more agents..."]
        end
    end

    %% SDD reads from and writes to Artifact Store
    SDD -->|reads/writes artifacts| ARTIFACTS

    %% Orchestrator resolves skills and injects into sub-agents
    SKILLS -->|loaded before each phase| SDD

    %% MCP provides tools to sub-agents during phases
    MCP -->|tools available during phases| SDD

    %% Agents run the orchestrator + sub-agents
    AGENTS -->|run orchestrator + sub-agents| SDD

    %% Engram MCP connects Engram memory to agents
    ENGRAM_MCP -.->|exposes| ENGRAM

    %% Style
    classDef sddPhase fill:#f9a8d4,stroke:#be185d,stroke-width:2px,color:#000
    classDef artifact fill:#a78bfa,stroke:#6d28d9,stroke-width:2px,color:#000
    classDef skill fill:#60a5fa,stroke:#1d4ed8,stroke-width:2px,color:#000
    classDef mcp fill:#34d399,stroke:#047857,stroke-width:2px,color:#000
    classDef agent fill:#fbbf24,stroke:#b45309,stroke-width:2px,color:#000

    class EXPLORE,PROPOSE,SPEC,DESIGN,TASKS,APPLY,VERIFY,ARCHIVE sddPhase
    class ENGRAM,OPENSPEC,HYBRID artifact
    class SDD_SKILLS,WORKFLOW_SKILLS,DOMAIN_SKILLS skill
    class CONTEXT7,ENGRAM_MCP,PROJECT_MCP mcp
    class CLAUDE,OPENCODE,CURSOR,OTHERS agent
```

---

## How the Pillars Connect to SDD

### Artifact Store — Persistence

Each SDD phase reads from and writes to the artifact store. The orchestrator decides the mode at session start:

```mermaid
flowchart LR
    SESSION_START[Session Preflight] --> MODE{Artifact Store Mode?}
    MODE -->|engram| ENGRAM_MODE["Engram only<br/>Fast, no files<br/>Upserts overwrite<br/>No iteration history"]
    MODE -->|openspec| OPENSPEC_MODE["OpenSpec only<br/>Files in repo<br/>Git history<br/>Team-shareable"]
    MODE -->|hybrid| HYBRID_MODE["Hybrid (both)<br/>Files + memory<br/>Cross-session recovery<br/>Higher token cost"]
    MODE -->|none| NONE_MODE["None<br/>Inline only<br/>No persistence<br/>Not recommended"]

    classDef mode fill:#a78bfa,stroke:#6d28d9,stroke-width:2px,color:#000
    class ENGRAM_MODE,OPENSPEC_MODE,HYBRID_MODE,NONE_MODE mode
```

Artifact routing by mode:

| Artifact | Engram topic_key | OpenSpec path |
|----------|-----------------|---------------|
| Exploration | `sdd/{change}/explore` | `openspec/changes/{change}/exploration.md` |
| Proposal | `sdd/{change}/proposal` | `openspec/changes/{change}/proposal.md` |
| Spec | `sdd/{change}/spec` | `openspec/changes/{change}/specs/{domain}/spec.md` |
| Design | `sdd/{change}/design` | `openspec/changes/{change}/design.md` |
| Tasks | `sdd/{change}/tasks` | `openspec/changes/{change}/tasks.md` |
| Apply progress | `sdd/{change}/apply-progress` | `openspec/changes/{change}/apply-progress.md` |
| Verify report | `sdd/{change}/verify-report` | `openspec/changes/{change}/verify-report.md` |
| Archive report | `sdd/{change}/archive-report` | `openspec/changes/archive/{date}-{change}/archive-report.md` |

### Skills — Execution Patterns

The orchestrator resolves skills from the registry ONCE per session and injects exact `SKILL.md` paths into each sub-agent's prompt. Sub-agents read those files BEFORE phase-specific work.

```mermaid
flowchart TB
    REGISTRY[".atl/skill-registry.md<br/>(project-local index)"] --> MATCH{Match by<br/>file context + task context}
    MATCH -->|sdd-explore phase| EXPLORE_SKILL["sdd-explore/SKILL.md"]
    MATCH -->|sdd-apply phase| APPLY_SKILL["sdd-apply/SKILL.md<br/>+ work-unit-commits/SKILL.md"]
    MATCH -->|sdd-verify phase| VERIFY_SKILL["sdd-verify/SKILL.md<br/>+ go-testing/SKILL.md"]
    MATCH -->|PR creation| PR_SKILL["branch-pr/SKILL.md<br/>+ chained-pr/SKILL.md"]

    EXPLORE_SKILL --> SUBAGENT1["Sub-agent<br/>(fresh context)"]
    APPLY_SKILL --> SUBAGENT2["Sub-agent<br/>(fresh context)"]
    VERIFY_SKILL --> SUBAGENT3["Sub-agent<br/>(fresh context)"]
    PR_SKILL --> SUBAGENT4["Sub-agent<br/>(fresh context)"]

    SUBAGENT1 --> REPORT1["skill_resolution: paths-injected"]
    SUBAGENT2 --> REPORT2["skill_resolution: paths-injected"]
    SUBAGENT3 --> REPORT3["skill_resolution: paths-injected"]
    SUBAGENT4 --> REPORT4["skill_resolution: paths-injected"]

    classDef registry fill:#60a5fa,stroke:#1d4ed8,stroke-width:2px,color:#000
    classDef skill fill:#93c5fd,stroke:#2563eb,stroke-width:1px,color:#000
    classDef subagent fill:#fbbf24,stroke:#b45309,stroke-width:2px,color:#000
    classDef report fill:#d1d5db,stroke:#4b5563,stroke-width:1px,color:#000

    class REGISTRY registry
    class EXPLORE_SKILL,APPLY_SKILL,VERIFY_SKILL,PR_SKILL skill
    class SUBAGENT1,SUBAGENT2,SUBAGENT3,SUBAGENT4 subagent
    class REPORT1,REPORT2,REPORT3,REPORT4 report
```

Skill resolution feedback loop: every sub-agent reports `skill_resolution` (paths-injected | fallback-registry | fallback-path | none). If the orchestrator sees anything other than `paths-injected`, it re-reads the registry and passes skill paths in subsequent delegations.

### MCP Servers — External Capabilities

MCP servers provide tools that sub-agents can call during any SDD phase:

```mermaid
flowchart LR
    subgraph MCP_SERVERS["MCP Servers"]
        CONTEXT7["Context7<br/>Live documentation<br/>lookup for any library"]
        ENGRAM_MCP["Engram MCP<br/>Memory access from<br/>any agent"]
        CUSTOM["Custom MCP<br/>Notion, Jira, databases,<br/>project-specific tools"]
    end

    subgraph SDD_PHASES["SDD Phases"]
        EXPLORE["explore"]
        SPEC["spec"]
        DESIGN["design"]
        APPLY["apply"]
        VERIFY["verify"]
    end

    CONTEXT7 -.->|"'what's the latest<br/>Next.js auth pattern?'"| EXPLORE
    CONTEXT7 -.->|"verify API usage<br/>against real docs"| SPEC
    CONTEXT7 -.->|"check framework<br/>conventions"| DESIGN
    ENGRAM_MCP -.->|"recall past<br/>decisions"| EXPLORE
    ENGRAM_MCP -.->|"save discoveries"| APPLY
    ENGRAM_MCP -.->|"recall bug fixes"| VERIFY
    CUSTOM -.->|"read Jira ticket"| SPEC
    CUSTOM -.->|"update Notion doc"| ARCHIVE["archive"]

    classDef mcp fill:#34d399,stroke:#047857,stroke-width:2px,color:#000
    classDef phase fill:#f9a8d4,stroke:#be185d,stroke-width:1px,color:#000

    class CONTEXT7,ENGRAM_MCP,CUSTOM mcp
    class EXPLORE,SPEC,DESIGN,APPLY,VERIFY,ARCHIVE phase
```

---

## Phase-by-Phase Ecosystem Usage

```mermaid
flowchart TB
    subgraph PHASES["SDD Phase Flow"]
        direction TB
        P1["1. explore<br/>Read: nothing<br/>Write: exploration<br/>Skills: sdd-explore<br/>MCP: Context7 (docs), Engram (recall)"]
        P2["2. propose<br/>Read: exploration<br/>Write: proposal<br/>Skills: sdd-propose<br/>MCP: Engram (recall context)"]
        P3["3. spec<br/>Read: proposal<br/>Write: spec<br/>Skills: sdd-spec<br/>MCP: Context7 (verify APIs)"]
        P4["4. design<br/>Read: proposal<br/>Write: design<br/>Skills: sdd-design<br/>MCP: Context7 (conventions)"]
        P5["5. tasks<br/>Read: spec + design<br/>Write: tasks<br/>Skills: sdd-tasks<br/>MCP: none typically"]
        P6["6. apply<br/>Read: tasks + spec + design<br/>Write: apply-progress<br/>Skills: sdd-apply + domain skills<br/>MCP: Context7 (implementation), Engram (save)"]
        P7["7. verify<br/>Read: spec + tasks + progress<br/>Write: verify-report<br/>Skills: sdd-verify + go-testing<br/>MCP: Engram (recall fixes)"]
        P8["8. archive<br/>Read: all artifacts<br/>Write: archive-report<br/>Skills: sdd-archive<br/>MCP: Custom (update external docs)"]

        P1 --> P2 --> P3 --> P4 --> P5 --> P6 --> P7 --> P8
    end

    classDef phase fill:#f9a8d4,stroke:#be185d,stroke-width:2px,color:#000
    class P1,P2,P3,P4,P5,P6,P7,P8 phase
```

---

## Delegation Model

The orchestrator is a COORDINATOR, not an executor. It delegates real work to sub-agents and synthesizes results. Each sub-agent gets a fresh context with no memory — the orchestrator controls context access.

```mermaid
flowchart TB
    ORCHESTRATOR["Orchestrator<br/>(thin thread, coordinates)"]

    ORCHESTRATOR -->|"delegates explore"| SUB_EXPLORE["sdd-explore sub-agent<br/>Fresh context"]
    ORCHESTRATOR -->|"delegates propose"| SUB_PROPOSE["sdd-propose sub-agent<br/>Fresh context"]
    ORCHESTRATOR -->|"delegates spec"| SUB_SPEC["sdd-spec sub-agent<br/>Fresh context"]
    ORCHESTRATOR -->|"delegates design"| SUB_DESIGN["sdd-design sub-agent<br/>Fresh context"]
    ORCHESTRATOR -->|"delegates tasks"| SUB_TASKS["sdd-tasks sub-agent<br/>Fresh context"]
    ORCHESTRATOR -->|"delegates apply"| SUB_APPLY["sdd-apply sub-agent<br/>Fresh context"]
    ORCHESTRATOR -->|"delegates verify"| SUB_VERIFY["sdd-verify sub-agent<br/>Fresh context"]
    ORCHESTRATOR -->|"delegates archive"| SUB_ARCHIVE["sdd-archive sub-agent<br/>Fresh context"]

    ORCHESTRATOR -->|"passes skill paths"| SKILL_INJECT["SKILL.md paths<br/>pre-resolved from registry"]
    ORCHESTRATOR -->|"passes artifact refs"| ARTIFACT_REFS["topic_keys or file paths<br/>(NOT content)"]
    ORCHESTRATOR -->|"gatekeeps between phases"| GATE["Validates: contract conformance,<br/>artifact existence, no hallucination,<br/>no drift, routing coherence"]

    SUB_EXPLORE -->|"returns"| RESULT["Result Contract:<br/>status, summary, artifacts,<br/>next_recommended, risks, skill_resolution"]
    GATE -->|"PASS → continue"| ORCHESTRATOR
    GATE -.->|"FAIL → re-run once"| ORCHESTRATOR

    classDef orch fill:#fbbf24,stroke:#b45309,stroke-width:2px,color:#000
    classDef sub fill:#f9a8d4,stroke:#be185d,stroke-width:2px,color:#000
    classDef inject fill:#60a5fa,stroke:#1d4ed8,stroke-width:1px,color:#000
    classDef gate fill:#f87171,stroke:#b91c1c,stroke-width:2px,color:#000
    classDef result fill:#d1d5db,stroke:#4b5563,stroke-width:1px,color:#000

    class ORCHESTRATOR orch
    class SUB_EXPLORE,SUB_PROPOSE,SUB_SPEC,SUB_DESIGN,SUB_TASKS,SUB_APPLY,SUB_VERIFY,SUB_ARCHIVE sub
    class SKILL_INJECT,ARTIFACT_REFS inject
    class GATE gate
    class RESULT result
```

---

## Dependency Graph

```mermaid
flowchart TB
    PROPOSAL["proposal"] --> SPEC["spec"]
    PROPOSAL --> DESIGN["design"]
    SPEC --> TASKS["tasks"]
    DESIGN --> TASKS
    TASKS --> APPLY["apply"]
    APPLY --> VERIFY["verify"]
    VERIFY --> ARCHIVE["archive"]

    classDef phase fill:#f9a8d4,stroke:#be185d,stroke-width:2px,color:#000
    class PROPOSAL,SPEC,DESIGN,TASKS,APPLY,VERIFY,ARCHIVE phase
```

---

## See Also

- [Skill Registry](skill-registry.md) — how skills are indexed and resolved
- [Engram Protocol](engram.md) — persistent memory protocol details
- [OpenSpec Config](openspec-config.md) — file-based artifact configuration
- [Intended Usage](intended-usage.md) — how SDD fits into the daily workflow
- [Components](components.md) — Gentle-AI internal component architecture
