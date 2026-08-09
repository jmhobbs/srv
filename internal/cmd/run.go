package cmd

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/jmhobbs/srv/internal/accesslogs"
	"github.com/jmhobbs/srv/internal/handler"
)

func Run(args []string, version string) {
	opts, err := ParseFlags(args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
	}
	if opts.ShowVersion {
		fmt.Fprintf(os.Stderr, "srv - %s\n", version)
		os.Exit(0)
	}
	logger := NewLogger(os.Stderr, *opts)

	absDir, err := GetWorkingDirectory(flag.Arg(0))
	if err != nil {
		logger.Fatal().Err(err).Msg("Could not resolve working directory")
	}

	headers, err := LoadHeadersFile(opts.HeadersFile)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Warn().Str("path", opts.HeadersFile).Msg("No headers file found")
		} else {
			logger.Fatal().Err(err).Msg("Error loading headers file")
		}
	}

	h, err := accesslogs.Setup(
		handler.New(logger, absDir, strings.Split(opts.DefaultDirectoryFiles, ","), headers),
		opts.AccessLogs,
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("Error setting up access logs.")
	}

	logger.Info().Msgf("Listening on http://%s/...", opts.Address)
	err = http.ListenAndServe(opts.Address, h)
	if err != nil {
		logger.Error().Err(err).Msg("Server shutdown unexpectedly")
		os.Exit(1)
	}

	logger.Info().Msg("Goodbye")
}
