package template

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"create-harness-app/internal/blueprint"
	"create-harness-app/internal/hub"
	"create-harness-app/templates"
)

// ResolvedTemplate holds loaded filesystem and blueprint manifest
type ResolvedTemplate struct {
	FS        fs.FS
	Blueprint *blueprint.Blueprint
}

// ResolveTemplate resolves blueprint using 4-step fallback chain:
// 1. Custom Flag Path / Local Path
// 2. User Home ~/.config/create-harness-app/templates/[name]
// 3. Embedded FS templates/[name]
// 4. HarnessHub Registry Auto-Pull (if not found locally, like docker run)
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
	if err == nil {
		if bpData, err := fs.ReadFile(subFS, "blueprint.json"); err == nil {
			if bp, err := blueprint.ParseBlueprint(bpData); err == nil {
				return &ResolvedTemplate{
					FS:        subFS,
					Blueprint: bp,
				}, nil
			}
		}
	}

	// Step 4: Auto-Pull from HarnessHub Registry (like docker run pulling missing image)
	fmt.Printf("📥 Template '%s' not found locally. Auto-pulling from HarnessHub...\n", templateName)
	client := hub.NewClient("")
	ctx := context.Background()

	// Query template info first to check latest version
	targetVer := "v1.0.0"
	detail, detailErr := client.GetTemplateDetail(ctx, templateName, "")
	if detailErr == nil && detail.Version != "" {
		targetVer = detail.Version
	}

	stream, err := client.DownloadTemplate(ctx, templateName, targetVer)
	if err == nil {
		defer stream.Close()
		destDir, installErr := InstallFromTarGz(stream, templateName)
		if installErr == nil {
			fmt.Printf("✅ Successfully auto-pulled template '%s' (%s) to %s\n", templateName, targetVer, destDir)
			return loadFromDisk(destDir)
		}
	}

	// Fallback to blueprint json if stream fails
	if detailErr == nil && detail.BlueprintJSON != "" {
		destDir, saveErr := SaveRawBlueprint(templateName, detail.BlueprintJSON)
		if saveErr == nil {
			fmt.Printf("✅ Successfully auto-pulled blueprint '%s' to %s\n", templateName, destDir)
			return loadFromDisk(destDir)
		}
	}

	return nil, fmt.Errorf("template '%s' not found in local cache, embedded templates, or HarnessHub registry", templateName)
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

// TemplateInfo contains metadata for a template
type TemplateInfo struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Author      string   `json:"author,omitempty"`
	License     string   `json:"license,omitempty"`
	Keywords    []string `json:"keywords,omitempty"`
	Source      string   `json:"source"` // "embedded" or "custom"
}

// ListTemplates scans embedded and user home custom template directories and returns available template infos
func ListTemplates() ([]TemplateInfo, error) {
	templatesMap := make(map[string]TemplateInfo)

	// 1. Scan EmbeddedFS
	entries, err := fs.ReadDir(templates.EmbeddedFS, ".")
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				name := entry.Name()
				bpData, err := fs.ReadFile(templates.EmbeddedFS, filepath.Join(name, "blueprint.json"))
				info := TemplateInfo{
					Name:   name,
					Source: "embedded",
				}
				if err == nil {
					if bp, parseErr := blueprint.ParseBlueprint(bpData); parseErr == nil {
						info.Name = bp.Name
						info.Version = bp.Version
						info.Description = bp.Description
						if bp.Author != nil {
							info.Author = bp.Author.Name
						}
						info.License = bp.License
						info.Keywords = bp.Keywords
					}
				}
				templatesMap[name] = info
			}
		}
	}

	// 2. Scan User Home Directory (~/.config/create-harness-app/templates)
	homeDir, err := os.UserHomeDir()
	if err == nil {
		userTemplatesDir := filepath.Join(homeDir, ".config", "create-harness-app", "templates")
		userEntries, err := os.ReadDir(userTemplatesDir)
		if err == nil {
			for _, entry := range userEntries {
				if entry.IsDir() {
					name := entry.Name()
					bpData, err := os.ReadFile(filepath.Join(userTemplatesDir, name, "blueprint.json"))
					info := TemplateInfo{
						Name:   name,
						Source: "custom",
					}
					if err == nil {
						if bp, parseErr := blueprint.ParseBlueprint(bpData); parseErr == nil {
							info.Name = bp.Name
							info.Version = bp.Version
							info.Description = bp.Description
							if bp.Author != nil {
								info.Author = bp.Author.Name
							}
							info.License = bp.License
							info.Keywords = bp.Keywords
						}
					}
					templatesMap[name] = info
				}
			}
		}
	}

	var result []TemplateInfo
	for _, info := range templatesMap {
		result = append(result, info)
	}

	return result, nil
}

