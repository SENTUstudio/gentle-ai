---
name: sdd-init
description: "Trigger: sdd init, iniciar sdd, openspec init. Initialize SDD context, testing capabilities, registry, and persistence."
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
> the dedicated `sdd-init` sub-agent using your platform's delegation primitive
> (e.g., `task(...)`, sub-agent invocation, etc.). This skill is for EXECUTORS
> only.

## Executor Override

If you ARE the `sdd-init` sub-agent (NOT the orchestrator), the gate above does NOT apply to you. Continue with the phase work below. Do NOT delegate. Do NOT call the Skill tool. You are the executor — execute.

## Language Domain Contract

Generated technical artifacts default to English. Do not inherit the user's conversational language or the active persona's regional voice for SDD artifacts unless the user explicitly requests that artifact language or the project convention requires it.

If Spanish technical artifacts are explicitly requested, use neutral/professional Spanish unless the user explicitly asks for a regional variant.

Public/contextual comments follow the target context language by default. Explicit user language or tone overrides win; Spanish comments default to neutral/professional Spanish unless the user or target context clearly calls for regional tone.

## Activation Contract

Run this phase when the orchestrator/user asks to initialize SDD in a project. You are the phase executor: do the work yourself, do not delegate, and do not behave like the orchestrator.

## Hard Rules

- Detect the real stack, conventions, architecture, testing tools, and persistence mode; never guess.
- In `engram` mode, do **not** create `openspec/`.
- In `openspec` mode, follow `../_shared/openspec-convention.md` and write file artifacts.
- In `hybrid` mode, write both openspec files and Engram observations.
- Always persist testing capabilities separately as `sdd/{project}/testing-capabilities` or `openspec/config.yaml` `testing:`.
- Always build `.atl/skill-registry.md`; also save `skill-registry` to Engram when available.
- Use `capture_prompt: false` for automated SDD/config saves when supported; omit it if the tool schema lacks it.
- If `openspec/` already exists, report what exists and ask before updating it.

## Decision Gates

| Input | Action |
|---|---|
| `mode=engram` | Save context and capabilities to Engram only. |
| `mode=openspec` | Create/update openspec bootstrap files only. |
| `mode=hybrid` | Do both Engram and openspec persistence. |
| `mode=none` | Return detected context only; write no SDD artifacts except registry if required. |
| strict TDD marker/config found | Use that value. |
| no marker/config but test runner exists | Default `strict_tdd: true`. |
| no test runner | Set `strict_tdd: false` and explain unavailable. |
| `openspec/config.yaml` already sets `domain:` | Respect it; skip re-detection. |
| no `domain:` in config | Run domain preflight (see Domain Preflight below). |

## Domain Preflight

SDD branches on an optional `domain` field in `openspec/config.yaml`. When
`domain` is absent the project is app-dev and every code path stays identical to
today. The preflight detects a data-engineering project, presents the hint for
confirmation, and writes the confirmed value so later phases can branch on it.

1. Run `gentle-ai sdd-config --detect --json` (or `--detect` for human output).
   It scans the project root for `template.yaml` AND `glue-jobs/*.py` and returns
   `domain`, `confidence`, and `evidence`.
   - confidence 0.8 (both markers) → strong signal.
   - confidence 0.5 (one marker) → weak hint.
   - confidence 0 (no markers) → app-dev; skip the rest of this preflight.
2. Present the hint to the user: domain, confidence, and evidence. Never apply a
   detected domain silently.
3. Confirm or override:
   - User confirms `data-engineering` → proceed to write.
   - User overrides (e.g. this is an app-dev project despite markers) → write
     the user's choice, or leave `domain` absent for app-dev.
   - User declines to set a domain now → leave `domain` absent (app-dev).
4. Write the confirmed `domain` to `openspec/config.yaml` (append or update the
   top-level `domain:` key; preserve all other keys and the `context: |` block).
5. Optionally run `gentle-ai sdd-config --validate-repos` to warn about missing
   `repos.infra` / `repos.carga` paths for data-engineering projects.

Skill patches in later phases (sdd-spec, sdd-design, sdd-apply, sdd-verify,
sdd-tasks) are gated on `domain: data-engineering`; when `domain` is unset they
follow the existing app-dev path unchanged.

## Execution Steps

1. Run the Domain Preflight above (skip if `domain` is already set in config).
2. Inspect project files (`package.json`, `go.mod`, `pyproject.toml`, CI, lint/test config) and summarize stack/conventions.
3. Detect test runner, test layers, coverage, linter, type checker, and formatter.
4. Resolve Strict TDD from agent marker, `openspec/config.yaml`, detected runner fallback, or no-runner fallback.
5. Initialize persistence for the resolved mode.
6. Build `.atl/skill-registry.md` using the skill-registry scan rules.
7. Persist testing capabilities and project context.
8. Return the structured initialization envelope.

## Output Contract

Return `status`, `executive_summary`, `artifacts`, `next_recommended`, and `risks`. Include project, stack, persistence mode, Strict TDD status, testing capability table, saved observation IDs/paths, registry path, and next `/sdd-explore` or `/sdd-new` step.

## References

- [references/init-details.md](references/init-details.md) — detection checklist, Engram payloads, config skeleton, and output templates.
- `../_shared/engram-convention.md` — Engram artifact naming.
- `../_shared/openspec-convention.md` — openspec layout and rules.
