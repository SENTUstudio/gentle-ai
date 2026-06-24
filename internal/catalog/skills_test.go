package catalog

import (
	"testing"

	"github.com/gentleman-programming/gentle-ai/internal/components/skills"
	"github.com/gentleman-programming/gentle-ai/internal/model"
)

// TestMVPSkillsCoverAllPresetSkills ensures every skill that presets.go would
// install is also registered in the catalog's mvpSkills allowlist. This
// prevents a future addition to sddSkills or foundationSkills from being
// silently rejected by normalizeSkills in cli/validate.go.
func TestMVPSkillsCoverAllPresetSkills(t *testing.T) {
	catalogSet := make(map[model.SkillID]bool)
	for _, s := range MVPSkills() {
		catalogSet[s.ID] = true
	}

	presetSkills := skills.AllSkillIDs()
	for _, id := range presetSkills {
		if !catalogSet[id] {
			t.Errorf("skill %q is in presets but missing from catalog mvpSkills", id)
		}
	}
}

// TestMVPSkillsNoDuplicates ensures no skill is listed twice in mvpSkills.
func TestMVPSkillsNoDuplicates(t *testing.T) {
	seen := make(map[model.SkillID]bool)
	for _, s := range MVPSkills() {
		if seen[s.ID] {
			t.Errorf("duplicate skill %q in mvpSkills", s.ID)
		}
		seen[s.ID] = true
	}
}

func TestMVPSkillsIncludeRequestedBundledSkillsWithCanonicalNames(t *testing.T) {
	required := map[model.SkillID]string{
		model.SkillCreator:       "skill-creator",
		model.SkillSkillRegistry: "skill-registry",
		model.SkillCognitiveDoc:  "cognitive-doc-design",
		model.SkillCommentWriter: "comment-writer",
		model.SkillJudgmentDay:   "judgment-day",
		model.SkillSDDInit:       "sdd-init",
		model.SkillImprover:      "skill-improver",
	}

	found := make(map[model.SkillID]string)
	for _, skill := range MVPSkills() {
		found[skill.ID] = skill.Name
		if skill.Name == "judgement-day" {
			t.Fatalf("catalog uses non-canonical spelling %q; want judgment-day", skill.Name)
		}
	}

	for id, wantName := range required {
		name, ok := found[id]
		if !ok {
			t.Fatalf("MVPSkills() missing requested bundled skill %q", id)
		}
		if name != wantName {
			t.Fatalf("MVPSkills() name for %q = %q, want %q", id, name, wantName)
		}
	}
}

// TestDataEngineerSkillsRegistered ensures the 8 data-engineer skills are
// registered in the catalog under category "data-engineering" and priority
// "p1". Category/priority are part of the contract surfaced to installers —
// drifting them silently would break domain-aware routing downstream.
func TestDataEngineerSkillsRegistered(t *testing.T) {
	required := map[model.SkillID]struct {
		name     string
		category string
		priority string
	}{
		model.SkillDataEngineerPatternDetect: {"data-engineer-pattern-detect", "data-engineering", "p1"},
		model.SkillDataEngineerStudyFile:     {"data-engineer-study-file", "data-engineering", "p1"},
		model.SkillDataEngineerETLS3:         {"data-engineer-etl-s3", "data-engineering", "p1"},
		model.SkillDataEngineerETLGlue:       {"data-engineer-etl-glue", "data-engineering", "p1"},
		model.SkillDataEngineerETLSharepoint: {"data-engineer-etl-sharepoint", "data-engineering", "p1"},
		model.SkillDataEngineerCreateTable:   {"data-engineer-create-table", "data-engineering", "p1"},
		model.SkillDataEngineerSQLFromLogic:  {"data-engineer-sql-from-logic", "data-engineering", "p1"},
		model.SkillDataEngineerIntegrate:     {"data-engineer-integrate", "data-engineering", "p1"},
	}

	found := make(map[model.SkillID]Skill)
	for _, s := range MVPSkills() {
		found[s.ID] = s
	}

	if len(required) != 8 {
		t.Fatalf("expected 8 required data-engineer skills, got %d", len(required))
	}

	for id, want := range required {
		s, ok := found[id]
		if !ok {
			t.Fatalf("MVPSkills() missing data-engineer skill %q", id)
		}
		if s.Name != want.name {
			t.Fatalf("MVPSkills() Name for %q = %q, want %q", id, s.Name, want.name)
		}
		if s.Category != want.category {
			t.Fatalf("MVPSkills() Category for %q = %q, want %q", id, s.Category, want.category)
		}
		if s.Priority != want.priority {
			t.Fatalf("MVPSkills() Priority for %q = %q, want %q", id, s.Priority, want.priority)
		}
	}
}
