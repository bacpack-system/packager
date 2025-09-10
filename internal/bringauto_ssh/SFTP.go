package bringauto_ssh

import (
	"github.com/bacpack-system/packager/internal/bringauto_const"
	"github.com/bacpack-system/packager/internal/bringauto_prerequisites"
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/mholt/archiver/v3"
	"github.com/pkg/sftp"
)

const (
	archiveName    string = "install_arch.tar"
	archiveNameSep string = string(os.PathSeparator) + archiveName
	// Size of the buffer used by bufio module
	bufferSize = 1024*1024
)

type SFTP struct {
	// Path to a directory on the remote machine
	RemoteDir string
	// Empty, existing local directory where the RemoteDir will be copy
	EmptyLocalDir  string
	SSHCredentials *SSHCredentials
	LogWriter      io.Writer
}

// DownloadDirectory
// Download directory from RemoteDir to EmptyLocalDir.
// EmptyLocalDir must be empty!
// Function returns error in case of problem or nil if succeeded.
func (sftpd *SFTP) DownloadDirectory() error {
	var err error

	tar, err := bringauto_prerequisites.CreateAndInitialize[Tar](archiveName, bringauto_const.DockerInstallDirConst)
	if err != nil {
		return fmt.Errorf("cannot initialize Tar - %w", err)
	}

	shellEvaluator := ShellEvaluator{
		Commands: tar.ConstructCMDLine(),
		StdOut:   sftpd.LogWriter,
	}

	err = shellEvaluator.RunOverSSH(*sftpd.SSHCredentials)

	if err != nil {
		return fmt.Errorf("cannot archive %s dir in docker container - %w", bringauto_const.DockerInstallDirConst, err)
	}

	sshSession := SSHSession{}
	err = sshSession.LoginMultipleAttempts(*sftpd.SSHCredentials)
	if err != nil {
		return fmt.Errorf("SFTP DownloadDirectory error - %w", err)
	}

	sftpClient, err := sftp.NewClient(sshSession.sshClient,
		sftp.MaxConcurrentRequestsPerFile(64),
		sftp.UseConcurrentReads(true),
		sftp.UseFstat(true),
		sftp.MaxPacket(1<<15),
	)
	if err != nil {
		return fmt.Errorf("SFTP DownloadDirectory problem - %w", err)
	}
	defer sftpClient.Close()

	if _, err = os.Stat(sftpd.EmptyLocalDir); os.IsNotExist(err) {
		return fmt.Errorf("EmptyLocalDir '%s' does not exist", sftpd.EmptyLocalDir)
	}

	localPathDirContent, _ := os.ReadDir(sftpd.EmptyLocalDir)
	localPathDirIsNotEmpty := len(localPathDirContent) != 0
	if localPathDirIsNotEmpty {
		return fmt.Errorf("local directory '%s' is not empty", sftpd.EmptyLocalDir)
	}

	localArchivePath := sftpd.EmptyLocalDir + archiveNameSep

	err = sftpd.copyFile(sftpClient, sftpd.RemoteDir+archiveNameSep, localArchivePath)
	if err != nil {
		return fmt.Errorf("cannot copy recursive  - %w", err)
	}

	tarArchive := archiver.Tar{
		OverwriteExisting:      false,
		MkdirAll:               false,
		ImplicitTopLevelFolder: false,
		ContinueOnError:        true,
	}

	err = tarArchive.Unarchive(localArchivePath, sftpd.EmptyLocalDir)
	if err != nil {
		return fmt.Errorf("cannot unarchive tar archive locally - %w", err)
	}

	err = os.Remove(localArchivePath)
	if err != nil {
		return fmt.Errorf("cannot remove local archive %s - %w", localArchivePath, err)
	}

	return nil
}

func (sftpd *SFTP) copyFile(sftpClient *sftp.Client, remoteFile string, localDir string) error {
	var err error
	remotePathStat, err := sftpClient.Lstat(remoteFile)
	if os.IsNotExist(err) {
		return fmt.Errorf("requested remote file %s does not exist", remoteFile)
	} else if err != nil {
		return fmt.Errorf("error retrieving %s remote file info - %w", remoteFile, err)
	}

	normalizedLocalDir, _ := normalizePath(localDir)
	sourceFile, err := sftpClient.Open(remoteFile)
	if err != nil {
		return fmt.Errorf("cannot open file %s for read - %w", remoteFile, err)
	}
	destFile, err := os.OpenFile(normalizedLocalDir, os.O_RDWR|os.O_CREATE, remotePathStat.Mode().Perm())
	if err != nil {
		return err
	}

	return copyIOFile(sourceFile, destFile)
}

func copyIOFile(sourceFile *sftp.File, destFile *os.File) error {
	sourceFileBuff := bufio.NewReaderSize(sourceFile, bufferSize)
	destFileBuff := bufio.NewWriterSize(destFile, bufferSize)

	var err error
	_, err = io.Copy(destFileBuff, sourceFileBuff)
	if err != nil {
		return fmt.Errorf("cannot copy remote IO files - %w", err)
	}

	err = destFileBuff.Flush()
	if err != nil {
		return fmt.Errorf("cannot flush destination buffer - %w", err)
	}

	err = destFile.Close()
	if err != nil {
		return fmt.Errorf("cannot close destination file - %w", err)
	}
	err = sourceFile.Close()
	if err != nil {
		return fmt.Errorf("cannot close source file - %w", err)
	}
	return nil
}

// normalizePath
func normalizePath(p string) (string, error) {
	regexp, regexpErr := regexp.CompilePOSIX("[/]{2,}")
	if regexpErr != nil {
		return "", fmt.Errorf("sftp cannot compile regex - %s", regexpErr)
	}
	normalizePath := regexp.ReplaceAllString(p, "/")
	return normalizePath, nil
}
