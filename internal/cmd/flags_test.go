package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jmhobbs/srv/internal/cmd"
)

func Test_Flags(t *testing.T) {
	opts, err := cmd.ParseFlags([]string{
		"-p", "5432",
		"-interface", "127.1.2.3",
		"-q",
		"-v",
		"-version",
		"-access-log", "access.log",
		"-default-dir-files", "dir.html",
		"-headers-file", "headers.txt",
	})
	require.NoError(t, err)
	assert.Equal(t, "127.1.2.3:5432", opts.Address)
	assert.True(t, opts.Quiet)
	assert.True(t, opts.Verbose)
	assert.True(t, opts.ShowVersion)
	assert.Equal(t, "access.log", opts.AccessLogs)
	assert.Equal(t, "dir.html", opts.DefaultDirectoryFiles)
	assert.Equal(t, "headers.txt", opts.HeadersFile)
}
