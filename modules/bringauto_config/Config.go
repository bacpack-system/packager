package bringauto_config

import (
	"bringauto/modules/bringauto_build"
	"bringauto/modules/bringauto_const"
	"bringauto/modules/bringauto_docker"
	"bringauto/modules/bringauto_git"
	"bringauto/modules/bringauto_package"
	"bringauto/modules/bringauto_prerequisites"
	"bringauto/modules/bringauto_sysroot"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jinzhu/copier"
)

// Build
// It stores configuration for given build system
// (CMake, autoconf, ...)
type Build struct {
	CMake *bringauto_build.CMake
}

type DockerMatrix struct {
	ImageNames []string
}

// Config
// Build configuration which stores how the package is build.
type Config struct {
	Env          map[string]string
	Git          bringauto_git.Git
	Build        Build
	Package      bringauto_package.Package
	DockerMatrix DockerMatrix
	DependsOn    []string
}

func (config *Config) FillDefault(*bringauto_prerequisites.Args) error {
	*config = Config{
		Env:       map[string]string{},
		Git:       bringauto_git.Git{},
		Build:     Build{},
		Package:   bringauto_package.Package{},
		DependsOn: []string{},
	}
	return nil
}

func (config *Config) FillDynamic(*bringauto_prerequisites.Args) error {
	return nil
}

func (config *Config) CheckPrerequisites(*bringauto_prerequisites.Args) error {
	return nil
}

func (config *Config) LoadJSONConfig(configPath string) error {
	mbytes, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(mbytes, config)
	if err != nil {
		return err
	}
	return nil
}

func (config *Config) SaveToJSONConfig(configPath string) error {
	mbytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	err = os.WriteFile(configPath, mbytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Returns array of builds structs for specific image name. The returned array will contain max one build.
// It is an array for simple handling of result using for loop.
func (config *Config) GetBuildStructure(imageName string, platformString *bringauto_package.PlatformString) []bringauto_build.Build {
	var buildConfigs []bringauto_build.Build
	for _, value := range config.DockerMatrix.ImageNames {
		if imageName != "" && imageName != value {
			continue
		}
		build := config.fillBuildStructure(imageName, platformString)
		defaultBuild := bringauto_prerequisites.CreateAndInitialize[bringauto_build.Build](imageName)
		err := copier.CopyWithOption(defaultBuild, build, copier.Option{DeepCopy: true, IgnoreEmpty: true})
		if err != nil {
			panic(fmt.Errorf("cannot merge default and real build config"))
		}
		buildConfigs = append(buildConfigs, *defaultBuild)
	}

	return buildConfigs
}

// fillBuildStructure
// Fills and returns Build structure.
func (config *Config) fillBuildStructure(dockerImageName string, platformString *bringauto_package.PlatformString) bringauto_build.Build {
	var err error
	defaultDocker := bringauto_prerequisites.CreateAndInitialize[bringauto_docker.Docker](dockerImageName)

	env := &bringauto_build.EnvironmentVariables{
		Env: config.Env,
	}
	err = bringauto_prerequisites.Initialize(env)
	if err != nil {
		panic(err)
	}
	err = bringauto_prerequisites.Initialize(&config.Git)
	if err != nil {
		panic(err)
	}
	err = bringauto_prerequisites.Initialize(config.Build.CMake)
	if err != nil {
		panic(err)
	}
	builtPackage := bringauto_prerequisites.CreateAndInitialize[bringauto_sysroot.BuiltPackage](
		config.Package.GetShortPackageName(),
		"",                                 // Will be filled later after build will have valid sysroot
		config.Git.URI,
		bringauto_const.EmptyGitCommitHash, // Will be filled later when the hash is retrieved from docker container
	)

	tmpPackage := config.Package
	err = bringauto_prerequisites.Initialize(&tmpPackage)
	if platformString != nil {
		tmpPackage.PlatformString = *platformString
	}
	if err != nil {
		panic(err)
	}

	build := bringauto_build.Build{
		Env:          env,
		Git:          &config.Git,
		CMake:        config.Build.CMake,
		Package:      &tmpPackage,
		BuiltPackage: builtPackage,
		Docker:       defaultDocker,
	}

	return build
}
