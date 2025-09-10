package repository

import "github.com/bacpack-system/packager/internal/bacpack_package"

// RepositoryInterface is an interface for every package repository.
// If you want to implement
type RepositoryInterface interface {
	// CopyToRepository copy package files represented by sourceDir to the repository.
	// Each repository has a different semantics for managing structure of th repository.
	//
	// Repository must not change the package name represented by pack.GetFullPackageName()
	CopyToRepository(pack bacpack_package.Package, sourceDir string) error
}
