package engine

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"create-harness-app/internal/config"
	"create-harness-app/internal/template"
)

type ScaffolderResult struct {
	TargetDir    string
	FilesCreated int
	Phases       []string
}

func Scaffold(cfg *config.Config, tmpl *template.ResolvedTemplate) (*ScaffolderResult, error) {
	targetAbsPath, err := filepath.Abs(cfg.AppName)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	if err := os.MkdirAll(targetAbsPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %w", targetAbsPath, err)
	}

	fmt.Printf("\n🚀 [create-harness-app] SDLC 하네스 스캐폴딩 시작: %s (블루프린트: %s)\n", targetAbsPath, tmpl.Blueprint.Name)

	result := &ScaffolderResult{
		TargetDir: targetAbsPath,
	}

	phaseMap := make(map[string]bool)

	for _, wf := range tmpl.Blueprint.Workflows {
		phaseDir := filepath.Join(targetAbsPath, wf.Dir)
		if err := os.MkdirAll(phaseDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create phase dir %s: %w", phaseDir, err)
		}
		phaseMap[wf.Dir] = true

		for _, node := range wf.Nodes {
			filePath := filepath.Join(wf.Dir, node.File)
			fullTargetPath := filepath.Join(targetAbsPath, filePath)

			// Ensure parent dir exists
			if err := os.MkdirAll(filepath.Dir(fullTargetPath), 0755); err != nil {
				return nil, fmt.Errorf("failed to create parent dir: %w", err)
			}

			// Read file content from template FS
			srcPath := filepath.Join("files", filePath)
			content, err := fs.ReadFile(tmpl.FS, srcPath)
			if err != nil {
				// Fallback to placeholder if template file isn't physically on disk yet
				content = []byte(fmt.Sprintf("# %s\n\n%s\n", node.File, node.Description))
			}

			if err := os.WriteFile(fullTargetPath, content, 0644); err != nil {
				return nil, fmt.Errorf("failed to write file %s: %w", fullTargetPath, err)
			}

			fmt.Printf("  ✓ 생성: %s\n", filePath)
			result.FilesCreated++
		}
	}

	for p := range phaseMap {
		result.Phases = append(result.Phases, p)
	}

	fmt.Printf("\n✅ 성공적으로 %d개의 SDLC 스펙 & 하네스 파일이 생성되었습니다.\n", result.FilesCreated)
	return result, nil
}
