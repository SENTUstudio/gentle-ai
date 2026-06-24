package model

import "testing"

func TestDefaultModelsForDomainAppDev(t *testing.T) {
	defaults := DefaultModelsForDomain("")
	if defaults["sdd-explore"].ModelID != "minimax-m3" {
		t.Errorf("app-dev explore = %s, want minimax-m3", defaults["sdd-explore"].ModelID)
	}
	if defaults["sdd-spec"].ModelID != "kimi-k2.7-code" {
		t.Errorf("app-dev spec = %s, want kimi-k2.7-code", defaults["sdd-spec"].ModelID)
	}
	if defaults["sdd-design"].ModelID != "minimax-m3" {
		t.Errorf("app-dev design = %s, want minimax-m3", defaults["sdd-design"].ModelID)
	}
}

func TestDefaultModelsForDomainDataEngineering(t *testing.T) {
	defaults := DefaultModelsForDomain("data-engineering")
	// Data-eng needs glm-5.2 for explore (not minimax-m3) — data profiling needs heavy reasoning
	if defaults["sdd-explore"].ModelID != "glm-5.2" {
		t.Errorf("data-eng explore = %s, want glm-5.2", defaults["sdd-explore"].ModelID)
	}
	// Data-eng needs minimax-m3 for spec (not kimi) — schema + DAG needs reasoning
	if defaults["sdd-spec"].ModelID != "minimax-m3" {
		t.Errorf("data-eng spec = %s, want minimax-m3", defaults["sdd-spec"].ModelID)
	}
	// Design stays minimax-m3 for both
	if defaults["sdd-design"].ModelID != "minimax-m3" {
		t.Errorf("data-eng design = %s, want minimax-m3", defaults["sdd-design"].ModelID)
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
	// explore MUST differ: app-dev=minimax-m3, data-eng=glm-5.2
	if appDev["sdd-explore"].ModelID == dataEng["sdd-explore"].ModelID {
		t.Error("explore model should differ between app-dev and data-engineering")
	}
	// spec MUST differ: app-dev=kimi, data-eng=minimax-m3
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
