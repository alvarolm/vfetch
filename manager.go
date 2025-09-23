package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func removeExisting(path string) error {
	if _, err := os.Lstat(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to check if path exists: %w", err)
	}

	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("failed to remove existing path: %w", err)
	}

	fmt.Printf("Removed existing: %s\n", path)
	return nil
}

func replaceVersionPlaceholders(url, version string) string {
	re := regexp.MustCompile(`(?i)\$version`)
	return re.ReplaceAllString(url, version)
}

func ProcessFetchItem(config *Config, item FetchItem) error {
	finalURL := replaceVersionPlaceholders(item.URL, item.Version)
	fmt.Printf("Downloading: %s\n", finalURL)
	downloadResult, err := DownloadFile(finalURL)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	fmt.Printf("Verifying hash...\n")
	if item.Hash != "" {
		if err := VerifyHash(downloadResult.Data, item.Hash); err != nil {
			return fmt.Errorf("hash verification failed: %w", err)
		}
	} else if len(item.Hashes) > 0 {
		if err := VerifyHashes(downloadResult.Data, item.Hashes); err != nil {
			return fmt.Errorf("hash verification failed: %w", err)
		}
	} else {
		return fmt.Errorf("no hash or hashes specified for verification")
	}

	var filesToWrite []FileToWrite

	if item.Extract {
		fmt.Printf("Extracting archive...\n")
		extractResult, err := ExtractArchive(downloadResult.Data, item.Name)
		if err != nil {
			return fmt.Errorf("extraction failed: %w", err)
		}

		for _, extractedFile := range extractResult.Files {
			// Create files under a directory named after the fetch item
			filePath := filepath.Join(item.Name, extractedFile.Name)
			filesToWrite = append(filesToWrite, FileToWrite{
				Name: filePath,
				Data: extractedFile.Data,
			})
		}
	} else {
		filesToWrite = append(filesToWrite, FileToWrite{
			Name: item.Name,
			Data: downloadResult.Data,
		})
	}

	outputDir := item.GetOutputDir(config.OutputDir)
	if outputDir != "" {
		if err := writeFiles(filesToWrite, outputDir, item.Extract, item.Name); err != nil {
			return fmt.Errorf("failed to write files: %w", err)
		}
	}

	binFile, shouldCreateSymlink := item.GetBinFileString()
	if shouldCreateSymlink {
		binDir := item.GetBinDir(config.BinsDir)
		if binDir == "" {
			return fmt.Errorf("bins-dir not specified for binary symlink")
		}

		var targetPath string
		if item.Extract {
			targetPath = filepath.Join(outputDir, item.Name, binFile)
		} else {
			if outputDir == "" {
				return fmt.Errorf("output-dir required when creating symlink for non-extracted file")
			}
			targetPath = filepath.Join(outputDir, item.Name)
		}

		symlinkName := filepath.Base(binFile)
		if err := createSymlink(targetPath, binDir, symlinkName); err != nil {
			return fmt.Errorf("failed to create symlink: %w", err)
		}
	}

	return nil
}

type FileToWrite struct {
	Name string
	Data []byte
}

func writeFiles(files []FileToWrite, outputDir string, isExtract bool, itemName string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if isExtract {
		extractDir := filepath.Join(outputDir, itemName)
		if err := removeExisting(extractDir); err != nil {
			return fmt.Errorf("failed to remove existing directory %s: %w", extractDir, err)
		}
	} else {
		singleFilePath := filepath.Join(outputDir, itemName)
		if err := removeExisting(singleFilePath); err != nil {
			return fmt.Errorf("failed to remove existing file %s: %w", singleFilePath, err)
		}
	}

	for _, file := range files {
		filePath := filepath.Join(outputDir, file.Name)

		fileDir := filepath.Dir(filePath)
		if err := os.MkdirAll(fileDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory for file %s: %w", filePath, err)
		}

		if err := os.WriteFile(filePath, file.Data, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}

		fmt.Printf("Written: %s\n", filePath)
	}

	return nil
}

func createSymlink(targetPath, binDir, symlinkName string) error {
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	symlinkPath := filepath.Join(binDir, symlinkName)

	if err := removeExisting(symlinkPath); err != nil {
		return fmt.Errorf("failed to remove existing file/directory at symlink location: %w", err)
	}

	if err := os.Symlink(targetPath, symlinkPath); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	if err := os.Chmod(targetPath, 0755); err != nil {
		fmt.Printf("Warning: failed to make target executable: %v\n", err)
	}

	fmt.Printf("Created symlink: %s -> %s\n", symlinkPath, targetPath)
	return nil
}
