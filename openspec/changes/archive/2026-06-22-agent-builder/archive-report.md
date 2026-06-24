# Archive Report: agent-builder

**Date**: 2026-06-22
**Change**: agent-builder
**Verify Status**: PASS WITH WARNINGS (new, static evidence)
**Archive Mode**: hybrid (OpenSpec + Engram)

---

## What Was Implemented

Agent Builder — TUI flow for creating custom AI sub-agents:
- `internal/agentbuilder/` package: types, engine interface (Claude/OpenCode/Gemini/Codex), prompt composition, output parser, JSON registry, atomic installer, SDD marker injection
- 8 TUI screens: engine selection, prompt input, SDD integration, SDD phase picker, generating (spinner), preview, installing, complete
- Router + Welcome menu wiring ("Create your own Agent" at index 5, disabled when no engines)
- Model Update/View/confirmSelection/goBack handlers for all 8 screens
- 122 tests across 13 files (unit + integration/TUI)
- PR #223 (v1.18.0) plus follow-ups for robustness

---

## Task Completion Gate

- **Implementation tasks**: All 29 tasks implemented per git evidence (PR #223, v1.18.0)
- **tasks.md checkboxes**: All 29 were unchecked — stale checkboxes
- **Stale-checkbox reconciliation**: APPROVED by orchestrator. All 29 reconciled to `[x]`.

---

## Verify Report

**Verdict**: PASS WITH WARNINGS
- Static-evidence-only verification (no runtime tests executed)
- 33/40 spec scenarios implemented as specified
- 4 not implemented as specified (deviations documented below)
- 3 partial/static-only

---

## Spec Deviations (Accepted by Orchestrator)

1. **5-minute timeout** instead of spec/design 120 seconds — generation timeout is `5*time.Minute`
2. **Missing Preview Edit** — `$EDITOR`/`vi` editing action not implemented in preview screen
3. **Silent registry replace** — custom-agent name conflict silently replaces instead of prompting user
4. **New-phase graph rewrite** — only appends marker block, does not rewrite orchestrator dependency graph string
5. **SDD injection method** — direct file I/O instead of `StrategyMarkdownSections` helper

These are documented deviations, not blockers. Orchestrator accepted all warnings.

---

## Spec Deltas

No delta specs to merge — no `specs/` subdirectory.

---

## Stale-Checkbox Reconciliation

**Reconciled**: 29 checkboxes (T-01 through T-29) from `[ ]` to `[x]`.
**Reason**: Tasks were implemented but checkboxes were not updated. Verify-report confirms all tasks complete with git evidence (commits `8a54d9b`, `e98e5af`, `c266163`).
**Orchestrator approval**: Explicitly approved for this change.

---

## Warnings Accepted

1. Static-evidence-only verification — no runtime test execution (orchestrator approved)
2. 29 stale checkboxes reconciled — all implemented per git evidence
3. 5-minute timeout deviation from 120s spec
4. Missing Preview Edit action
5. Silent custom-agent replace instead of confirmation dialog
6. New-phase graph not rewritten (marker block only)
7. SDD injection via direct file I/O
8. No `apply-progress` artifact

---

## Archive Contents

- proposal.md ✅
- spec.md ✅ (no delta, full spec)
- design.md ✅
- tasks.md ✅ (29/29 reconciled)
- verify-report.md ✅
