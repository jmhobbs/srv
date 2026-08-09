package handler

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	headers "github.com/jmhobbs/cloudflare-headers-file"
	"github.com/jmhobbs/srv/internal/listing"
	"github.com/rs/zerolog"
)

type Handler struct {
	root                  string
	logger                zerolog.Logger
	defaultDirectoryFiles []string
	fileServer            http.Handler
	headers               *headers.File
}

func New(logger zerolog.Logger, dir string, defaultDirectoryFiles []string, headers *headers.File) *Handler {
	return &Handler{dir, logger, defaultDirectoryFiles, http.FileServer(http.Dir(dir)), headers}
}

func (h *Handler) SetHeaders(w http.ResponseWriter, r *http.Request) {
	if h.headers != nil {
		headersToSet := h.headers.Match(*r.URL)
		for _, header := range headersToSet {
			split := strings.SplitN(header, ":", 2)
			if len(split) == 2 {
				w.Header().Add(split[0], split[1])
			}
		}
	}
}

func (h *Handler) DirectoryListing(w http.ResponseWriter, clean, abs string) {
	for _, defaultFile := range h.defaultDirectoryFiles {
		f, err := os.Open(filepath.Join(abs, defaultFile))
		if err == nil {
			defer func() {
				if err := f.Close(); err != nil {
					h.logger.Error().Err(err).Str("path", defaultFile).Msg("error closing file")
				}
			}()
			h.writeFile(w, f, abs, clean)
			return
		}
	}

	dirents, err := os.ReadDir(abs)
	if err != nil {
		h.logger.Error().Err(err).Str("path", clean).Msg("error reading directory")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = Template.ExecuteTemplate(w, "index", listing.New(clean, dirents))
	if err != nil {
		h.logger.Error().Err(err).Str("path", clean).Msg("error rendering directory listing")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clean := filepath.Clean(r.URL.Path)

	if clean != r.URL.Path {
		http.Redirect(w, r, clean, http.StatusMovedPermanently)
		return
	}

	h.SetHeaders(w, r)

	abs := filepath.Join(h.root, clean)
	finfo, err := os.Stat(abs)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		h.logger.Error().Err(err).Str("path", abs).Msg("failed to stat file")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if finfo.IsDir() {
		h.DirectoryListing(w, clean, abs)
		return
	}

	if contentType := DetectContentType(abs); contentType != nil {
		w.Header().Set("content-type", *contentType)
	}

	h.fileServer.ServeHTTP(w, r)
}

func (h *Handler) writeFile(w http.ResponseWriter, f io.Reader, abs, clean string) {
	if contentType := DetectContentType(abs); contentType != nil {
		w.Header().Set("content-type", *contentType)
	}

	_, err := io.Copy(w, f)
	if err != nil {
		h.logger.Error().Err(err).Str("path", clean).Msg("error writing file")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func DetectContentType(path string) *string {
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

var knownContentTypes = map[string]string{
	"css": "text/css",
	"js":  "application/javascript",
	"svg": "image/svg+xml",
}
