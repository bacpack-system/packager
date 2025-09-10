package main

import (
	"github.com/bacpack-system/packager/internal/bringauto_build"
	"github.com/bacpack-system/packager/internal/bringauto_config"
	"github.com/bacpack-system/packager/internal/bringauto_const"
	"github.com/bacpack-system/packager/internal/bringauto_context"
	"github.com/bacpack-system/packager/internal/bringauto_docker"
	"github.com/bacpack-system/packager/internal/bringauto_log"
	"github.com/bacpack-system/packager/internal/bringauto_package"
	"github.com/bacpack-system/packager/internal/bringauto_prerequisites"
	"github.com/bacpack-system/packager/internal/bringauto_process"
	"github.com/bacpack-system/packager/internal/bringauto_repository"
	"github.com/bacpack-system/packager/internal/bringauto_ssh"
	"github.com/bacpack-system/packager/internal/bringauto_sysroot"
	"github.com/bacpack-system/packager/internal/bringauto_error"
	"fmt"
	"strconv"
	"slices"
)

type buildDepList struct {
	dependsMap map[string]*map[string]bool
}

func removeDuplicates(configList *[]bringauto_config.Config) []bringauto_config.Config {
	var newConfigList []bringauto_config.Config
	packageMap := make(map[string]bool)
	for _, cconfig := range *configList {
		packageName := cconfig.Package.Name + ":" + strconv.FormatBool(cconfig.Package.IsDebug)
		exist, _ := packageMap[packageName]
		if exist {
			continue
		}
		packageMap[packageName] = true
		newConfigList = append(newConfigList, cconfig)
	}
	return newConfigList
}

func (list *buildDepList) TopologicalSort(buildMap bringauto_context.ConfigMapType) ([]bringauto_config.Config, error) {

	// Map represents 'PackageName: []DependsOnPackageNames'
	var dependsMap map[string]*map[string]bool
	var allDependencies map[string]bool

	dependsMap, allDependencies, err := bringauto_context.CreateDependsMap(&buildMap)
	if err != nil {
		return []bringauto_config.Config{}, err
	}

	dependsMapCopy := make(map[string]*map[string]bool, len(dependsMap))
	for key, value := range dependsMap {
		dependsMapCopy[key] = value
	}
	var rootList []string
	for dependencyName, _ := range allDependencies {
		delete(dependsMapCopy, dependencyName)
	}
	for key, _ := range dependsMapCopy {
		rootList = append(rootList, key)
	}

	var sortedDependencies []string
	for key, _ := range dependsMapCopy {
		sortedDependencies = append(sortedDependencies, *list.sortDependencies(key, &dependsMap)...)
	}

	sortedLen := len(sortedDependencies)
	sortedReverse := make([]string, sortedLen)
	for i := sortedLen - 1; i >= 0; i-- {
		sortedReverse[i] = sortedDependencies[sortedLen-i-1]
	}

	var sortedDependenciesConfig []bringauto_config.Config
	for _, packageName := range sortedReverse {
		sortedDependenciesConfig = append(sortedDependenciesConfig, buildMap[packageName]...)
	}

	return removeDuplicates(&sortedDependenciesConfig), nil
}

func (list *buildDepList) sortDependencies(rootName string, dependsMap *map[string]*map[string]bool) *[]string {
	sorted := []string{rootName}
	rootDeps, found := (*dependsMap)[rootName]

	if !found || len(*rootDeps) == 0 {
		return &sorted
	}

	for packageName, _ := range *rootDeps {
		packageDeps := list.sortDependencies(packageName, dependsMap)
		sorted = append(sorted, *packageDeps...)
	}

	return &sorted
}

// performPreBuildChecks
// Performs Git lfs and sysroot consistency checks. This should be called before builds.
func performPreBuildChecks(
	repo           *bringauto_repository.GitLFSRepository,
	contextManager *bringauto_context.ContextManager,
	platformString *bringauto_package.PlatformString,
	imageName      string,
) error {
	logger := bringauto_log.GetLogger()
	logger.Info("Checking Git Lfs directory consistency")
	err := repo.CheckGitLfsConsistency(contextManager, platformString, imageName)
	if err != nil {
		logger.Error("Git Lfs consistency error - %s", err)
		return bringauto_error.GitLfsErr
	}
	logger.Info("Checking Sysroot directory consistency")
	err = checkSysrootDirs(platformString)
	if err != nil {
		return err
	}
	return nil
}

// BuildPackage
// process Package mode of the program
func BuildPackage(cmdLine *BuildPackageCmdLineArgs, contextPath string) error {
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
		ForPackage: true,
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
		return buildAllPackages(cmdLine, &contextManager, platformString, repo)
	} else {
		return buildSinglePackage(cmdLine, &contextManager, platformString, repo)
	}
}

// buildAllPackages
// Builds all packages specified in contextPath. Also takes care of building all deps for all
// packages in correct order. It returns nil if everything is ok, or not nil in case of error.
func buildAllPackages(
	cmdLine        *BuildPackageCmdLineArgs,
	contextManager *bringauto_context.ContextManager,
	platformString *bringauto_package.PlatformString,
	repo           bringauto_repository.GitLFSRepository,
) error {
	configMap := contextManager.GetAllConfigsMap()

	depsList := buildDepList{}
	configList, err := depsList.TopologicalSort(configMap)
	if err != nil {
		return err
	}

	count := int32(0)
	for _, config := range configList {
		buildConfigs, err := config.GetBuildStructure(
			*cmdLine.DockerImageName,
			platformString,
			uint16(*cmdLine.Port),
			false,
			"",
		)
		if err != nil {
			return err
		}
		if len(buildConfigs) == 0 {
			continue
		}
		count++
		err = buildAndCopyPackage(&buildConfigs, platformString, repo, bringauto_const.PackageDirName)
		if err != nil {
			return fmt.Errorf("cannot build package '%s' - %w", config.Package.Name, err)
		}
	}
	if count == 0 {
		return fmt.Errorf("no Packages to build for %s image", *cmdLine.DockerImageName)
	}

	return nil
}

// prepareConfigs
// Returns sorted packageConfigs list.
func sortConfigs(packageConfigs []bringauto_config.Config) ([]bringauto_config.Config, error) {
	var configList []bringauto_config.Config
	defsMap := make(bringauto_context.ConfigMapType)
	addConfigsToDefsMap(&defsMap, packageConfigs)
	depList := buildDepList{}
	configList, err := depList.TopologicalSort(defsMap)
	if err != nil {
		return []bringauto_config.Config{}, err
	}
	return configList, nil
}

// prepareConfigsNoBuildDeps
// Returns Config structures only for given Package or App (depends on packageOrApp).
func prepareConfigsNoBuildDeps(
	packageName    string,
	contextManager *bringauto_context.ContextManager,
	platformString *bringauto_package.PlatformString,
	packageOrApp   string,
) ([]bringauto_config.Config, error) {
	logger := bringauto_log.GetLogger()

	if packageOrApp == bringauto_const.PackageDirName {
		value, err := isPackageDepsInSysroot(packageName, contextManager, platformString, false)
		if err != nil {
			return []bringauto_config.Config{}, err
		}
		if !value {
			logger.Error("Package dependencies are not in sysroot")
			return []bringauto_config.Config{}, bringauto_error.PackageMissingDependencyErr
		}
	}

	packageConfigs, err := contextManager.GetPackageConfigs(packageName)
	if err != nil {
		return []bringauto_config.Config{}, err
	}

	return packageConfigs, nil
}

// prepareConfigsBuildDepsOrBuildDepsOn
// Returns Config structures based on --build-deps and --build-deps-on flags.
func prepareConfigsBuildDepsOrBuildDepsOn(
	cmdLine        *BuildPackageCmdLineArgs,
	packageName    string,
	contextManager *bringauto_context.ContextManager,
	platformString *bringauto_package.PlatformString,
) ([]bringauto_config.Config, error) {
	var packageConfigs []bringauto_config.Config

	logger := bringauto_log.GetLogger()

	if *cmdLine.BuildDeps {
		configs, err := contextManager.GetPackageWithDepsConfigs(packageName)
		if err != nil {
			return []bringauto_config.Config{}, err
		}
		packageConfigs = append(packageConfigs, configs...)
	} else if *cmdLine.BuildDepsOn || *cmdLine.BuildDepsOnRecursive {
		value, err := isPackageDepsInSysroot(packageName, contextManager, platformString, true)
		if err != nil {
			return []bringauto_config.Config{}, err
		}
		if !value {
			logger.Error("--build-deps-on(-recursive) set but base package or its dependencies are not in sysroot")
			return []bringauto_config.Config{}, bringauto_error.PackageMissingDependencyErr
		}
	}
	if *cmdLine.BuildDepsOn || *cmdLine.BuildDepsOnRecursive {
		configs, err := contextManager.GetPackageWithDepsOnConfigs(packageName, *cmdLine.BuildDepsOnRecursive)
		if err != nil {
			return []bringauto_config.Config{}, err
		}
		if len(configs) == 0 {
			logger.Warn("No package depends on %s", packageName)
		}
		packageConfigs = append(packageConfigs, configs...)
	}
	return sortConfigs(packageConfigs)
}

// buildSinglePackage
// Builds single package specified by name in cmdLine. Also takes care of building all deps for
// given package in correct order. It returns nil if everything is ok, or not nil in case of error.
func buildSinglePackage(
	cmdLine        *BuildPackageCmdLineArgs,
	contextManager *bringauto_context.ContextManager,
	platformString *bringauto_package.PlatformString,
	repo           bringauto_repository.GitLFSRepository,
) error {
	packageName := *cmdLine.Name
	var err error
	var configList []bringauto_config.Config

	if *cmdLine.BuildDeps || *cmdLine.BuildDepsOn || *cmdLine.BuildDepsOnRecursive {
		configList, err = prepareConfigsBuildDepsOrBuildDepsOn(cmdLine, packageName, contextManager, platformString)
	} else {
		configList, err = prepareConfigsNoBuildDeps(packageName, contextManager, platformString, bringauto_const.PackageDirName)
	}
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
			false,
			"",
		)
		if err != nil {
			return err
		}
		err = buildAndCopyPackage(&buildConfigs, platformString, repo, bringauto_const.PackageDirName)
		if err != nil {
			return fmt.Errorf("cannot build package '%s' - %w", config.Package.Name, err)
		}
	}
	return nil
}

// addConfigsToDefsMap
// Adds Configs in packageConfigs to defsMap.
func addConfigsToDefsMap(defsMap *bringauto_context.ConfigMapType, packageConfigs []bringauto_config.Config) {
	for _, config := range packageConfigs {
		packageName := config.Package.Name
		_, found := (*defsMap)[packageName]
		if !found {
			(*defsMap)[packageName] = []bringauto_config.Config{}
		}
		(*defsMap)[packageName] = append((*defsMap)[packageName], config)
	}
}

// buildAndCopyPackage
// Builds single Package or App (depends on packageOrApp), takes care of every step of build for
// single package.
func buildAndCopyPackage(
	build          *[]bringauto_build.Build,
	platformString *bringauto_package.PlatformString,
	repo           bringauto_repository.GitLFSRepository,
	packageOrApp   string,
) error {
	var err error
	var removeHandler func()

	logger := bringauto_log.GetLogger()

	for _, buildConfig := range *build {
		logger.Info("Build %s", buildConfig.Package.GetFullPackageName())

		sysroot := bringauto_sysroot.Sysroot{
			IsDebug:        buildConfig.Package.IsDebug,
			PlatformString: platformString,
		}
		err = bringauto_prerequisites.Initialize(&sysroot)
		buildConfig.SetSysroot(&sysroot)

		removeHandler = bringauto_process.SignalHandlerAddHandler(buildConfig.CleanUp)
		var buildPerformed bool
		err, buildPerformed = buildConfig.RunBuild()
		if err != nil {
			return bringauto_error.BuildErr
		}

		if buildPerformed {
			logger.InfoIndent("Copying to local sysroot directory")
			err = sysroot.CopyToSysroot(buildConfig.GetLocalInstallDirPath(), *buildConfig.BuiltPackage)
			if err != nil {
				break
			}
			
			logger.InfoIndent("Copying to Git repository")
			err = repo.CopyToRepository(*buildConfig.Package, buildConfig.GetLocalInstallDirPath(), packageOrApp)
			if err != nil {
				break
			}
		}

		removeHandler()
		removeHandler = nil
		logger.InfoIndent("Build OK")
	}
	if removeHandler != nil {
		removeHandler()
	}
	return err
}

// determinePlatformString
// Will construct platform string suitable for sysroot.
func determinePlatformString(dockerImageName string, dockerPort uint16) (*bringauto_package.PlatformString, error) {
	defaultDocker, err := bringauto_prerequisites.CreateAndInitialize[bringauto_docker.Docker](dockerImageName, dockerPort)
	if err != nil {
		return nil, err
	}
	defaultDocker.ImageName = dockerImageName

	sshCreds, err := bringauto_prerequisites.CreateAndInitialize[bringauto_ssh.SSHCredentials]()
	if err != nil {
		return nil, err
	}
	sshCreds.Port = uint16(defaultDocker.Port)

	platformString := bringauto_package.PlatformString{
		Mode: bringauto_package.ModeAuto,
	}

	err = bringauto_prerequisites.Initialize[bringauto_package.PlatformString](&platformString, sshCreds, defaultDocker)
	return &platformString, err
}

// checkSysrootDirs
// Checks if sysroot release and debug directories are empty. If not, prints a warning.
func checkSysrootDirs(platformString *bringauto_package.PlatformString) (error) {
	sysroot := bringauto_sysroot.Sysroot{
		IsDebug:        false,
		PlatformString: platformString,
	}
	err := bringauto_prerequisites.Initialize(&sysroot)
	if err != nil {
		return err
	}

	logger := bringauto_log.GetLogger()
	if !sysroot.IsSysrootDirectoryEmpty() {
		logger.Warn("Sysroot release directory is not empty - the package build may fail")
	}
	sysroot.IsDebug = true
	if !sysroot.IsSysrootDirectoryEmpty() {
		logger.Warn("Sysroot debug directory is not empty - the package build may fail")
	}
	return nil
}

// isPackageWithDepsInSysroot
// Returns true if packageName an its dependencies are in sysroot, else returns false. If the
// checkPackageItself is true, it also checks for presence of Package itself. 
func isPackageDepsInSysroot(
	packageName        string,
	contextManager     *bringauto_context.ContextManager,
	platformString     *bringauto_package.PlatformString,
	checkPackageItself bool,
) (bool, error) {
	configMap, err := contextManager.GetPackageWithDepsConfigs(packageName)
	if err != nil {
		return false, err
	}

	sysroot := bringauto_sysroot.Sysroot{
		IsDebug:        false,
		PlatformString: platformString,
	}
	err = bringauto_prerequisites.Initialize(&sysroot)
	if err != nil {
		return false, err
	}

	for _, config := range configMap {
		if !checkPackageItself && config.Package.Name == packageName {
			continue
		}
		sysroot.IsDebug = config.Package.IsDebug
		builtPackage := bringauto_sysroot.BuiltPackage {
			Name: config.Package.GetShortPackageName(),
			DirName: sysroot.GetDirNameInSysroot(),
			GitUri: config.Git.URI,
			GitCommitHash: bringauto_const.EmptyGitCommitHash,
		}
		if !sysroot.IsPackageInSysroot(builtPackage) {
			return false, nil
		}
	}

	return true, nil
}
