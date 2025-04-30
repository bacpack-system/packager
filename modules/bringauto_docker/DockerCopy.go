package bringauto_docker

import (
	"bringauto/modules/bringauto_process"
	"bytes"
	"fmt"
)

type DockerCopy Docker

// Copy
// Copies file from container with filePath path to localPath.
func (args *DockerCopy) Copy(filePath string, localPath string) error {

	var extraArgs []string
	extraArgs = append(extraArgs, "cp")
	extraArgs = append(extraArgs, "-L")
	extraArgs = append(extraArgs, fmt.Sprintf("%s:%s", args.containerId, filePath))
	extraArgs = append(extraArgs, localPath)

	var errBuff bytes.Buffer
	process := bringauto_process.Process{
		CommandAbsolutePath: DockerExecutablePathConst,
		Args: bringauto_process.ProcessArgs{
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
