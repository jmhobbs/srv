package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Handler struct {
  root string
	logger *Logger
}

func newHandler(logger *Logger, dir string) *Handler {
	return &Handler{dir, logger}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clean := filepath.Clean(r.URL.Path)

	if clean != r.URL.Path {
		http.Redirect(w, r, clean, 301)
		return
	}

	abs := filepath.Join(h.root, clean)
	finfo, err := os.Stat(abs) // todo: stat cache?
	if err != nil {
		if os.IsNotExist(err) {
			// 404
			w.WriteHeader(404)
			w.Write([]byte("Not Found"))
			return
		}
	}

	if finfo.IsDir() {
		dirents, err := os.ReadDir(abs)
		if err != nil {
			h.logger.Error("error reading directory %q: %v", clean, err)
			w.WriteHeader(500)
			w.Write([]byte("Internal Server Error"))
			return
		}

		err = Template.ExecuteTemplate(w, "index", newDirectoryListing(clean, dirents))
		if err != nil {
			h.logger.Error("error rendering directory listing %q: %v", clean, err)
			w.WriteHeader(500)
			w.Write([]byte("Internal Server Error"))
			return
		}
		return
	}

	f, err := os.Open(abs)
	if err != nil {
		h.logger.Error("error opening file %q: %v", clean, err)
		w.WriteHeader(500)
		w.Write([]byte("Internal Server Error"))
		return
	}
	defer f.Close()

	// todo: etags
	_, err = io.Copy(w, f)
	if err != nil {
		h.logger.Error("error copyinf file %q: %v", clean, err)
		w.WriteHeader(500)
		w.Write([]byte("Internal Server Error"))
		return
	}
}

type directoryListing struct {
	Path string
	Parent string
	Entries []directoryListingEntry
}

type directoryListingEntry struct {
	Name string
	Path string
	IsDir bool
} 

func newDirectoryListing(path string, dirents []os.DirEntry) directoryListing {
	listing := directoryListing{
		Path: path,
		Parent: findParentPath(path),
		Entries: []directoryListingEntry{},
	}

	prefix := path
	if path == "/" {
		prefix = ""
	}

	for _, ent := range dirents {
		listing.Entries = append(listing.Entries, directoryListingEntry{
			Name: ent.Name(),
			Path: strings.Join([]string{prefix, ent.Name()}, "/"),
			IsDir: ent.IsDir(),
		})
	}
	return listing
}

// todo: fairly naive implementation here
func findParentPath(path string) string {
	segments := strings.Split(strings.TrimRight(path, "/"), "/")
	if len(segments) <= 1 {
		return ""
	}
	if len(segments) == 2 {
		return "/"
	}
	return strings.Join(segments[:len(segments)-1], "/")
}