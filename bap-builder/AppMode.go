package main

import (
	"bringauto/modules/bringauto_error"
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
	contextManager := bringauto_context.ContextManager{
		ContextPath: contextPath,
		ForPackage: false,
	}
	err = bringauto_prerequisites.Initialize(&contextManager)
	if err != nil {
		logger := bringauto_log.GetLogger()
		logger.Error("Context consistency error - %s", err)
		return bringauto_error.ContextErr
	}
	err = performPreBuildChecks(&repo, &contextManager, platformString, *cmdLine.DockerImageName)
	if err != nil {
		return err
	}

	handleRemover := bringauto_process.SignalHandlerAddHandler(repo.RestoreAllChanges)
	defer handleRemover()

	if *cmdLine.All {
		return buildAllApps(cmdLine, &contextManager, platformString, repo)
	} else {
		return buildSingleApp(cmdLine, &contextManager, platformString, repo)
	}
}

// buildAllApps
// Builds all Apps specified in contextPath. Returns nil if everything is ok, else returns error.
func buildAllApps(
	cmdLine        *BuildAppCmdLineArgs,
	contextManager *bringauto_context.ContextManager,
	platformString *bringauto_package.PlatformString,
	repo           bringauto_repository.GitLFSRepository,
) error {
	configMap := contextManager.GetAllConfigsMap()

	count := int32(0)
	for appName := range configMap {
		for _, config := range configMap[appName] {
			buildConfigs, err := config.GetBuildStructure(
				*cmdLine.DockerImageName,
				platformString,
				uint16(*cmdLine.Port),
				*cmdLine.UseLocalRepo,
				repo.GitRepoPath,
			)
			if err != nil {
				return err
			}
			if len(buildConfigs) == 0 {
				continue
			}
			count++
			err = buildAndCopyPackage(&buildConfigs, platformString, repo, bringauto_const.AppDirName)
			if err != nil {
				return fmt.Errorf("cannot build App '%s' - %w", config.Package.Name, err)
			}
		}
		err := bringauto_sysroot.RemoveInstallSysroot()
		if err != nil {
			return fmt.Errorf("cannot remove install sysroot directory")
		}
	}
	if count == 0 {
		return fmt.Errorf("no Apps to build for %s image", *cmdLine.DockerImageName)
	}

	return nil
}

// buildSingleApp
// Builds single App specified by name in cmdLine. Returns nil if everything is ok, else returns error.
func buildSingleApp(
	cmdLine        *BuildAppCmdLineArgs,
	contextManager *bringauto_context.ContextManager,
	platformString *bringauto_package.PlatformString,
	repo           bringauto_repository.GitLFSRepository,
) error {
	configList, err := prepareConfigsNoBuildDeps(*cmdLine.Name, contextManager, platformString, bringauto_const.AppDirName)
	if err != nil {
		return err
	}
	if len(configList) == 0 {
		return fmt.Errorf("nothing to build")
	}
	for _, config := range configList {
		if !slices.Contains(config.DockerMatrix.ImageNames, *cmdLine.DockerImageName) {
			return fmt.Errorf("'%s' does not support %s image", config.Package.Name, *cmdLine.DockerImageName)
		}
		buildConfigs, err := config.GetBuildStructure(
			*cmdLine.DockerImageName,
			platformString,
			uint16(*cmdLine.Port),
			*cmdLine.UseLocalRepo,
			repo.GitRepoPath,
		)
		if err != nil {
			return err
		}
		err = buildAndCopyPackage(&buildConfigs, platformString, repo, bringauto_const.AppDirName)
		if err != nil {
			return fmt.Errorf("cannot build App '%s' - %w", *cmdLine.Name, err)
		}
	}
	return nil
}
