package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"syscall"
)

const VERSION string = "0.0.1"

func main () {
	var (
		port *int = flag.Int("p", 5000, "Port to listen on")
		iface *string = flag.String("interface", "127.0.0.1", "Network interface to listen on")
		quiet *bool = flag.Bool("q", false, "Quiet mode, disable most logging")
		verbose *bool = flag.Bool("v", false, "Verbose mode, enable debug logging")
		showVersion *bool = flag.Bool("version", false, "Show version and exit")
	)
	flag.Usage = func ()  {
		fmt.Fprint(flag.CommandLine.Output(), "usage: srv [options] [directory]\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *showVersion {
		fmt.Fprintf(flag.CommandLine.Output(), "srv - v%s\n", VERSION)
		os.Exit(0)
	}

	logLevel := INFO
	if *quiet {
		logLevel = ERROR
	} else if *verbose {
		logLevel = DEBUG
	}
	logger := newLogger(os.Stderr, logLevel)

	dir := flag.Arg(0)
	if "" == dir {
		cwd, err := syscall.Getwd()
		if err != nil {
			logger.Error("unable to determine current working directory: %v", err)
			os.Exit(1)
		}
		dir = cwd
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		logger.Error("unable to determine path to host: %v", err)
		os.Exit(1)
	}

	addr := fmt.Sprintf("%s:%d", *iface, *port)

	fmt.Printf("Listening on %s...\n", addr)

	var handler http.Handler = newHandler(logger, absDir)
	if ! *quiet {
		handler = LoggingMiddleware(logger, handler)
	}

	err = http.ListenAndServe(addr, handler)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}

	fmt.Println("Goodbye.")
}
