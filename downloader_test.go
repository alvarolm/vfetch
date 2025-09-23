package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/sha3"
)

func TestDownloadFile(t *testing.T) {
	testData := []byte("test file content")

	tests := []struct {
		name         string
		statusCode   int
		responseBody []byte
		expectError  bool
		expectedData []byte
	}{
		{
			name:         "successful download",
			statusCode:   http.StatusOK,
			responseBody: testData,
			expectError:  false,
			expectedData: testData,
		},
		{
			name:        "404 not found",
			statusCode:  http.StatusNotFound,
			expectError: true,
		},
		{
			name:        "500 server error",
			statusCode:  http.StatusInternalServerError,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				if tt.responseBody != nil {
					w.Write(tt.responseBody)
				}
			}))
			defer server.Close()

			result, err := DownloadFile(server.URL + "/testfile.txt")

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

			if string(result.Data) != string(tt.expectedData) {
				t.Errorf("Expected data %q, got %q", string(tt.expectedData), string(result.Data))
			}
		})
	}
}

func TestDownloadFileInvalidURL(t *testing.T) {
	_, err := DownloadFile("invalid-url")
	if err == nil {
		t.Errorf("Expected error for invalid URL, but got none")
	}
}

func TestVerifySHA256(t *testing.T) {
	testData := []byte("test data")
	hasher := sha256.New()
	hasher.Write(testData)
	correctHash := fmt.Sprintf("%x", hasher.Sum(nil))

	tests := []struct {
		name        string
		data        []byte
		hash        string
		expectError bool
	}{
		{
			name:        "correct hash",
			data:        testData,
			hash:        correctHash,
			expectError: false,
		},
		{
			name:        "incorrect hash",
			data:        testData,
			hash:        "incorrect_hash",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifySHA256(tt.data, tt.hash)
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

func TestVerifySHA512(t *testing.T) {
	testData := []byte("test data")
	hasher := sha512.New()
	hasher.Write(testData)
	correctHash := fmt.Sprintf("%x", hasher.Sum(nil))

	tests := []struct {
		name        string
		data        []byte
		hash        string
		expectError bool
	}{
		{
			name:        "correct hash",
			data:        testData,
			hash:        correctHash,
			expectError: false,
		},
		{
			name:        "incorrect hash",
			data:        testData,
			hash:        "incorrect_hash",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifySHA512(tt.data, tt.hash)
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

func TestVerifySHA3(t *testing.T) {
	testData := []byte("test data")
	hasher := sha3.New256()
	hasher.Write(testData)
	correctHash := fmt.Sprintf("%x", hasher.Sum(nil))

	tests := []struct {
		name        string
		data        []byte
		hash        string
		expectError bool
	}{
		{
			name:        "correct hash",
			data:        testData,
			hash:        correctHash,
			expectError: false,
		},
		{
			name:        "incorrect hash",
			data:        testData,
			hash:        "incorrect_hash",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifySHA3(tt.data, tt.hash)
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

func TestVerifyBLAKE2b(t *testing.T) {
	testData := []byte("test data")
	hasher, _ := blake2b.New256(nil)
	hasher.Write(testData)
	correctHash := fmt.Sprintf("%x", hasher.Sum(nil))

	tests := []struct {
		name        string
		data        []byte
		hash        string
		expectError bool
	}{
		{
			name:        "correct hash",
			data:        testData,
			hash:        correctHash,
			expectError: false,
		},
		{
			name:        "incorrect hash",
			data:        testData,
			hash:        "incorrect_hash",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifyBLAKE2b(tt.data, tt.hash)
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

func TestVerifyBLAKE2s(t *testing.T) {
	testData := []byte("test data")
	hasher, _ := blake2s.New256(nil)
	hasher.Write(testData)
	correctHash := fmt.Sprintf("%x", hasher.Sum(nil))

	tests := []struct {
		name        string
		data        []byte
		hash        string
		expectError bool
	}{
		{
			name:        "correct hash",
			data:        testData,
			hash:        correctHash,
			expectError: false,
		},
		{
			name:        "incorrect hash",
			data:        testData,
			hash:        "incorrect_hash",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifyBLAKE2s(tt.data, tt.hash)
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

func TestVerifyHash(t *testing.T) {
	testData := []byte("test data")

	sha256Hasher := sha256.New()
	sha256Hasher.Write(testData)
	sha256Hash := fmt.Sprintf("%x", sha256Hasher.Sum(nil))

	tests := []struct {
		name         string
		data         []byte
		expectedHash string
		expectError  bool
	}{
		{
			name:         "valid sha256 hash",
			data:         testData,
			expectedHash: "sha256:" + sha256Hash,
			expectError:  false,
		},
		{
			name:         "invalid hash format",
			data:         testData,
			expectedHash: "invalid_format",
			expectError:  true,
		},
		{
			name:         "unsupported hash type",
			data:         testData,
			expectedHash: "md5:abcd1234",
			expectError:  true,
		},
		{
			name:         "wrong hash value",
			data:         testData,
			expectedHash: "sha256:wronghash",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyHash(tt.data, tt.expectedHash)
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

func TestVerifyHashes(t *testing.T) {
	testData := []byte("test data")

	sha256Hasher := sha256.New()
	sha256Hasher.Write(testData)
	sha256Hash := fmt.Sprintf("%x", sha256Hasher.Sum(nil))

	sha512Hasher := sha512.New()
	sha512Hasher.Write(testData)
	sha512Hash := fmt.Sprintf("%x", sha512Hasher.Sum(nil))

	tests := []struct {
		name           string
		data           []byte
		expectedHashes []string
		expectError    bool
	}{
		{
			name:           "valid hashes - all correct",
			data:           testData,
			expectedHashes: []string{"sha256:" + sha256Hash, "sha512:" + sha512Hash},
			expectError:    false,
		},
		{
			name:           "valid hashes - one correct",
			data:           testData,
			expectedHashes: []string{"sha256:" + sha256Hash, "sha256:wronghash"},
			expectError:    false,
		},
		{
			name:           "all hashes wrong",
			data:           testData,
			expectedHashes: []string{"sha256:wronghash1", "sha512:wronghash2"},
			expectError:    true,
		},
		{
			name:           "empty hashes array",
			data:           testData,
			expectedHashes: []string{},
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyHashes(tt.data, tt.expectedHashes)
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
