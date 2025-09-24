package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tidwall/jsonc"
)

type Config struct {
	OutputDir string      `json:"output-dir"`
	BinsDir   string      `json:"bins-dir"`
	Fetch     []FetchItem `json:"fetch"`
}

type FetchItem struct {
	Name       string      `json:"name"`
	URL        string      `json:"url"`
	Version    string      `json:"version"`
	Hash       string      `json:"hash"`
	Hashes     []string    `json:"hashes"`
	Extract    bool        `json:"extract"`
	BinFile    interface{} `json:"bin-file"`
	BinDir     string      `json:"bin-dir"`
	OutputDir  string      `json:"output-dir"`
	HomeURL    string      `json:"home-url,omitempty"`
	SourceURL  string      `json:"source-url,omitempty"`
	LicenseURL string      `json:"license-url,omitempty"`
	AuthorURL  string      `json:"author-url,omitempty"`
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(jsonc.ToJSONInPlace(data), &config); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	// Resolve relative paths in the config relative to the config file's directory
	configDir := filepath.Dir(configPath)
	if config.OutputDir != "" && !filepath.IsAbs(config.OutputDir) {
		config.OutputDir = filepath.Join(configDir, config.OutputDir)
	}
	if config.BinsDir != "" && !filepath.IsAbs(config.BinsDir) {
		config.BinsDir = filepath.Join(configDir, config.BinsDir)
	}

	// Resolve relative output-dir paths for individual items
	for i := range config.Fetch {
		if config.Fetch[i].OutputDir != "" && !filepath.IsAbs(config.Fetch[i].OutputDir) {
			config.Fetch[i].OutputDir = filepath.Join(configDir, config.Fetch[i].OutputDir)
		}
		if config.Fetch[i].BinDir != "" && !filepath.IsAbs(config.Fetch[i].BinDir) {
			config.Fetch[i].BinDir = filepath.Join(configDir, config.Fetch[i].BinDir)
		}
	}

	return &config, nil
}

func ValidateConfig(config *Config) error {
	if len(config.Fetch) == 0 {
		return fmt.Errorf("no fetch items specified")
	}

	for i, item := range config.Fetch {
		if err := validateFetchItem(item, i); err != nil {
			return err
		}
	}

	return nil
}

func validateFetchItem(item FetchItem, index int) error {
	if item.Name == "" {
		return fmt.Errorf("fetch item %d: name is required", index)
	}

	if item.URL == "" {
		return fmt.Errorf("fetch item %d: URL is required", index)
	}

	if item.Version == "" {
		return fmt.Errorf("fetch item %d: version is required", index)
	}

	// Check if version field contains placeholders (which would cause infinite recursion)
	if strings.Contains(item.Version, "$version") || strings.Contains(item.Version, "$VERSION") {
		return fmt.Errorf("fetch item %d: version field cannot contain version placeholders ($version or $VERSION)", index)
	}

	// Ensure only one of hash or hashes is provided
	hasHash := item.Hash != ""
	hasHashes := len(item.Hashes) > 0

	if !hasHash && !hasHashes {
		return fmt.Errorf("fetch item %d: either 'hash' or 'hashes' field is required", index)
	}

	if hasHash && hasHashes {
		return fmt.Errorf("fetch item %d: cannot specify both 'hash' and 'hashes' fields, use only one", index)
	}

	validHashTypes := []string{"sha256:", "sha512:", "sha3:", "blake2b:", "blake2s:", "blake3:"}

	// Validate single hash
	if hasHash {
		if err := validateHashFormat(item.Hash, validHashTypes, index); err != nil {
			return err
		}
	}

	// Validate multiple hashes
	if hasHashes {
		for i, hash := range item.Hashes {
			if err := validateHashFormat(hash, validHashTypes, index); err != nil {
				return fmt.Errorf("fetch item %d, hash %d: %w", index, i, err)
			}
		}
	}

	if item.BinFile != nil {
		switch binFile := item.BinFile.(type) {
		case bool:
			if binFile && item.Extract {
				return fmt.Errorf("fetch item %d: bin-file cannot be true when extract is true", index)
			}
		case string:
		default:
			return fmt.Errorf("fetch item %d: bin-file must be a string or boolean", index)
		}
	}

	return nil
}

func validateHashFormat(hash string, validHashTypes []string, index int) error {
	validHash := false
	for _, hashType := range validHashTypes {
		if strings.HasPrefix(hash, hashType) {
			validHash = true
			break
		}
	}
	if !validHash {
		return fmt.Errorf("hash must be in format 'type:value' where type is one of: sha256, sha512, sha3, blake2b, blake2s, blake3")
	}
	return nil
}

func (item FetchItem) GetBinFileString() (string, bool) {
	if item.BinFile == nil {
		return "", false
	}

	switch binFile := item.BinFile.(type) {
	case bool:
		if binFile {
			return item.Name, true
		}
		return "", false
	case string:
		return binFile, true
	default:
		return "", false
	}
}

func (item FetchItem) GetOutputDir(globalOutputDir string) string {
	if item.OutputDir != "" {
		return item.OutputDir
	}
	return globalOutputDir
}

func (item FetchItem) GetBinDir(globalBinDir string) string {
	if item.BinDir != "" {
		return item.BinDir
	}
	return globalBinDir
}
