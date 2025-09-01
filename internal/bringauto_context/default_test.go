package bringauto_context

import (
	"bringauto/internal/bringauto_const"
	"bringauto/internal/bringauto_package"
	"bringauto/internal/bringauto_prerequisites"
	"path/filepath"
	"testing"
	"os"
)

const (
	TestDataDirName = "test_data"
	Set1DirName = "set1"
	Set2DirName = "set2"
	Set3DirName = "set3"
	Set4DirName = "set4"
	Set1DirPath = TestDataDirName + "/" + Set1DirName
	Set2DirPath = TestDataDirName + "/" + Set2DirName
	Set3DirPath = TestDataDirName + "/" + Set3DirName
	Set4DirPath = TestDataDirName + "/" + Set4DirName

	Pack1Name = "pack1"
	Pack2Name = "pack2"
	Pack3Name = "pack3"
	Pack4Name = "pack4"
	Pack5Name = "pack5"
	Pack6Name = "pack6"
	Image1Name = "image1"
	Image2Name = "image2"
	DockerfileName = "Dockerfile"
)

var defaultPlatformString bringauto_package.PlatformString

func TestMain(m *testing.M) {
	stringExplicit := bringauto_package.PlatformStringExplicit {
		DistroName: "distro",
		DistroRelease: "1.0",
		Machine: "machine",
	}

	defaultPlatformString = bringauto_package.PlatformString{
		Mode: bringauto_package.ModeExplicit,
		String: stringExplicit,
	}
	os.Exit(m.Run())
}

func initContext(contextPath string) (*ContextManager, error) {
	context := ContextManager {
		ContextPath: contextPath,
	}
	err := bringauto_prerequisites.Initialize(&context)
	if err != nil {
		return nil, err
	}
	return &context, nil
}

func TestGetAllConfigsMap(t *testing.T) {
	context, err := initContext(Set1DirPath)
	if err != nil {
		t.Fatalf("Cannot initialize context - %s", err)
	} 

	configMap := context.GetAllConfigsMap()

	_, ok1 := configMap[Pack1Name]
	_, ok2 := configMap[Pack2Name]
	_, ok3 := configMap[Pack3Name]

	if !ok1 || !ok2 || !ok3 {
		t.Fatalf("some Config was not returned")
	}
}

func TestGetPackageConfigs(t *testing.T) {
	context, err := initContext(Set1DirPath)
	if err != nil {
		t.Fatalf("Cannot initialize context - %s", err)
	} 

	pack1Configs, err := context.GetPackageConfigs(Pack1Name)
	if err != nil {
		t.Fatalf("GetPackageConfigs failed - %s", err)
	}

	if len(pack1Configs) != 2 {
		t.Fatalf("wrong number of returned Configs")
	}

	for _, config := range pack1Configs {
		if config.Package.Name != Pack1Name {
			t.Error("wrong Config content")
		}
	}
}

func TestGetAllPackageConfigsArray(t *testing.T) {
	context, err := initContext(Set1DirPath)
	if err != nil {
		t.Fatalf("Cannot initialize context - %s", err)
	} 

	configs := context.GetAllPackageConfigsArray(&defaultPlatformString)

	if len(configs) != 4 {
		t.Fatal("wrong number of returned configs")
	}

	// Checking some properties
	for _, config := range configs {
		if config.Package.Name == Pack1Name && config.Package.IsDebug {
			if config.Package.VersionTag != "v1.0.0" {
				t.Error("wrong config content")
			}
		} else if config.Package.Name == Pack1Name {
			if (config.DockerMatrix.ImageNames[0] != Image1Name ||
				len(config.DependsOn) != 0) {
				t.Error("wrong config content")
			}
		} else if config.Package.Name == Pack2Name {
			if config.Build.CMake.Defines["BRINGAUTO_INSTALL"] != "ON" {
				t.Error("wrong config content")
			}
		} else if config.Package.Name == Pack3Name {
			if (config.DependsOn[0] != "pack1" ||
				config.DependsOn[1] != "pack2") {
				t.Error("wrong config content")
			}
		} else {
			t.Error("returned config for unknown package")
		}
	}
}

func TestGetAllImagesDockerfilePaths(t *testing.T) {
	context, err := initContext(Set1DirPath)
	if err != nil {
		t.Fatalf("Cannot initialize context - %s", err)
	} 

	paths := context.GetAllImagesDockerfilePaths()

	image1PathToCheck, ok1 := paths[Image1Name]
	image2PathToCheck, ok2 := paths[Image2Name]

	if !ok1 || !ok2 {
		t.Fatalf("some image was not returned")
	}

	commonPath := filepath.Join(Set1DirPath, bringauto_const.DockerDirName)
	image1Path := filepath.Join(commonPath, Image1Name, DockerfileName)
	image2Path := filepath.Join(commonPath, Image2Name, DockerfileName)

	if ((image1PathToCheck != image1Path) || (image2PathToCheck != image2Path)) {
		t.Fatalf("wrong returned paths")
	}
}

func TestGetPackageWithDepsConfigs(t *testing.T) {
	context, err := initContext(Set2DirPath)
	if err != nil {
		t.Fatalf("Cannot initialize context - %s", err)
	} 

	configs, err := context.GetPackageWithDepsConfigs(Pack3Name)
	if err != nil {
		t.Fatalf("GetPackageWithDepsConfigs failed - %s", err)
	}

	if len(configs) != 4 {
		t.Fatalf("wrong number of returned configs")
	}

	for _, config := range configs {
		if (config.Package.Name != Pack1Name &&
			config.Package.Name != Pack2Name &&
			config.Package.Name != Pack3Name &&
			config.Package.Name != Pack4Name) {
			t.Error("wrong returned configs")
		}
	}
}

func TestGetPackageWithDepsConfigsNoDepWithBuildType(t *testing.T) {
	context, err := initContext(Set3DirPath)
	if err != nil {
		t.Fatalf("Cannot initialize context - %s", err)
	} 

	_, err = context.GetPackageWithDepsConfigs(Pack2Name)
	if err == nil {
		t.Error("GetPackageWithDepsConfigs didn't returned error")
	}
}

func TestGetPackageWithDepsConfigsCircularDependency(t *testing.T) {
	context, err := initContext(Set4DirPath)
	if err != nil {
		t.Fatalf("Cannot initialize context - %s", err)
	} 

	configs, err := context.GetPackageWithDepsConfigs(Pack1Name)
	if err != nil {
		t.Fatalf("GetPackageWithDepsConfigs failed - %s", err)
	}

	if len(configs) != 3 {
		t.Fatalf("wrong number of returned configs")
	}

	for _, config := range configs {
		if (config.Package.Name != Pack1Name &&
			config.Package.Name != Pack2Name &&
			config.Package.Name != Pack3Name) {
			t.Error("wrong returned configs")
		}
	}
}

func TestGetPackageWithDepsOnConfigs(t *testing.T) {
	context, err := initContext(Set2DirPath)
	if err != nil {
		t.Fatalf("Cannot initialize context - %s", err)
	} 

	configs, err := context.GetPackageWithDepsOnConfigs(Pack1Name, false)
	if err != nil {
		t.Fatalf("GetPackageWithDepsOnConfigs failed - %s", err)
	}

	if len(configs) != 3 {
		t.Fatalf("wrong number of returned configs")
	}

	for _, config := range configs {
		if (config.Package.Name != Pack4Name &&
			config.Package.Name != Pack5Name &&
			config.Package.Name != Pack6Name) {
			t.Error("wrong returned configs")
		}
	}
}

func TestGetPackageWithDepsOnConfigsRecursively(t *testing.T) {
	context, err := initContext(Set2DirPath)
	if err != nil {
		t.Fatalf("Cannot initialize context - %s", err)
	} 

	configs, err := context.GetPackageWithDepsOnConfigs(Pack1Name, true)
	if err != nil {
		t.Fatalf("GetPackageWithDepsOnConfigs failed - %s", err)
	}

	if len(configs) != 4 {
		t.Fatalf("wrong number of returned configs")
	}

	for _, config := range configs {
		if (config.Package.Name != Pack3Name &&
			config.Package.Name != Pack4Name &&
			config.Package.Name != Pack5Name &&
			config.Package.Name != Pack6Name) {
			t.Error("wrong returned configs")
		}
	}
}

func TestGetPackageWithDepsOnConfigsRecursivelyCircularDependency(t *testing.T) {
	context, err := initContext(Set4DirPath)
	if err != nil {
		t.Fatalf("Cannot initialize context - %s", err)
	} 

	configs, err := context.GetPackageWithDepsOnConfigs(Pack1Name, true)
	if err != nil {
		t.Fatalf("GetPackageWithDepsOnConfigs failed - %s", err)
	}

	if len(configs) != 1 {
		t.Fatalf("wrong number of returned configs")
	}

	for _, config := range configs {
		if (config.Package.Name != Pack4Name) {
			t.Error("wrong returned config")
		}
	}
}
