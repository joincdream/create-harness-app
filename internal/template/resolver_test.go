package template

import (
	"testing"
)

func TestListTemplates(t *testing.T) {
	templates, err := ListTemplates()
	if err != nil {
		t.Fatalf("ListTemplates failed: %v", err)
	}

	if len(templates) == 0 {
		t.Fatalf("expected at least 1 template, got 0")
	}

	foundDefault := false
	for _, tmpl := range templates {
		if tmpl.Name == "default" {
			foundDefault = true
			if tmpl.Source != "embedded" {
				t.Errorf("expected default template source to be 'embedded', got '%s'", tmpl.Source)
			}
		}
	}

	if !foundDefault {
		t.Errorf("expected to find 'default' template in list")
	}
}
