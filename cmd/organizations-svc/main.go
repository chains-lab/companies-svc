package main

import (
	"os"

	"github.com/chains-lab/organizations-svc/cmd/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
