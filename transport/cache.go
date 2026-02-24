package transport

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CachePaths computes the cache file paths for a given prompt file.
// Returns (jsonlPath, hashPath, tmpPath).
func CachePaths(promptFile string) (jsonl, hash, tmp string) {
	dir := filepath.Dir(promptFile)
	base := strings.TrimSuffix(filepath.Base(promptFile), filepath.Ext(promptFile))
	jsonl = filepath.Join(dir, base+".jsonl")
	hash = filepath.Join(dir, "."+base+".hash")
	tmp = filepath.Join(dir, base+".jsonl.tmp")
	return
}

// PromptHash computes the SHA256 hex digest of file content.
func PromptHash(content []byte) string {
	h := sha256.Sum256(content)
	return fmt.Sprintf("%x", h[:])
}

// CacheValid returns true if the cached JSONL exists and its hash file
// matches the current prompt content hash.
func CacheValid(promptFile string) bool {
	content, err := os.ReadFile(promptFile)
	if err != nil {
		return false
	}

	jsonlPath, hashPath, _ := CachePaths(promptFile)

	// JSONL must exist
	if _, err := os.Stat(jsonlPath); err != nil {
		return false
	}

	// Hash file must exist and match
	stored, err := os.ReadFile(hashPath)
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(stored)) == PromptHash(content)
}

// WriteHashFile writes the SHA256 of the prompt content to the hash file.
func WriteHashFile(promptFile string) error {
	content, err := os.ReadFile(promptFile)
	if err != nil {
		return err
	}

	_, hashPath, _ := CachePaths(promptFile)
	return os.WriteFile(hashPath, []byte(PromptHash(content)+"\n"), 0644)
}
