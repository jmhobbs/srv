package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

const VERSION string = "0.0.4"

type options struct {
	address               string
	quiet                 bool
	verbose               bool
	showVersion           bool
	accessLogs            string
	defaultDirectoryFiles string
	headersFile           string
}

func flags() options {
	var (
		port                  = flag.Int("p", 5050, "Port to listen on")
		iface                 = flag.String("interface", "127.0.0.1", "Network interface to listen on")
		quiet                 = flag.Bool("q", false, "Quiet mode, disable most logging")
		verbose               = flag.Bool("v", false, "Verbose mode, enable debug logging")
		showVersion           = flag.Bool("version", false, "Show version and exit")
		accessLogs            = flag.String("access-log", "-", "Where to write access logs, default is STDOUT. Pass empty string to disable.")
		defaultDirectoryFiles = flag.String("default-dir-files", "index.html,index.htm", "Default files to show for directory, when present.")
		headersFile           = flag.String("headers-file", "_headers", "Path to _headers file to apply.")
	)
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), "usage: srv [options] [directory]\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	return options{
		address:               fmt.Sprintf("%s:%d", *iface, *port),
		quiet:                 *quiet,
		verbose:               *verbose,
		showVersion:           *showVersion,
		accessLogs:            *accessLogs,
		defaultDirectoryFiles: *defaultDirectoryFiles,
		headersFile:           *headersFile,
	}
}

func getWorkingDirectory(dir string) (string, error) {
	if dir == "" {
		cwd, err := syscall.Getwd()
		if err != nil {
			return "", fmt.Errorf("unable to determine current working directory: %v", err)
		}
		dir = cwd
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("unable to determine absolute path for %q: %v", dir, err)
	}

	return absDir, nil
}

func main() {
	opts := flags()
	if opts.showVersion {
		fmt.Fprintf(flag.CommandLine.Output(), "srv - v%s\n", VERSION)
		os.Exit(0)
	}

	logLevel := INFO
	if opts.quiet {
		logLevel = ERROR
	} else if opts.verbose {
		logLevel = DEBUG
	}
	logger := newLogger(os.Stderr, logLevel)

	absDir, err := getWorkingDirectory(flag.Arg(0))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	headers, err := loadHeadersFile(opts.headersFile)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Warn("No headers file found at %q", opts.headersFile)
		} else {
			logger.Error("Error loading headers file: %v", err)
			os.Exit(1)
		}
	}

	handler, err := setupAccessLogs(
		newHandler(logger, absDir, strings.Split(opts.defaultDirectoryFiles, ","), headers),
		opts.accessLogs,
	)
	if err != nil {
		logger.Error("Error setting up access logs: %v", err)
		os.Exit(1)
	}

	logger.Info("Listening on http://%s/...", opts.address)
	err = http.ListenAndServe(opts.address, handler)
	if err != nil {
		logger.Error("Error serving: %v", err)
		os.Exit(1)
	}

	logger.Info("Goodbye")
}

func setupAccessLogs(handler http.Handler, accessLogs string) (http.Handler, error) {
	if accessLogs != "" {
		var sink io.Writer
		if accessLogs == "-" {
			sink = os.Stdout
		} else {
			f, err := os.OpenFile(accessLogs, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return nil, fmt.Errorf("error opening access log: %v", err)
			}
			//nolint:errcheck
			defer f.Close()
			sink = f
		}
		return LoggingMiddleware(sink, handler), nil
	}
	return handler, nil
}
