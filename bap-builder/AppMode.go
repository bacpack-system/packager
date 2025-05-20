package main

import (
	"bringauto/modules/bringauto_config"
	"bringauto/modules/bringauto_const"
	"bringauto/modules/bringauto_context"
	"bringauto/modules/bringauto_log"
	"bringauto/modules/bringauto_package"
	"bringauto/modules/bringauto_prerequisites"
	"bringauto/modules/bringauto_process"
	"bringauto/modules/bringauto_repository"
	"bringauto/modules/bringauto_sysroot"
	"fmt"
	"slices"
)

// BuildApp
func BuildApp(cmdLine *BuildAppCmdLineArgs, contextPath string) error {
	platformString, err := determinePlatformString(*cmdLine.DockerImageName, uint16(*cmdLine.Port))
	if err != nil {
		return err
	}
	repo := bringauto_repository.GitLFSRepository{
		GitRepoPath: *cmdLine.OutputDir,
	}
	err = bringauto_prerequisites.Initialize(&repo)
	if err != nil {
		return err
	}
	err = performPreBuildChecks(contextPath, &repo, platformString, *cmdLine.DockerImageName)
	if err != nil {
		return err
	}

	handleRemover := bringauto_process.SignalHandlerAddHandler(repo.RestoreAllChanges)
	defer handleRemover()

	if *cmdLine.All {
		err = buildAllApps(*cmdLine.DockerImageName, contextPath, platformString, repo, uint16(*cmdLine.Port))
	} else {
		err = buildSingleApp(cmdLine, contextPath, platformString, repo, uint16(*cmdLine.Port))
	}
	if err != nil {
		return err
	}
	
	return nil
}

// buildAllApps
// Builds all Apps specified in contextPath. Returns nil if everything is ok, else returns error.
func buildAllApps(
	imageName      string,
	contextPath    string,
	platformString *bringauto_package.PlatformString,
	repo           bringauto_repository.GitLFSRepository,
	dockerPort     uint16,
) error {
	contextManager := bringauto_context.ContextManager{
		ContextPath: contextPath,
	}
	appJsonPathMap, err := contextManager.GetAllConfigJsonPaths(bringauto_const.AppDirName)
	if err != nil {
		return err
	}

	defsMap := make(ConfigMapType)
	for _, appJsonPathList := range appJsonPathMap {
		addConfigsToDefsMap(&defsMap, appJsonPathList)
	}

	logger := bringauto_log.GetLogger()

	count := int32(0)
	for appName := range defsMap {
		for _, config := range defsMap[appName] {
			if isDepsInConfig(config) {
				return fmt.Errorf("App has non-empty DependsOn")
			}
			buildConfigs := config.GetBuildStructure(imageName, platformString, dockerPort)
			if len(buildConfigs) == 0 {
				continue
			}
			count++
			err := buildAndCopyPackage(&buildConfigs, platformString, repo, bringauto_const.AppDirName)
			if err != nil {
				return fmt.Errorf("cannot build App '%s' - %s", config.Package.Name, err)
			}
		}
		err = bringauto_sysroot.RemoveInstallSysroot()
		if err != nil {
			return fmt.Errorf("cannot remove install sysroot directory")
		}
	}
	if count == 0 {
		logger.Warn("Nothing to build. Did you enter correct image name?")
	}

	return nil
}

// buildSingleApp
// Builds single App specified by name in cmdLine. Returns nil if everything is ok, else returns error.
func buildSingleApp(
	cmdLine        *BuildAppCmdLineArgs,
	contextPath    string,
	platformString *bringauto_package.PlatformString,
	repo           bringauto_repository.GitLFSRepository,
	dockerPort     uint16,
) error{
	contextManager := bringauto_context.ContextManager{
		ContextPath: contextPath,
	}

	configList, err := prepareConfigsNoBuildDeps(*cmdLine.Name, &contextManager, bringauto_const.AppDirName)
	if err != nil {
		return err
	}
	if len(configList) == 0 {
		return fmt.Errorf("nothing to build")
	}
	for _, config := range configList {
		if isDepsInConfig(config) {
			return fmt.Errorf("App has non-empty DependsOn")
		}
		if !slices.Contains(config.DockerMatrix.ImageNames, *cmdLine.DockerImageName) {
			return fmt.Errorf("'%s' does not support %s image", config.Package.Name, *cmdLine.DockerImageName)
		}
		buildConfigs := config.GetBuildStructure(*cmdLine.DockerImageName, platformString, dockerPort)
		err := buildAndCopyPackage(&buildConfigs, platformString, repo, bringauto_const.AppDirName)
		if err != nil {
			return fmt.Errorf("cannot build App '%s' - %s", *cmdLine.Name, err)
		}
	}
	return nil
}

// isDepsInConfig
// Returns true if given config has non-empty DependsOn array, else returns false.
func isDepsInConfig(config *bringauto_config.Config) bool {
	return len(config.DependsOn) > 0
}
