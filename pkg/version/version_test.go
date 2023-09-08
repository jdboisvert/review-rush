package version

import (
	"fmt"
	"runtime"
	"testing"
)

func TestGitCommit(t *testing.T) {
	if GitCommit != "" {
		t.Error("Expected GitCommit to be empty on initial build")
	}
}

func TestVersion(t *testing.T) {
	if Version == "" {
		t.Error("Expected Version to be non-empty")
	}
}

func TestBuildDate(t *testing.T) {
	if BuildDate != "" {
		t.Error("Expected BuildDate to be empty on initial build")
	}
}

func TestGoVersion(t *testing.T) {
	if GoVersion == "" {
		t.Error("Expected GoVersion to be non-empty")
	}
}

func TestOsArch(t *testing.T) {
	expected := fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)
	if OsArch != expected {
		t.Errorf("Expected OsArch to be %q, but got %q", expected, OsArch)
	}
}
