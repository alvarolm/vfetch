package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/sha3"
)

type DownloadResult struct {
	Data     []byte
	Filename string
}

func DownloadFile(url string) (*DownloadResult, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed with status: %d %s", resp.StatusCode, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	filename, err := getFilenameFromURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to extract filename from URL: %w", err)
	}

	return &DownloadResult{
		Data:     data,
		Filename: filename,
	}, nil
}

func getFilenameFromURL(url_ string) (string, error) {
	u, err := url.Parse(url_)
	if err != nil {
		return "", err
	}

	filename := path.Base(u.Path)
	if filename == "/" || filename == "." {
		return "download", nil
	}

	return filename, nil
}

func verifySHA256(data []byte, hash string) error {
	hasher := sha256.New()
	hasher.Write(data)
	actualHash := fmt.Sprintf("%x", hasher.Sum(nil))

	if actualHash != hash {
		return fmt.Errorf("hash mismatch: expected %s, got %s", hash, actualHash)
	}

	return nil
}

func verifySHA512(data []byte, hash string) error {
	hasher := sha512.New()
	hasher.Write(data)
	actualHash := fmt.Sprintf("%x", hasher.Sum(nil))

	if actualHash != hash {
		return fmt.Errorf("hash mismatch: expected %s, got %s", hash, actualHash)
	}

	return nil
}

func verifySHA3(data []byte, hash string) error {
	hasher := sha3.New256()
	hasher.Write(data)
	actualHash := fmt.Sprintf("%x", hasher.Sum(nil))

	if actualHash != hash {
		return fmt.Errorf("hash mismatch: expected %s, got %s", hash, actualHash)
	}

	return nil
}

func verifyBLAKE2b(data []byte, hash string) error {
	hasher, err := blake2b.New256(nil)
	if err != nil {
		return fmt.Errorf("failed to create BLAKE2b hasher: %w", err)
	}
	hasher.Write(data)
	actualHash := fmt.Sprintf("%x", hasher.Sum(nil))

	if actualHash != hash {
		return fmt.Errorf("hash mismatch: expected %s, got %s", hash, actualHash)
	}

	return nil
}

func verifyBLAKE2s(data []byte, hash string) error {
	hasher, err := blake2s.New256(nil)
	if err != nil {
		return fmt.Errorf("failed to create BLAKE2s hasher: %w", err)
	}
	hasher.Write(data)
	actualHash := fmt.Sprintf("%x", hasher.Sum(nil))

	if actualHash != hash {
		return fmt.Errorf("hash mismatch: expected %s, got %s", hash, actualHash)
	}

	return nil
}

func VerifyHash(data []byte, expectedHash string) error {
	parts := strings.Split(expectedHash, ":")

	if len(parts) != 2 {
		return fmt.Errorf("invalid hash format: %s", expectedHash)
	}

	hashType := parts[0]
	hashValue := parts[1]

	switch hashType {
	case "sha256":
		return verifySHA256(data, hashValue)
	case "sha512":
		return verifySHA512(data, hashValue)
	case "sha3":
		return verifySHA3(data, hashValue)
	case "blake2b":
		return verifyBLAKE2b(data, hashValue)
	case "blake2s":
		return verifyBLAKE2s(data, hashValue)
	default:
		return fmt.Errorf("unsupported hash type: %s", hashType)
	}
}

func VerifyHashes(data []byte, expectedHashes []string) error {
	if len(expectedHashes) == 0 {
		return fmt.Errorf("no hashes provided for verification")
	}

	var errors []string
	verified := false

	for _, expectedHash := range expectedHashes {
		if err := VerifyHash(data, expectedHash); err != nil {
			errors = append(errors, fmt.Sprintf("hash %s failed: %v", expectedHash, err))
		} else {
			verified = true
		}
	}

	if !verified {
		return fmt.Errorf("all hash verifications failed:\n%s", strings.Join(errors, "\n"))
	}

	return nil
}
