package buildinfo

import "testing"

func TestCurrentReturnsCanonicalMetadata(t *testing.T) {
	t.Setenv("DINARIA_ENVIRONMENT", "sandbox")
	oldVersion, oldCommit, oldBuiltAt := Version, Commit, BuiltAt
	Version, Commit, BuiltAt = "1.2.3", "abc123", "2026-09-22T00:00:00Z"
	t.Cleanup(func() { Version, Commit, BuiltAt = oldVersion, oldCommit, oldBuiltAt })
	got := Current("connector-example")
	if got.Service != "connector-example" || got.Version != "1.2.3" || got.Commit != "abc123" || got.Environment != "sandbox" {
		t.Fatalf("unexpected metadata: %+v", got)
	}
}
