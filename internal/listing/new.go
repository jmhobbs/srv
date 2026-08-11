package listing

import (
	"os"
	"strings"

	"github.com/jmhobbs/srv/internal/humanize"
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
			Size:  humanize.Bytes(finfo.Size()),
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
