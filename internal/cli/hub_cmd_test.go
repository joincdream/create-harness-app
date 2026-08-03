package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSubstituteTemplateVariables(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "harness-var-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test description file
	descPath := filepath.Join(tmpDir, "desc.md")
	descContent := "# Test Template Description\n\nThis is a test markdown."
	if err := os.WriteFile(descPath, []byte(descContent), 0644); err != nil {
		t.Fatalf("failed to write test desc file: %v", err)
	}

	stateMap := map[string]interface{}{
		"app_name":    "var-app",
		"version":     "v1.0.0",
		"description": "${file:desc.md}",
		"nodes": map[string]interface{}{
			"test-node": map[string]interface{}{
				"summary": "${file:desc.md}",
			},
		},
	}

	substituteTemplateVariables(stateMap, tmpDir)

	if stateMap["description"] != descContent {
		t.Errorf("expected description %q, got %q", descContent, stateMap["description"])
	}

	nodesMap := stateMap["nodes"].(map[string]interface{})
	testNode := nodesMap["test-node"].(map[string]interface{})
	if testNode["summary"] != descContent {
		t.Errorf("expected node summary %q, got %q", descContent, testNode["summary"])
	}
}
