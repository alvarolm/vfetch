package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"testing"
)

func TestExtractArchive(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		expectError bool
	}{
		{
			name:        "unsupported format",
			filename:    "file.rar",
			expectError: true,
		},
		{
			name:        "zip file",
			filename:    "file.zip",
			expectError: false,
		},
		{
			name:        "tar.gz file",
			filename:    "file.tar.gz",
			expectError: false,
		},
		{
			name:        "tar file",
			filename:    "file.tar",
			expectError: false,
		},
		{
			name:        "gz file",
			filename:    "file.gz",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testData []byte
			var err error

			switch {
			case tt.filename == "file.zip":
				testData, err = createTestZip()
			case tt.filename == "file.tar.gz":
				testData, err = createTestTarGz()
			case tt.filename == "file.tar":
				testData, err = createTestTar()
			case tt.filename == "file.gz":
				testData, err = createTestGzip()
			default:
				testData = []byte("invalid data")
			}

			if err != nil && !tt.expectError {
				t.Fatalf("Failed to create test data: %v", err)
			}

			result, err := ExtractArchive(testData, tt.filename)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Errorf("Expected result, but got nil")
				return
			}

			if len(result.Files) == 0 {
				t.Errorf("Expected files in result, but got none")
			}
		})
	}
}

func TestExtractZip(t *testing.T) {
	zipData, err := createTestZip()
	if err != nil {
		t.Fatalf("Failed to create test zip: %v", err)
	}

	result, err := extractZip(zipData)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if result == nil {
		t.Errorf("Expected result, but got nil")
		return
	}

	if len(result.Files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(result.Files))
		return
	}

	expectedContent := "test file content"
	if string(result.Files[0].Data) != expectedContent {
		t.Errorf("Expected content %q, got %q", expectedContent, string(result.Files[0].Data))
	}

	if result.Files[0].Name != "testfile.txt" {
		t.Errorf("Expected filename 'testfile.txt', got %q", result.Files[0].Name)
	}
}

func TestExtractZipInvalidData(t *testing.T) {
	invalidData := []byte("not a zip file")
	_, err := extractZip(invalidData)
	if err == nil {
		t.Errorf("Expected error for invalid zip data, but got none")
	}
}

func TestExtractTar(t *testing.T) {
	tarData, err := createTestTar()
	if err != nil {
		t.Fatalf("Failed to create test tar: %v", err)
	}

	result, err := extractTar(tarData)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if result == nil {
		t.Errorf("Expected result, but got nil")
		return
	}

	if len(result.Files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(result.Files))
		return
	}

	expectedContent := "test file content"
	if string(result.Files[0].Data) != expectedContent {
		t.Errorf("Expected content %q, got %q", expectedContent, string(result.Files[0].Data))
	}

	if result.Files[0].Name != "testfile.txt" {
		t.Errorf("Expected filename 'testfile.txt', got %q", result.Files[0].Name)
	}
}

func TestExtractTarGz(t *testing.T) {
	tarGzData, err := createTestTarGz()
	if err != nil {
		t.Fatalf("Failed to create test tar.gz: %v", err)
	}

	result, err := extractTarGz(tarGzData)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if result == nil {
		t.Errorf("Expected result, but got nil")
		return
	}

	if len(result.Files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(result.Files))
		return
	}

	expectedContent := "test file content"
	if string(result.Files[0].Data) != expectedContent {
		t.Errorf("Expected content %q, got %q", expectedContent, string(result.Files[0].Data))
	}
}

func TestExtractGzip(t *testing.T) {
	gzipData, err := createTestGzip()
	if err != nil {
		t.Fatalf("Failed to create test gzip: %v", err)
	}

	result, err := extractGzip(gzipData, "testfile.txt.gz")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if result == nil {
		t.Errorf("Expected result, but got nil")
		return
	}

	if len(result.Files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(result.Files))
		return
	}

	expectedContent := "test file content"
	if string(result.Files[0].Data) != expectedContent {
		t.Errorf("Expected content %q, got %q", expectedContent, string(result.Files[0].Data))
	}

	if result.Files[0].Name != "testfile.txt" {
		t.Errorf("Expected filename 'testfile.txt', got %q", result.Files[0].Name)
	}
}

func TestExtractGzipInvalidData(t *testing.T) {
	invalidData := []byte("not gzip data")
	_, err := extractGzip(invalidData, "test.gz")
	if err == nil {
		t.Errorf("Expected error for invalid gzip data, but got none")
	}
}

func TestExtractGzipFallbackFilename(t *testing.T) {
	gzipData, err := createTestGzip()
	if err != nil {
		t.Fatalf("Failed to create test gzip: %v", err)
	}

	result, err := extractGzip(gzipData, "notgzfile")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if result.Files[0].Name != "decompressed_file" {
		t.Errorf("Expected fallback filename 'decompressed_file', got %q", result.Files[0].Name)
	}
}

// Helper functions to create test archive data

func createTestZip() ([]byte, error) {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	file, err := zipWriter.Create("testfile.txt")
	if err != nil {
		return nil, err
	}

	_, err = file.Write([]byte("test file content"))
	if err != nil {
		return nil, err
	}

	err = zipWriter.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func createTestTar() ([]byte, error) {
	var buf bytes.Buffer
	tarWriter := tar.NewWriter(&buf)

	header := &tar.Header{
		Name: "testfile.txt",
		Mode: 0644,
		Size: int64(len("test file content")),
	}

	err := tarWriter.WriteHeader(header)
	if err != nil {
		return nil, err
	}

	_, err = tarWriter.Write([]byte("test file content"))
	if err != nil {
		return nil, err
	}

	err = tarWriter.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func createTestTarGz() ([]byte, error) {
	tarData, err := createTestTar()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)

	_, err = gzipWriter.Write(tarData)
	if err != nil {
		return nil, err
	}

	err = gzipWriter.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func createTestGzip() ([]byte, error) {
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)

	_, err := gzipWriter.Write([]byte("test file content"))
	if err != nil {
		return nil, err
	}

	err = gzipWriter.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}