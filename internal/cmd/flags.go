package cmd

import (
	"flag"
	"fmt"
)

type Options struct {
	Address               string
	Quiet                 bool
	Verbose               bool
	ShowVersion           bool
	AccessLogs            string
	DefaultDirectoryFiles string
	HeadersFile           string
}

func ParseFlags(args []string) (*Options, error) {
	var (
		fs                    = flag.NewFlagSet("srv", flag.ExitOnError)
		port                  = fs.Int("p", 5050, "Port to listen on")
		iface                 = fs.String("interface", "127.0.0.1", "Network interface to listen on")
		quiet                 = fs.Bool("q", false, "Quiet mode, disable most logging")
		verbose               = fs.Bool("v", false, "Verbose mode, enable debug logging")
		showVersion           = fs.Bool("version", false, "Show version and exit")
		accessLogs            = fs.String("access-log", "-", "Where to write access logs, default is STDOUT. Pass empty string to disable.")
		defaultDirectoryFiles = fs.String("default-dir-files", "index.html,index.htm", "Default files to show for directory, when present.")
		headersFile           = fs.String("headers-file", "_headers", "Path to _headers file to apply.")
	)
	fs.Usage = func() {
		fmt.Fprint(fs.Output(), "usage: srv [options] [directory]\n\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	return &Options{
		Address:               fmt.Sprintf("%s:%d", *iface, *port),
		Quiet:                 *quiet,
		Verbose:               *verbose,
		ShowVersion:           *showVersion,
		AccessLogs:            *accessLogs,
		DefaultDirectoryFiles: *defaultDirectoryFiles,
		HeadersFile:           *headersFile,
	}, nil
}
