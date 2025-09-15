package sysroot

import (
	"github.com/bacpack-system/packager/internal/constants"
	"github.com/bacpack-system/packager/internal/log"
	"encoding/json"
	"fmt"
	"os"
	"path"
)

const (
	jsonFileName = "built_packages.json"
	indent = "\x20\x20\x20\x20" // four spaces
)

// Contains built Packages in sysroot and has functions for Json encoding and decoding of built
// Packages.
type BuiltPackages struct {
	Packages []BuiltPackage
}

// AddToBuiltPackages
// Adds packageName to built Packages.
func (builtPackages *BuiltPackages) AddToBuiltPackages(pack BuiltPackage) error {
	err := builtPackages.updateBuiltPackages()
	if err != nil {
		return fmt.Errorf("can't update builtPackages from json - %w", err)
	}
	builtPackages.Packages = append(builtPackages.Packages, pack)
	bytes, err := json.MarshalIndent(builtPackages.Packages, "", indent)
	if err != nil {
		return err
	}
	err = os.WriteFile(path.Join(sysrootDirectoryName, jsonFileName), bytes, 0644)
	return err
}

// UpdateBuiltPackages
// Updates builtPackages struct based on built_packages.json.
func (builtPackages *BuiltPackages) updateBuiltPackages() error {
	bytes, err := os.ReadFile(path.Join(sysrootDirectoryName, jsonFileName))
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to read built packages file - %w", err)
	}

	err = json.Unmarshal(bytes, &builtPackages.Packages)
	if err != nil {
		return fmt.Errorf("failed to parse built packages file - %w", err)
	}
	return nil
}

// Contains
// Returns true if given Package is in builtPackages, else false. All fields of BuiltPackage struct
// are compared. Only if the pack has empty GitCommitHash, the GitCommitHash is not compared.
func (builtPackages *BuiltPackages) Contains(pack BuiltPackage) bool {
	err := builtPackages.updateBuiltPackages()
	if err != nil {
		logger := log.GetLogger()
		logger.Error("Can't update builtPackages from json - %s", err)
	}
	for _, p := range builtPackages.Packages {
		condition := p.Name == pack.Name && p.DirName == pack.DirName && pack.GitUri == p.GitUri
		if pack.GitCommitHash != constants.EmptyGitCommitHash {
			condition = condition && pack.GitCommitHash == p.GitCommitHash
		}
		if condition {
			return true
		}
	}
	return false
}
