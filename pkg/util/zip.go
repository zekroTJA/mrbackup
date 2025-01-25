package util

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func ZipDirectory(sourceDir string, targetFile string) error {
	zipFile, err := os.Create(targetFile)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk source directory: %w", err)
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		zipFile, err := zipWriter.Create(relPath)
		if err != nil {
			return fmt.Errorf("failed to create zip file entry: %w", err)
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open source file: %w", err)
		}
		defer srcFile.Close()

		_, err = io.Copy(zipFile, srcFile)
		if err != nil {
			return fmt.Errorf("failed to copy file to zip archive: %w", err)
		}

		return nil
	})

	return err
}
