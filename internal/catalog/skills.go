package catalog

import "github.com/gentleman-programming/gentle-ai/internal/model"

type Skill struct {
	ID       model.SkillID
	Name     string
	Category string
	Priority string
}

var mvpSkills = []Skill{
	// SDD skills
	{ID: model.SkillSDDInit, Name: "sdd-init", Category: "sdd", Priority: "p0"},

	{ID: model.SkillSDDApply, Name: "sdd-apply", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDVerify, Name: "sdd-verify", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDExplore, Name: "sdd-explore", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDPropose, Name: "sdd-propose", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDSpec, Name: "sdd-spec", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDDesign, Name: "sdd-design", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDTasks, Name: "sdd-tasks", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDArchive, Name: "sdd-archive", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDOnboard, Name: "sdd-onboard", Category: "sdd", Priority: "p0"},
	// Foundation skills
	{ID: model.SkillGoTesting, Name: "go-testing", Category: "testing", Priority: "p0"},
	{ID: model.SkillCreator, Name: "skill-creator", Category: "workflow", Priority: "p0"},
	{ID: model.SkillJudgmentDay, Name: "judgment-day", Category: "workflow", Priority: "p0"},
	{ID: model.SkillBranchPR, Name: "branch-pr", Category: "workflow", Priority: "p0"},
	{ID: model.SkillIssueCreation, Name: "issue-creation", Category: "workflow", Priority: "p0"},
	{ID: model.SkillSkillRegistry, Name: "skill-registry", Category: "workflow", Priority: "p0"},
	// Data engineering skills
	{ID: model.SkillDataEngineerStudyFile, Name: "data-engineer-study-file", Category: "data-engineering", Priority: "p0"},
	{ID: model.SkillDataEngineerSQLFromLogic, Name: "data-engineer-sql-from-logic", Category: "data-engineering", Priority: "p0"},
	{ID: model.SkillDataEngineerETLS3, Name: "data-engineer-etl-s3", Category: "data-engineering", Priority: "p0"},
	{ID: model.SkillDataEngineerETLGlue, Name: "data-engineer-etl-glue", Category: "data-engineering", Priority: "p0"},
	{ID: model.SkillDataEngineerETLSharePoint, Name: "data-engineer-etl-sharepoint", Category: "data-engineering", Priority: "p0"},
	{ID: model.SkillDataEngineerCreateTable, Name: "data-engineer-create-table", Category: "data-engineering", Priority: "p0"},
	{ID: model.SkillDataEngineerIntegrate, Name: "data-engineer-integrate", Category: "data-engineering", Priority: "p0"},
}

func MVPSkills() []Skill {
	skills := make([]Skill, len(mvpSkills))
	copy(skills, mvpSkills)
	return skills
}
