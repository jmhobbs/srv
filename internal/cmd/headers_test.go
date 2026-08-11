package cmd_test

import (
	"os"
	"path"
	"testing"

	"github.com/jmhobbs/srv/internal/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LoadHeadersFile_Valid(t *testing.T) {
	dir := t.TempDir()
	file := path.Join(dir, "headers.txt")
	err := os.WriteFile(file, []byte("/\n\tx-foo: bar\n"), 0644)
	require.NoError(t, err)
	h, err := cmd.LoadHeadersFile(file)
	require.NoError(t, err)
	require.NotNil(t, h)
	assert.Len(t, *h, 1)
}

func Test_LoadHeadersFile_Missing(t *testing.T) {
	dir := t.TempDir()
	file := path.Join(dir, "doesnt_exist.txt")
	_, err := cmd.LoadHeadersFile(file)
	assert.Error(t, err)
}

func Test_LoadHeadersFile_Invalid(t *testing.T) {
	dir := t.TempDir()
	file := path.Join(dir, "invalid_headers.txt")
	err := os.WriteFile(file, []byte("\tinvalid"), 0644)
	require.NoError(t, err)
	_, err = cmd.LoadHeadersFile(file)
	assert.Error(t, err)
}
