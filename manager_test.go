package main

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "verifetch-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	files := []FileToWrite{
		{
			Name: "file1.txt",
			Data: []byte("content of file 1"),
		},
		{
			Name: "subdir/file2.txt",
			Data: []byte("content of file 2"),
		},
	}

	err = writeFiles(files, tmpDir, false, "test-item")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	file1Path := filepath.Join(tmpDir, "file1.txt")
	if _, err := os.Stat(file1Path); os.IsNotExist(err) {
		t.Errorf("Expected file1.txt to exist")
	}

	file1Content, err := os.ReadFile(file1Path)
	if err != nil {
		t.Errorf("Failed to read file1.txt: %v", err)
	} else if string(file1Content) != "content of file 1" {
		t.Errorf("Expected content 'content of file 1', got %q", string(file1Content))
	}

	file2Path := filepath.Join(tmpDir, "subdir", "file2.txt")
	if _, err := os.Stat(file2Path); os.IsNotExist(err) {
		t.Errorf("Expected subdir/file2.txt to exist")
	}

	file2Content, err := os.ReadFile(file2Path)
	if err != nil {
		t.Errorf("Failed to read subdir/file2.txt: %v", err)
	} else if string(file2Content) != "content of file 2" {
		t.Errorf("Expected content 'content of file 2', got %q", string(file2Content))
	}
}

func TestCreateSymlink(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "verifetch-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	targetFile := filepath.Join(tmpDir, "target", "binary")
	binDir := filepath.Join(tmpDir, "bin")

	if err := os.MkdirAll(filepath.Dir(targetFile), 0755); err != nil {
		t.Fatalf("Failed to create target dir: %v", err)
	}

	if err := os.WriteFile(targetFile, []byte("#!/bin/bash\necho hello"), 0644); err != nil {
		t.Fatalf("Failed to create target file: %v", err)
	}

	err = createSymlink(targetFile, binDir, "mybinary")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	symlinkPath := filepath.Join(binDir, "mybinary")
	if _, err := os.Lstat(symlinkPath); os.IsNotExist(err) {
		t.Errorf("Expected symlink to exist at %s", symlinkPath)
		return
	}

	linkTarget, err := os.Readlink(symlinkPath)
	if err != nil {
		t.Errorf("Failed to read symlink: %v", err)
		return
	}

	if linkTarget != targetFile {
		t.Errorf("Expected symlink target %q, got %q", targetFile, linkTarget)
	}

	targetInfo, err := os.Stat(targetFile)
	if err != nil {
		t.Errorf("Failed to stat target file: %v", err)
		return
	}

	if targetInfo.Mode()&0111 == 0 {
		t.Errorf("Expected target file to be executable")
	}
}

func TestCreateSymlinkReplaceExisting(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "verifetch-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	targetFile := filepath.Join(tmpDir, "target", "binary")
	binDir := filepath.Join(tmpDir, "bin")
	symlinkPath := filepath.Join(binDir, "mybinary")

	if err := os.MkdirAll(filepath.Dir(targetFile), 0755); err != nil {
		t.Fatalf("Failed to create target dir: %v", err)
	}

	if err := os.WriteFile(targetFile, []byte("#!/bin/bash\necho hello"), 0644); err != nil {
		t.Fatalf("Failed to create target file: %v", err)
	}

	if err := os.MkdirAll(binDir, 0755); err != nil {
		t.Fatalf("Failed to create bin dir: %v", err)
	}

	if err := os.Symlink("/old/target", symlinkPath); err != nil {
		t.Fatalf("Failed to create existing symlink: %v", err)
	}

	err = createSymlink(targetFile, binDir, "mybinary")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	linkTarget, err := os.Readlink(symlinkPath)
	if err != nil {
		t.Errorf("Failed to read symlink: %v", err)
		return
	}

	if linkTarget != targetFile {
		t.Errorf("Expected symlink target %q, got %q", targetFile, linkTarget)
	}
}

func TestProcessFetchItemDownloadOnly(t *testing.T) {
	testData := []byte("test file content")
	hasher := sha256.New()
	hasher.Write(testData)
	expectedHash := fmt.Sprintf("sha256:%x", hasher.Sum(nil))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(testData)
	}))
	defer server.Close()

	tmpDir, err := os.MkdirTemp("", "verifetch-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	config := &Config{
		OutputDir: tmpDir,
	}

	item := FetchItem{
		Name: "test-item",
		URL:  server.URL + "/testfile.txt",
		Hash: expectedHash,
	}

	err = ProcessFetchItem(config, item)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	filePath := filepath.Join(tmpDir, "test-item")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("Expected downloaded file to exist at %s", filePath)
		return
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Errorf("Failed to read downloaded file: %v", err)
		return
	}

	if string(content) != string(testData) {
		t.Errorf("Expected content %q, got %q", string(testData), string(content))
	}
}

func TestProcessFetchItemWithExtraction(t *testing.T) {
	zipData, err := createTestZipForManager()
	if err != nil {
		t.Fatalf("Failed to create test zip: %v", err)
	}

	hasher := sha256.New()
	hasher.Write(zipData)
	expectedHash := fmt.Sprintf("sha256:%x", hasher.Sum(nil))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(zipData)
	}))
	defer server.Close()

	tmpDir, err := os.MkdirTemp("", "verifetch-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	config := &Config{
		OutputDir: tmpDir,
	}

	item := FetchItem{
		Name:    "test-item.zip",
		URL:     server.URL + "/testfile.zip",
		Hash:    expectedHash,
		Extract: true,
	}

	err = ProcessFetchItem(config, item)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	extractedFilePath := filepath.Join(tmpDir, "test-item.zip", "extracted.txt")
	if _, err := os.Stat(extractedFilePath); os.IsNotExist(err) {
		t.Errorf("Expected extracted file to exist at %s", extractedFilePath)
		return
	}

	content, err := os.ReadFile(extractedFilePath)
	if err != nil {
		t.Errorf("Failed to read extracted file: %v", err)
		return
	}

	expectedContent := "content of extracted file"
	if string(content) != expectedContent {
		t.Errorf("Expected content %q, got %q", expectedContent, string(content))
	}
}

func TestProcessFetchItemWithBinFile(t *testing.T) {
	testData := []byte("#!/bin/bash\necho hello")
	hasher := sha256.New()
	hasher.Write(testData)
	expectedHash := fmt.Sprintf("sha256:%x", hasher.Sum(nil))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(testData)
	}))
	defer server.Close()

	tmpDir, err := os.MkdirTemp("", "verifetch-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	outputDir := filepath.Join(tmpDir, "output")
	binDir := filepath.Join(tmpDir, "bin")

	config := &Config{
		OutputDir: outputDir,
		BinsDir:   binDir,
	}

	item := FetchItem{
		Name:    "test-item",
		URL:     server.URL + "/mybinary",
		Hash:    expectedHash,
		BinFile: true,
	}

	err = ProcessFetchItem(config, item)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	symlinkPath := filepath.Join(binDir, "test-item")
	if _, err := os.Lstat(symlinkPath); os.IsNotExist(err) {
		t.Errorf("Expected symlink to exist at %s", symlinkPath)
		return
	}

	targetFile := filepath.Join(outputDir, "test-item")
	linkTarget, err := os.Readlink(symlinkPath)
	if err != nil {
		t.Errorf("Failed to read symlink: %v", err)
		return
	}

	if linkTarget != targetFile {
		t.Errorf("Expected symlink target %q, got %q", targetFile, linkTarget)
	}
}

func TestProcessFetchItemHashVerificationFailure(t *testing.T) {
	testData := []byte("test file content")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(testData)
	}))
	defer server.Close()

	config := &Config{}

	item := FetchItem{
		Name: "test-item",
		URL:  server.URL + "/testfile.txt",
		Hash: "sha256:wronghash",
	}

	err := ProcessFetchItem(config, item)
	if err == nil {
		t.Errorf("Expected error for hash verification failure, but got none")
	}
}

func TestProcessFetchItemNoHash(t *testing.T) {
	testData := []byte("test file content")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(testData)
	}))
	defer server.Close()

	config := &Config{}

	item := FetchItem{
		Name: "test-item",
		URL:  server.URL + "/testfile.txt",
	}

	err := ProcessFetchItem(config, item)
	if err == nil {
		t.Errorf("Expected error for missing hash, but got none")
	}
}

func createTestZipForManager() ([]byte, error) {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	file, err := zipWriter.Create("extracted.txt")
	if err != nil {
		return nil, err
	}

	_, err = file.Write([]byte("content of extracted file"))
	if err != nil {
		return nil, err
	}

	err = zipWriter.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func TestReplaceVersionPlaceholders(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		version  string
		expected string
	}{
		{
			name:     "lowercase $version placeholder",
			url:      "https://github.com/user/repo/releases/download/v$version/file.tar.gz",
			version:  "1.2.3",
			expected: "https://github.com/user/repo/releases/download/v1.2.3/file.tar.gz",
		},
		{
			name:     "uppercase $VERSION placeholder",
			url:      "https://example.com/releases/$VERSION/archive.zip",
			version:  "2.0.0",
			expected: "https://example.com/releases/2.0.0/archive.zip",
		},
		{
			name:     "multiple placeholders",
			url:      "https://api.example.com/v$version/download/$VERSION/file.bin",
			version:  "1.5.0",
			expected: "https://api.example.com/v1.5.0/download/1.5.0/file.bin",
		},
		{
			name:     "no placeholders",
			url:      "https://example.com/static/file.tar.gz",
			version:  "1.0.0",
			expected: "https://example.com/static/file.tar.gz",
		},
		{
			name:     "empty version",
			url:      "https://example.com/v$version/file.zip",
			version:  "",
			expected: "https://example.com/v/file.zip",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceVersionPlaceholders(tt.url, tt.version)
			if result != tt.expected {
				t.Errorf("replaceVersionPlaceholders() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestProcessFetchItemWithVersionPlaceholder(t *testing.T) {
	testData := []byte("test file content with version")
	hasher := sha256.New()
	hasher.Write(testData)
	expectedHash := fmt.Sprintf("sha256:%x", hasher.Sum(nil))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The URL should contain the replaced version (1.2.3)
		if r.URL.Path != "/releases/v1.2.3/file.txt" {
			t.Errorf("Expected URL path '/releases/v1.2.3/file.txt', got '%s'", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(testData)
	}))
	defer server.Close()

	tmpDir, err := os.MkdirTemp("", "verifetch-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	config := &Config{
		OutputDir: tmpDir,
	}

	item := FetchItem{
		Name:    "test-version-item",
		URL:     server.URL + "/releases/v$version/file.txt",
		Version: "1.2.3",
		Hash:    expectedHash,
	}

	err = ProcessFetchItem(config, item)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	filePath := filepath.Join(tmpDir, "test-version-item")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("Expected downloaded file to exist at %s", filePath)
		return
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Errorf("Failed to read downloaded file: %v", err)
		return
	}

	if string(content) != string(testData) {
		t.Errorf("Expected content %q, got %q", string(testData), string(content))
	}
}
func TestRemoveExisting(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "verifetch-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("remove existing file", func(t *testing.T) {
		filePath := filepath.Join(tmpDir, "test-file.txt")
		if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		if err := removeExisting(filePath); err != nil {
			t.Errorf("Unexpected error removing file: %v", err)
		}

		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			t.Errorf("Expected file to be removed")
		}
	})

	t.Run("remove existing directory", func(t *testing.T) {
		dirPath := filepath.Join(tmpDir, "test-dir")
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		testFile := filepath.Join(dirPath, "nested-file.txt")
		if err := os.WriteFile(testFile, []byte("nested content"), 0644); err != nil {
			t.Fatalf("Failed to create nested file: %v", err)
		}

		if err := removeExisting(dirPath); err != nil {
			t.Errorf("Unexpected error removing directory: %v", err)
		}

		if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
			t.Errorf("Expected directory to be removed")
		}
	})

	t.Run("remove existing symlink", func(t *testing.T) {
		targetPath := filepath.Join(tmpDir, "target-file.txt")
		if err := os.WriteFile(targetPath, []byte("target content"), 0644); err != nil {
			t.Fatalf("Failed to create target file: %v", err)
		}

		symlinkPath := filepath.Join(tmpDir, "test-symlink")
		if err := os.Symlink(targetPath, symlinkPath); err != nil {
			t.Fatalf("Failed to create symlink: %v", err)
		}

		if err := removeExisting(symlinkPath); err != nil {
			t.Errorf("Unexpected error removing symlink: %v", err)
		}

		if _, err := os.Lstat(symlinkPath); !os.IsNotExist(err) {
			t.Errorf("Expected symlink to be removed")
		}

		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			t.Errorf("Expected target file to remain")
		}
	})

	t.Run("remove non-existent path", func(t *testing.T) {
		nonExistentPath := filepath.Join(tmpDir, "non-existent-file")

		if err := removeExisting(nonExistentPath); err != nil {
			t.Errorf("Unexpected error for non-existent path: %v", err)
		}
	})
}

func TestWriteFilesRemoval(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "verifetch-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("remove existing file before writing single file", func(t *testing.T) {
		itemName := "existing-file.txt"
		existingFilePath := filepath.Join(tmpDir, itemName)

		if err := os.WriteFile(existingFilePath, []byte("old content"), 0644); err != nil {
			t.Fatalf("Failed to create existing file: %v", err)
		}

		files := []FileToWrite{
			{Name: itemName, Data: []byte("new content")},
		}

		if err := writeFiles(files, tmpDir, false, itemName); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		content, err := os.ReadFile(existingFilePath)
		if err != nil {
			t.Errorf("Failed to read file: %v", err)
		}

		if string(content) != "new content" {
			t.Errorf("Expected 'new content', got %q", string(content))
		}
	})

	t.Run("remove existing directory before extracting", func(t *testing.T) {
		itemName := "existing-dir"
		existingDirPath := filepath.Join(tmpDir, itemName)

		if err := os.MkdirAll(existingDirPath, 0755); err != nil {
			t.Fatalf("Failed to create existing directory: %v", err)
		}

		oldFilePath := filepath.Join(existingDirPath, "old-file.txt")
		if err := os.WriteFile(oldFilePath, []byte("old content"), 0644); err != nil {
			t.Fatalf("Failed to create old file: %v", err)
		}

		files := []FileToWrite{
			{Name: filepath.Join(itemName, "new-file.txt"), Data: []byte("new content")},
		}

		if err := writeFiles(files, tmpDir, true, itemName); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if _, err := os.Stat(oldFilePath); !os.IsNotExist(err) {
			t.Errorf("Expected old file to be removed")
		}

		newFilePath := filepath.Join(tmpDir, itemName, "new-file.txt")
		content, err := os.ReadFile(newFilePath)
		if err != nil {
			t.Errorf("Failed to read new file: %v", err)
		}

		if string(content) != "new content" {
			t.Errorf("Expected 'new content', got %q", string(content))
		}
	})
}

func TestCreateSymlinkRemoval(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "verifetch-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	targetFile := filepath.Join(tmpDir, "target")
	binDir := filepath.Join(tmpDir, "bin")
	symlinkName := "mybinary"
	symlinkPath := filepath.Join(binDir, symlinkName)

	if err := os.WriteFile(targetFile, []byte("#!/bin/bash\necho hello"), 0644); err != nil {
		t.Fatalf("Failed to create target file: %v", err)
	}

	t.Run("remove existing file before creating symlink", func(t *testing.T) {
		if err := os.MkdirAll(binDir, 0755); err != nil {
			t.Fatalf("Failed to create bin dir: %v", err)
		}

		if err := os.WriteFile(symlinkPath, []byte("regular file"), 0644); err != nil {
			t.Fatalf("Failed to create existing file: %v", err)
		}

		if err := createSymlink(targetFile, binDir, symlinkName); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		linkTarget, err := os.Readlink(symlinkPath)
		if err != nil {
			t.Errorf("Failed to read symlink: %v", err)
		}

		if linkTarget != targetFile {
			t.Errorf("Expected symlink target %q, got %q", targetFile, linkTarget)
		}
	})

	t.Run("remove existing directory before creating symlink", func(t *testing.T) {
		if err := os.RemoveAll(symlinkPath); err != nil {
			t.Fatalf("Failed to clean up: %v", err)
		}

		if err := os.MkdirAll(symlinkPath, 0755); err != nil {
			t.Fatalf("Failed to create existing directory: %v", err)
		}

		if err := createSymlink(targetFile, binDir, symlinkName); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		linkTarget, err := os.Readlink(symlinkPath)
		if err != nil {
			t.Errorf("Failed to read symlink: %v", err)
		}

		if linkTarget != targetFile {
			t.Errorf("Expected symlink target %q, got %q", targetFile, linkTarget)
		}
	})
}
