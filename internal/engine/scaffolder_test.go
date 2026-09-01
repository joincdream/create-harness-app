package engine_test

import (
	"os"
	"path/filepath"
	"testing"

	"create-harness-app/internal/config"
	"create-harness-app/internal/engine"
	"create-harness-app/internal/template"
)

func TestScaffold(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "harness_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	targetAppDir := filepath.Join(tempDir, "test-app")

	cfg := &config.Config{
		AppName:      targetAppDir,
		TemplateName: "default",
		Force:        true,
	}

	tmpl, err := template.ResolveTemplate("default", "")
	if err != nil {
		t.Fatalf("failed to resolve default template: %v", err)
	}

	result, err := engine.Scaffold(cfg, tmpl)
	if err != nil {
		t.Fatalf("Scaffold failed: %v", err)
	}

	if result.FilesCreated == 0 {
		t.Errorf("expected files to be created, got 0")
	}

	// Verify required planning file exists
	reqFile := filepath.Join(targetAppDir, "harness", "docs", "01_planning", "01_user_requirements.md")
	if _, err := os.Stat(reqFile); os.IsNotExist(err) {
		t.Errorf("expected file %s to exist", reqFile)
	}
}
