package build

import (
	"fmt"
	"runtime"
)

var (
	GitBranch string
	GitCommit string

	// Build time in UTC with format %Y.%m.%d.%H%M
	BuildTime string
)

// Git commit and branch at build time
func GitInfo() string {
	return fmt.Sprintf("%s-%s", GitCommit, GitBranch)
}

// Go version at build time
func GoVersion() string {
	return runtime.Version()
}

func Version() string {
	return fmt.Sprintf("%s-%s", BuildTime, GitInfo())
}

func Print() {
	fmt.Println("=== Build Info ===")
	fmt.Printf("Go Version = %s\n", GoVersion())
	fmt.Printf("App Version = %s\n", Version())
	fmt.Println()
}
