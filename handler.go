package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Handler struct {
	root                  string
	logger                *Logger
	defaultDirectoryFiles []string
}

func newHandler(logger *Logger, dir string, defaultDirectoryFiles []string) *Handler {
	return &Handler{dir, logger, defaultDirectoryFiles}
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
		for _, defaultFile := range h.defaultDirectoryFiles {
			f, err := os.Open(filepath.Join(abs, defaultFile))
			if err == nil {
				defer f.Close()
				h.writeFile(w, f, abs, clean)
				return
			}
		}

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

	h.writeFile(w, f, abs, clean)
}

func (h *Handler) writeFile(w http.ResponseWriter, f io.Reader, abs, clean string) {
	if contentType := detectContentType(abs); contentType != nil {
		w.Header().Set("content-type", *contentType)
	}

	_, err := io.Copy(w, f)
	if err != nil {
		h.logger.Error("error copying file %q: %v", clean, err)
		w.WriteHeader(500)
		w.Write([]byte("Internal Server Error"))
		return
	}
}

type directoryListing struct {
	Path    string
	Parent  string
	Entries []directoryListingEntry
}

type directoryListingEntry struct {
	Name  string
	Path  string
	IsDir bool
	Size  string
}

func newDirectoryListing(path string, dirents []os.DirEntry) directoryListing {
	listing := directoryListing{
		Path:    path,
		Parent:  findParentPath(path),
		Entries: []directoryListingEntry{},
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

		listing.Entries = append(listing.Entries, directoryListingEntry{
			Name:  ent.Name(),
			Path:  strings.Join([]string{prefix, ent.Name()}, "/"),
			IsDir: ent.IsDir(),
			Size:  humanize(finfo.Size()),
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

func humanize(size int64) string {
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

func detectContentType(path string) *string {
	i := strings.LastIndexByte(path, '.')
	if i == -1 {
		return nil
	}

	suffix := path[i+1:]

	if contentType, ok := knownContentTypes[suffix]; ok {
		return &contentType
	}

	return nil
}

var knownContentTypes map[string]string

func init() {
	knownContentTypes = map[string]string{
		"css": "text/css",
		"js":  "application/javascript",
		"svg": "image/svg+xml",
	}
}
