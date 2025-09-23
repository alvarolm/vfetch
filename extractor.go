package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

type ExtractedFile struct {
	Name string
	Data []byte
}

type ExtractionResult struct {
	Files []ExtractedFile
}

func ExtractArchive(data []byte, filename string) (*ExtractionResult, error) {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".zip":
		return extractZip(data)
	case ".gz":
		if strings.HasSuffix(strings.ToLower(filename), ".tar.gz") {
			return extractTarGz(data)
		}
		return extractGzip(data, filename)
	case ".tgz":
		return extractTarGz(data)
	case ".tar":
		return extractTar(data)
	default:
		return nil, fmt.Errorf("unsupported archive format: %s", ext)
	}
}

func extractZip(data []byte) (*ExtractionResult, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to open zip archive: %w", err)
	}

	var files []ExtractedFile

	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			continue
		}

		rc, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s in zip: %w", file.Name, err)
		}

		fileData, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s from zip: %w", file.Name, err)
		}

		files = append(files, ExtractedFile{
			Name: file.Name,
			Data: fileData,
		})
	}

	return &ExtractionResult{Files: files}, nil
}

func extractTarGz(data []byte) (*ExtractionResult, error) {
	gzReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to open gzip reader: %w", err)
	}
	defer gzReader.Close()

	decompressed, err := io.ReadAll(gzReader)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress gzip data: %w", err)
	}

	return extractTar(decompressed)
}

func extractTar(data []byte) (*ExtractionResult, error) {
	tarReader := tar.NewReader(bytes.NewReader(data))

	var files []ExtractedFile

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read tar header: %w", err)
		}

		if header.Typeflag != tar.TypeReg {
			continue
		}

		fileData, err := io.ReadAll(tarReader)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s from tar: %w", header.Name, err)
		}

		files = append(files, ExtractedFile{
			Name: header.Name,
			Data: fileData,
		})
	}

	return &ExtractionResult{Files: files}, nil
}

func extractGzip(data []byte, filename string) (*ExtractionResult, error) {
	gzReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to open gzip reader: %w", err)
	}
	defer gzReader.Close()

	decompressed, err := io.ReadAll(gzReader)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress gzip data: %w", err)
	}

	outputName := strings.TrimSuffix(filename, ".gz")
	if outputName == filename {
		outputName = "decompressed_file"
	}

	return &ExtractionResult{
		Files: []ExtractedFile{
			{
				Name: outputName,
				Data: decompressed,
			},
		},
	}, nil
}