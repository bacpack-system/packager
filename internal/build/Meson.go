package build

import (
	"github.com/bacpack-system/packager/internal/prerequisites"
	"fmt"
	"strings"
)

const (
	mesonBuildDirConst = "build"
)

type Meson struct {
	BuildSystem *BuildSystem
	Options      map[string]string
	Defines      map[string]string
}

func (meson *Meson) FillDefault(*prerequisites.Args) error {
	meson.Options = map[string]string{}
	meson.Defines = map[string]string{}
	return nil
}

func (meson *Meson) FillDynamic(*prerequisites.Args) error {
	return nil
}

func (meson *Meson) CheckPrerequisites(*prerequisites.Args) error {
	for key := range meson.Options {
		if !validateOptionName(key) {
			return fmt.Errorf("invalid Meson option: %s", key)
		}
	}
	for key := range meson.Defines {
		if !validateDefineName(key) {
			return fmt.Errorf("invalid Meson define: %s", key)
		}
	}
	if meson.BuildSystem == nil {
		return fmt.Errorf("BuildSystem is required for Meson")
	}
	_, found := meson.Options["prefix"]
	if found {
		return fmt.Errorf("do not specify prefix option")
	}
	_, found = meson.Options["cmake-prefix-path"]
	if found {
		return fmt.Errorf("do not specify cmake-prefix-path option")
	}
	return nil
}

func (meson *Meson) ConstructCMDLine() []string {
	meson.UpdateOptions()
	var cmdSetup []string
	cmdSetup = append(cmdSetup, "meson")
	cmdSetup = append(cmdSetup, "setup")
	for key, value := range meson.Options {
		valuePair := "--" + key + "=" + value
		cmdSetup = append(cmdSetup, valuePair)
	}
	for key, value := range meson.Defines {
		valuePair := "-D" + key + "=" + escapeDefineValue(value)
		cmdSetup = append(cmdSetup, valuePair)
	}
	cmdSetup = append(cmdSetup, mesonBuildDirConst)
	cmdSetup = append(cmdSetup, meson.BuildSystem.SourceDir)

	cmdInstall := []string{"meson", "install", "-C", mesonBuildDirConst}
	
	return []string{
		strings.Join(cmdSetup, " "),
		strings.Join(cmdInstall, " "),
	}
}

func (meson *Meson) UpdateOptions() {
	meson.Options["prefix"] = meson.BuildSystem.InstallPrefix
	meson.Options["cmake-prefix-path"] = meson.BuildSystem.PrefixPath
}
