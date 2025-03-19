package bringauto_sysroot

import (
	"bringauto/modules/bringauto_prerequisites"
)

// Represents one built Package in sysroot.
type BuiltPackage struct {
	Name string
	DirName string
	GitUrl string
	GitCommitHash string
}

type builtPackageInitArgs struct {
	Name string
	DirName string
	GitUrl string
	GitCommitHash string
}

func (builtPackage *BuiltPackage) FillDefault(*bringauto_prerequisites.Args) error {
	builtPackage.Name = ""
	builtPackage.DirName = ""
	builtPackage.GitUrl = ""
	builtPackage.GitCommitHash = ""
	return nil
}

func (builtPackage *BuiltPackage) FillDynamic(args *bringauto_prerequisites.Args) error {
	if !bringauto_prerequisites.IsEmpty(args) {
		var argsStruct builtPackageInitArgs
		bringauto_prerequisites.GetArgs(args, &argsStruct)
		builtPackage.Name = argsStruct.Name
		builtPackage.DirName = argsStruct.DirName
		builtPackage.GitUrl = argsStruct.GitUrl
		builtPackage.GitCommitHash = argsStruct.GitCommitHash
	}
	return nil
}

func (builtPackage *BuiltPackage) CheckPrerequisites(*bringauto_prerequisites.Args) error {
	return nil
}
