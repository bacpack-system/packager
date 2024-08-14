package main

import (
	"bringauto/modules/bringauto_log"
	"bringauto/modules/bringauto_build"
	"bringauto/modules/bringauto_config"
	"bringauto/modules/bringauto_package"
	"bringauto/modules/bringauto_prerequisites"
	"bringauto/modules/bringauto_repository"
	"bringauto/modules/bringauto_sysroot"
	"fmt"
	"strconv"
	"os"
	"io"
)

type (
	dependsMapType      map[string]*map[string]bool
	allDependenciesType map[string]bool
	ConfigMapType       map[string][]*bringauto_config.Config
)

type buildDepList struct {
	dependsMap map[string]*map[string]bool
}

func removeDuplicates(configList *[]*bringauto_config.Config) []*bringauto_config.Config {
	var newConfigList []*bringauto_config.Config
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

func (list *buildDepList) TopologicalSort(buildMap ConfigMapType) []*bringauto_config.Config {

	// Map represents 'PackageName: []DependsOnPackageNames'
	var dependsMap map[string]*map[string]bool
	var allDependencies map[string]bool

	dependsMap, allDependencies = list.createDependsMap(&buildMap)

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

	var sortedDependenciesConfig []*bringauto_config.Config
	for _, packageName := range sortedReverse {
		sortedDependenciesConfig = append(sortedDependenciesConfig, buildMap[packageName]...)
	}

	return removeDuplicates(&sortedDependenciesConfig)
}

func (list *buildDepList) createDependsMap(buildMap *ConfigMapType) (dependsMapType, allDependenciesType) {
	allDependencies := make(map[string]bool)
	dependsMap := make(map[string]*map[string]bool)

	for _, configArray := range *buildMap {
		if len(configArray) == 0 {
			panic("invalid entry in dependency map")
		}
		packageName := configArray[0].Package.Name
		item, found := dependsMap[packageName]
		if !found {
			item = &map[string]bool{}
			dependsMap[packageName] = item
		}
		for _, config := range configArray {
			if len(config.DependsOn) == 0 {
				continue
			}
			for _, v := range config.DependsOn {
				(*item)[v] = true
				allDependencies[v] = true
			}
		}
	}
	return dependsMap, allDependencies
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

// isDirEmpty
// Returns true if specified dir do not exists or exists but is empty, otherwise returns false
func isDirEmpty(dirPath string) (bool, error) {
	f, err := os.Open(dirPath)
	if err != nil { // The directory do not exists
		return true, nil
	}
	defer f.Close()

	_, err = f.Readdirnames(1)

	if err == io.EOF { // The directory exists, but is empty
		return true, nil
	}
	return false, err
}

func checkSysrootDirectory() {
	isDirEmpty, err := isDirEmpty(bringauto_sysroot.SysrootDirectoryName)
	if !isDirEmpty {
		logger := bringauto_log.GetLogger()
		if err != nil {
			logger.Warn("Cannot read in sysroot directory: %s", err)
		} else {
			logger.Warn("Sysroot directory is not empty - the package build may fail")
		}
	}
}

// BuildPackage
// process Package mode of the program
func BuildPackage(cmdLine *BuildPackageCmdLineArgs, contextPath string) error {
	checkSysrootDirectory()
	buildAll := cmdLine.All
	if *buildAll {
		return buildAllPackages(cmdLine, contextPath)
	}
	return buildSinglePackage(cmdLine, contextPath)
}

// buildAllPackages
// builds all docker images in the given contextPath.
// It returns nil if everything is ok, or not nil in case of error
func buildAllPackages(cmdLine *BuildPackageCmdLineArgs, contextPath string) error {
	contextManager := ContextManager{
		ContextPath: contextPath,
	}
	packagesDefs, err := contextManager.GetAllPackagesJsonDefPaths()
	if err != nil {
		return err
	}

	defsList := make(map[string][]*bringauto_config.Config)
	for _, packageJsonDef := range packagesDefs {
		for _, defdef := range packageJsonDef {
			var config bringauto_config.Config
			err = config.LoadJSONConfig(defdef)
			packageName := config.Package.Name
			_, found := defsList[packageName]
			if !found {
				defsList[packageName] = []*bringauto_config.Config{}
			}
			defsList[packageName] = append(defsList[packageName], &config)
		}

	}
	depsList := buildDepList{}
	configList := depsList.TopologicalSort(defsList)

	logger := bringauto_log.GetLogger()

	count := int32(0)
	for _, config := range configList {
		buildConfigs := config.GetBuildStructure(*cmdLine.DockerImageName)
		if len(buildConfigs) == 0 {
			continue
		}
		count++
		logger.Info("Build %s", buildConfigs[0].Package.GetFullPackageName())
		err = buildAndCopyPackage(cmdLine, &buildConfigs)
		if err != nil {
			logger.Fatal("cannot build package '%s' - %s", config.Package.Name, err)
		}
	}
	if count == 0 {
		logger.Warn("Nothing to build. Did you enter correct image name?")
	}

	return nil
}

// buildSinglePackage
// build single package specified by a name
// It returns nil if everything is ok, or not nil in case of error
func buildSinglePackage(cmdLine *BuildPackageCmdLineArgs, contextPath string) error {
	contextManager := ContextManager{
		ContextPath: contextPath,
	}
	packageName := *cmdLine.Name
	packageJsonDefsList, err := contextManager.GetPackageJsonDefPaths(packageName)
	if err != nil {
		return err
	}

	logger := bringauto_log.GetLogger()

	for _, packageJsonDef := range packageJsonDefsList {
		var config bringauto_config.Config
		err = config.LoadJSONConfig(packageJsonDef)
		if err != nil {
			logger.Warn("package '%s' JSON config def problem - %s\n", packageName, err)
			continue
		}

		buildConfigs := config.GetBuildStructure(*cmdLine.DockerImageName)
		logger.Info("Build %s", buildConfigs[0].Package.GetFullPackageName())
		err = buildAndCopyPackage(cmdLine, &buildConfigs)
		if err != nil {
			logger.Error("cannot build package '%s' - %s\n", packageName, err)
			continue
		}
	}
	return nil
}

func buildAndCopyPackage(cmdLine *BuildPackageCmdLineArgs, build *[]bringauto_build.Build) error {
	if *cmdLine.OutputDirMode != OutputDirModeGitLFS {
		return fmt.Errorf("invalid OutputDirmode. Only GitLFS is supported")
	}

	var err error

	repo := bringauto_repository.GitLFSRepository{
		GitRepoPath: *cmdLine.OutputDir,
	}
	err = bringauto_prerequisites.Initialize(&repo)
	if err != nil {
		return err
	}

	for _, buildConfig := range *build {
		platformString, err := determinePlatformString(&buildConfig)
		if err != nil {
			return err
		}

		sysroot := bringauto_sysroot.Sysroot{
			IsDebug:        buildConfig.Package.IsDebug,
			PlatformString: platformString,
		}
		err = bringauto_prerequisites.Initialize(&sysroot)

		buildConfig.SetSysroot(&sysroot)
		err = buildConfig.RunBuild()
		if err != nil {
			return err
		}

		err = repo.CopyToRepository(*buildConfig.Package, buildConfig.GetLocalInstallDirPath())
		if err != nil {
			return err
		}

		err = sysroot.CopyToSysroot(buildConfig.GetLocalInstallDirPath())
		if err != nil {
			return err
		}

		err = buildConfig.CleanUp()
		if err != nil {
			return err
		}
	}
	return nil
}

// determinePlatformString will construct platform string suitable
// for sysroot.
// For example: the any_machine platformString must be copied to all machine-specific sysroot for
// a given image.
func determinePlatformString(build *bringauto_build.Build) (*bringauto_package.PlatformString, error) {
	platformStringSpecialized := build.Package.PlatformString
	if build.Package.PlatformString.Mode == bringauto_package.ModeAnyMachine {
		platformStringStruct := bringauto_package.PlatformString{
			Mode: bringauto_package.ModeAuto,
		}
		platformStringStruct.Mode = bringauto_package.ModeAuto
		err := bringauto_prerequisites.Initialize[bringauto_package.PlatformString](&platformStringStruct,
			build.SSHCredentials, build.Docker,
		)
		if err != nil {
			return nil, err
		}
		platformStringSpecialized.String.Machine = platformStringStruct.String.Machine
	}
	return &platformStringSpecialized, nil
}
