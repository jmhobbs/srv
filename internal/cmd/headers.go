package cmd

import (
	"os"

	headers "github.com/jmhobbs/cloudflare-headers-file"
)

func LoadHeadersFile(path string) (*headers.File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	//nolint:errcheck
	defer f.Close()

	return headers.Parse(f)
}
