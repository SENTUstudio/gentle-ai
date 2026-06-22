## Verification Report

**Change**: update-experience
**Version**: v1.41.0
**Mode**: Strict TDD (static-evidence-only verify; no test execution per orchestrator)

### Completeness
| Metric | Value |
|--------|-------|
| Tasks total | 45 (slices 2–7 task checkboxes) |
| Tasks complete | 45 |
| Tasks incomplete | 0 |

> **Reconciliation note**: `tasks.md` had all checkboxes unchecked when this verify started. Git-log evidence shows all five implementation PRs (#875, #876, #877, #879, #880) are merged into the feature branch, and static inspection confirms every implementation checkbox is satisfied. The orchestrator pre-approved stale-checkbox reconciliation for this cleanup batch; task boxes were reconciled to checked.

### Build & Tests Execution
**Build**: ➖ Not executed (static-evidence-only verify)

**Tests**: ➖ Not executed (static-evidence-only verify)

> Strict TDD is active and test runner `go test ./...` is configured, but this verification was explicitly scoped to static evidence: git log for PRs #875–#880 and source/test-file inspection. The next full verify run should execute the test suite to convert these to runtime evidence.

### Spec Compliance Matrix
| Requirement | Scenario | Test | Result |
|-------------|----------|------|--------|
| **update-check-cache** — Cooldown gate | Cache fresh → no network call | `internal/update/cooldown_test.go > TestCheckAllWithCooldown_FreshCacheSkipsNetwork` | ✅ COMPLIANT (static) |
| **update-check-cache** — Cooldown gate | Cache stale → refresh from GitHub | `internal/update/cooldown_test.go > TestCheckAllWithCooldown_StaleCacheRefreshes` | ✅ COMPLIANT (static) |
| **update-check-cache** — Cooldown gate | Cache missing → first run checks | `internal/update/cooldown_test.go > TestCheckAllWithCooldown_MissingCache` | ✅ COMPLIANT (static) |
| **update-check-cache** — Failure resilience | Rate-limited/network error does not advance timestamp | `internal/update/cooldown_test.go > TestCheckAllWithCooldown_FailedCheckDoesNotAdvanceTimestamp` | ✅ COMPLIANT (static) |
| **update-check-cache** — State persistence | `last_update_check` round-trips and is backward-compatible | `internal/state/state_test.go > TestLastUpdateCheck_*`, `TestMergeAgents_PreservesLastUpdateCheck` | ✅ COMPLIANT (static) |
| **upgrade-channel** — Channel-aware routing | Stable upgrade (default) | `internal/update/upgrade/strategy_test.go > TestEngramBinaryUpgrade_StableChannelCallsDownloadFn` | ✅ COMPLIANT (static) |
| **upgrade-channel** — Channel-aware routing | Beta upgrade (`GENTLE_AI_CHANNEL=beta`) | `internal/update/upgrade/strategy_test.go > TestEngramBinaryUpgrade_BetaChannelUsesGoInstallMain` | ✅ COMPLIANT (static) |
| **upgrade-channel** — Channel-aware routing | Unknown channel value falls back to stable with warning | `internal/cli/channel_test.go > TestResolveInstallChannel` (invalid case) + `strategy.go` warning | ✅ COMPLIANT (static) |
| **upgrade-channel** — Channel-aware routing | Empty env string treated as unset/stable | `internal/cli/channel_test.go > TestResolveInstallChannel` (empty env string case) | ✅ COMPLIANT (static) |
| **upgrade-sync** — Sync completes across self-upgrade | Upgrade without self-upgrade → inline sync | `internal/tui/model_test.go > TestStartUpgradeSync_DoesNotSetPendingSyncWhenGentleAINotUpgraded` | ✅ COMPLIANT (static) |
| **upgrade-sync** — Sync completes across self-upgrade | Upgrade WITH self-upgrade → `pending_sync=true` written | `internal/app/selfupdate_test.go > TestSelfUpdate_SetsPendingSyncOnSuccess`, `internal/tui/model_test.go > TestStartUpgradeSync_SetsPendingSyncWhenGentleAIUpgraded` | ✅ COMPLIANT (static) |
| **upgrade-sync** — Sync completes across self-upgrade | Deferred sync runs on next launch and clears flag | `internal/app/app_test.go > TestRunArgs_PendingSync_RunsSyncAndClearsFlag` | ✅ COMPLIANT (static) |
| **upgrade-sync** — Sync completes across self-upgrade | Deferred sync fails → flag stays set for retry | `internal/app/app_test.go > TestRunArgs_PendingSync_LeavesSetOnFailure` | ✅ COMPLIANT (static) |
| **self-update** — CLI prompt default | Default prompt shown unconditionally | `internal/app/selfupdate_test.go > TestSelfUpdate_PromptAlwaysShown` | ✅ COMPLIANT (static) |
| **self-update** — CLI prompt default | User accepts (Y/Enter) | `internal/app/selfupdate_test.go > TestSelfUpdate_ConfirmUpdate_UserAccepts`, `TestSelfUpdate_ConfirmUpdateTable` | ✅ COMPLIANT (static) |
| **self-update** — CLI prompt default | User declines (N) | `internal/app/selfupdate_test.go > TestSelfUpdate_ConfirmUpdate_UserDeclines`, `TestSelfUpdate_ConfirmUpdateTable` | ✅ COMPLIANT (static) |
| **self-update** — CLI prompt default | `--yes` / `GENTLE_AI_YES=1` auto-accepts | `internal/app/selfupdate_test.go > TestSelfUpdate_YesFlag_AutoAccepts`, `TestSelfUpdate_YesEnvVar_AutoAccepts` | ✅ COMPLIANT (static) |
| **update-prompt** — Launch-time prompt presence | Update available → TUI pre-Welcome screen | `internal/tui/model_test.go > TestUpdatePromptScreen_ShownWhenUpdateAvailable` | ✅ COMPLIANT (static) |
| **update-prompt** — Launch-time prompt presence | No update → skip prompt | `internal/tui/model_test.go > TestUpdatePromptScreen_SkippedWhenNoUpdate` | ✅ COMPLIANT (static) |
| **update-prompt** — Launch-time prompt presence | Check failed/offline → skip prompt | `internal/tui/model_test.go > TestUpdatePromptScreen_SkippedWhenCheckFailed` | ✅ COMPLIANT (static) |
| **update-prompt** — Update action | TUI “Update” runs upgrade and quits | `internal/tui/model_test.go > TestUpdatePromptScreen_KeyU_RunsUpgradeThenQuits` | ✅ COMPLIANT (static) |
| **update-prompt** — Update action | TUI “Keep current version” → Welcome | `internal/tui/model_test.go > TestUpdatePromptScreen_KeyC_TransitionsToWelcome`, `TestUpdatePromptScreen_KeyEnter_TransitionsToWelcome` | ✅ COMPLIANT (static) |
| **update-prompt** — View changes | TUI “View changes” opens browser / prints URL and stays | `internal/tui/model_test.go > TestUpdatePromptScreen_KeyV_CallsOpenBrowser`, `TestUpdatePromptScreen_KeyV_FallsBackWhenBrowserFails` | ✅ COMPLIANT (static) |
| **update-prompt** — Ask-every-launch cadence | No skip/snooze state persisted; `ScreenUpdatePrompt` shown every launch with updates | `internal/tui/model.go` — no snooze fields; `UpdateCheckResultMsg` handler transitions on every result | ✅ COMPLIANT (static) |
| **advisory-manifest** — Manifest fetch at launch | Valid message → displayed | `internal/update/advisory_test.go > TestFetchAdvisory_ValidJSON`, `internal/tui/model_test.go > TestWelcomeView_ContainsAdvisoryMessage` | ✅ COMPLIANT (static) |
| **advisory-manifest** — Manifest fetch at launch | Empty/absent message → silent | `internal/update/advisory_test.go > TestFetchAdvisory_EmptyMessage`, `internal/tui/model_test.go > TestAdvisoryMsg_EmptyAdvisoryNoChange` | ✅ COMPLIANT (static) |
| **advisory-manifest** — Manifest fetch at launch | Unreachable / non-200 / timeout → silent fail-open | `internal/update/advisory_test.go > TestFetchAdvisory_HTTP500`, `TestFetchAdvisory_HTTP404`, `TestFetchAdvisory_Timeout` | ✅ COMPLIANT (static) |
| **advisory-manifest** — Malformed payload | Malformed JSON / unexpected schema → silent | `internal/update/advisory_test.go > TestFetchAdvisory_MalformedJSON`, `TestFetchAdvisory_OversizedBody` | ✅ COMPLIANT (static) |
| **advisory-manifest** — No version gating | Version fields ignored, never blocks | `internal/update/advisory.go` — `Advisory` struct has no version fields; `internal/tui/model.go` displays only `Message` | ✅ COMPLIANT (static) |
| **version-resolution** — Engram always-latest | Stable tag filter `^v[0-9]+\.[0-9]+\.[0-9]+$` used for update check and download | `internal/update/registry.go` (`ReleaseTagPattern`), `internal/components/engram/download.go` (`engramCoreTagPattern`, `fetchLatestEngramVersion`) | ✅ COMPLIANT (static) |

**Compliance summary**: 30/30 scenarios covered by static test evidence.

### Correctness (Static Evidence)
| Requirement | Status | Notes |
|------------|--------|-------|
| `last_update_check` persisted in state | ✅ Implemented | `internal/state/state.go` — `*time.Time` with `omitempty`; carried in `MergeAgents` |
| 6h cooldown gate | ✅ Implemented | `internal/update/cooldown.go` — `UpdateCheckTTL = 6 * time.Hour`; skip when elapsed < TTL; refresh only on success |
| Cooldown wired into TUI init and CLI self-update | ✅ Implemented | `internal/tui/model.go:675`, `internal/app/selfupdate.go:110` both call `update.CheckAllWithCooldown` |
| `GENTLE_AI_CONFIRM_UPDATE` removed | ✅ Implemented | `internal/app/selfupdate.go` — constant removed; prompt called unconditionally; tests confirm env ignored |
| `[Y/n]` default prompt | ✅ Implemented | `internal/app/selfupdate.go:58` prints `[Y/n]`; empty/Y/Yes → true; N/No → false |
| `--yes` / `GENTLE_AI_YES=1` auto-accept | ✅ Implemented | `internal/app/selfupdate.go:41`, `:136` — replaces `promptFn` with auto-accept stub |
| `GENTLE_AI_CHANNEL` honored for engram | ✅ Implemented | `internal/update/upgrade/strategy.go:598` resolves channel; beta → `engramBetaInstallFn` (`go install @main`), stable → `engramDownloadFn` |
| `pending_sync` flag | ✅ Implemented | `internal/state/state.go:93`; set in `selfupdate.go:184` and TUI `startUpgradeSync`; cleared in `app.go:150` |
| Deferred sync on next launch | ✅ Implemented | `internal/app/app.go:145` — reads state, runs `deferredSyncFn`, clears flag on success, leaves on failure |
| Converged close-and-restart | ✅ Implemented | `internal/app/selfupdate.go:201` — prints restart message and returns on every OS; no re-exec |
| TUI pre-Welcome update prompt | ✅ Implemented | `internal/tui/model.go:365` `ScreenUpdatePrompt`; transitions in `UpdateCheckResultMsg`; dedicated `handleUpdatePromptKey` |
| Advisory manifest fetch | ✅ Implemented | `internal/update/advisory.go` — 2s timeout dedicated client, fail-open, URL var for tests; `internal/tui/model.go:686` launched in `Init` |
| Advisory display on Welcome | ✅ Implemented | `internal/tui/model.go:961` appends `"Advisory: " + message` to banner; sanitized on store |
| No version gating anywhere | ✅ Implemented | No version comparison in advisory path; no forced update gate |

### Coherence (Design)
| Decision | Followed? | Notes |
|----------|-----------|-------|
| TUI pre-Welcome prompt as new `Screen` enum | ✅ Yes | `ScreenUpdatePrompt` added after other screens; `View()`, `Update()`, key handlers, `optionCount()`, `router.go` all extended |
| Converge both OSes on close-and-reopen | ✅ Yes | `restartAfterGentleAIUpgrade` prints guidance and returns on all platforms; Unix re-exec removed |
| CLI prompt default, TTY-gated, `--yes` opt-in | ✅ Yes | `defaultPromptForUpdate` declines on non-TTY; `selfUpdateYesFn` reads `GENTLE_AI_YES=1`; no `GENTLE_AI_CONFIRM_UPDATE` |
| 6h TTL in `state.json`, refresh-on-success-only | ✅ Yes | `UpdateCheckTTL = 6h`; `LastUpdateCheck` updated only when `checkSucceeded(results)` is true |
| `pending_sync` flag drives deferred sync | ✅ Yes | Set before exit on self-upgrade; cleared after successful deferred sync; retry on failure |
| Advisory manifest = release asset, async, fail-open | ✅ Partial | Endpoint uses dedicated `advisory` tag asset (`.../download/advisory/advisory.json`) rather than `latest` release asset as design initially suggested. Implementation matches the open question resolution in design §Open Questions and is functionally equivalent (owner-controlled, CDN-backed). |
| Channel-honoring engram upgrade | ✅ Yes | `ResolveInstallChannel` is single source of truth; `DownloadLatestBinary(profile, isBeta)` supports beta `@main` vs stable release download |
| 7-slice order, TUI prompt stands alone | ✅ Yes | Slices shipped as independent PRs #875–#880; TUI prompt was its own PR (#880) |

### TDD Compliance (Strict TDD active)
| Check | Result | Details |
|-------|--------|---------|
| TDD Evidence reported | ⚠️ Skipped | `apply-progress` artifact does not exist for this change; RED/GREEN/TRIANGULATE/SAFETY-NET/REFACTOR table cannot be verified. |
| All tasks have tests | ✅ Yes | Every implementation task has corresponding test files (verified statically). |
| RED confirmed (tests exist) | ✅ Yes | Test files for all slices exist: `cooldown_test.go`, `strategy_test.go`, `selfupdate_test.go`, `model_test.go`, `advisory_test.go`, `state_test.go`, `app_test.go`, `download_test.go`, `channel_test.go`. |
| GREEN confirmed (tests pass) | ➖ Not executed | Static verify did not run `go test ./...`; runtime evidence pending. |
| Triangulation adequate | ✅ Yes | Multiple distinct test cases per behavior (e.g., fresh/stale/missing/failed/future-timestamp for cooldown). |
| Safety net for modified files | ➖ Not verified | No `apply-progress` artifact to cross-reference. |

### Assertion Quality
✅ All assertions verify real behavior — no tautologies, ghost loops, or smoke-test-only patterns observed in the reviewed test files.

### Quality Metrics
**Linter**: ➖ Not executed (static-evidence-only verify)
**Type Checker**: ➖ Not executed (static-evidence-only verify)

### Issues Found
**CRITICAL**: None

**WARNING**:
- `apply-progress` artifact is missing. Strict TDD mode expects a TDD Cycle Evidence table to validate RED/GREEN/SAFETY-NET columns; this verify could not do so. This is a process gap, not an implementation gap.
- Design document originally proposed hosting `advisory.json` on the `latest` release asset; implementation uses a dedicated `advisory` tag/release. This is a documented open-question resolution and is functionally equivalent, but it is a design deviation from the first-draft wording.
- This verify did not execute `go test ./...`; static compliance must be confirmed by a runtime test run before archive readiness.

**SUGGESTION**:
- Run `go test ./...` and, if available, `go test ./... -cover` to convert the static matrix to runtime evidence and report changed-file coverage.
- Add an `apply-progress` artifact in future apply phases when Strict TDD is active so the verify phase can cross-reference TDD cycle evidence.

### Verdict
**PASS WITH WARNINGS**

All implementation tasks are satisfied by static evidence (merged PRs #875–#880 and source/test inspection), all spec scenarios have covering tests, and the design is followed with one minor endpoint-URL deviation that is functionally equivalent. Warnings are procedural (missing `apply-progress`, no runtime test execution). A follow-up runtime verify is recommended to reach full Strict TDD compliance before archiving.
