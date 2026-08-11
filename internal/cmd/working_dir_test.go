package cmd_test

import (
	"path"
	"syscall"
	"testing"

	"github.com/jmhobbs/srv/internal/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetWorkingDirectory(t *testing.T) {
	expected, err := syscall.Getwd()
	require.NoError(t, err)
	cwd, err := cmd.GetWorkingDirectory("")
	require.NoError(t, err)
	assert.Equal(t, expected, cwd)
}

func Test_GetWorkingDirectory_Override(t *testing.T) {
	dir := t.TempDir()
	cwd, err := cmd.GetWorkingDirectory(dir)
	require.NoError(t, err)
	assert.Equal(t, dir, cwd)
}

func Test_GetWorkingDirectory_OverrideAbsolute(t *testing.T) {
	dir := t.TempDir()
	cwd, err := cmd.GetWorkingDirectory(dir + "/subdir/../subdir/")
	require.NoError(t, err)
	assert.Equal(t, path.Join(dir, "subdir"), cwd)
}
