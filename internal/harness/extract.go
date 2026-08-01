package harness

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ExtractTemplate copies .agents/ and docs/ assets from sampleAppDir into templates/[templateName]/files/
func ExtractTemplate(sampleAppDir string, templateName string, harnessRepoDir string) error {
	sampleAppAbs, err := filepath.Abs(sampleAppDir)
	if err != nil {
		return fmt.Errorf("invalid sample app dir: %w", err)
	}

	targetTemplateDir := filepath.Join(harnessRepoDir, "templates", templateName, "files")
	if err := os.MkdirAll(targetTemplateDir, 0755); err != nil {
		return fmt.Errorf("failed to create target template dir: %w", err)
	}

	dirsToExtract := []string{".agents", "docs"}
	extractedCount := 0

	for _, dirName := range dirsToExtract {
		srcDir := filepath.Join(sampleAppAbs, dirName)
		if _, err := os.Stat(srcDir); os.IsNotExist(err) {
			continue
		}

		destDir := filepath.Join(targetTemplateDir, dirName)
		if err := copyDir(srcDir, destDir); err != nil {
			return fmt.Errorf("failed to extract dir %s: %w", dirName, err)
		}
		fmt.Printf("  ✓ 역추출 동기화 완료: %s -> templates/%s/files/%s\n", dirName, templateName, dirName)
		extractedCount++
	}

	if extractedCount == 0 {
		return fmt.Errorf("no extractable harness assets (.agents, docs) found in %s", sampleAppDir)
	}

	fmt.Printf("\n🎉 성공적으로 '%s'의 하네스 자산이 'templates/%s'로 역추출 자산화되었습니다.\n", sampleAppDir, templateName)
	return nil
}

func copyDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		return copyFile(path, targetPath)
	})
}

func copyFile(srcFile, dstFile string) error {
	in, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dstFile), 0755); err != nil {
		return err
	}

	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
