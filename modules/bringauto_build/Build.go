package bringauto_build

import (
	"bringauto/modules/bringauto_const"
	"bringauto/modules/bringauto_docker"
	"bringauto/modules/bringauto_git"
	"bringauto/modules/bringauto_log"
	"bringauto/modules/bringauto_package"
	"bringauto/modules/bringauto_prerequisites"
	"bringauto/modules/bringauto_process"
	"bringauto/modules/bringauto_ssh"
	"bringauto/modules/bringauto_sysroot"
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"time"
	"strconv"
)

type Build struct {
	Env            *EnvironmentVariables
	Docker         *bringauto_docker.Docker
	Git            *bringauto_git.Git
	CMake          *CMake
	GNUMake        *GNUMake
	SSHCredentials *bringauto_ssh.SSHCredentials
	Package        *bringauto_package.Package
	BuiltPackage   *bringauto_sysroot.BuiltPackage
	sysroot        *bringauto_sysroot.Sysroot
}

type buildInitArgs struct {
	DockerImageName string
	DockerPort uint16
}

// FillDefault
// It fills up defaults for all members in the Build structure.
func (build *Build) FillDefault(args *bringauto_prerequisites.Args) error {
	var argsStruct buildInitArgs
	var err error
	bringauto_prerequisites.GetArgs(args, &argsStruct)
	if build.Git == nil {
		build.Git, err = bringauto_prerequisites.CreateAndInitialize[bringauto_git.Git]()
		if err != nil {
			return err
		}
	}
	if build.Docker == nil {
		build.Docker, err = bringauto_prerequisites.CreateAndInitialize[bringauto_docker.Docker](argsStruct.DockerImageName, argsStruct.DockerPort)
		if err != nil {
			return err
		}
	}
	if build.SSHCredentials == nil {
		build.SSHCredentials, err = bringauto_prerequisites.CreateAndInitialize[bringauto_ssh.SSHCredentials]()
		if err != nil {
			return err
		}
		build.SSHCredentials.Port = build.Docker.Port
	}
	if build.CMake == nil {
		build.CMake, err = bringauto_prerequisites.CreateAndInitialize[CMake]()
		if err != nil {
			return err
		}
	}
	if build.GNUMake == nil {
		build.GNUMake, err = bringauto_prerequisites.CreateAndInitialize[GNUMake]()
		if err != nil {
			return err
		}
	}
	if build.Env == nil {
		build.Env, err = bringauto_prerequisites.CreateAndInitialize[EnvironmentVariables]()
		if err != nil {
			return err
		}
	}

	if build.Package == nil {
		build.Package, err = bringauto_prerequisites.CreateAndInitialize[bringauto_package.Package]()
		if err != nil {
			return err
		}
	}

	return nil
}

func (build *Build) FillDynamic(*bringauto_prerequisites.Args) error {
	return nil
}

func (build *Build) CheckPrerequisites(*bringauto_prerequisites.Args) error {
	copyDir := build.GetLocalInstallDirPath()
	if _, err := os.Stat(copyDir); !os.IsNotExist(err) {
		return fmt.Errorf("package directory exist. Please delete it: %s", copyDir)
	}

	return nil
}

// performPreBuildTasks
// Downloads Package files for build in docker container. Clones the repository and updates all
// submodules.
func (build *Build) performPreBuildTasks(shellEvaluator *bringauto_ssh.ShellEvaluator) error {
	gitClone := bringauto_git.GitClone{Git: *build.Git}
	gitCheckout := bringauto_git.GitCheckout{Git: *build.Git}
	gitSubmoduleUpdate := bringauto_git.GitSubmoduleUpdate{Git: *build.Git}
	startupScript, err := bringauto_prerequisites.CreateAndInitialize[StartupScript]()
	if err != nil {
		return err
	}

	startupChain := BuildChain{
		Chain: []CMDLineInterface{
			startupScript,
		},
	}	
	preparePackageChain := BuildChain{
		Chain: []CMDLineInterface{
			build.Env,
			&gitClone,
			&gitCheckout,
			&gitSubmoduleUpdate,
		},
	}

	shellEvaluator.PreparingCommands = startupChain.GenerateCommands()
	shellEvaluator.Commands = preparePackageChain.GenerateCommands()

	err = shellEvaluator.RunOverSSH(*build.SSHCredentials)
	if err != nil {
		logger := bringauto_log.GetLogger()
		logger.Error("Failed to clone or checkout git repository, check the log file, is the git URI and revision correct?")
		return err
	}

	return nil
}
// prepareForBuild
// Prepares some fields of Build struct for build and makes pre build checks.
func (build *Build) prepareForBuild() error {
	err := build.CheckPrerequisites(nil)
	if err != nil {
		return err
	}

	if build.BuiltPackage == nil {
		return fmt.Errorf("BuiltPackage is nil")
	}

	build.Git.ClonePath = dockerGitCloneDirConst
	build.CMake.SourceDir = dockerGitCloneDirConst

	_, found := build.CMake.Defines["CMAKE_INSTALL_PREFIX"]
	if found {
		return fmt.Errorf("do not specify CMAKE_INSTALL_PREFIX")
	}
	build.CMake.Defines["CMAKE_INSTALL_PREFIX"] = bringauto_const.DockerInstallDirConst

	if build.sysroot != nil {
		build.sysroot.CreateSysrootDir()
		sysPath := build.sysroot.GetSysrootPath()
		err = build.Docker.SetVolume(sysPath, "/sysroot")
		if err != nil {
			return err
		}
		build.CMake.SetDefine("CMAKE_PREFIX_PATH", "/sysroot")
	}

	return nil
}

// RunBuild
// Creates a Docker container, performs a build in it. After build, the files are downloaded to
// local directory and Docker container is stopped and removed. Returns bool which indicates if
// the build was performed succesfully.
func (build *Build) RunBuild() (error, bool) { // Long function - it is hard to refactor to make readability better
	err := build.prepareForBuild()
	if err != nil {
		return err, false
	}

	shellEvaluator := bringauto_ssh.ShellEvaluator{
		Commands: []string{},
	}

	logger := bringauto_log.GetLogger()
	packBuildChainLogger := logger.CreateContextLogger(build.Docker.ImageName, build.Package.GetShortPackageName(), bringauto_log.BuildChainContext)
	if packBuildChainLogger != nil {
		file, err := packBuildChainLogger.GetFile()
		if err != nil {
			logger.Error("Failed to open file - %s", err)
			return err, false
		}
		defer file.Close()

		shellEvaluator.StdOut = file
	}

	dockerRun := (*bringauto_docker.DockerRun)(build.Docker)
	removeHandler := bringauto_process.SignalHandlerAddHandler(func() error {
		// Waiting for docker run command to get container id
		time.Sleep(300 * time.Millisecond)
		return build.stopAndRemoveContainer()
	})
	defer removeHandler()

	logger.InfoIndent("Starting docker container")

	err = dockerRun.Run()
	if err != nil {
		return err, false
	}
	build.SSHCredentials.Port = build.Docker.Port

	logger.InfoIndent("Cloning Package git repository inside docker container")

	err = build.performPreBuildTasks(&shellEvaluator)
	if err != nil {
		return err, false
	}

	build.BuiltPackage.GitCommitHash, err = build.getGitCommitHash()
	if err != nil {
		return fmt.Errorf("can't get git commit hash from container - %w", err), false
	}
	build.BuiltPackage.DirName = build.sysroot.GetDirNameInSysroot()

	if build.sysroot.IsPackageInSysroot(*build.BuiltPackage) {
		logger.InfoIndent("Package already built in sysroot - skipping build")
		return nil, false
	}
	startupScript, err := bringauto_prerequisites.CreateAndInitialize[StartupScript]()
	if err != nil {
		return err, false
	}

	startupChain := BuildChain{
		Chain: []CMDLineInterface{
			startupScript,
		},
	}
	buildChain := BuildChain{
		Chain: []CMDLineInterface{
			build.Env,
			build.CMake,
			build.GNUMake,
		},
	}

	shellEvaluator.PreparingCommands = startupChain.GenerateCommands()
	shellEvaluator.Commands = buildChain.GenerateCommands()

	logger.InfoIndent("Running build inside container")

	err = shellEvaluator.RunOverSSH(*build.SSHCredentials)
	if err != nil {
		return fmt.Errorf("build failed inside docker container, check the log file"), false
	}

	logger.InfoIndent("Copying install files from container to local directory")

	err = build.downloadInstalledFiles()
	if err != nil {
		return fmt.Errorf("can't download files from container to local directory"), false
	}

	return nil, true
}

func (build *Build) SetSysroot(sysroot *bringauto_sysroot.Sysroot) {
	build.sysroot = sysroot
}

func (build *Build) GetLocalInstallDirPath() string {
	workingDir, err := os.Getwd()
	if err != nil {
		logger := bringauto_log.GetLogger()
		logger.Fatal("cannot call Getwd - %s", err)
	}
	suffix := ""
	if build.Docker.Port != bringauto_const.DefaultSSHPort {
		suffix = strconv.Itoa(int(build.Docker.Port) - bringauto_const.DefaultSSHPort)
	}
	copyBaseDir := filepath.Join(workingDir, localInstallDirNameConst + suffix)
	return copyBaseDir
}

func (build *Build) stopAndRemoveContainer() error {
	var err error

	dockerStop := (*bringauto_docker.DockerStop)(build.Docker)
	dockerRm := (*bringauto_docker.DockerRm)(build.Docker)
	logger := bringauto_log.GetLogger()
	err = dockerStop.Stop()
	if err != nil {
		logger.Error("Can't stop container - %s", err)
	}
	err = dockerRm.RemoveContainer()
	if err != nil {
		logger.Error("Can't remove container - %s", err)
	}
	return nil
}

func (build *Build) CleanUp() error {
	var err error
	copyDir := build.GetLocalInstallDirPath()
	if _, err = os.Stat(copyDir); os.IsNotExist(err) {
		return nil
	}
	err = os.RemoveAll(copyDir)
	if err != nil {
		return err
	}
	return nil
}

func (build *Build) downloadInstalledFiles() error {
	var err error

	copyDir := build.GetLocalInstallDirPath()
	if _, err = os.Stat(copyDir); os.IsNotExist(err) {
		err = os.MkdirAll(copyDir, 0766)
		if err != nil {
			return fmt.Errorf("cannot create directory %s", copyDir)
		}
	}

	sftpClient := bringauto_ssh.SFTP{
		RemoteDir:      bringauto_const.DockerInstallDirConst,
		EmptyLocalDir:  copyDir,
		SSHCredentials: build.SSHCredentials,
	}

	packTarLogger := bringauto_log.GetLogger().CreateContextLogger(build.Docker.ImageName, build.Package.GetShortPackageName(), bringauto_log.TarContext)
	if packTarLogger != nil {
		logFile, err := packTarLogger.GetFile()
		if err != nil {
			return fmt.Errorf("failed to open file - %w", err)
		}
		defer logFile.Close()

		sftpClient.LogWriter = logFile
	}

	err = sftpClient.DownloadDirectory()
	return err
}

func (build *Build) getGitCommitHash() (string, error) {
	pipeReader, pipeWriter := io.Pipe()
	defer pipeReader.Close()
	defer pipeWriter.Close()
	gitGetHash := bringauto_git.GitGetHash{Git: *build.Git}
	shellEvaluator := bringauto_ssh.ShellEvaluator{
		Commands: gitGetHash.ConstructCMDLine(),
		StdOut:   pipeWriter,
	}

	err := shellEvaluator.RunOverSSH(*build.SSHCredentials)
	if err != nil {
		return "", err
	}

	buf := bufio.NewReader(pipeReader)
	var line string

	for {
		line, err = buf.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", err
		}
		if err == nil { // The newline character is present
			line = line[:len(line)-1]
		}

		hash := getGitCommitHashFromLine(line)
		if hash != "" {
			return hash, nil
		}

		if err == io.EOF {
			break
		}
	}

	return "", fmt.Errorf("no commit hash in output")
}

func getGitCommitHashFromLine(line string) string {
	// regexp for long git commit hash, it must be used, because the ssh output has several commands and it is long
	re := regexp.MustCompile("[a-f0-9]{40}")
	match := re.FindStringSubmatch(line)
	if len(match) > 0 {
		return match[0]
	}

	return ""
}
