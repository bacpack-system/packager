package git

import "github.com/bacpack-system/packager/internal/prerequisites"

func (git *Git) FillDefault(*prerequisites.Args) error {
	return nil
}

func (git *Git) FillDynamics(*prerequisites.Args) error {
	return nil
}

// CheckPrerequisites
// Function should check if the git can be run and if not it returns error
// (not nil value)
func (git *Git) CheckPrerequisites(*prerequisites.Args) error {
	// Git server as cmdline constructor for remote server is not good
	return nil
}
