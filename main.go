package main

import (
	"os"

	"github.com/jmhobbs/srv/internal/cmd"
)

var version string = "v0.0.0-dev"

func main() {
	cmd.Run(os.Args, version)
}
