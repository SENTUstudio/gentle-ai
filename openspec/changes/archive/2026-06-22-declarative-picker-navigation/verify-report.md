## Verification Report

**Change**: declarative-picker-navigation
**Version**: v1.41.0
**Mode**: Strict TDD (test runner: `go test ./...`) — **static-evidence-only verify; no test execution performed per orchestrator instruction**

---

### Completeness

| Metric | Value |
|--------|-------|
| Tasks total | 35 |
| Tasks complete | 35 |
| Tasks incomplete | 0 |

All checkboxes in `tasks.md` are marked complete; no stale-checkbox issue reported.

---

### Build & Tests Execution

**Build**: ✅ Passed
```text
$ go build ./...
(no errors)
```

**Tests**: ➖ Not executed
```text
Test runner: go test ./...
Instruction: static-evidence-only verification; no test runs needed.
```

**Coverage**: ➖ Not available (no coverage run performed)

> Note: Because tests were not executed, this report cannot assert runtime pass/fail for any scenario. Compliance below is reported as "COVERED (static)" when a covering test or regression test exists in source, not as a runtime PASS.

---

### Spec Compliance Matrix

| Requirement | Scenario | Test | Result |
|-------------|----------|------|--------|
| INV-1 Forward order | Scenario 1 — All agents full forward | `TestInstallNavigationRoundTrips` (all 11 cases) | ⚠️ COVERED (static) — test exists; runtime not verified |
| INV-1 Forward skip | Scenario 2 — Kiro excluded | `TestInstallNavigationRoundTrips` | ⚠️ COVERED (static) |
| INV-2 Backward parity | Scenario 3 — SDDMode back to Codex (INV-2a) | `TestInstallNavigationRoundTrips` + `TestPickerBackRowRegression` | ⚠️ COVERED (static) |
| INV-2 Backward parity | Scenario 4 — Codex "← Back" not inert (INV-2b) | `TestInstallNavigationRoundTrips` + `TestPickerBackRowRegression` | ⚠️ COVERED (static) |
| INV-6 ModelPicker gate | Scenario 5 — ModelPicker excluded when SDDMode != Multi | `TestPickerFlowSlice` | ⚠️ COVERED (static) |
| INV-6 ModelPicker gate | Scenario 6 — ModelPicker excluded when cache absent | `TestPickerFlowSlice` | ⚠️ COVERED (static) |
| INV-3 Custom inversion | Scenario 7 — PresetCustom ordering inversion | `TestPickerFlowSlice` + `TestCustomPresetPostComponentFlowMatrix` | ⚠️ COVERED (static) |
| INV-4 ModelConfigMode exit | Scenario 8 — ModelConfigMode exit-ramp | Existing ModelConfigMode tests | ⚠️ COVERED (static) |
| INV-7 Round-trip symmetry | Scenario 9 — Full symmetry (all agents) | `TestInstallNavigationRoundTrips` | ⚠️ COVERED (static) |
| INV-7 Round-trip symmetry | Scenario 10 — Full symmetry (partial subset) | `TestInstallNavigationRoundTrips` | ⚠️ COVERED (static) |
| INV-5 OpenCodePluginsStandalone guard | (no scenario) | Existing standalone guard tests + `TestStrictTDDForward` | ⚠️ COVERED (static) |

**Compliance summary**: 11/11 spec invariants/scenarios have covering tests or existing test enforcement in source. Runtime execution was not performed.

---

### Correctness (Static Evidence)

| Requirement | Status | Notes |
|------------|--------|-------|
| INV-1 — Forward navigation order | ✅ Implemented | `pickerFlowSlice()` builds `ScreenPreset → [Claude]* → [Kiro]* → [Codex]* → [SDDMode]* → [ModelPicker]** → [StrictTDD]* → [DependencyTree]` in `internal/tui/model.go:3660-3692`. |
| INV-2 — Backward navigation parity | ✅ Implemented | `pickerPreviousScreen()` in `model.go:3710-3718` scans the same slice; all Back-row call sites in `confirmSelection` and `goBack` route through it. |
| INV-2a — SDDMode back to Codex | ✅ Implemented | `goBack` `ScreenSDDMode` block (`model.go:2895-2900`) and SDDMode Back-row (`model.go:1933-1936`) both use `pickerPreviousScreen`. |
| INV-2b — Codex Back row not inert | ✅ Implemented | `confirmSelection` `ScreenCodexModelPicker` Back row (`model.go:1896-1907`) uses `pickerPreviousScreen`; `goBack` `ScreenCodexModelPicker` (`model.go:2916-2921`) mirrors it. |
| INV-3 — PresetCustom inversion | ✅ Implemented | `pickerFlowSlice()` appends `ScreenDependencyTree` second when `Preset == PresetCustom` (`model.go:3663-3666`) and skips the trailing anchor. |
| INV-4 — ModelConfigMode exit-ramp | ✅ Implemented | Early returns preserved in HandleNav forward (`model.go:1131-1139`, `1160-1167`, `1208-1227`, `1984-1993`), Back rows (`model.go:1874-1882`, `1886-1895`, `1898-1907`, `2005-2009`), and `goBack` (`model.go:2822-2827`). |
| INV-5 — OpenCodePluginsStandalone guard | ✅ Implemented | Guard remains outside the slice in `confirmSelection` (`model.go:2018-2020`, `2097-2103`) and `goBack` (`model.go:2891-2893`, plus references in `model.go:3351`, `3368`, `3393`). |
| INV-6 — ModelPicker uses injectable stat | ✅ Implemented | `pickerFlowSlice()` calls `osStatModelCache(opencode.DefaultCachePath())` (`model.go:3679`); `osStatModelCache` remains a package var. |
| INV-7 — Round-trip symmetry | ✅ Implemented | Slice is deterministic per `m.Selection`; forward/backward scanners are exact inverses. |

---

### Coherence (Design)

| Decision | Followed? | Notes |
|----------|-----------|-------|
| Source of order: one `pickerFlowSlice()` | ✅ Yes | Implemented in `model.go:3660-3692`; `linearRoutes` in `router.go` was not touched. |
| Slice endpoints include Preset + DependencyTree anchors | ✅ Yes | Both anchors always present (custom/non-custom). |
| Custom inversion encoded in slice | ✅ Yes | Single `if custom` branch controls `ScreenDependencyTree` position. |
| ModelPicker gate via `SDDMode==Multi` + `osStatModelCache` | ✅ Yes | Matches design contract exactly. |
| Cross-cutting guards kept outside slice | ✅ Yes | ModelConfigMode, OpenCodePluginsStandalone, Upgrade-Esc remain early returns. |
| Allocation: rebuild per call | ✅ Yes | Helpers call `pickerFlowSlice()` each invocation; no cache. |
| `applyPickerEntry` initializes all picker targets | ✅ Yes | Switch on target initializes Claude/Kiro/Codex/ModelPicker and calls `setScreen`. |
| StrictTDD forward preserves custom Review fallback | ✅ Yes | `model.go:2018-2038` keeps guard order: OpenCodePlugins → custom SkillPicker/Review → non-custom `pickerNextScreen`. |

---

### TDD Compliance

> `apply-progress` artifact does not exist for this change. TDD compliance is inferred from git history only.

| Check | Result | Details |
|-------|--------|---------|
| TDD Evidence reported | ⚠️ Inferred | No `apply-progress` file found; commit history shows RED/GREEN/REFACTOR-style units. |
| All tasks have tests | ✅ Yes | Unit tests `TestPickerFlowSlice`, `TestPickerNextScreen`, `TestPickerPreviousScreen`, `TestApplyPickerEntry`, `TestPickerBackRowRegression`, and `TestStrictTDDForward` exist in `model_test.go`. |
| RED confirmed (tests exist before methods) | ⚠️ Partially visible | Git history shows helper implementation commits (`0c556b5`, `8aa47aa`) precede their dedicated test commits (`76cdc7f`), which is the reverse of strict RED-first. However, `d1e8f38` (Back-row regression tests) clearly precedes the call-site rewrites (`3f6ac4a`, `c56985b`). |
| GREEN confirmed (tests pass on execution) | ➖ Not verified | Tests were not run. |
| Triangulation adequate | ✅ Yes | `TestPickerFlowSlice` has 7 cases; `TestPickerNextScreen` 11 cases; `TestPickerPreviousScreen` 11 cases; `TestApplyPickerEntry` 9 cases; `TestPickerBackRowRegression` 8 cases; `TestStrictTDDForward` 4 cases. |
| Safety net for modified files | ➖ Not verified | No runtime execution. |

**TDD Compliance**: Inferred from git history; strict RED-first ordering is partially visible. Runtime pass status unknown.

---

### Test Layer Distribution

| Layer | Tests | Files | Tools |
|-------|-------|-------|-------|
| Unit | `TestPickerFlowSlice`, `TestPickerNextScreen`, `TestPickerPreviousScreen`, `TestApplyPickerEntry`, `TestPickerBackRowRegression`, `TestStrictTDDForward` | `internal/tui/model_test.go` | `go test` |
| Integration | `TestInstallNavigationRoundTrips`, flow-matrix goldens | `internal/tui/preset_flow_test.go` | `go test` |
| E2E | 0 | — | — |
| **Total** | **6 new unit tests + existing integration tests** | **2 files** | |

---

### Changed File Coverage

➖ Coverage analysis skipped — no coverage run performed (static-evidence-only verify).

---

### Assertion Quality

✅ No trivial assertions found in the new test functions during static inspection. Tests assert concrete screen transitions and initialized picker state (non-empty presets, non-nil maps).

---

### Quality Metrics

**Linter**: ➖ Not run
**Type Checker**: ✅ `go build ./...` succeeded (no type errors)

---

### Issues Found

**CRITICAL**: None

**WARNING**:
1. **Static-only verification**: No tests were executed. Runtime compliance of all spec scenarios is unproven. This report cannot certify that the suite passes.
2. **PR review budget exceeded**: Gross changed lines for the change (implementation + tests) are ~1109 lines across `internal/tui/model.go` and `internal/tui/model_test.go` per `git diff --stat 41cdb1b^..0a9ffdd`, well above the 400-line single-PR budget forecast in `tasks.md` (~375 lines). The change was delivered as a single PR despite the overshoot.
3. **TDD evidence from apply-progress missing**: No `apply-progress` artifact exists, so strict TDD cycle verification relies on git history inference only.
4. **Helper tests committed after helper implementation**: `TestPickerNextScreen`/`TestPickerPreviousScreen` (`76cdc7f`) were committed after the helpers (`0c556b5`), which is GREEN-before-RED for those units.

**SUGGESTION**:
1. `ScreenPreset` forward entry (`model.go:1833-1867`) still uses a manual predicate ladder to decide the first picker. While out of the explicit refactor scope, it duplicates the order encoded in `pickerFlowSlice()`. Consider converging it onto `pickerNextScreen()` in a follow-up to fully eliminate desync risk.
2. Run the full test suite (`go test ./...`) and, if available, coverage before archiving to convert static "COVERED" statuses into runtime "COMPLIANT" evidence.

---

### Verdict

**PASS WITH WARNINGS**

Implementation aligns with the spec and design based on static inspection and git history; all 35 tasks are checked; `go build ./...` succeeds; no golden files were modified. The verdict is conditional because tests were not executed and the change exceeded the 400-line PR budget.
