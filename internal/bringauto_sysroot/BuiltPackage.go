package bringauto_sysroot

import (
	"bringauto/internal/bringauto_prerequisites"
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

func (builtPackage *BuiltPackage) FillDefault(*bringauto_prerequisites.Args) error {
	builtPackage.Name = ""
	builtPackage.DirName = ""
	builtPackage.GitUri = ""
	builtPackage.GitCommitHash = ""
	return nil
}

func (builtPackage *BuiltPackage) FillDynamic(args *bringauto_prerequisites.Args) error {
	if !bringauto_prerequisites.IsEmpty(args) {
		var argsStruct builtPackageInitArgs
		bringauto_prerequisites.GetArgs(args, &argsStruct)
		builtPackage.Name = argsStruct.Name
		builtPackage.DirName = argsStruct.DirName
		builtPackage.GitUri = argsStruct.GitUri
		builtPackage.GitCommitHash = argsStruct.GitCommitHash
	}
	return nil
}

func (builtPackage *BuiltPackage) CheckPrerequisites(*bringauto_prerequisites.Args) error {
	return nil
}
