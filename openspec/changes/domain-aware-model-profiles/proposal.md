# Proposal: Domain-Aware Model Profiles

## Intent

Model profiles (`model.Profile`) are domain-blind: a `cheap` profile assigns
Haiku to `sdd-explore` regardless of whether the project is app-dev or
data-engineering. Per `docs/sdd-profiles-domain-model.md`, data-engineering
needs Sonnet for `sdd-explore` (data profiling = judgment) and Opus for
`sdd-spec` (schema + DAG = precision). The domain config (`sddconfig.Domain`)
already exists from PR #2 but is not read by the profile system. This change
bridges the two so model assignments vary by domain.

## Scope

### In Scope

- Add `Domain string` to `model.Profile` (empty = app-dev, backward-compat).
- New pure function `DefaultModelsForDomain(domain)` sourced from the doc table.
- TUI: auto-detect domain from `sddconfig`, pre-fill model picker, show in footer.
- CLI: new `--profile-domain name:domain` flag + auto-detect from `sddconfig`.
- Collision guard: same profile name + different domain = error/disambiguation.
- Spec delta: ADDED "Profile Domain" requirement + round-trip scenario.

### Out of Scope

- Encoding `Domain` in `opencode.json` overlay (Q3=A: Go-only state for v1).
- YAML-declared per-domain defaults (Q2=B: premature scope).
- New TUI step for domain selection (Q5: no new screen).
- Profile rename / domain migration flow (follow-up).

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `sdd-profiles`: ADDED `Domain` field on Profile + per-domain default models + collision guard + round-trip preservation scenario.

## Approach

Domain is **metadata on the struct, not a name suffix or JSON field**. All 6
exploration decisions adopted as a coherent set (Q1=C, Q2=C, Q3=A, Q4=C+A,
Q5=C+B, Q6=C). The overlay JSON schema is untouched; the domain lives only in
Go state (`Selection.Profiles[*].Domain`). Both CLI and TUI read domain from
`sddconfig.LoadConfig` by default; the explicit `--profile-domain` flag exists
for CI/headless override. The single recommendation table is the shared source
of truth for both `DefaultModelsForDomain` and orchestrator prompt rendering,
preventing drift.

**Delivery**: single PR (~500 lines), 800-line review budget, landing on
`feature/aws-dataengineer` (PR #2 branch). User-confirmed override of the
exploration's 3-PR recommendation.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/model/types.go` | Modified | Add `Domain string` to `Profile` + backward-compat comment. |
| `internal/model/profiles_default.go` | New | `DefaultModelsForDomain(domain)` pure function + table tests. |
| `internal/components/sdd/profiles.go` | Modified | Orchestrator prompt table renders domain column; overlay loop absorbs domain. |
| `internal/cli/sync.go` | Modified | `--profile-domain` flag + `parseProfileFlags` sddconfig auto-detect + collision error. |
| `internal/tui/screens/profile_create.go` | Modified | Footer domain detection + picker pre-fill + step-2 collision warning. |
| `internal/tui/screens/model_picker.go` | Modified | Domain-aware tier legend (cosmetic). |
| `openspec/specs/sdd-profiles/spec.md` | Delta | ADDED Profile Domain requirement + round-trip scenario. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Defaults drift from doc table | Low | Pure function + test asserts each (phase, model) pair against the doc. |
| Name collision in `opencode.json` (one `sdd-orchestrator-{name}` slot) | Medium | TUI warning at step 2 + CLI hard error on same-name/different-domain. |
| Domain lost on hand-edit of `opencode.json` | Low | Accepted v1 trade-off; profiles are managed artifacts, next sync re-sets domain. |
| Headless sync without `openspec/config.yaml` | Low | Auto-detect returns empty = app-dev (documented safe default). |
| Single PR exceeds 800-line budget | Low | Exploration estimates ~500 lines; tasks phase re-forecasts before apply. |

## Rollback Plan

Revert the single PR. The `Domain` field is additive and zero-value for all
existing profiles; revert restores today's domain-blind behavior with no data
loss. No migration, no schema change, no JSON overlay change to undo.

## Dependencies

- PR #2 (`data-engineering-domain-profile`) merged — provides `sddconfig.Config.Domain`.

## Success Criteria

- [ ] Existing `makeHaikuProfile()` tests pass unchanged (empty Domain = app-dev).
- [ ] `DefaultModelsForDomain("data-engineering")` returns Sonnet for explore, Opus for spec/design.
- [ ] TUI in a data-eng project pre-fills the data-eng defaults; user can override any row.
- [ ] CLI `--profile cheap:haiku` in a data-eng project produces `Profile.Domain = "data-engineering"`.
- [ ] Same name + different domain → TUI warns, CLI errors.
- [ ] Round-trip (create → serialize → re-detect) preserves domain in Go state.
