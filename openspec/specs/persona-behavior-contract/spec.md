# persona-behavior-contract Specification

## Requirements

### Requirement: Neutral Mentor Behavior Parity

The system MUST treat `neutral` as a level-neutral variant of the Gentleman mentor behavior contract. Neutral persona content MUST preserve the same senior mentor expectations as Gentleman, including concise answers, direct correction after verification, concept-first teaching, careful technical reasoning, and user-growth-oriented guidance, while MUST NOT include Rioplatense Spanish, regional slang, voseo, Gentleman branding, or any persona-specific regional voice.

#### Scenario: Neutral receives the same mentor contract without regional voice

- GIVEN an agent persona asset is rendered with persona `neutral`
- WHEN the generated instruction content is inspected
- THEN it includes the same mentor behavior expectations as Gentleman for brevity, verification, concept-first explanation, and constructive correction
- AND it does not include Rioplatense Spanish, regional slang, voseo, Gentleman branding, or regional persona voice instructions

#### Scenario: Gentleman keeps regional mentor behavior when explicitly selected

- GIVEN an agent persona asset is rendered with persona `gentleman`
- WHEN the generated instruction content is inspected
- THEN it preserves the Gentleman mentor behavior contract
- AND it preserves the Gentleman regional voice constraints

---

### Requirement: Neutral Interaction Discipline

The neutral persona contract MUST require disciplined interaction defaults across supported agent consumers: short default answers, at most one question at a time, stopping after asking a question, no option menus or exhaustive alternatives unless a real tradeoff exists, and verification before accepting or correcting a user claim.

#### Scenario: Neutral defaults to brief replies

- GIVEN a neutral persona instruction is installed for an agent
- WHEN the agent answers a normal user request that does not require extensive detail
- THEN the instruction requires the minimum useful response
- AND it permits expansion only when the user asks or the task genuinely requires it

#### Scenario: Neutral asks one question and stops

- GIVEN a neutral persona instruction needs clarification from the user
- WHEN it asks a question
- THEN it asks at most one question
- AND it instructs the agent to stop and wait for the user's answer

#### Scenario: Neutral avoids unnecessary menus

- GIVEN a neutral persona instruction describes how to present alternatives
- WHEN there is no real fork with meaningful tradeoffs
- THEN it prohibits option menus, exhaustive lists, and multiple approaches by default

#### Scenario: Neutral verifies before agreeing or correcting

- GIVEN a user makes a technical claim
- WHEN a neutral persona instruction governs the response
- THEN it requires verification against code, docs, or other evidence before agreeing with the claim
- AND it requires explaining why the claim is wrong when evidence disproves it

---

### Requirement: Direct conversation follows the active persona

Direct user and orchestrator conversation MUST be governed by the active persona. Persona rules MUST apply to conversational replies, clarification prompts, and user-facing orchestration status, but MUST NOT be treated as the default language or regional style for generated technical artifacts.

#### Scenario: Gentleman governs direct conversation

- GIVEN the active persona is `gentleman`
- WHEN the agent replies directly to the user or orchestrator
- THEN the reply MUST preserve the Gentleman teaching voice
- AND the reply MUST use the expected Rioplatense senior-architect style when Spanish is used
- AND the reply MAY use voseo, warm/direct phrasing, and concept-before-code explanations

#### Scenario: Neutral governs direct conversation

- GIVEN the active persona is `neutral`
- WHEN the agent replies directly to the user or orchestrator
- THEN the reply MUST keep the same teaching core and clear technical guidance
- AND the reply MUST NOT use Rioplatense regional expressions by default

#### Scenario: Persona does not cross the artifact boundary

- GIVEN the active persona has regional Spanish conversation rules
- WHEN the agent generates a technical artifact
- THEN the artifact language MUST be selected from the artifact language contract
- AND the artifact MUST NOT inherit regional persona phrasing unless the user explicitly requests that regional style for the artifact

---

### Requirement: Technical artifacts default to English

Generated technical artifacts MUST default to English regardless of active persona or conversation language. Technical artifacts SHALL include OpenSpec proposal, spec, design, task, verification, and archive artifacts; SDD phase artifacts; prompt-generated technical files; generated code comments; UI copy; tests; fixtures; and other repository-facing files unless an explicit user request or project convention requires another language.

#### Scenario: Spanish conversation produces English OpenSpec artifacts

- GIVEN the user and agent are conversing in Spanish
- AND no explicit Spanish artifact request exists
- WHEN the agent writes an OpenSpec artifact
- THEN the artifact MUST be written in English
- AND the artifact MUST NOT include Rioplatense conversational expressions

#### Scenario: Persona-specific voice is excluded from generated technical files

- GIVEN the active persona is `gentleman`
- WHEN the agent writes specs, designs, tasks, generated code comments, UI copy, tests, fixtures, or prompt-generated technical files
- THEN those files MUST default to English
- AND they MUST NOT use voseo or regional Spanish terms from the Gentleman persona unless explicitly requested for that artifact

#### Scenario: Project convention can require a non-English artifact

- GIVEN a repository convention clearly requires a technical artifact in Spanish
- WHEN the agent writes that artifact
- THEN the artifact MAY be written in Spanish
- AND the artifact MUST follow the Spanish technical artifact contract

---

### Requirement: Spanish technical artifacts use neutral professional Spanish

When technical artifacts are explicitly requested in Spanish, or when a project convention requires Spanish artifacts, the artifact MUST use neutral/professional Spanish by default. Regional Spanish variants MAY be used only when the user explicitly requests the regional style for that artifact.

#### Scenario: Explicit Spanish artifact request defaults neutral

- GIVEN the user asks for a technical artifact in Spanish
- AND the user does not request a regional variant
- WHEN the agent writes the artifact
- THEN the artifact MUST use neutral/professional Spanish
- AND it MUST avoid Rioplatense expressions such as voseo by default

#### Scenario: Explicit regional artifact request is honored

- GIVEN the user asks for a Spanish technical artifact
- AND the user explicitly requests Rioplatense tone for that artifact
- WHEN the agent writes the artifact
- THEN the artifact MAY use Rioplatense phrasing
- AND the artifact MUST still remain professional and technically precise

---

### Requirement: Comment writer follows target context language

The `comment-writer` skill MUST be context-reactive. It MUST write public or contextual comments in the target context language by default, and an explicit user language or tone override MUST take precedence over inferred context.

#### Scenario: Spanish context yields Spanish comment

- GIVEN the target issue, pull request, review thread, or message is in Spanish
- AND the user does not override the language
- WHEN `comment-writer` drafts the comment
- THEN the comment MUST be written in Spanish

#### Scenario: English context yields English comment

- GIVEN the target issue, pull request, review thread, or message is in English
- AND the user does not override the language
- WHEN `comment-writer` drafts the comment
- THEN the comment MUST be written in English

#### Scenario: Mixed context follows target message language

- GIVEN the surrounding thread is mixed language
- AND a specific target message or audience language is identifiable
- WHEN `comment-writer` drafts the comment
- THEN the comment MUST use the target message or audience language

#### Scenario: User override wins

- GIVEN the target context language is identifiable
- AND the user explicitly requests a different language or tone
- WHEN `comment-writer` drafts the comment
- THEN the comment MUST follow the explicit user request

---

### Requirement: Spanish comments default neutral professional

Spanish comments produced by `comment-writer` MUST default to neutral/professional Spanish. Regional Spanish tone MAY be used only when the user explicitly asks for it or when the surrounding target context clearly calls for that regional tone.

#### Scenario: Spanish comment without regional signal is neutral

- GIVEN `comment-writer` is drafting a Spanish comment
- AND neither the user request nor the target context calls for a regional tone
- WHEN the comment is produced
- THEN the comment MUST use neutral/professional Spanish
- AND it MUST NOT force Rioplatense wording or voseo

#### Scenario: Regional Spanish comment requires a clear signal

- GIVEN `comment-writer` is drafting a Spanish comment
- AND the user explicitly requests Rioplatense tone or the target context clearly uses and expects that tone
- WHEN the comment is produced
- THEN the comment MAY use Rioplatense phrasing
- AND the comment MUST remain appropriate for the public or contextual target

#### Scenario: Root and embedded skill contracts stay aligned

- GIVEN the root `skills/comment-writer/SKILL.md` and embedded `internal/assets/skills/comment-writer/SKILL.md` exist
- WHEN their language behavior rules are inspected
- THEN both skill sources MUST require context-reactive comment language
- AND both skill sources MUST require neutral/professional Spanish by default
- AND neither source MUST force Rioplatense Spanish for all Spanish comments

---

### Requirement: All supported SDD agent assets implement the contract

Every supported SDD orchestrator asset MUST codify the persona/artifact/comment language boundary. Covered assets SHALL include OpenCode, Kilocode through the OpenCode asset path, Claude, Kimi, Codex, Gemini, Qwen, Cursor, Windsurf, Antigravity, Kiro, the generic fallback, and any additional supported agent-specific SDD orchestrator assets discovered in `internal/assets` or the supported agent registry.

#### Scenario: Known supported asset set is covered

- GIVEN the SDD orchestrator assets are reviewed
- WHEN language-contract coverage is evaluated
- THEN coverage MUST include `internal/assets/opencode/sdd-orchestrator.md`
- AND coverage MUST include Kilocode behavior that uses the OpenCode SDD asset path
- AND coverage MUST include `internal/assets/claude/sdd-orchestrator.md`
- AND coverage MUST include `internal/assets/kimi/sdd-orchestrator.md`
- AND coverage MUST include `internal/assets/codex/sdd-orchestrator.md`
- AND coverage MUST include `internal/assets/gemini/sdd-orchestrator.md`
- AND coverage MUST include `internal/assets/qwen/sdd-orchestrator.md`
- AND coverage MUST include `internal/assets/cursor/sdd-orchestrator.md`
- AND coverage MUST include `internal/assets/windsurf/sdd-orchestrator.md`
- AND coverage MUST include `internal/assets/antigravity/sdd-orchestrator.md`
- AND coverage MUST include `internal/assets/kiro/sdd-orchestrator.md`
- AND coverage MUST include `internal/assets/generic/sdd-orchestrator.md`

#### Scenario: Newly discovered supported assets are not skipped

- GIVEN an additional supported agent-specific SDD orchestrator asset exists
- WHEN implementation and tests enumerate language-contract coverage
- THEN the asset MUST be included in the same contract checks
- AND the generic fallback MUST remain covered for unsupported or unrecognized agents

#### Scenario: Persona-specific direct conversation assets remain allowed

- GIVEN asset language guards inspect supported assets
- WHEN Gentleman direct-conversation persona assets are inspected
- THEN tests MAY allow intentional Rioplatense direct-conversation wording in those persona assets
- AND tests MUST still prevent those persona rules from becoming the default for technical artifacts or comments

---

### Requirement: Install and sync preserve the updated language contract

Install and sync flows MUST install, refresh, and regenerate only assets that comply with the updated language contract. They MUST NOT regenerate old persona leaks from stale embedded sources, overlays, shared prompt files, generated agent files, or root skill copies.

#### Scenario: Fresh install does not regenerate old leaks

- GIVEN a fresh install writes SDD and skill assets
- WHEN generated assets are inspected after install
- THEN they MUST contain the updated artifact/comment language contract
- AND they MUST NOT contain the known old leak terms in persona-agnostic SDD assets

#### Scenario: Sync does not restore stale wording

- GIVEN an installation already contains outdated SDD or `comment-writer` wording
- WHEN sync refreshes SDD and skill assets
- THEN the refreshed assets MUST contain the updated language contract
- AND sync MUST NOT restore stale Rioplatense defaults for persona-agnostic artifacts or comments

#### Scenario: OpenCode and Kilocode leak path is guarded

- GIVEN OpenCode or Kilocode SDD assets are installed or synced
- WHEN generated orchestrator prompts, overlays, or shared prompt files are inspected
- THEN the generated content MUST NOT include the old OpenCode leak wording
- AND the generated content MUST preserve the updated artifact/comment contract

---

### Requirement: Delegated prompts forward the artifact and comment contract

Delegated SDD phase prompts and subagent instructions MUST forward the artifact and comment language contract. Delegated agents MUST know that direct conversation may be persona-governed, generated technical artifacts default to English, Spanish technical artifacts default to neutral/professional Spanish, and comments follow the target context language by default.

#### Scenario: Delegated phase prompt receives artifact defaults

- GIVEN the orchestrator delegates an SDD phase to a phase executor or subagent
- WHEN the delegated prompt is constructed
- THEN the prompt MUST state that generated technical artifacts default to English
- AND the prompt MUST state that Spanish technical artifacts use neutral/professional Spanish unless explicitly regional

#### Scenario: Delegated comment work receives comment defaults

- GIVEN the orchestrator delegates comment-writing or review-comment drafting work
- WHEN the delegated prompt is constructed
- THEN the prompt MUST forward that comments use the target context language by default
- AND the prompt MUST forward that Spanish comments default to neutral/professional Spanish unless user or context clearly calls for regional tone

#### Scenario: Delegation does not convert persona into artifact language

- GIVEN the active persona uses Rioplatense direct conversation
- WHEN the orchestrator delegates artifact-writing work
- THEN the delegated prompt MUST NOT instruct the phase executor to write artifacts in Rioplatense Spanish by default
- AND the delegated prompt MUST preserve the artifact language contract

---

### Requirement: Known language leaks are prevented

Persona-agnostic SDD assets, generated technical artifacts, delegated prompts, install/sync outputs, and root or embedded comment-writer skill contracts MUST prevent recurrence of the known leaks: `elegí`, `Respondé`, the hardcoded `¿Querés ajustar algo o continuamos?` continuation in persona-agnostic SDD flows, and root `comment-writer` forcing Rioplatense Spanish.

#### Scenario: Persona-agnostic SDD assets reject known leak terms

- GIVEN persona-agnostic SDD orchestrator assets are inspected
- WHEN language guard tests run
- THEN the assets MUST NOT contain `elegí`
- AND the assets MUST NOT contain `Respondé`
- AND the assets MUST NOT contain `¿Querés ajustar algo o continuamos?`

#### Scenario: Generated artifacts reject known leak terms

- GIVEN install, sync, or SDD artifact generation produces prompt or artifact files
- WHEN generated files are inspected
- THEN persona-agnostic generated outputs MUST NOT contain the known leak terms
- AND any allowed mentions MUST be limited to explicit regression-test assertions or boundary documentation that names the leak as prohibited

#### Scenario: Comment writer no longer forces Rioplatense Spanish

- GIVEN `comment-writer` drafts a Spanish comment without a regional override or clear regional context
- WHEN the comment is produced
- THEN the comment MUST NOT use Rioplatense Spanish solely because the root skill requested it
- AND the root skill MUST NOT contain a rule that forces Rioplatense Spanish for all Spanish comments

---

### Requirement: Claude Neutral Output Style Contract

Claude-specific neutral output-style content MUST be meaningful and MUST NOT fall back to a generic default assistant character. It MUST encode the neutral mentor behavior contract, interaction discipline, verification-first rule, and artifact language independence without regional voice.

#### Scenario: Claude neutral output-style is not default assistant behavior

- GIVEN Claude assets are generated with persona `neutral`
- WHEN the neutral output-style content is inspected
- THEN it contains explicit neutral mentor behavior instructions
- AND it contains brevity, one-question, no-menu, verification-first, and artifact-language constraints
- AND it does not describe or imply an unstyled default assistant character

#### Scenario: Claude explicit Gentleman output-style remains honored

- GIVEN Claude assets are generated with persona `gentleman`
- WHEN the output-style content is inspected
- THEN it preserves Gentleman-specific mentor and regional voice instructions
- AND it is not replaced by neutral output-style content

---

### Requirement: Kimi Neutral Output Style Content

Kimi neutral output-style module content MUST be meaningful, non-empty, and semantically aligned with the generic neutral behavior contract. Empty files, placeholder text, or whitespace-only injected output-style content MUST NOT be accepted for neutral.

#### Scenario: Kimi neutral output-style is meaningful

- GIVEN Kimi assets are generated or injected with persona `neutral`
- WHEN the `output-style.md` content is inspected
- THEN it is non-empty after trimming whitespace
- AND it includes neutral mentor behavior, interaction discipline, verification-first, and artifact-language constraints
- AND it excludes regional Gentleman voice instructions

#### Scenario: Kimi neutral output-style rejects placeholder-only content

- GIVEN the Kimi neutral output-style source contains only placeholder or whitespace-only content
- WHEN the asset is prepared for injection
- THEN the system treats the content as invalid for neutral parity
- AND implementation MUST provide meaningful neutral output-style content instead

---

### Requirement: Generic Neutral Asset Parity

All neutral consumers that are not covered by an agent-specific override MUST receive parity through the generic neutral persona or output-style asset. Agent-specific assets MAY adapt wording to platform mechanics, but MUST NOT weaken the neutral behavior contract.

#### Scenario: Non-agent-specific consumers receive generic neutral parity

- GIVEN an agent or surface consumes the generic neutral persona asset
- WHEN neutral instructions are rendered for that consumer
- THEN the rendered content includes the neutral mentor behavior contract
- AND it includes brevity, one-question, no-menu, verification-first, and artifact-language constraints

#### Scenario: Agent-specific neutral assets do not weaken generic behavior

- GIVEN an agent has its own neutral persona or output-style asset
- WHEN that asset is compared against the generic neutral behavior contract
- THEN it preserves all generic neutral requirements
- AND any agent-specific differences are limited to platform-accurate wording or installation mechanics

---

### Requirement: Safe Persona Fallback Semantics

When persisted persona state is missing, empty, unreadable, or invalid, sync and persona resolution MUST NOT silently select or reactivate `gentleman`. The fallback MUST be neutral/default-safe behavior that does not introduce Gentleman regional voice unless the user explicitly selected Gentleman.

#### Scenario: Missing persisted persona does not reactivate Gentleman

- GIVEN persisted persona state is absent
- WHEN sync resolves the persona to apply
- THEN it does not select `gentleman` implicitly
- AND it applies neutral/default-safe persona behavior without regional voice

#### Scenario: Invalid persisted persona does not reactivate Gentleman

- GIVEN persisted persona state contains an unknown or invalid value
- WHEN sync resolves the persona to apply
- THEN it does not select `gentleman` implicitly
- AND it applies neutral/default-safe persona behavior without regional voice

#### Scenario: Unreadable persisted persona does not reactivate Gentleman

- GIVEN persisted persona state cannot be read
- WHEN sync resolves the persona to apply
- THEN it does not select `gentleman` implicitly
- AND it applies neutral/default-safe persona behavior without regional voice
- AND it may surface a warning if the sync command already reports recoverable configuration issues

---

### Requirement: Explicit Persona Selection Preservation

Explicit persona selections MUST remain authoritative. When the user explicitly selects Gentleman, the system MUST apply Gentleman behavior and regional voice; when the user explicitly selects neutral, the system MUST apply neutral parity behavior without regional voice.

#### Scenario: Explicit Gentleman selection remains honored during sync

- GIVEN the user has explicitly selected persona `gentleman`
- WHEN sync resolves and applies persona assets
- THEN Gentleman persona assets are selected
- AND Gentleman regional voice instructions remain present

#### Scenario: Explicit neutral selection remains honored during sync

- GIVEN the user has explicitly selected persona `neutral`
- WHEN sync resolves and applies persona assets
- THEN neutral persona assets are selected
- AND the rendered content includes neutral parity behavior without regional voice

#### Scenario: Fallback does not override an explicit selection

- GIVEN a valid explicit persona selection exists
- WHEN sync applies persona assets
- THEN fallback logic is not used to replace that explicit selection
- AND the selected persona remains the source of truth for rendered persona content
