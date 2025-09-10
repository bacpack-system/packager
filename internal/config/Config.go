package config

import (
	"github.com/bacpack-system/packager/internal/build"
	"github.com/bacpack-system/packager/internal/constants"
	"github.com/bacpack-system/packager/internal/docker"
	"github.com/bacpack-system/packager/internal/git"
	"github.com/bacpack-system/packager/internal/bacpack_package"
	"github.com/bacpack-system/packager/internal/prerequisites"
	"github.com/bacpack-system/packager/internal/sysroot"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jinzhu/copier"
)

// Build
// It stores configuration for given build system
// (CMake, autoconf, ...)
type Build struct {
	CMake *build.CMake
}

type DockerMatrix struct {
	ImageNames []string
}

// Config
// Build configuration which stores how the package is build.
type Config struct {
	Env          map[string]string
	Git          git.Git
	Build        Build
	Package      bacpack_package.Package
	DockerMatrix DockerMatrix
	DependsOn    []string
}

func (config *Config) FillDefault(*prerequisites.Args) error {
	*config = Config{
		Env:       map[string]string{},
		Git:       git.Git{},
		Build:     Build{},
		Package:   bacpack_package.Package{},
		DependsOn: []string{},
	}
	return nil
}

func (config *Config) FillDynamic(*prerequisites.Args) error {
	return nil
}

func (config *Config) CheckPrerequisites(*prerequisites.Args) error {
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
func (config *Config) GetBuildStructure(
	imageName      string,
	platformString *bacpack_package.PlatformString,
	dockerPort     uint16,
	useLocalRepo   bool,
	repoPath       string,
) ([]build.Build, error) {
	var buildConfigs []build.Build
	for _, value := range config.DockerMatrix.ImageNames {
		if imageName != "" && imageName != value {
			continue
		}
		build_obj, err := config.fillBuildStructure(imageName, platformString, dockerPort, useLocalRepo, repoPath)
		if err != nil {
			return []build.Build{}, err
		}
		defaultBuild, err := prerequisites.CreateAndInitialize[build.Build](imageName, dockerPort)
		if err != nil {
			return []build.Build{}, fmt.Errorf("cannot initialize Build struct - %w", err)
		}
		err = copier.CopyWithOption(defaultBuild, build_obj, copier.Option{DeepCopy: true, IgnoreEmpty: true})
		if err != nil {
			return []build.Build{}, fmt.Errorf("cannot initialize Build struct - %w", err)
		}
		buildConfigs = append(buildConfigs, *defaultBuild)
	}

	return buildConfigs, nil
}

// fillBuildStructure
// Fills and returns Build structure.
func (config *Config) fillBuildStructure(
	dockerImageName string,
	platformString  *bacpack_package.PlatformString,
	dockerPort      uint16,
	useLocalRepo    bool,
	repoPath        string,
) (build.Build, error) {
	var err error
	defaultDocker, err := prerequisites.CreateAndInitialize[docker.Docker](dockerImageName, dockerPort)
	if err != nil {
		return build.Build{}, err
	}
	if useLocalRepo {
		err := defaultDocker.SetVolume(repoPath, constants.ContainerPackageRepoPath)
		if err != nil {
			return build.Build{}, err
		}
	}

	env := &build.EnvironmentVariables{
		Env: config.Env,
	}
	err = prerequisites.Initialize(env)
	if err != nil {
		return build.Build{}, err
	}
	err = prerequisites.Initialize(&config.Git)
	if err != nil {
		return build.Build{}, err
	}
	err = prerequisites.Initialize(config.Build.CMake)
	if err != nil {
		return build.Build{}, err
	}
	builtPackage, _ := prerequisites.CreateAndInitialize[sysroot.BuiltPackage](
		config.Package.GetShortPackageName(),
		"",                                 // Will be filled later after build will have valid sysroot
		config.Git.URI,
		constants.EmptyGitCommitHash, // Will be filled later when the hash is retrieved from docker container
	)

	tmpPackage := config.Package
	err = prerequisites.Initialize(&tmpPackage)
	if platformString != nil {
		tmpPackage.PlatformString = *platformString
	}
	if err != nil {
		return build.Build{}, err
	}

	build := build.Build{
		Env:          env,
		Git:          &config.Git,
		CMake:        config.Build.CMake,
		Package:      &tmpPackage,
		BuiltPackage: builtPackage,
		Docker:       defaultDocker,
		UseLocalRepo: useLocalRepo,
	}

	return build, nil
}
