package cmd

import (
	"io"

	"github.com/rs/zerolog"
)

func NewLogger(out io.Writer, opts Options) zerolog.Logger {
	level := zerolog.InfoLevel
	if opts.Quiet {
		level = zerolog.WarnLevel
	} else if opts.Verbose {
		level = zerolog.DebugLevel
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	return zerolog.New(zerolog.ConsoleWriter{Out: out}).With().Timestamp().Logger().Level(level)
}
