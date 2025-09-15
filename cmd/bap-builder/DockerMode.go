package main

import (
	"github.com/bacpack-system/packager/internal/log"
	"github.com/bacpack-system/packager/internal/docker"
	"github.com/bacpack-system/packager/internal/context"
	"github.com/bacpack-system/packager/internal/packager_error"
	"github.com/bacpack-system/packager/internal/prerequisites"
	"path"
)

// BuildDockerImage
// process Docker mode of cmd line
//
func BuildDockerImage(cmdLine *BuildImageCmdLineArgs, contextPath string) error {
	contextManager := context.ContextManager{
		ContextPath: contextPath,
	}
	err := prerequisites.Initialize(&contextManager)
	if err != nil {
		logger := log.GetLogger()
		logger.Error("Context consistency error - %s", err)
		return packager_error.ContextErr
	}
	buildAll := cmdLine.All
	if *buildAll {
		return buildAllDockerImages(contextManager)
	}

	dockerfilePath, err := contextManager.GetImageDockerfilePath(*cmdLine.Name)
	if err != nil {
		return err
	}
	return buildSingleDockerImage(*cmdLine.Name, dockerfilePath, contextPath)
}

// buildAllDockerImages
// builds all docker images in the given contextPath.
// It returns nil if everything is ok, or not nil in case of error
//
func buildAllDockerImages(contextManager context.ContextManager) error {
	dockerfilePathList := contextManager.GetAllImagesDockerfilePaths()

	for imageName, dockerfilePath := range dockerfilePathList {
		err := buildSingleDockerImage(imageName, dockerfilePath, contextManager.ContextPath)
		if err != nil {
			return err
		}
	}
	return nil
}

// buildSingleDockerImage
// builds a single docker image specified by an image name and a path to Dockerfile.
//
func buildSingleDockerImage(imageName string, dockerfilePath string, contextPath string) error {
	logger := log.GetLogger()
	dockerfileDir := path.Dir(dockerfilePath)
	dockerBuild := docker.DockerBuild{
		DockerfileDir: dockerfileDir,
		Tag:           imageName,
		Context:       contextPath,
	}
	logger.Info("Build Docker Image: %s", imageName)

	// Building image does not require any handler when SIGINT is received. 'docker build' creates
	// image after all steps from Dockerfile are successfully executed.
	err := dockerBuild.Build()
	if err != nil {
		logger.ErrorIndent("Can't build image - %s", err)
		return packager_error.BuildErr
	}
	logger.InfoIndent("Build OK")
	return nil
}
