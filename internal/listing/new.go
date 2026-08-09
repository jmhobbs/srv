package listing

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type DirectoryListing struct {
	Path    string
	Parent  string
	Entries []DirectoryListingEntry
}

type DirectoryListingEntry struct {
	Name  string
	Path  string
	IsDir bool
	Size  string
}

func New(path string, dirents []os.DirEntry) DirectoryListing {
	listing := DirectoryListing{
		Path:    path,
		Parent:  FindParentPath(path),
		Entries: []DirectoryListingEntry{},
	}

	prefix := path
	if path == "/" {
		prefix = ""
	}

	for _, ent := range dirents {
		finfo, err := ent.Info()
		if err != nil {
			// todo: log or handle
			continue
		}

		listing.Entries = append(listing.Entries, DirectoryListingEntry{
			Name:  ent.Name(),
			Path:  strings.Join([]string{prefix, ent.Name()}, "/"),
			IsDir: ent.IsDir(),
			Size:  Humanize(finfo.Size()),
		})
	}
	return listing
}

// todo: fairly naive implementation here
func FindParentPath(path string) string {
	segments := strings.Split(strings.TrimRight(path, "/"), "/")
	if len(segments) <= 1 {
		return ""
	}
	if len(segments) == 2 {
		return "/"
	}
	return strings.Join(segments[:len(segments)-1], "/")
}

func Humanize(size int64) string {
	if size < 1024 {
		return strconv.FormatInt(size, 10)
	}
	fsize := float64(size)
	if size < 1048576 {
		return fmt.Sprintf("%0.2f kb", fsize/1024.0)
	}
	if size < 1073741824 {
		return fmt.Sprintf("%0.2f mb", fsize/1048576.0)
	}
	return fmt.Sprintf("%0.2f gb", fsize/1073741824.0)
}
