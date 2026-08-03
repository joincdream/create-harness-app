package template

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndRemoveCachedTemplate(t *testing.T) {
	templateName := "test-remove-template"
	targetDir, err := SaveRawBlueprint(templateName, `{"name":"test-remove-template","version":"1.0.0"}`)
	if err != nil {
		t.Fatalf("SaveRawBlueprint failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(targetDir, "blueprint.json")); err != nil {
		t.Fatalf("expected blueprint.json to exist: %v", err)
	}

	if err := RemoveCachedTemplate(templateName); err != nil {
		t.Fatalf("RemoveCachedTemplate failed: %v", err)
	}

	if _, err := os.Stat(targetDir); !os.IsNotExist(err) {
		t.Fatalf("expected directory to be removed, got err: %v", err)
	}
}
