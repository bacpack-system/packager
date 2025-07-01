package bringauto_context

import (
	"bringauto/modules/bringauto_config"
	"bringauto/modules/bringauto_const"
	"bringauto/modules/bringauto_log"
	"bringauto/modules/bringauto_package"
	"bringauto/modules/bringauto_prerequisites"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

type (
	ImagesPathType      map[string]string
	DependsMapType      map[string]*map[string]bool
	AllDependenciesType map[string]bool
	ConfigMapType       map[string][]bringauto_config.Config
	dependsMapNodeType struct {
		DependsOn  []string
		HasDebug   bool
		HasRelease bool
	}
	dependsMapTagsType  map[string]dependsMapNodeType
)

// ContextManager
// Manages all operations on the given Context. After initialization the Configs from Context are
// loaded, so the ContextPath and ForPackage can't be changed.
type ContextManager struct {
	ContextPath string
	// ForPackage boolean value if the Context is used for Packages or Apps
	ForPackage     bool
	images         ImagesPathType
	configs        *ConfigMapType
	appConfigs     ConfigMapType
	packageConfigs ConfigMapType
}

func (context *ContextManager) FillDefault(*bringauto_prerequisites.Args) error {
	context.ContextPath = ""
	context.ForPackage = true
	return nil
}

func (context *ContextManager) FillDynamic(*bringauto_prerequisites.Args) error {
	return nil
}

func (context *ContextManager) CheckPrerequisites(*bringauto_prerequisites.Args) error {
	logger := bringauto_log.GetLogger()
	logger.Info("Checking Context (%s) consistency", context.ContextPath)

	err := context.validateContextPath()
	if err != nil {
		return err
	}

	err = context.loadImagesDockerfilePaths()
	if err != nil {
		return err
	}

	packageConfigs, appConfigs, err := context.loadConfigs()
	if err != nil {
		return err
	}

	err = checkPackageConfigs(&packageConfigs)
	if err != nil {
		return err
	}

	err = checkAppConfigs(&appConfigs)
	if err != nil {
		return err
	}
	context.packageConfigs = packageConfigs
	context.appConfigs = appConfigs

	if context.ForPackage {
		context.configs = &context.packageConfigs
	} else {
		context.configs = &context.appConfigs
	}

	err = context.checkAllConfigs()
	if err != nil {
		return err
	}

	return nil
}

// loadConfigs
// Loads Configs from Context into Context Manager structs. Checks if all directories in contextPath
// have same name as Package names from JSON definitions inside this directory. If not, returns
// error with description, else returns nil. Also returns error if the Package JSON definition
// can't be loaded.
func (context *ContextManager) loadConfigs() (ConfigMapType, ConfigMapType, error) {
	packageConfigs := make(ConfigMapType)
	appConfigs := make(ConfigMapType)
	err := filepath.WalkDir(context.ContextPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			section := filepath.Base(filepath.Dir(filepath.Dir(path)))
			if section != bringauto_const.PackageDirName && section != bringauto_const.AppDirName {
				return nil
			}
			var config bringauto_config.Config
			err = config.LoadJSONConfig(path)
			if err != nil {
				return fmt.Errorf("can't load JSON config from %s path - %w", path, err)
			}
			dirName := filepath.Base(filepath.Dir(path))
			if config.Package.Name != dirName {
				return fmt.Errorf("directory name (%s) is different from package name (%s)", dirName, config.Package.Name)
			}

			if section == bringauto_const.PackageDirName {
				packageConfigs[config.Package.Name] = append(packageConfigs[config.Package.Name], config)
			} else {
				appConfigs[config.Package.Name] = append(appConfigs[config.Package.Name], config)
			}
		}
		return nil
	})

	if err != nil {
		return ConfigMapType{}, ConfigMapType{}, err
	}

	return packageConfigs, appConfigs, err
}

// checkAllConfigs
// Checks supported Images validity for all Configs in Context.
func (context *ContextManager) checkAllConfigs() error {
	err := context.checkImagesInConfigs(&context.packageConfigs)
	if err != nil {
		return err
	}

	return context.checkImagesInConfigs(&context.appConfigs)
}

// checkImagesInConfigs
// Checks if Images supported by Configs in configsMap were defined in Context.
func (context *ContextManager) checkImagesInConfigs(configsMap *ConfigMapType) error {
	for packName, configs := range *configsMap {
		for _, config := range configs {
			if len(config.DockerMatrix.ImageNames) == 0 {
				return fmt.Errorf("Package/App %s does not support any image", packName)
			}
			for _, dockerName := range config.DockerMatrix.ImageNames {
				_, exists := context.images[dockerName]
				if !exists {
					return fmt.Errorf("Package/App %s supports unknown image %s", packName, dockerName)
				}
			}
		}
	}
	return nil
}

// checkPackageConfigs
// Checks circular dependencies between Package Configs.
func checkPackageConfigs(configsMap *ConfigMapType) error {
	dependsMap, _, err := CreateDependsMap(configsMap)
	if err != nil {
		return err
	}
	err = checkForCircularDependency(dependsMap)
	if err != nil {
		return err
	}
	dependsMapType , err := createDependsMapWithTags(configsMap, &dependsMap)
	if err != nil {
		return err
	}
	return checkDebugReleaseTrees(dependsMapType)
}

// createDependsMapWithTags
// Creates dependency map with information if the Package has Debug and/or Release Config.
func createDependsMapWithTags(configsMap *ConfigMapType, dependsMap *DependsMapType) (dependsMapTagsType, error) {
	dependsMapType := make(dependsMapTagsType)
	for packageName, deps := range *dependsMap {
		dependsMapType[packageName] = dependsMapNodeType{
			DependsOn: []string{},
			HasDebug: false,
			HasRelease: false,
		}
		entry := dependsMapType[packageName]
		for depPackageName := range *deps {
			
			entry.DependsOn = append(entry.DependsOn, depPackageName)
		}
		for _, config := range (*configsMap)[packageName] {
			if config.Package.IsDebug {
				entry.HasDebug = true
			} else {
				entry.HasRelease = true
			}
		}

		dependsMapType[packageName] = entry
	}
	return dependsMapType, nil
}

// checkDebugReleaseTrees
// Checks if all dependencies of Packages in dependsMap have the same Debug and Release Configs as
// depended Package.
func checkDebugReleaseTrees(dependsMap dependsMapTagsType) error {
	for packageName, entry := range dependsMap {
		if entry.HasRelease {
			err := checkDebugReleaseTree(packageName, dependsMap, true)
			if err != nil {
				return err
			}
		}
		if entry.HasDebug {
			err := checkDebugReleaseTree(packageName, dependsMap, false)
			if err != nil {
				return err
			}
		}		
	}
	return nil
}

// checkDebugReleaseTree
// Checks if all dependencies of Package packageName have the same Debug or Release Configs as
// depended Package.
func checkDebugReleaseTree(packageName string, dependsMap dependsMapTagsType, isRelease bool) error {
	entry := dependsMap[packageName]
	if isRelease && !entry.HasRelease {
		return fmt.Errorf("no release Config for Package %s", packageName)
	}
	if !isRelease && !entry.HasDebug {
		return fmt.Errorf("no debug Config for Package %s", packageName)
	}
	for _, depPackageName := range entry.DependsOn {
		err := checkDebugReleaseTree(depPackageName, dependsMap, isRelease)
		if err != nil {
			return err
		}
	}
	return nil
}

// checkAppConfigs
// Checks DependsOn field in App Configs, which must be empty.
func checkAppConfigs(configsMap *ConfigMapType) error {
	for _, configsArray := range *configsMap {
		for _, config := range configsArray {
			if len(config.DependsOn) > 0 {
				return fmt.Errorf("App %s has non-empty DependsOn", config.Package.Name)
			}
		}
	}
	return nil
}

// checkForCircularDependency
// Checks for circular dependency in dependsMap. If there is any, returns error with message
// and problematic packages, else returns nil.
func checkForCircularDependency(dependsMap DependsMapType) error {
	visited := make(map[string]bool)

	for packageName := range dependsMap {
		cycleDetected, cycleString := detectCycle(packageName, dependsMap, visited)
		if cycleDetected {
			return fmt.Errorf("circular dependency detected - %s", packageName+" -> "+cycleString)
		}
		// Clearing recursion stack after one path through graph was checked
		for visitedPackage := range visited {
			visited[visitedPackage] = false
		}
	}
	return nil
}

// detectCycle
// Detects cycle between package dependencies in one path through graph. visited is current
// recursion stack and dependsMap is whole graph representation. packageName is root node where
// cycle detection should start.
func detectCycle(packageName string, dependsMap DependsMapType, visited map[string]bool) (bool, string) {
	visited[packageName] = true
	depsMap, found := dependsMap[packageName]
	if found {
		for depPackageName := range *depsMap {
			if visited[depPackageName] {
				return true, depPackageName
			} else {
				cycleDetected, cycleString := detectCycle(depPackageName, dependsMap, visited)
				if cycleDetected {
					return cycleDetected, depPackageName + " -> " + cycleString
				}
			}
		}
	}
	visited[packageName] = false
	return false, ""
}

// CreateDependsMap
// Creates dependency map (DependsMapType) from configsMap.
func CreateDependsMap(configsMap *ConfigMapType) (DependsMapType, AllDependenciesType, error) {
	dependsMap := make(DependsMapType)
	allDependencies := make(AllDependenciesType)

	for _, configArray := range *configsMap {
		if len(configArray) == 0 {
			return nil, nil, fmt.Errorf("invalid entry in dependency map")
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
	return dependsMap, allDependencies, nil
}

// GetAllConfigsMap
// Returns all Package/App Configs in the Context directory as a Config map.
func (context *ContextManager) GetAllConfigsMap() ConfigMapType {
	mapConfig := make(ConfigMapType)
	for packageName, configsArray := range *context.configs {
		mapConfig[packageName] = configsArray
	}
	return mapConfig
}

// GetAllConfigsArray
// Returns Config structs of all Package/App Configs as an array. If platformString is not nil, it
// is added to all Packages/Apps.
func (context *ContextManager) GetAllConfigsArray(platformString *bringauto_package.PlatformString) []*bringauto_config.Config {
	return getConfigsArray(platformString, context.configs)
}

// GetAllPackageConfigsArray
// Returns Config structs of all Package Configs as an array. If platformString is not nil, it
// is added to all Packages.
func (context *ContextManager) GetAllPackageConfigsArray(platformString *bringauto_package.PlatformString) []*bringauto_config.Config {
	return getConfigsArray(platformString, &context.packageConfigs)
}

// GetAllAppConfigsArray
// Returns Config structs of all App Configs as an array. If platformString is not nil, it
// is added to all Apps.
func (context *ContextManager) GetAllAppConfigsArray(platformString *bringauto_package.PlatformString) []*bringauto_config.Config {
	return getConfigsArray(platformString, &context.appConfigs)
}

// getConfigsArray
// Returns Config structs from configMap as an array. If platformString is not nil, it is added to
// all Packages/Apps.
func getConfigsArray(platformString *bringauto_package.PlatformString, configMap *ConfigMapType) []*bringauto_config.Config {
	var configs []*bringauto_config.Config
	for _, configsArray := range *configMap {
		for _, config := range configsArray {
			configCopy := config
			if platformString != nil {
				configCopy.Package.PlatformString = *platformString
			}
			configs = append(configs, &configCopy)
		}
	}
	return configs
}

// GetAllPackagesStructs
// Returns Package structs of all Package Configs. If platformString is not nil, it is added to all
// Packages.
func (context *ContextManager) GetAllPackagesStructs(platformString *bringauto_package.PlatformString) ([]bringauto_package.Package, error) {
	packConfigs := context.GetAllConfigsArray(platformString)

	var packages []bringauto_package.Package
	for _, packConfig := range packConfigs {
		packages = append(packages, packConfig.Package)
	}

	return packages, nil
}

// GetPackageConfigs
// Returns Configs for given Package/App.
func (context *ContextManager) GetPackageConfigs(packageName string) ([]bringauto_config.Config, error) {
	configs, exists := (*context.configs)[packageName]
	if !exists {
		return []bringauto_config.Config{}, fmt.Errorf("Package/App %s does not exist, please check the name", packageName)
	}

	return configs, nil
}

// getAllDepConfigs
// Returns all Configs for given Package (specified with config) and all Configs for its
// dependencies recursively. For tracking of circular dependencies, the visited map must be
// initialized before function call.
func (context *ContextManager) getAllDepConfigs(config bringauto_config.Config, visited map[string]struct{}) ([]bringauto_config.Config, error) {
	visited[config.Package.GetShortPackageName()] = struct{}{}
	addedPackages := 0
	var configListWithDeps []bringauto_config.Config
	for _, packageDep := range config.DependsOn {
		packageDepConfigs, err := context.GetPackageConfigs(packageDep)
		if err != nil {
			return []bringauto_config.Config{}, fmt.Errorf("cant't get Config of %s package", packageDep)
		}
		for _, depConfig := range packageDepConfigs {
			if depConfig.Package.IsDebug != config.Package.IsDebug {
				continue
			}
			addedPackages++
			_, packageVisited := visited[depConfig.Package.GetShortPackageName()]
			if packageVisited {
				continue
			}
			configListWithDeps = append(configListWithDeps, depConfig)
			configListWithDepsTmp, err := context.getAllDepConfigs(depConfig, visited)
			if err != nil {
				return []bringauto_config.Config{}, err
			}
			configListWithDeps = append(configListWithDeps, configListWithDepsTmp...)
		}
	}

	if addedPackages < len(config.DependsOn) {
		return []bringauto_config.Config{}, fmt.Errorf("package %s dependencies do not have package with same build type", config.Package.Name)
	}

	return configListWithDeps, nil
}

// getAllDepOnConfigs
// Returns all Configs of Packages which depends on Package specified with config. If recursively
// is set to true, it is done recursively. For tracking of circular dependencies, the visited map
// must be initialized before function call.
func (context *ContextManager) getAllDepOnConfigs(config bringauto_config.Config, visited map[string]struct{}, recursively bool) ([]bringauto_config.Config, error) {
	packConfigs := context.GetAllConfigsArray(nil)
	visited[config.Package.Name] = struct{}{}
	var packsToBuild []bringauto_config.Config
	for _, packConfig := range packConfigs {
		if packConfig.Package.Name == config.Package.Name ||
			packConfig.Package.IsDebug != config.Package.IsDebug {
			continue
		}
		for _, dep := range packConfig.DependsOn {
			if dep == config.Package.Name {
				_, packageVisited := visited[packConfig.Package.Name]
				if packageVisited {
					break
				}
				err := context.addDependsOnPackagesToBuild(&packsToBuild, packConfig, visited, recursively)
				if err != nil {
					return []bringauto_config.Config{}, err
				}
				break
			}
		}
	}

	return packsToBuild, nil
}

// addDependsOnPackagesToBuild
// Adds packConfig Package nad its dependecies to packsToBuild. It is done recursively if
// recursively is set.
func (context *ContextManager) addDependsOnPackagesToBuild(packsToBuild *[]bringauto_config.Config, packConfig *bringauto_config.Config, visited map[string]struct{}, recursively bool) error {
	packWithDeps, err := context.GetPackageWithDepsConfigs(packConfig.Package.Name)
	if err != nil {
		return err
	}
	*packsToBuild = append(*packsToBuild, packWithDeps...)
	if recursively {
		packsDepsOnRecursive, err := context.getAllDepOnConfigs(*packConfig, visited, true)
		if err != nil {
			return err
		}
		*packsToBuild = append(*packsToBuild, packsDepsOnRecursive...)
	}
	return nil
}

// GetPackageWithDepsConfigs
// Returns all Configs for given Package and all its dependencies Configs recursively.
func (context *ContextManager) GetPackageWithDepsConfigs(packageName string) ([]bringauto_config.Config, error) {
	packageConfigs, err := context.GetPackageConfigs(packageName)
	if err != nil {
		return []bringauto_config.Config{}, fmt.Errorf("can't get Configs for Package '%s' - %w", packageName, err)
	}
	var packageDeps []bringauto_config.Config
	visitedPackages := make(map[string]struct{})
	for _, config := range packageConfigs {
		packageDepsTmp, err := context.getAllDepConfigs(config, visitedPackages)
		if err != nil {
			return []bringauto_config.Config{}, err
		}
		packageDeps = append(packageDeps, packageDepsTmp...)
	}

	packageConfigs = append(packageConfigs, packageDeps...)

	return packageConfigs, nil
}

// GetPackageWithDepsOnConfigs
// Returns all Configs which depends on given Package and all its dependencies Configs without Package
// (packageName) itself and its dependencies. If recursively is set to true, it is done recursively.
func (context *ContextManager) GetPackageWithDepsOnConfigs(packageName string, recursively bool) ([]bringauto_config.Config, error) {
	packageConfigs, err := context.GetPackageConfigs(packageName)
	if err != nil {
		return []bringauto_config.Config{}, err
	}
	var packsToBuild []bringauto_config.Config
	visitedPackages := make(map[string]struct{})
	for _, config := range packageConfigs {
		packageDepsTmp, err := context.getAllDepOnConfigs(config, visitedPackages, recursively)
		if err != nil {
			return []bringauto_config.Config{}, err
		}
		packsToBuild = append(packsToBuild, packageDepsTmp...)
	}

	packsToRemove, err := context.GetPackageWithDepsConfigs(packageName)
	if err != nil {
		return []bringauto_config.Config{}, err
	}
	packsToBuild = removeConfigs(packsToBuild, packsToRemove)
	return packsToBuild, nil
}

// removeConfigs
// Removes list2 Configs from list1.
func removeConfigs(list1 []bringauto_config.Config, list2 []bringauto_config.Config) []bringauto_config.Config {
	for _, config2 := range list2 {
		list1 = removeConfig(list1, config2)
	}
	return list1
}

// removeConfig
// Removes config from list1.
func removeConfig(list1 []bringauto_config.Config, config bringauto_config.Config) []bringauto_config.Config {
	i := 0
	for _, config1 := range list1 {
		if config1.Package.GetShortPackageName() != config.Package.GetShortPackageName() {
			list1[i] = config1
			i++
		}
	}
	return list1[:i]
}

// loadImagesDockerfilePaths
// Loads Image Dockerfile paths from Context into Context Manager struct.
func (context *ContextManager) loadImagesDockerfilePaths() error {
	imageDir := path.Join(context.ContextPath, bringauto_const.DockerDirName)

	// Resolve symlink if imageDir is a symlink
	resolvedImageDir, err := filepath.EvalSymlinks(imageDir)
	if err != nil {
		return fmt.Errorf("cannot resolve symlink %s for image directory", imageDir)
	}

	reg, err := regexp.CompilePOSIX("^Dockerfile$")
	if err != nil {
		return fmt.Errorf("cannot compile regexp for matchiing Dockerfile")
	}

	dockerfileList, err := getAllFilesInSubdirByRegexp(resolvedImageDir, reg)
	if err != nil {
		return err
	}

	imagePaths := make(ImagesPathType)
	for imageName, pathList := range dockerfileList {
		if len(pathList) != 1 {
			return fmt.Errorf("wrong number of Dockerfiles for %s image (should be 1)", imageName)
		}
		imagePaths[imageName] = pathList[0]
	}

	context.images = imagePaths
	return nil
}

// GetAllImagesDockerfilePaths
// Returns all Dockerfile paths located in the Context directory.
func (context *ContextManager) GetAllImagesDockerfilePaths() ImagesPathType {
	return context.images
}

// GetImageDockerfilePath
// Returns Dockerfile path for the given Image locate in the given Context.
func (context *ContextManager) GetImageDockerfilePath(imageName string) (string, error) {
	path, exists := context.images[imageName]
	if !exists {
		return "", fmt.Errorf("docker image definition does not exist, please check the name")
	}

	return path, nil
}

// validateContextPath
// Validates Context path if the structure in the Context directory works
// Return nil if structure is valid, error if the structure is invalid.
func (context *ContextManager) validateContextPath() error {
	var err error
	ContextStat, err := os.Stat(context.ContextPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("context path does not exist - %s\n", context.ContextPath)
	}
	if !ContextStat.IsDir() {
		return fmt.Errorf("context path is not a directory - %s\n", context.ContextPath)
	}

	dockerDirPath := path.Join(context.ContextPath, bringauto_const.DockerDirName)
	DockerStat, err := os.Stat(dockerDirPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("docker dir path does not exist - %s\n", dockerDirPath)
	}
	if !DockerStat.IsDir() {
		return fmt.Errorf("docker path is not a directory - %s\n", dockerDirPath)
	}

	packageDirPath := path.Join(context.ContextPath, bringauto_const.PackageDirName)
	packageStat, err := os.Stat(packageDirPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("package path does not exist - %s\n", packageDirPath)
	}
	if !packageStat.IsDir() {
		return fmt.Errorf("package path is not a directory - %s\n", packageDirPath)
	}

	appDirPath := path.Join(context.ContextPath, bringauto_const.AppDirName)
	appStat, err := os.Stat(appDirPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("App path does not exist - %s\n", appDirPath)
	}
	if !appStat.IsDir() {
		return fmt.Errorf("App path is not a directory - %s\n", appDirPath)
	}

	return nil
}

// getAllFilesInDirByRegexp
// Get all file in subdirs of rootDir which matches given regexp.
func getAllFilesInSubdirByRegexp(rootDir string, reg *regexp.Regexp) (map[string][]string, error) {
	acceptedFileList := map[string][]string{}
	walkError := filepath.WalkDir(rootDir, func(item string, d fs.DirEntry, err error) error {
		if d.Name() == path.Base(rootDir) {
			return nil
		}
		packageName := d.Name()
		packageBaseDir := path.Join(rootDir, d.Name())
		packageFileDefs, err := getAllFilesInDirByRegexp(packageBaseDir, reg)
		if err != nil {
			return nil
		}
		acceptedFileList[packageName] = packageFileDefs
		return nil
	},
	)
	return acceptedFileList, walkError
}

// getAllFilesInDirByRegexp
// Get all files from given rootDir which matches given regexp.
func getAllFilesInDirByRegexp(rootDir string, reg *regexp.Regexp) ([]string, error) {
	var acceptedFileList []string
	dirEntryList, err := os.ReadDir(rootDir)
	if err != nil {
		return []string{}, fmt.Errorf("cannot list dir %s", rootDir)
	}

	for _, dirEntry := range dirEntryList {
		packageNameOk := reg.MatchString(dirEntry.Name())
		if !packageNameOk {
			continue
		}
		acceptedFileList = append(acceptedFileList, path.Join(rootDir, dirEntry.Name()))
	}
	return acceptedFileList, nil
}
