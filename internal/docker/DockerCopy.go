package docker

import (
	"github.com/bacpack-system/packager/internal/process"
	"bytes"
	"fmt"
)

type DockerCopy Docker

// Copy
// Copies file from container with filePath path to localPath.
func (args *DockerCopy) Copy(filePath string, localPath string) error {
	if args.containerId == "" {
		return fmt.Errorf("dockerCopy copy error - container ID is empty")
	}

	var extraArgs []string
	extraArgs = append(extraArgs, "cp")
	extraArgs = append(extraArgs, "-L")
	extraArgs = append(extraArgs, fmt.Sprintf("%s:%s", args.containerId, filePath))
	extraArgs = append(extraArgs, localPath)

	var errBuff bytes.Buffer
	process := process.Process{
		CommandAbsolutePath: DockerExecutablePathConst,
		Args: process.ProcessArgs{
			ExtraArgs: &extraArgs,
		},
		StdErr: &errBuff,
	}
	err := process.Run()

	if err != nil {
		return fmt.Errorf("dockerCopy copy error - %s", errBuff.String())
	}

	return nil
}
