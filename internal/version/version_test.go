package version

import (
	"testing"
)

func TestGetVersion(t *testing.T) {
	version := GetVersion()
	if version == "" {
		t.Error("GetVersion() should not return empty string")
	}
}

func TestGetFullVersion(t *testing.T) {
	fullVersion := GetFullVersion()
	if fullVersion == "" {
		t.Error("GetFullVersion() should not return empty string")
	}

	// フォーマットが正しいか確認（Version (GitCommit)）
	expected := Version + " (" + GitCommit + ")"
	if fullVersion != expected {
		t.Errorf("GetFullVersion() = %q, want %q", fullVersion, expected)
	}
}

func TestVersionVariables(t *testing.T) {
	// デフォルト値が設定されていることを確認
	if Version == "" {
		t.Error("Version should have a default value")
	}
	if GitCommit == "" {
		t.Error("GitCommit should have a default value")
	}
	if BuildDate == "" {
		t.Error("BuildDate should have a default value")
	}
}
