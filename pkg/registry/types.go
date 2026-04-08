package registry

import (
	"fmt"
	"strings"
	"time"
)

const DefaultHost = "github.com"

// PackageRef is a namespaced package reference like "github.com/owner/repo".
type PackageRef struct {
	Host  string // e.g. "github.com"
	Owner string // e.g. "artpar"
	Repo  string // e.g. "calculator"
}

// ParsePackageRef parses a package reference string.
// Accepts "github.com/owner/repo" (full) or "owner/repo" (assumes github.com).
func ParsePackageRef(s string) (PackageRef, error) {
	parts := strings.Split(s, "/")
	switch len(parts) {
	case 3:
		// github.com/owner/repo
		if parts[0] == "" || parts[1] == "" || parts[2] == "" {
			return PackageRef{}, fmt.Errorf("invalid package reference: %q", s)
		}
		return PackageRef{Host: parts[0], Owner: parts[1], Repo: parts[2]}, nil
	case 2:
		// owner/repo → assume github.com
		if parts[0] == "" || parts[1] == "" {
			return PackageRef{}, fmt.Errorf("invalid package reference: %q", s)
		}
		return PackageRef{Host: DefaultHost, Owner: parts[0], Repo: parts[1]}, nil
	default:
		return PackageRef{}, fmt.Errorf("invalid package reference: %q (expected github.com/owner/repo)", s)
	}
}

// String returns the full namespaced reference: "github.com/owner/repo".
func (r PackageRef) String() string {
	return r.Host + "/" + r.Owner + "/" + r.Repo
}

// OwnerRepo returns "owner/repo" for GitHub API calls.
func (r PackageRef) OwnerRepo() string {
	return r.Owner + "/" + r.Repo
}

// PackageType is the kind of Canopy package.
type PackageType string

const (
	TypeApp       PackageType = "app"
	TypeComponent PackageType = "component"
	TypeTheme     PackageType = "theme"
	TypeFFIConfig PackageType = "ffi-config"
)

// Manifest is the canopy.json schema at the root of a package repo.
type Manifest struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Type         PackageType       `json:"type"`
	Description  string            `json:"description,omitempty"`
	Author       string            `json:"author,omitempty"`
	License      string            `json:"license,omitempty"`
	Icon         string            `json:"icon,omitempty"`
	Entry        string            `json:"entry"`
	Prompt       string            `json:"prompt,omitempty"`
	BundleID     string            `json:"bundleId,omitempty"`
	Keywords     []string          `json:"keywords,omitempty"`
	Dependencies map[string]string `json:"dependencies,omitempty"`
}

// PackageEntry tracks an installed package in the local registry.
type PackageEntry struct {
	Name        string      `json:"name"`
	Type        PackageType `json:"type"`
	Repo        string      `json:"repo"`
	Version     string      `json:"version"`
	Commit      string      `json:"commit"`
	InstalledAt time.Time   `json:"installed_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Path        string      `json:"path"`
	Entry       string      `json:"entry"`
}

// RegistryFile is the on-disk format of ~/.canopy/registry.json.
type RegistryFile struct {
	Version  int                      `json:"version"`
	Packages map[string]*PackageEntry `json:"packages"`
}

// Semver is a parsed semantic version.
type Semver struct {
	Major int
	Minor int
	Patch int
	Pre   string
}

// UpdateInfo describes an available update for an installed package.
type UpdateInfo struct {
	Name           string `json:"name"`
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
	Repo           string `json:"repo"`
}
