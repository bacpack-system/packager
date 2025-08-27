package bringauto_docker_test

import (
	"bringauto/modules/bringauto_docker"
	"bringauto/modules/bringauto_prerequisites"
	"reflect"
	"testing"
)

func TestDockerRm_GenerateCmdLine(t *testing.T) {
	docker, err := bringauto_prerequisites.CreateAndInitialize[bringauto_docker.Docker]()
	if err != nil {
		t.Fatalf("cannot create Docker instance - %s", err)
	}
	dockerRm := (*bringauto_docker.DockerRm)(docker)
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
	docker, err := bringauto_prerequisites.CreateAndInitialize[bringauto_docker.Docker]()
	if err != nil {
		t.Fatalf("cannot create Docker instance - %s", err)
	}
	dockerStop := (*bringauto_docker.DockerStop)(docker)
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

	docker, err := bringauto_prerequisites.CreateAndInitialize[bringauto_docker.Docker]()
	if err != nil {
		t.Fatalf("cannot create Docker instance - %s", err)
	}
	dockerRun := (*bringauto_docker.DockerRun)(docker)

	dockerRun.RunAsDaemon = true

	validCmdLine := []string{
		"run",
		"-d",
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

	dockerRun.Port = 1212
	validCmdLine = validCmdLine[:len(validCmdLine)-1]
	validCmdLine = append(validCmdLine, "-p", "1212:125", dockerRun.ImageName)
	cmdLine, err = dockerRun.GenerateCmdLine()
	if err != nil {
		t.Errorf("cannot generate reference cmd line")
		return
	}
	cmdLineValid = reflect.DeepEqual(cmdLine, validCmdLine)
	if !cmdLineValid {
		t.Errorf("invalid Docker Run cmd line with ports!")
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
