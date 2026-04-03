package main

import (
	"fmt"
	"os"

	"github.com/zarishsphere/zs-core-fhir-engine/cmd/fhir-engine/internal/cli"
)

// Version information (injected at build time via ldflags)
var (
	Version = "develop"
	Commit  = ""
	Date    = ""
)

func main() {
	if err := cli.Run(Version, Commit, Date); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
