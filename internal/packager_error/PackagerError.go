package packager_error

import (
	"errors"
)

const (
	DEFAULT_ERROR                    = 1 // Not specified errors
	CMD_LINE_ERROR                   = 2 // Cmd line parsing errors
	CONTEXT_ERROR                    = 3 // Context consistency errors
	GIT_LFS_ERROR                    = 4 // Git Lfs and Context comparison consistency errors
	BUILD_ERROR                      = 5 // Build image, App or Package errors 
	PACKAGE_MISSING_DEPENDENCY_ERROR = 6 // Package dependency is not on sysroot error
	CREATING_SYSROOT_ERROR           = 7 // Creating sysroot errors
	OVERWRITE_FILE_IN_SYSROOT_ERROR  = 8 // Overwriting files in sysroot error
)

var CmdLineErr = errors.New("cmd parse error")
var ContextErr = errors.New("context consistency error")
var GitLfsErr = errors.New("git lfs consistency error")
var BuildErr = errors.New("build error")
var PackageMissingDependencyErr = errors.New("package missing dependency in sysroot error")
var CreatingSysrootErr = errors.New("creating sysroot error")
var OverwriteFileInSysrootErr = errors.New("trying to overwrite file in sysroot error")

func GetReturnCode(err error) int {
	if errors.Is(err, CmdLineErr) {
		return CMD_LINE_ERROR
	} else if errors.Is(err, ContextErr) {
		return CONTEXT_ERROR
	} else if errors.Is(err, GitLfsErr) {
		return GIT_LFS_ERROR
	} else if errors.Is(err, PackageMissingDependencyErr) {
		return PACKAGE_MISSING_DEPENDENCY_ERROR
	} else if errors.Is(err, BuildErr) {
		return BUILD_ERROR
	} else if errors.Is(err, CreatingSysrootErr) {
		return CREATING_SYSROOT_ERROR
	} else if errors.Is(err, OverwriteFileInSysrootErr) {
		return OVERWRITE_FILE_IN_SYSROOT_ERROR
	}
	return DEFAULT_ERROR
}
