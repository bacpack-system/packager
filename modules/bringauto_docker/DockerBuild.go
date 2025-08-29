package bringauto_docker

import (
	"bringauto/modules/bringauto_log"
	"fmt"
	"os/exec"
)

// DockerBuild
// options needed for a build a docker image
type DockerBuild struct {
	// directory where the Dockerfile is located
	DockerfileDir string
	// tag which will be used to tag image after build
	Tag string
	// build context
	Context string
}

// Build given docker image
func (dockerBuild *DockerBuild) Build() error {
	if dockerBuild.DockerfileDir == "" {
		return fmt.Errorf("DockerBuild - DockerfileDir is empty")
	}

	var ok = dockerBuild.prepareAndRun(prepareBuildArgs)
	if !ok {
		return fmt.Errorf("DockerBuild build error for %s image", dockerBuild.Tag)
	}
	return nil
}

// prepareAndRun
// Runs docker build command. Returns bool value which indicates if the command was succesfull or
// not.
func (dockerBuild *DockerBuild) prepareAndRun(f func(build *DockerBuild) []string) bool {
	logger := bringauto_log.GetLogger()
	contextLogger := logger.CreateContextLogger(dockerBuild.Tag, "", bringauto_log.ImageBuildContext)

	var cmd exec.Cmd
	cmdArgs := f(dockerBuild)
	cmdArgs = append([]string{DockerExecutablePathConst}, cmdArgs...)
	cmd.Args = cmdArgs
	cmd.Path = DockerExecutablePathConst
	if contextLogger != nil {
		file, err := contextLogger.GetFile()
		if err != nil {
			logger.Error("Failed to open file for logs - %s", err)
			return false
		}
		cmd.Stderr = file
		cmd.Stdout = file
	}
	err := cmd.Run()
	if err != nil {
		return false
	}
	if cmd.ProcessState.ExitCode() > 0 {
		return false
	}
	return true
}

func prepareBuildArgs(dockerBuild *DockerBuild) []string {
	cmdArgs := make([]string, 0)
	cmdArgs = append(cmdArgs, "build")
	cmdArgs = append(cmdArgs, dockerBuild.DockerfileDir)
	if dockerBuild.DockerfileDir != "" {
		cmdArgs = append(cmdArgs, "--tag", dockerBuild.Tag)
	}

	if dockerBuild.Context != "" {
		cmdArgs = append(cmdArgs, "--build-context", "package-context=" + dockerBuild.Context)
	}
	return cmdArgs
}
