package bringauto_docker

import (
	"bringauto/modules/bringauto_prerequisites"
	"bringauto/modules/bringauto_process"
	"bringauto/modules/bringauto_const"
	"fmt"
	"os"
	"bytes"
	"strconv"
)

const (
	sshPort = 22
	defaultImageNameConst = "unknown"
)

// Docker
type Docker struct {
	// ImageName - tag or image hash
	ImageName string
	// Used port
	Port uint16
	// Volumes map a host directory (represented by absolute path)
	// to the directory inside the docker container
	// in manner map[string]string { <host_volume_abs_path>:<> }
	Volumes map[string]string `json:"-"`
	// If true docker command will run in non-blocking mode - as a daemon.
	RunAsDaemon bool `json:"-"`
	containerId string
}

type dockerInitArgs struct {
	ImageName string
	Port uint16
}

func (docker *Docker) FillDefault(*bringauto_prerequisites.Args) error {
	*docker = Docker{
		Volumes:     map[string]string{},
		RunAsDaemon: true,
		ImageName:   defaultImageNameConst,
		Port: bringauto_const.DefaultSSHPort,
	}
	return nil
}

func (docker *Docker) FillDynamic(args *bringauto_prerequisites.Args) error {
	var argsStruct dockerInitArgs
	bringauto_prerequisites.GetArgs(args, &argsStruct)
	if bringauto_prerequisites.IsEmpty(args) {
		return nil
	}
	if argsStruct.ImageName != "" {
		docker.ImageName = argsStruct.ImageName
	}
	docker.Port = argsStruct.Port
	return nil
}

// CheckPrerequisites
// It checks if the docker is installed and can be run by given user.
// Function returns nil if Docker installation is ok, not nil of the problem is recognized
func (docker *Docker) CheckPrerequisites(*bringauto_prerequisites.Args) error {
	portAvailable, err := isPortAvailable(docker.Port)
	if err != nil {
		return err
	} else if !portAvailable {
		return fmt.Errorf("port %d not available", docker.Port)
	}
	var outBuff bytes.Buffer
	process := bringauto_process.Process{
		CommandAbsolutePath: DockerExecutablePathConst,
		Args: bringauto_process.ProcessArgs{
			ExtraArgs: &[]string{
				"images",
				"-q",
				docker.ImageName,
			},
		},
		StdOut: &outBuff,
	}
	err = process.Run()
	if err != nil {
		return err
	}

	if outBuff.Len() == 0 {
		return fmt.Errorf("image %s does not exist", docker.ImageName)
	}

	for hostVolume, _ := range docker.Volumes {
		if _, err = os.Stat(hostVolume); os.IsNotExist(err) {
			return fmt.Errorf("connot mount non existent directory as volume: '%s'", hostVolume)
		}
	}

	return nil
}

// SetVolume set volume mapping for a Docker container.
// It's not possible to overwrite volume mapping that already exists (panic occure)
func (docker *Docker) SetVolume(hostDirectory string, containerDirectory string) {
	_, hostFound := docker.Volumes[hostDirectory]
	if hostFound {
		panic(fmt.Errorf("volume mapping is already set: '%s' --> '%s'", hostDirectory, containerDirectory))
	}
	docker.Volumes[hostDirectory] = containerDirectory
}


// isPortAvailable
// Returns true if port for docker is available, else returns false.
// When false is returned, the error contains message from the docker command.
func isPortAvailable(port uint16) (bool, error) {
	var outBuff, errBuff bytes.Buffer

	process := bringauto_process.Process{
		CommandAbsolutePath: DockerExecutablePathConst,
		Args: bringauto_process.ProcessArgs{
			ExtraArgs: &[]string{
				"container",
				"ls",
				"--filter",
				"publish=" + strconv.Itoa(int(port)),
				"--format",
				"{{.ID}}{{.Ports}}",
			},
		},
		StdOut: &outBuff,
		StdErr: &errBuff,
	}

	err := process.Run()
	if err != nil {
		return false, fmt.Errorf(errBuff.String())
	}

	return outBuff.Len() == 0, nil
}
