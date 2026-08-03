package template

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// GetUserTemplatesDir returns local user templates cache path ~/.config/create-harness-app/templates
func GetUserTemplatesDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, ".config", "create-harness-app", "templates"), nil
}

// InstallFromTarGz extracts .tar.gz stream into ~/.config/create-harness-app/templates/<templateName>
func InstallFromTarGz(r io.Reader, templateName string) (string, error) {
	baseDir, err := GetUserTemplatesDir()
	if err != nil {
		return "", err
	}

	targetDir := filepath.Join(baseDir, templateName)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create target directory %s: %w", targetDir, err)
	}

	gzr, err := gzip.NewReader(r)
	if err != nil {
		// Fallback: If it's not a gzip stream, we don't treat it as tar.gz
		return "", fmt.Errorf("failed to open gzip stream: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("error reading tar archive: %w", err)
		}

		// Sanitize path to prevent Zip Slip / Tar Slip vulnerability
		cleanedPath := filepath.Clean(header.Name)
		if strings.HasPrefix(cleanedPath, "..") {
			continue
		}

		destPath := filepath.Join(targetDir, cleanedPath)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return "", fmt.Errorf("failed to create directory %s: %w", destPath, err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return "", fmt.Errorf("failed to create parent dir for %s: %w", destPath, err)
			}
			outFile, err := os.Create(destPath)
			if err != nil {
				return "", fmt.Errorf("failed to create file %s: %w", destPath, err)
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return "", fmt.Errorf("failed to write file content %s: %w", destPath, err)
			}
			outFile.Close()
		}
	}

	return targetDir, nil
}

// SaveRawBlueprint creates a minimal blueprint template when downloading sample tar in dev mode
func SaveRawBlueprint(templateName, blueprintJSON string) (string, error) {
	baseDir, err := GetUserTemplatesDir()
	if err != nil {
		return "", err
	}

	targetDir := filepath.Join(baseDir, templateName)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create target directory %s: %w", targetDir, err)
	}

	bpPath := filepath.Join(targetDir, "blueprint.json")
	if err := os.WriteFile(bpPath, []byte(blueprintJSON), 0644); err != nil {
		return "", fmt.Errorf("failed to write blueprint.json: %w", err)
	}

	return targetDir, nil
}

// RemoveCachedTemplate deletes downloaded template from ~/.config/create-harness-app/templates/<templateName>
func RemoveCachedTemplate(templateName string) error {
	baseDir, err := GetUserTemplatesDir()
	if err != nil {
		return err
	}

	targetDir := filepath.Join(baseDir, templateName)
	info, err := os.Stat(targetDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("template '%s' is not installed in local cache (%s)", templateName, targetDir)
	}
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("path %s is not a directory", targetDir)
	}

	if err := os.RemoveAll(targetDir); err != nil {
		return fmt.Errorf("failed to remove template directory %s: %w", targetDir, err)
	}

	return nil
}
