package main

import (
	"os"

	headers "github.com/jmhobbs/cloudflare-headers-file"
)

func loadHeadersFile(path string) (*headers.File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return headers.Parse(f)
}
