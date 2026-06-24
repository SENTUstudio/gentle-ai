package model

type Selection struct {
	Agents                      []AgentID
	Components                  []ComponentID
	Skills                      []SkillID
	Persona                     PersonaID
	Preset                      PresetID
	SDDMode                     SDDModeID
	SDDProfileStrategy          SDDProfileStrategyID
	StrictTDD                   bool
	CodexMultiAgent             bool                             // deprecated: Codex now always writes features.multi_agent = true; retained for state/back-compat
	ModelAssignments            map[string]ModelAssignment       // key = sub-agent name (e.g., "sdd-init")
	ClaudeModelAssignments      map[string]ClaudeModelAlias      // key = phase name; value = fable|opus|sonnet|haiku
	ClaudePhaseAssignments      map[string]ClaudePhaseAssignment // key = phase name; value = Claude model+effort
	KiroModelAssignments        map[string]KiroModelAlias        // key = phase name; value = Kiro-native model alias
	CodexModelAssignments       map[string]CodexEffort           // key = phase name; value = low|medium|high|xhigh
	CodexCarrilModelAssignments map[string]string                // key = carril profile (sdd-strong|sdd-mid|sdd-cheap); value = model id
	CodexPhaseModelAssignments  map[string]string                // key = phase name; value = model id (Custom per-phase picker only)
	Profiles                    []Profile                        // named SDD profiles to generate/update during sync
	OpenCodePlugins             []OpenCodeCommunityPluginID      // optional community OpenCode TUI plugins
}

func (s Selection) HasAgent(agent AgentID) bool {
	for _, current := range s.Agents {
		if current == agent {
			return true
		}
	}

	return false
}

// DefaultModelsForDomain returns recommended model assignments per SDD phase for a given domain.
// Empty domain or "app-dev" returns app-dev defaults. "data-engineering" returns DE defaults.
// Data-engineering needs higher tiers because data profiling and schema design require more judgment.
//
// Model tiers (based on user's opencode-go providers):
//   Light:  mimo-v2.5 (high)         — mechanical tasks (archive)
//   Mid:    minimax-m3 (thinking)    — reasoning tasks (explore, design)
//           kimi-k2.7-code           — coding specialist (spec, tasks, verify)
//   High:   glm-5.2 (max)            — heavy reasoning (propose, apply)
//           deepseek-v4-pro (high)   — init/onboard
func DefaultModelsForDomain(domain string) map[string]ModelAssignment {
	switch domain {
	case "data-engineering":
		return map[string]ModelAssignment{
			"sdd-explore":  {ProviderID: "opencode-go", ModelID: "glm-5.2", Effort: "max"},        // upgraded: data profiling needs heavy reasoning
			"sdd-propose":  {ProviderID: "opencode-go", ModelID: "glm-5.2", Effort: "max"},
			"sdd-spec":     {ProviderID: "opencode-go", ModelID: "minimax-m3", Effort: "thinking"}, // upgraded: schema + DAG needs reasoning
			"sdd-design":   {ProviderID: "opencode-go", ModelID: "minimax-m3", Effort: "thinking"},
			"sdd-tasks":    {ProviderID: "opencode-go", ModelID: "kimi-k2.7-code"},
			"sdd-apply":    {ProviderID: "opencode-go", ModelID: "glm-5.2", Effort: "max"},
			"sdd-verify":   {ProviderID: "opencode-go", ModelID: "kimi-k2.7-code"},
			"sdd-archive":  {ProviderID: "opencode-go", ModelID: "mimo-v2.5", Effort: "high"},
		}
	default: // app-dev or empty — matches user's current opencode.json config
		return map[string]ModelAssignment{
			"sdd-explore":  {ProviderID: "opencode-go", ModelID: "minimax-m3", Effort: "thinking"},
			"sdd-propose":  {ProviderID: "opencode-go", ModelID: "glm-5.2", Effort: "max"},
			"sdd-spec":     {ProviderID: "opencode-go", ModelID: "kimi-k2.7-code"},
			"sdd-design":   {ProviderID: "opencode-go", ModelID: "minimax-m3", Effort: "thinking"},
			"sdd-tasks":    {ProviderID: "opencode-go", ModelID: "kimi-k2.7-code"},
			"sdd-apply":    {ProviderID: "opencode-go", ModelID: "glm-5.2", Effort: "max"},
			"sdd-verify":   {ProviderID: "opencode-go", ModelID: "kimi-k2.7-code"},
			"sdd-archive":  {ProviderID: "opencode-go", ModelID: "mimo-v2.5", Effort: "high"},
		}
	}
}

func (s Selection) HasComponent(component ComponentID) bool {
	for _, current := range s.Components {
		if current == component {
			return true
		}
	}

	return false
}

// SyncOverrides holds optional overrides applied to the sync selection.
// Used when the TUI "Configure Models" flow needs to persist model assignments
// without re-running the full install pipeline.
//
// Nil fields mean "no override" — the sync uses defaults from BuildSyncSelection.
// A non-nil but empty map means "reset to defaults" (explicit clear).
type SyncOverrides struct {
	// TargetAgents forces TUI sync to run the adapter(s) affected by the
	// override, even when persisted install state omits them. This is used by
	// model/profile configurators, where the user picked a concrete target agent.
	TargetAgents                []AgentID
	ModelAssignments            map[string]ModelAssignment       // nil = no override; empty map = reset to defaults
	ClaudeModelAssignments      map[string]ClaudeModelAlias      // nil = no override; empty map = reset to defaults
	ClaudePhaseAssignments      map[string]ClaudePhaseAssignment // nil = no override; empty map = reset to defaults
	KiroModelAssignments        map[string]KiroModelAlias        // nil = no override; empty map = reset to defaults
	CodexModelAssignments       map[string]CodexEffort           // nil = no override; empty map = reset to defaults
	CodexCarrilModelAssignments map[string]string                // nil = no override; empty map = reset to defaults
	CodexPhaseModelAssignments  map[string]string                // nil = no override (partial sync); non-nil empty = clear (preset selected); non-nil non-empty = custom per-phase assignments
	SDDMode                     SDDModeID                        // "" = no override; when non-empty, overrides the sync's default SDD mode
	SDDProfileStrategy          SDDProfileStrategyID             // "" = auto; otherwise explicit sync profile strategy
	StrictTDD                   *bool                            // nil = no override; non-nil = override strict TDD mode
	Profiles                    []Profile                        // NEW: profile creation/updates during sync
}
