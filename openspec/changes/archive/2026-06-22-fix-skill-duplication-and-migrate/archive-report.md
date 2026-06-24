# Archive Report: fix-skill-duplication-and-migrate

**Date**: 2026-06-22
**Change**: fix-skill-duplication-and-migrate
**Verify Status**: PASS (pre-existing)
**Archive Mode**: hybrid (OpenSpec + Engram)

---

## What Was Implemented

**Part 1 — LLM-first skills refactor (commits `99f7062` + `f070845`):**
- Slimmed 5 heavy SKILL.md files to ≤60 lines each (chained-pr, judgment-day, sdd-init, go-testing, sdd-verify)
- Extracted verbose content to `references/*.md` companion files
- Updated both injectors (`internal/components/skills/inject.go`, `internal/components/sdd/inject.go`) to copy subdirectories recursively
- Shipped `docs/skill-style-guide.md`
- Token cost on activation drops ~85% for the 5 refactored skills

**Part 2 — SDD picker frontmatter flags (commit `7a3bff9`):**
- Added `user-invocable: false` and `disable-model-invocation: true` to all 11 SDD SKILL.md files
- Widened frontmatter linter allowlist by 2 keys
- Regenerated 16 golden files across adapters
- Closes #457

---

## Task Completion Gate

- **Implementation tasks**: All checked (Part 1 and Part 2)
- **Phase 5 manual verification steps** (5.1–5.6): UNCHECKED — these are human-in-the-loop post-merge verification steps, not implementation tasks. They remain as a reviewer checklist for manual confirmation against a real Claude Code v2.1.131+ environment.
- **Gate verdict**: PASS — unchecked items are manual verification, not implementation. Orchestrator approved direct archive.

---

## Verify Report

**Verdict**: PASS
- All in-scope CI-verifiable scenarios (A–E, H–L) pass
- Full test suite green across 40+ packages
- WARNING-1 (inject.go contamination flag) resolved — changes are intentional Part 1 scope
- `size:exception` applied — Part 1 was pre-existing local commits bundled in PR #458

---

## Spec Deltas

No delta specs to merge — the spec is a full specification at the change root level (not in a `specs/` subdirectory). No corresponding domain exists in `openspec/specs/`. The spec documents both Part 1 and Part 2 scope.

---

## Stale-Checkbox Reconciliation

None required — all implementation tasks were checked. Phase 5 manual verification items are intentionally unchecked (human-in-the-loop post-merge steps).

---

## Warnings Accepted

None — verify report has no open warnings.

---

## Archive Contents

- proposal.md ✅
- spec.md ✅ (full spec, not delta)
- design.md ✅
- tasks.md ✅ (Phase 5 manual steps intentionally unchecked)
- verify-report.md ✅
