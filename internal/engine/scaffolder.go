package engine

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"create-harness-app/internal/config"
	"create-harness-app/internal/harness"
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

	// 1. Initialize Harness State (harness/state.json)
	if err := harness.InitState(tmpl.Blueprint, targetAbsPath); err != nil {
		return nil, fmt.Errorf("failed to initialize harness state: %w", err)
	}
	fmt.Println("  ✓ 생성: harness/state.json (결정적 동적 상태 원본)")

	// 2. Clean 1:1 Copy of all template files and directories into target project
	_ = fs.WalkDir(tmpl.FS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || path == "." || path == "blueprint.json" || path == "files" {
			return nil
		}

		targetRelPath := strings.TrimPrefix(path, "files/")
		fullTargetPath := filepath.Join(targetAbsPath, targetRelPath)

		if d.IsDir() {
			return os.MkdirAll(fullTargetPath, 0755)
		}

		if err := os.MkdirAll(filepath.Dir(fullTargetPath), 0755); err != nil {
			return err
		}

		content, readErr := fs.ReadFile(tmpl.FS, path)
		if readErr == nil {
			if writeErr := os.WriteFile(fullTargetPath, content, 0644); writeErr == nil {
				fmt.Printf("  ✓ 생성: %s\n", targetRelPath)
				result.FilesCreated++
			}
		}
		return nil
	})

	fmt.Printf("\n✅ 성공적으로 %d개의 SDLC 스펙 & 하네스 파일이 생성되었습니다.\n", result.FilesCreated)
	return result, nil
}
