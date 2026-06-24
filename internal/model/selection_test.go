package model

import "testing"

func TestDefaultModelsForDomainAppDev(t *testing.T) {
	defaults := DefaultModelsForDomain("")
	if defaults["sdd-explore"].ModelID != "claude-haiku-4-5-20250315" {
		t.Errorf("app-dev explore = %s, want claude-haiku-4-5-20250315", defaults["sdd-explore"].ModelID)
	}
	if defaults["sdd-spec"].ModelID != "claude-sonnet-4-20250514" {
		t.Errorf("app-dev spec = %s, want claude-sonnet-4-20250514", defaults["sdd-spec"].ModelID)
	}
	if defaults["sdd-design"].ModelID != "claude-opus-4-20250514" {
		t.Errorf("app-dev design = %s, want claude-opus-4-20250514", defaults["sdd-design"].ModelID)
	}
}

func TestDefaultModelsForDomainDataEngineering(t *testing.T) {
	defaults := DefaultModelsForDomain("data-engineering")
	// Data-eng needs Sonnet for explore (not Haiku) — data profiling needs judgment
	if defaults["sdd-explore"].ModelID != "claude-sonnet-4-20250514" {
		t.Errorf("data-eng explore = %s, want claude-sonnet-4-20250514", defaults["sdd-explore"].ModelID)
	}
	// Data-eng needs Opus for spec (not Sonnet) — schema + DAG needs precision
	if defaults["sdd-spec"].ModelID != "claude-opus-4-20250514" {
		t.Errorf("data-eng spec = %s, want claude-opus-4-20250514", defaults["sdd-spec"].ModelID)
	}
	// Design stays Opus for both
	if defaults["sdd-design"].ModelID != "claude-opus-4-20250514" {
		t.Errorf("data-eng design = %s, want claude-opus-4-20250514", defaults["sdd-design"].ModelID)
	}
}

func TestDefaultModelsForDomainHasAllPhases(t *testing.T) {
	defaults := DefaultModelsForDomain("data-engineering")
	expectedPhases := []string{
		"sdd-explore", "sdd-propose", "sdd-spec", "sdd-design",
		"sdd-tasks", "sdd-apply", "sdd-verify", "sdd-archive",
	}
	for _, phase := range expectedPhases {
		if _, ok := defaults[phase]; !ok {
			t.Errorf("missing phase %s in data-engineering defaults", phase)
		}
	}
}

func TestDefaultModelsForDomainDataEngDiffersFromAppDev(t *testing.T) {
	appDev := DefaultModelsForDomain("")
	dataEng := DefaultModelsForDomain("data-engineering")
	// explore MUST differ: app-dev=Haiku, data-eng=Sonnet
	if appDev["sdd-explore"].ModelID == dataEng["sdd-explore"].ModelID {
		t.Error("explore model should differ between app-dev and data-engineering")
	}
	// spec MUST differ: app-dev=Sonnet, data-eng=Opus
	if appDev["sdd-spec"].ModelID == dataEng["sdd-spec"].ModelID {
		t.Error("spec model should differ between app-dev and data-engineering")
	}
}

func TestProfileDomainFieldBackwardCompat(t *testing.T) {
	// A profile with empty Domain must behave as app-dev
	p := Profile{Name: "cheap"}
	if p.Domain != "" {
		t.Error("new Profile should have empty Domain by default")
	}
}
