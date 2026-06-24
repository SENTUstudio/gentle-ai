package etl

// RepoPrefix returns the task annotation for which company repo a task targets.
// "glue-tables" → "infra", "glue-jobs" → "carga", both or unknown → "both".
func RepoPrefix(target string) string {
	switch {
	case contains(target, "glue-tables") && contains(target, "glue-jobs"):
		return "both"
	case contains(target, "glue-tables"):
		return "infra"
	case contains(target, "glue-jobs"):
		return "carga"
	default:
		return "both"
	}
}

// GitFlowForRepo returns the git workflow for a given repo prefix.
// "master" → "github" (GitHub Flow: feature → PR → main).
// "infra" or "carga" → "bitbucket" (GitFlow variant: feature → develop → release).
func GitFlowForRepo(repo string) string {
	switch repo {
	case "master":
		return "github"
	case "infra", "carga":
		return "bitbucket"
	default:
		return "github"
	}
}

// contains is a simple substring check without importing strings (keeps the package lean).
func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
