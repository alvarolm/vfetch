package main

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		configJSON  string
		expectError bool
	}{
		{
			name: "valid config",
			configJSON: `{
				"output-dir": "/tmp/test",
				"bins-dir": "/tmp/bins",
				"fetch": [
					{
						"name": "test-item",
						"url": "https://example.com/file.zip",
						"hash": "sha256:abcd1234",
						"extract": true
					}
				]
			}`,
			expectError: false,
		},
		{
			name: "empty fetch array",
			configJSON: `{
				"output-dir": "/tmp/test",
				"fetch": []
			}`,
			expectError: false,
		},
		{
			name:        "invalid JSON",
			configJSON:  `{invalid json}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "config-*.json")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.WriteString(tt.configJSON); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			tmpFile.Close()

			config, err := LoadConfig(tmpFile.Name())
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

			if config == nil {
				t.Errorf("Expected config, but got nil")
			}
		})
	}
}

func TestLoadConfigFileNotFound(t *testing.T) {
	_, err := LoadConfig("nonexistent-file.json")
	if err == nil {
		t.Errorf("Expected error for nonexistent file, but got none")
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectError bool
	}{
		{
			name: "valid config with single hash",
			config: Config{
				Fetch: []FetchItem{
					{
						Name:    "test",
						URL:     "https://example.com/file.zip",
						Version: "1.0.0",
						Hash:    "sha256:abcd1234",
					},
				},
			},
			expectError: false,
		},
		{
			name: "valid config with multiple hashes",
			config: Config{
				Fetch: []FetchItem{
					{
						Name:    "test",
						URL:     "https://example.com/file.zip",
						Version: "1.0.0",
						Hashes:  []string{"sha256:abcd1234", "sha512:efgh5678"},
					},
				},
			},
			expectError: false,
		},
		{
			name: "empty fetch array",
			config: Config{
				Fetch: []FetchItem{},
			},
			expectError: true,
		},
		{
			name: "missing name",
			config: Config{
				Fetch: []FetchItem{
					{
						URL:  "https://example.com/file.zip",
						Hash: "sha256:abcd1234",
					},
				},
			},
			expectError: true,
		},
		{
			name: "missing URL",
			config: Config{
				Fetch: []FetchItem{
					{
						Name: "test",
						Hash: "sha256:abcd1234",
					},
				},
			},
			expectError: true,
		},
		{
			name: "missing hash and hashes",
			config: Config{
				Fetch: []FetchItem{
					{
						Name: "test",
						URL:  "https://example.com/file.zip",
					},
				},
			},
			expectError: true,
		},
		{
			name: "both hash and hashes provided",
			config: Config{
				Fetch: []FetchItem{
					{
						Name:   "test",
						URL:    "https://example.com/file.zip",
						Hash:   "sha256:abcd1234",
						Hashes: []string{"sha512:efgh5678"},
					},
				},
			},
			expectError: true,
		},
		{
			name: "invalid hash format",
			config: Config{
				Fetch: []FetchItem{
					{
						Name: "test",
						URL:  "https://example.com/file.zip",
						Hash: "invalid-hash",
					},
				},
			},
			expectError: true,
		},
		{
			name: "bin-file true with extract true",
			config: Config{
				Fetch: []FetchItem{
					{
						Name:    "test",
						URL:     "https://example.com/file.zip",
						Hash:    "sha256:abcd1234",
						Extract: true,
						BinFile: true,
					},
				},
			},
			expectError: true,
		},
		{
			name: "bin-file string with extract true",
			config: Config{
				Fetch: []FetchItem{
					{
						Name:    "test",
						URL:     "https://example.com/file.zip",
						Version: "1.0.0",
						Hash:    "sha256:abcd1234",
						Extract: true,
						BinFile: "binary",
					},
				},
			},
			expectError: false,
		},
		{
			name: "version field contains $version placeholder",
			config: Config{
				Fetch: []FetchItem{
					{
						Name:    "test",
						URL:     "https://example.com/v$version/file.zip",
						Version: "$version",
						Hash:    "sha256:abcd1234",
					},
				},
			},
			expectError: true,
		},
		{
			name: "version field contains $VERSION placeholder",
			config: Config{
				Fetch: []FetchItem{
					{
						Name:    "test",
						URL:     "https://example.com/v$VERSION/file.zip",
						Version: "v$VERSION",
						Hash:    "sha256:abcd1234",
					},
				},
			},
			expectError: true,
		},
		{
			name: "version field contains both placeholders",
			config: Config{
				Fetch: []FetchItem{
					{
						Name:    "test",
						URL:     "https://example.com/v$version/file.zip",
						Version: "$version-$VERSION",
						Hash:    "sha256:abcd1234",
					},
				},
			},
			expectError: true,
		},
		{
			name: "valid version field with placeholder in URL",
			config: Config{
				Fetch: []FetchItem{
					{
						Name:    "test",
						URL:     "https://example.com/v$version/file.zip",
						Version: "1.2.3",
						Hash:    "sha256:abcd1234",
					},
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(&tt.config)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateHashFormat(t *testing.T) {
	validHashTypes := []string{"sha256:", "sha512:", "sha3:", "blake2b:", "blake2s:", "blake3:"}

	tests := []struct {
		name        string
		hash        string
		expectError bool
	}{
		{
			name:        "valid sha256",
			hash:        "sha256:abcd1234",
			expectError: false,
		},
		{
			name:        "valid sha512",
			hash:        "sha512:efgh5678",
			expectError: false,
		},
		{
			name:        "valid sha3",
			hash:        "sha3:ijkl9012",
			expectError: false,
		},
		{
			name:        "valid blake2b",
			hash:        "blake2b:mnop3456",
			expectError: false,
		},
		{
			name:        "valid blake2s",
			hash:        "blake2s:qrst7890",
			expectError: false,
		},
		{
			name:        "valid blake3",
			hash:        "blake3:uvwx1234",
			expectError: false,
		},
		{
			name:        "invalid hash type",
			hash:        "md5:abcd1234",
			expectError: true,
		},
		{
			name:        "no colon separator",
			hash:        "sha256abcd1234",
			expectError: true,
		},
		{
			name:        "empty hash",
			hash:        "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateHashFormat(tt.hash, validHashTypes, 0)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestFetchItem_GetBinFileString(t *testing.T) {
	tests := []struct {
		name           string
		item           FetchItem
		expectedString string
		expectedBool   bool
	}{
		{
			name: "nil bin-file",
			item: FetchItem{
				BinFile: nil,
			},
			expectedString: "",
			expectedBool:   false,
		},
		{
			name: "bin-file true",
			item: FetchItem{
				Name:    "myfile",
				URL:     "https://example.com/myfile.tar.gz",
				BinFile: true,
			},
			expectedString: "myfile",
			expectedBool:   true,
		},
		{
			name: "bin-file true with query params",
			item: FetchItem{
				Name:    "myfile",
				URL:     "https://example.com/myfile.tar.gz?version=1.0",
				BinFile: true,
			},
			expectedString: "myfile",
			expectedBool:   true,
		},
		{
			name: "bin-file false",
			item: FetchItem{
				BinFile: false,
			},
			expectedString: "",
			expectedBool:   false,
		},
		{
			name: "bin-file string",
			item: FetchItem{
				BinFile: "custom-binary",
			},
			expectedString: "custom-binary",
			expectedBool:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultString, resultBool := tt.item.GetBinFileString()
			if resultString != tt.expectedString {
				t.Errorf("Expected string %q, got %q", tt.expectedString, resultString)
			}
			if resultBool != tt.expectedBool {
				t.Errorf("Expected bool %v, got %v", tt.expectedBool, resultBool)
			}
		})
	}
}

func TestFetchItem_GetOutputDir(t *testing.T) {
	tests := []struct {
		name              string
		item              FetchItem
		globalOutputDir   string
		expectedOutputDir string
	}{
		{
			name: "item has output dir",
			item: FetchItem{
				OutputDir: "/custom/output",
			},
			globalOutputDir:   "/global/output",
			expectedOutputDir: "/custom/output",
		},
		{
			name:              "item has no output dir",
			item:              FetchItem{},
			globalOutputDir:   "/global/output",
			expectedOutputDir: "/global/output",
		},
		{
			name:              "both empty",
			item:              FetchItem{},
			globalOutputDir:   "",
			expectedOutputDir: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.item.GetOutputDir(tt.globalOutputDir)
			if result != tt.expectedOutputDir {
				t.Errorf("Expected %q, got %q", tt.expectedOutputDir, result)
			}
		})
	}
}

func TestFetchItem_GetBinDir(t *testing.T) {
	tests := []struct {
		name           string
		item           FetchItem
		globalBinDir   string
		expectedBinDir string
	}{
		{
			name: "item has bin dir",
			item: FetchItem{
				BinDir: "/custom/bin",
			},
			globalBinDir:   "/global/bin",
			expectedBinDir: "/custom/bin",
		},
		{
			name:           "item has no bin dir",
			item:           FetchItem{},
			globalBinDir:   "/global/bin",
			expectedBinDir: "/global/bin",
		},
		{
			name:           "both empty",
			item:           FetchItem{},
			globalBinDir:   "",
			expectedBinDir: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.item.GetBinDir(tt.globalBinDir)
			if result != tt.expectedBinDir {
				t.Errorf("Expected %q, got %q", tt.expectedBinDir, result)
			}
		})
	}
}