package build

import (
	"github.com/bacpack-system/packager/internal/prerequisites"
	"fmt"
	"regexp"
)

// BuildSystem
// Represents build system. Its main task is to create command line for specific build system.
type BuildSystem struct {
	SourceDir     string
	InstallPrefix string
	PrefixPath    string
	CMake         *CMake
	GNUMake       *GNUMake
	Meson         *Meson
}

var optionRegexp *regexp.Regexp = regexp.MustCompilePOSIX("^[0-9a-zA-Z-]+$")

// FillDefault
// It fills up defaults for all members in the Build structure.
func (buildSystem *BuildSystem) FillDefault(args *prerequisites.Args) error {
	return nil
}

func (buildSystem *BuildSystem) FillDynamic(*prerequisites.Args) error {
	var err error
	if buildSystem.CMake != nil {
		buildSystem.GNUMake, err = prerequisites.CreateAndInitialize[GNUMake]()
		if err != nil {
			return err
		}
	}
	buildSystem.UpdateBuildSystemPointers()
	return nil
}

func (buildSystem *BuildSystem) CheckPrerequisites(*prerequisites.Args) error {
	if buildSystem.CMake != nil && buildSystem.Meson != nil {
		return fmt.Errorf("more than one build system specified")
	} else if buildSystem.CMake != nil {
		return buildSystem.CMake.CheckPrerequisites(nil)
	} else if buildSystem.Meson != nil {
		return buildSystem.Meson.CheckPrerequisites(nil)
	}
	return nil
}

func (buildSystem *BuildSystem) ConstructCMDLine() []string {
	if buildSystem.SourceDir == "" {
		panic(fmt.Errorf("source directory is empty"))
	}
	if buildSystem.CMake != nil {
		commands := buildSystem.CMake.ConstructCMDLine()
		commands = append(commands, buildSystem.GNUMake.ConstructCMDLine()...)
		return commands
	} else if buildSystem.Meson != nil {
		return buildSystem.Meson.ConstructCMDLine()
	} else {
		return []string{}
	}
}

// UpdateBuildSystemPointers updates pointers to BuildSystem in all build system specific structs
func (buildSystem *BuildSystem) UpdateBuildSystemPointers() {
	if buildSystem.CMake != nil {
		buildSystem.CMake.BuildSystem = buildSystem
	}
	if buildSystem.Meson != nil {
		buildSystem.Meson.BuildSystem = buildSystem
	}
}

func escapeDefineValue(varValue string) string {
	return "\"" + varValue + "\""
}

func validateOptionName(varName string) bool {
	return optionRegexp.MatchString(varName)
}