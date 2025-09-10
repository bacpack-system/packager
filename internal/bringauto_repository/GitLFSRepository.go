package bringauto_repository

import (
	"github.com/bacpack-system/packager/internal/bringauto_package"
	"github.com/bacpack-system/packager/internal/bringauto_prerequisites"
	"github.com/bacpack-system/packager/internal/bringauto_log"
	"github.com/bacpack-system/packager/internal/bringauto_context"
	"github.com/bacpack-system/packager/internal/bringauto_config"
	"github.com/bacpack-system/packager/internal/bringauto_const"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"slices"
)

// GitLFSRepository represents Package/App repository based on Git LFS
type GitLFSRepository struct {
	GitRepoPath string
}

const (
	gitExecutablePath = "/usr/bin/git"
	// Count of files which will be list in warnings
	listFileCount = 10
)

func (lfs *GitLFSRepository) FillDefault(args *bringauto_prerequisites.Args) error {
	return nil
}

func (lfs *GitLFSRepository) FillDynamic(args *bringauto_prerequisites.Args) error {
	return nil
}

func (lfs *GitLFSRepository) CheckPrerequisites(*bringauto_prerequisites.Args) error {
	if _, err := os.Stat(lfs.GitRepoPath); os.IsNotExist(err) {
		return fmt.Errorf("package repository '%s' does not exist", lfs.GitRepoPath)
	}
	if _, err := os.Stat(lfs.GitRepoPath + "/.git"); os.IsNotExist(err) {
		return fmt.Errorf("package repository '%s' is not a git repository", lfs.GitRepoPath)
	}

	isStatusEmpty := lfs.gitIsStatusEmpty()
	if !isStatusEmpty {
		return fmt.Errorf("sorry, but the given git root does not have empty `git status`. clean up changes and try again")
	}
	return nil
}

// commitPackage
// Adds all changes to staged and then makes a commit with packageName description.
func (lfs *GitLFSRepository) commitPackage(packageName string) error {
	err := lfs.gitAddAll()
	if err != nil {
		return err
	}
	err = lfs.gitCommit(packageName)
	if err != nil {
		return err
	}

	return nil
}

// RestoreAllChanges
// Restores all changes in repository and cleans all untracked changes.
func (lfs *GitLFSRepository) RestoreAllChanges() error {
	var err error
	if !lfs.isRepoEmpty() {
		err = lfs.gitRestoreAll()
		if err != nil {
			return err
		}
	}

	err = lfs.gitCleanAll()
	if err != nil {
		return err
	}
	return nil
}

// dividePackagesForCurrentImage
// Divides allConfigs to packages for imageName and not for imageName and returns 2 slices.
func dividePackagesForCurrentImage(allConfigs []*bringauto_config.Config, imageName string) ([]bringauto_package.Package, []bringauto_package.Package) {
	var packagesForImage []bringauto_package.Package
	var packagesNotForImage []bringauto_package.Package

	for _, config := range allConfigs {
		if slices.Contains(config.DockerMatrix.ImageNames, imageName) {
			packagesForImage = append(packagesForImage, config.Package)
		} else {
			packagesNotForImage = append(packagesNotForImage, config.Package)
		}
	}

	return packagesForImage, packagesNotForImage
}

// compareConfigsAndGitLfs
// Compares Packages/Apps (depends on packageOrApp string) in Context and in Git Lfs and returns
// three arrays of strings - error Paths, expected Paths for imageName and expected paths not for
// imageName.
func (lfs *GitLFSRepository) compareConfigsAndGitLfs(
	platformString *bringauto_package.PlatformString,
	configs []*bringauto_config.Config,
	imageName string,
	packageOrApp string,
) (error, []string, []string, []string) {
	packagesForImage, packagesNotForImage := dividePackagesForCurrentImage(configs, imageName)

	var errorPaths, expectedPathsForImage, expectedPathsNotForImage []string
	for _, pack := range packagesForImage {
		packPath := filepath.Join(lfs.CreatePath(pack, packageOrApp) + "/" + pack.GetFullPackageName() + ".zip")
		expectedPathsForImage = append(expectedPathsForImage, packPath)
	}
	for _, pack := range packagesNotForImage {
		packPath := filepath.Join(lfs.CreatePath(pack, packageOrApp) + "/" + pack.GetFullPackageName() + ".zip")
		expectedPathsNotForImage = append(expectedPathsNotForImage, packPath)
	}

	lookupPath := filepath.Join(
		lfs.GitRepoPath,
		packageOrApp,
		platformString.String.DistroName,
		platformString.String.DistroRelease,
		platformString.String.Machine,
	)

	_, err := os.Stat(lookupPath)
	if !os.IsNotExist(err) {
		err = filepath.WalkDir(lookupPath, func(path string, d fs.DirEntry, err error) error {
			if d.Name() == ".git" && d.IsDir() {
				return filepath.SkipDir
			}
			if !d.IsDir() {
				if !slices.Contains(expectedPathsForImage, path) {
					errorPaths = append(errorPaths, path)
				} else {
					// Remove element from expected package paths
					index := slices.Index(expectedPathsForImage, path)
					expectedPathsForImage[index] = expectedPathsForImage[len(expectedPathsForImage) - 1]
					expectedPathsForImage = expectedPathsForImage[:len(expectedPathsForImage) - 1]
				}
			}
			return nil
		})
		if err != nil {
			return err, []string{}, []string{}, []string{}
		}
	}

	return nil, errorPaths, expectedPathsForImage, expectedPathsNotForImage
}

// comparePackagesAndGitLfs
// Compares Packages in Context and in Git Lfs and returns three arrays of strings - error Paths,
// expected Paths for imageName and expected paths not for imageName.
func (lfs *GitLFSRepository) comparePackagesAndGitLfs(
	contextManager *bringauto_context.ContextManager,
	platformString *bringauto_package.PlatformString,
	imageName string,
) (error, []string, []string, []string) {
	configs := contextManager.GetAllPackageConfigsArray(platformString)
	return lfs.compareConfigsAndGitLfs(platformString, configs, imageName, bringauto_const.PackageDirName)
}

// compareAppsAndGitLfs
// Compares Apps in Context and in Git Lfs and returns three arrays of strings - error Paths,
// expected Paths for imageName and expected paths not for imageName.
func (lfs *GitLFSRepository) compareAppsAndGitLfs(
	contextManager *bringauto_context.ContextManager,
	platformString *bringauto_package.PlatformString,
	imageName string,
) (error, []string, []string, []string) {
	configs := contextManager.GetAllAppConfigsArray(platformString)
	return lfs.compareConfigsAndGitLfs(platformString, configs, imageName, bringauto_const.AppDirName)
}

// CheckGitLfsConsistency
// Checks Git Lfs consistency based on Context. Prints and returns errors if in Git Lfs is any
// Package/App which is not in Context. Prints warnings if the Git Lfs is missing any Packages/Apps
// present in Context and prints warnings if any Package/App won't build for current imageName.
func (lfs *GitLFSRepository) CheckGitLfsConsistency(contextManager *bringauto_context.ContextManager, platformString *bringauto_package.PlatformString, imageName string) error {
	err, errorPaths, expectedPathsForImage, expectedPathsNotForImage := lfs.comparePackagesAndGitLfs(contextManager, platformString, imageName)
	if err != nil {
		return err
	}

	err, errorPaths_Apps, expectedPathsForImage_Apps, expectedPathsNotForImage_Apps := lfs.compareAppsAndGitLfs(contextManager, platformString, imageName)

	errorPaths = append(errorPaths, errorPaths_Apps...)
	expectedPathsForImage = append(expectedPathsForImage, expectedPathsForImage_Apps...)
	expectedPathsNotForImage = append(expectedPathsNotForImage, expectedPathsNotForImage_Apps...)

	err = printErrors(errorPaths, expectedPathsForImage, expectedPathsNotForImage)
	if err != nil {
		return err
	}
	return nil
}

// printErrors
// Prints errors and warnings for Git Lfs consistency check.
func printErrors(errorPaths []string, expectedPathsForImage []string, expectedPathsNotForImage []string) error {
	logger := bringauto_log.GetLogger()
	if len(errorPaths) > 0 {
		logger.Error("%d Packages/Apps are not in Json definitions but are in Git Lfs (listing first %d):", len(errorPaths), listFileCount)
		for i, errorPath := range errorPaths {
			if i > listFileCount - 1 {
				break
			}
			logger.ErrorIndent("%s", errorPath)
		}
		return fmt.Errorf("Packages/Apps in Git Lfs are not subset of Packages/Apps in Json definitions")
	}

	if len(expectedPathsForImage) > 0 {
		logger.Warn("Expected %d Packages/Apps (built for target image) to be in Git Lfs (listing first %d):", len(expectedPathsForImage), listFileCount)
		for i, expectedPathForImage := range expectedPathsForImage {
			if i > listFileCount - 1 {
				break
			}
			logger.WarnIndent("%s", expectedPathForImage)
		}
	}
	if len(expectedPathsNotForImage) > 0 {
		logger.Warn("%d Packages/Apps are in context but are not built for target image (listing first %d):", len(expectedPathsNotForImage), listFileCount)
		for i, expectedPathNotForImage := range expectedPathsNotForImage {
			if i > listFileCount - 1 {
				break
			}
			logger.WarnIndent("%s", expectedPathNotForImage)
		}
	}
	return nil
}

// CreatePath
// Returns path for specific pack inside Git Lfs. The path depends on packageOrApp string which
// should be either "package" or "app".
func (lfs *GitLFSRepository) CreatePath(pack bringauto_package.Package, packageOrApp string) string {
	repositoryPath := path.Join(
		pack.PlatformString.String.DistroName,
		pack.PlatformString.String.DistroRelease,
		pack.PlatformString.String.Machine,
		pack.Name,
	)
	return path.Join(lfs.GitRepoPath, packageOrApp, repositoryPath)
}

// CopyToRepository
// Copies the pack to the Git LFS repository. packageOrApp is a string representing type of
// Package, it should be either "package" or "app". Each Package/App is stored in different
// directory structure represented by
// packageOrApp / PlatformString.DistroName / PlatformString.DistroRelease / PlatformString.Machine / <package>
func (lfs *GitLFSRepository) CopyToRepository(pack bringauto_package.Package, sourceDir string, packageOrApp string) error {
	archiveDirectory := lfs.CreatePath(pack, packageOrApp)

	var err error
	err = os.MkdirAll(archiveDirectory, 0755)
	if err != nil {
		return err
	}

	err = pack.CreatePackage(sourceDir, archiveDirectory)
	if err != nil {
		return err
	}

	err = lfs.commitPackage(pack.GetFullPackageName())
	if err != nil {
		return err
	}

	return nil
}

// gitIsStatusEmpty
// Returns true, if the git status in Git Lfs is empty, else returns false.
func (lfs *GitLFSRepository) gitIsStatusEmpty() bool {
	var ok, buffer = lfs.prepareAndRun([]string{
		"status",
		"-s",
	},
	)
	if !ok {
		return false
	}
	if buffer.Len() != 0 {
		return false
	}
	return true
}

// gitAddAll
// Adds all to staged in Git Lfs.
func (lfs *GitLFSRepository) gitAddAll() error {
	var ok, _ = lfs.prepareAndRun([]string{
		"add",
		".",
	},
	)
	if !ok {
		return fmt.Errorf("cannot add changes in Git Lfs")
	}
	return nil
}

// gitCommit
// Commits Git Lfs with packageName description.
func (lfs *GitLFSRepository) gitCommit(packageName string) error {
	var ok, _ = lfs.prepareAndRun([]string{
		"commit",
		"-m",
		"Build package " + packageName,
	},
	)
	if !ok {
		return fmt.Errorf("cannot commit changes in Git Lfs")
	}
	return nil
}

// gitRestoreAll
// Restores all changes in Git Lfs.
func (lfs *GitLFSRepository) gitRestoreAll() error {
	var ok, _ = lfs.prepareAndRun([]string{
		"restore",
		".",
	},
	)
	if !ok {
		return fmt.Errorf("cannot restore changes in Git Lfs")
	}
	return nil
}

// gitCleanAll
// Cleans all changes in Git Lfs.
func (lfs *GitLFSRepository) gitCleanAll() error {
	var ok, _ = lfs.prepareAndRun([]string{
		"clean",
		"-f",
		".",
	},
	)
	if !ok {
		return fmt.Errorf("cannot clean changes in Git Lfs")
	}
	return nil
}

// isRepoEmpty
// Returns true, if the Git Lfs is empty (no commits), else returns false.
func (lfs *GitLFSRepository) isRepoEmpty() bool {
	var ok, _ = lfs.prepareAndRun([]string{
		"log",
	},
	)
	return !ok
}

func (repo *GitLFSRepository) prepareAndRun(cmdline []string) (bool, *bytes.Buffer) {
	var cmd exec.Cmd
	var outBuffer bytes.Buffer
	var err error

	repoPath := repo.GitRepoPath
	if !filepath.IsAbs(repoPath) {
		workingDir, err := os.Getwd()
		if err != nil {
			return false, nil
		}
		repoPath = path.Join(workingDir, repoPath)
	}

	cmd.Dir = repoPath
	cmdArgs := append([]string{gitExecutablePath}, cmdline...)
	cmd.Args = cmdArgs
	cmd.Path = gitExecutablePath
	cmd.Stdout = &outBuffer
	err = cmd.Run()
	if err != nil {
		return false, &outBuffer
	}
	if cmd.ProcessState.ExitCode() > 0 {
		return false, &outBuffer
	}
	return true, &outBuffer
}
