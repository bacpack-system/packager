package build

import (
	"github.com/bacpack-system/packager/internal/prerequisites"
	"fmt"
	"os"
	"path"
	"strings"
)

type CMake struct {
	BuildSystem  *BuildSystem
	Defines      map[string]string
	CMakeListDir string
}

func (cmake *CMake) FillDefault(*prerequisites.Args) error {
	cmake.CMakeListDir = "." + string(os.PathSeparator)
	cmake.Defines = map[string]string{}

	return nil
}

func (cmake *CMake) FillDynamic(*prerequisites.Args) error {
	return nil
}

func (cmake *CMake) CheckPrerequisites(*prerequisites.Args) error {
	for key := range cmake.Defines {
		if !validateDefineName(key) {
			return fmt.Errorf("invalid CMake define: %s", key)
		}
	}
	if cmake.BuildSystem != nil {
		_, found := cmake.Defines["CMAKE_INSTALL_PREFIX"]
		if found {
			return fmt.Errorf("do not specify CMAKE_INSTALL_PREFIX define")
		}
		_, found = cmake.Defines["CMAKE_PREFIX_PATH"]
		if found {
			return fmt.Errorf("do not specify CMAKE_PREFIX_PATH define")
		}
	}
	return nil
}

func (cmake *CMake) ConstructCMDLine() []string {
	cmake.UpdateDefines()
	var cmdLine []string
	cmdLine = append(cmdLine, "cmake")
	for key, value := range cmake.Defines {
		valuePair := "-D" + key + "=" + escapeDefineValue(value)
		cmdLine = append(cmdLine, valuePair)
	}
	cmdLine = append(cmdLine, path.Join(cmake.BuildSystem.SourceDir, cmake.CMakeListDir))
	return []string{strings.Join(cmdLine, " ")}
}

func (cmake *CMake) UpdateDefines() {
	cmake.Defines["CMAKE_INSTALL_PREFIX"] = cmake.BuildSystem.InstallPrefix
	cmake.Defines["CMAKE_PREFIX_PATH"] = cmake.BuildSystem.PrefixPath
}
