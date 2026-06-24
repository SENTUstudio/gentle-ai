# Delta: sdd-profiles (MODIFIED)

## MODIFIED Requirements

### Requirement: Domain-Aware Profile

A Profile MUST carry an optional `Domain` field. When `Domain` is empty, the profile behaves as app-dev (today's default). When `Domain` is "data-engineering", `DefaultModelsForDomain` MUST return domain-appropriate model recommendations per phase. The TUI MUST auto-detect the domain from `sdd-config` and prefill recommendations. `DetectProfiles` MUST reject two profiles with the same name but different domains.

#### Scenario: Data-engineering profile gets higher-tier defaults

- GIVEN a project with `domain: data-engineering` in openspec/config.yaml
- WHEN the TUI creates a profile named "cheap"
- THEN `DefaultModelsForDomain("data-engineering")` returns Sonnet for explore (not Haiku)
- AND returns Opus for spec (not Sonnet)

#### Scenario: App-dev profile unchanged

- GIVEN a project with no domain field
- WHEN the TUI creates a profile named "cheap"
- THEN `DefaultModelsForDomain("")` returns Haiku for explore
- AND returns Sonnet for spec

#### Scenario: Collision guard

- GIVEN opencode.json contains `sdd-orchestrator-cheap` with domain "app-dev"
- WHEN DetectProfiles finds another `sdd-orchestrator-cheap` with domain "data-engineering"
- THEN DetectProfiles returns an error

#### Scenario: Backward compatibility

- GIVEN an existing opencode.json with profiles that have no domain
- WHEN DetectProfiles reads them
- THEN each profile gets Domain="" (app-dev)
- AND no behavior changes
