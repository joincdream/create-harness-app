package template

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"create-harness-app/internal/blueprint"
	"create-harness-app/templates"
)

// ResolvedTemplate holds loaded filesystem and blueprint manifest
type ResolvedTemplate struct {
	FS        fs.FS
	Blueprint *blueprint.Blueprint
}

// ResolveTemplate resolves blueprint using 3-step fallback chain:
// 1. Custom Flag Path / Local Path
// 2. User Home ~/.config/create-harness-app/templates/[name]
// 3. Embedded FS templates/[name]
func ResolveTemplate(templateName string, customPath string) (*ResolvedTemplate, error) {
	// Step 1: Explicit Custom Path
	if customPath != "" {
		if _, err := os.Stat(customPath); err == nil {
			return loadFromDisk(customPath)
		}
	}

	// Step 2: User Home Directory (~/.config/create-harness-app/templates/[name])
	homeDir, err := os.UserHomeDir()
	if err == nil {
		userCustomPath := filepath.Join(homeDir, ".config", "create-harness-app", "templates", templateName)
		if _, err := os.Stat(userCustomPath); err == nil {
			return loadFromDisk(userCustomPath)
		}
	}

	// Step 3: Default Embedded FS
	subPath := fmt.Sprintf("%s", templateName)
	subFS, err := fs.Sub(templates.EmbeddedFS, subPath)
	if err != nil {
		return nil, fmt.Errorf("embedded template '%s' not found: %w", templateName, err)
	}

	bpData, err := fs.ReadFile(subFS, "blueprint.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read blueprint.json from embedded template: %w", err)
	}

	bp, err := blueprint.ParseBlueprint(bpData)
	if err != nil {
		return nil, err
	}

	return &ResolvedTemplate{
		FS:        subFS,
		Blueprint: bp,
	}, nil
}

func loadFromDisk(dirPath string) (*ResolvedTemplate, error) {
	diskFS := os.DirFS(dirPath)
	bpData, err := fs.ReadFile(diskFS, "blueprint.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read blueprint.json from %s: %w", dirPath, err)
	}

	bp, err := blueprint.ParseBlueprint(bpData)
	if err != nil {
		return nil, err
	}

	return &ResolvedTemplate{
		FS:        diskFS,
		Blueprint: bp,
	}, nil
}
