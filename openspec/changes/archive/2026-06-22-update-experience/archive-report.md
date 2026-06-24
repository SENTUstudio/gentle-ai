# Archive Report: update-experience

**Date**: 2026-06-22
**Change**: update-experience
**Verify Status**: PASS WITH WARNINGS (new, static evidence)
**Archive Mode**: hybrid (OpenSpec + Engram)

---

## What Was Implemented

Update experience overhaul across 6 slices (2–7), delivered as 5 PRs (#875–#880):

**Slice 2 — Update-Check Cooldown**: 6h TTL in `state.json`, refresh-on-success-only, clock injection for test determinism.

**Slice 3 — Channel-Honoring Upgrade**: `GENTLE_AI_CHANNEL` honored for engram download (beta → `@main`, stable → pinned version).

**Slice 4 — Upgrade+Sync Deferred**: `pending_sync` flag written before exit on self-upgrade; deferred sync runs on next launch; clears on success, retries on failure.

**Slice 5 — CLI Prompt Default**: Removed `GENTLE_AI_CONFIRM_UPDATE` env gate; prompt shown unconditionally; `[Y/n]` default; `--yes`/`GENTLE_AI_YES=1` auto-accept.

**Slice 6 — TUI Pre-Welcome Update Prompt**: `ScreenUpdatePrompt` with Update/View changes/Keep current options; spinner while checking; cooldown gate integrated.

**Slice 7 — Advisory Manifest**: `FetchAdvisory` with 2s timeout, fail-open; displayed on Welcome screen; no version gating.

45/45 tasks checked, 30/30 spec scenarios covered (static evidence).

---

## Task Completion Gate

- **Implementation tasks**: All 45 tasks checked (`[x]`)
- **Gate verdict**: PASS — all tasks complete (reconciled by verify-report before this archive run)

---

## Verify Report

**Verdict**: PASS WITH WARNINGS
- Static-evidence-only verification (no runtime tests executed per orchestrator directive)
- 30/30 spec scenarios covered by static test evidence
- Warnings: no `apply-progress` artifact, no runtime test execution, advisory endpoint uses dedicated `advisory` tag instead of `latest` release asset (functionally equivalent)

---

## Spec Deltas

**Created 7 new main specs** in `openspec/specs/`:
| Domain | Action |
|--------|--------|
| `advisory-manifest` | Created |
| `self-update` | Created |
| `update-check-cache` | Created |
| `update-prompt` | Created |
| `upgrade-channel` | Created |
| `upgrade-sync` | Created |
| `version-resolution` | Created |

All 7 were full specs (no existing main specs to merge into).

---

## Stale-Checkbox Reconciliation

None required for this archive run — verify-report had already reconciled all 45 checkboxes before this archive phase.

---

## Warnings Accepted

1. Static-evidence-only verification — no runtime test execution (orchestrator approved)
2. No `apply-progress` artifact — TDD cycle evidence missing
3. Advisory endpoint uses dedicated `advisory` tag instead of `latest` release asset (functionally equivalent, documented deviation)

---

## Archive Contents

- proposal.md ✅
- specs/ (7 domains, all copied to main specs) ✅
- design.md ✅
- tasks.md ✅ (45/45 tasks complete)
- verify-report.md ✅
