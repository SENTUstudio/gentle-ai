package etl

import "testing"

func TestRepoPrefix(t *testing.T) {
	tests := []struct {
		name   string
		target string
		want   string
	}{
		{"infra only", "update glue-tables/mi_tabla.yaml", "infra"},
		{"carga only", "update glue-jobs/etl_foo.py", "carga"},
		{"both", "update glue-tables/x.yaml and glue-jobs/y.py", "both"},
		{"unknown", "update README.md", "both"},
		{"empty", "", "both"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RepoPrefix(tt.target); got != tt.want {
				t.Errorf("RepoPrefix(%q) = %q, want %q", tt.target, got, tt.want)
			}
		})
	}
}

func TestGitFlowForRepo(t *testing.T) {
	tests := []struct {
		repo string
		want string
	}{
		{"master", "github"},
		{"infra", "bitbucket"},
		{"carga", "bitbucket"},
		{"both", "github"},
		{"unknown", "github"},
	}
	for _, tt := range tests {
		t.Run(tt.repo, func(t *testing.T) {
			if got := GitFlowForRepo(tt.repo); got != tt.want {
				t.Errorf("GitFlowForRepo(%q) = %q, want %q", tt.repo, got, tt.want)
			}
		})
	}
}
