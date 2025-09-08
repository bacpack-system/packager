package docker_test

import (
	"github.com/bacpack-system/packager/internal/docker"
	"github.com/bacpack-system/packager/internal/prerequisites"
	"github.com/bacpack-system/packager/internal/constants"
	"reflect"
	"testing"
	"strconv"
)

func TestDockerRm_GenerateCmdLine(t *testing.T) {
	dock, err := prerequisites.CreateAndInitialize[docker.Docker]("debian12", uint16(constants.DefaultSSHPort))
	if err != nil {
		t.Fatalf("cannot create Docker instance - %s", err)
	}
	dockerRm := (*docker.DockerRm)(dock)
	cmdLine, err := dockerRm.GenerateCmdLine()
	if err != nil {
		t.Errorf("cannot generate reference cmd line")
	}

	validCmdLine := []string{
		"rm",
		"",
	}

	cmdLineValid := reflect.DeepEqual(cmdLine, validCmdLine)
	if !cmdLineValid {
		t.Errorf("invalid Docker RM cmd line!")
	}
}

func TestDockerStop_GenerateCmdLine(t *testing.T) {
	dock, err := prerequisites.CreateAndInitialize[docker.Docker]("debian12", uint16(constants.DefaultSSHPort))
	if err != nil {
		t.Fatalf("cannot create Docker instance - %s", err)
	}
	dockerStop := (*docker.DockerStop)(dock)
	cmdLine, err := dockerStop.GenerateCmdLine()
	if err != nil {
		t.Errorf("cannot generate reference cmd line")
	}

	validCmdLine := []string{
		"stop",
		"",
	}

	cmdLineValid := reflect.DeepEqual(cmdLine, validCmdLine)
	if !cmdLineValid {
		t.Errorf("invalid Docker Stop cmd line!")
	}
}

func TestDockerRun_GenerateCmdLine(t *testing.T) {
	var cmdLine []string
	var cmdLineValid bool
	var err error

	dock, err := prerequisites.CreateAndInitialize[docker.Docker]("debian12", uint16(constants.DefaultSSHPort))
	if err != nil {
		t.Fatalf("cannot create Docker instance - %s", err)
	}
	dockerRun := (*docker.DockerRun)(dock)

	dockerRun.RunAsDaemon = true

	validCmdLine := []string{
		"run",
		"-d",
		"-p",
		strconv.Itoa(constants.DefaultSSHPort) + ":22",
		dockerRun.ImageName,
	}
	cmdLine, err = dockerRun.GenerateCmdLine()
	if err != nil {
		t.Errorf("cannot generate reference cmd line")
		return
	}
	cmdLineValid = reflect.DeepEqual(cmdLine, validCmdLine)
	if !cmdLineValid {
		t.Errorf("invalid Docker Run cmd line as a daemon!")
		return
	}

	dockerRun.Volumes = map[string]string{
		"A": "A",
		"B": "BVol",
	}
	validCmdLine = validCmdLine[:len(validCmdLine)-1]
	validCmdLine = append(validCmdLine, "-v", "A:A", "-v", "B:BVol", dockerRun.ImageName)
	cmdLine, err = dockerRun.GenerateCmdLine()
	if err != nil {
		t.Errorf("cannot generate reference cmd line")
		return
	}
	cmdLineValid = reflect.DeepEqual(cmdLine, validCmdLine)
	if !cmdLineValid {
		t.Errorf("invalid Docker Run cmd line with volumes!")
		return
	}
}
