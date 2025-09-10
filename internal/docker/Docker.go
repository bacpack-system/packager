package docker

import (
	"github.com/bacpack-system/packager/internal/prerequisites"
	"github.com/bacpack-system/packager/internal/process"
	"github.com/bacpack-system/packager/internal/constants"
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

func (docker *Docker) FillDefault(*prerequisites.Args) error {
	*docker = Docker{
		Volumes:     map[string]string{},
		RunAsDaemon: true,
		ImageName:   defaultImageNameConst,
		Port: constants.DefaultSSHPort,
	}
	return nil
}

func (docker *Docker) FillDynamic(args *prerequisites.Args) error {
	var argsStruct dockerInitArgs
	prerequisites.GetArgs(args, &argsStruct)
	if prerequisites.IsEmpty(args) {
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
func (docker *Docker) CheckPrerequisites(*prerequisites.Args) error {
	err := checkIfDockerIsUsable()
	if err != nil {
		return fmt.Errorf(
			"Docker cannot be used, it is not installed, the Docker daemon is not running " +
			"or current user is not in Docker group - %w", err,
		)
	}
	portAvailable, err := isPortAvailable(docker.Port)
	if err != nil {
		return err
	} else if !portAvailable {
		return fmt.Errorf("port %d not available", docker.Port)
	}

	err = checkForImageExistence(docker.ImageName)
	if err != nil {
		return err
	}

	for hostVolume, _ := range docker.Volumes {
		if _, err = os.Stat(hostVolume); os.IsNotExist(err) {
			return fmt.Errorf("cannot mount non existent directory as volume: '%s'", hostVolume)
		}
	}

	return nil
}

// SetVolume set volume mapping for a Docker container.
// It's not possible to overwrite volume mapping that already exists
func (docker *Docker) SetVolume(hostDirectory string, containerDirectory string) error {
	_, hostFound := docker.Volumes[hostDirectory]
	if hostFound {
		return fmt.Errorf("volume mapping is already set: '%s' --> '%s'", hostDirectory, containerDirectory)
	}
	docker.Volumes[hostDirectory] = containerDirectory

	return nil
}

// checkIfDockerIsUsable
// Checks if the Docker can be used. If not, returns error, else nil.
func checkIfDockerIsUsable() error {
	var errBuff bytes.Buffer
	process := process.Process{
		CommandAbsolutePath: DockerExecutablePathConst,
		Args: process.ProcessArgs{
			ExtraArgs: &[]string{
				"info",
			},
		},
		StdErr: &errBuff,
	}

	err := process.Run()
	if err != nil {
		return fmt.Errorf(errBuff.String())
	}
	return nil	
}

// checkForImageExistence
// Checks if the Docker image exists. If not, returns error, else nil.
func checkForImageExistence(imageName string) error {
	var outBuff bytes.Buffer
	process := process.Process{
		CommandAbsolutePath: DockerExecutablePathConst,
		Args: process.ProcessArgs{
			ExtraArgs: &[]string{
				"images",
				"-q",
				imageName,
			},
		},
		StdOut: &outBuff,
	}
	err := process.Run()
	if err != nil {
		return err
	}

	if outBuff.Len() == 0 {
		return fmt.Errorf("image %s does not exist", imageName)
	}
	return nil
}

// isPortAvailable
// Returns true if port for docker is available, else returns false.
// When false is returned, the error contains message from the docker command.
func isPortAvailable(port uint16) (bool, error) {
	var outBuff, errBuff bytes.Buffer

	process := process.Process{
		CommandAbsolutePath: DockerExecutablePathConst,
		Args: process.ProcessArgs{
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
