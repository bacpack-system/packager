package sysroot

import (
	"github.com/bacpack-system/packager/internal/prerequisites"
)

// Represents one built Package in sysroot. Before build of a Package, this struct is created and
// compared to already built Packages in sysroot (a json file in sysroot directory), if any field
// is changed, the build is performed, else it is skipped. At the end the struct in form of a json
// is saved to sysroot directory.
type BuiltPackage struct {
	Name string
	DirName string
	GitUri string
	GitCommitHash string
}

type builtPackageInitArgs BuiltPackage

func (builtPackage *BuiltPackage) FillDefault(*prerequisites.Args) error {
	builtPackage.Name = ""
	builtPackage.DirName = ""
	builtPackage.GitUri = ""
	builtPackage.GitCommitHash = ""
	return nil
}

func (builtPackage *BuiltPackage) FillDynamic(args *prerequisites.Args) error {
	if !prerequisites.IsEmpty(args) {
		var argsStruct builtPackageInitArgs
		prerequisites.GetArgs(args, &argsStruct)
		builtPackage.Name = argsStruct.Name
		builtPackage.DirName = argsStruct.DirName
		builtPackage.GitUri = argsStruct.GitUri
		builtPackage.GitCommitHash = argsStruct.GitCommitHash
	}
	return nil
}

func (builtPackage *BuiltPackage) CheckPrerequisites(*prerequisites.Args) error {
	return nil
}
