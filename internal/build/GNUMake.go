package build

import (
	"github.com/bacpack-system/packager/internal/prerequisites"
	"strconv"
	"strings"
)

const (
	makeJobsCount = 10
)

// GNUMake cmd line interface for standard GNU Make utility
type GNUMake struct {}

func (make *GNUMake) FillDefault(*prerequisites.Args) error {
	return nil
}

func (make *GNUMake) FillDynamic(*prerequisites.Args) error {
	return nil
}

func (make *GNUMake) CheckPrerequisites(*prerequisites.Args) error {
	return nil
}

func (make *GNUMake) ConstructCMDLine() []string {
	cmdBuild := []string{"make", "-j", strconv.Itoa(makeJobsCount)}
	cmdInstall := []string{"make", "install"}
	return []string{
		strings.Join(cmdBuild, " "),
		strings.Join(cmdInstall, " "),
	}
}
