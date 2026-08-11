package accesslogs_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/jmhobbs/srv/internal/accesslogs"
	"github.com/stretchr/testify/assert"
)

var durationRegex = regexp.MustCompile(`\(\d+(\.\d+)?[µnm]?s\)`)

func Test_LoggingMiddleware(t *testing.T) {
	var buf bytes.Buffer
	mw := accesslogs.LoggingMiddleware(
		&buf,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("x-foo", "bar")
			w.WriteHeader(http.StatusAccepted)
			//nolint:errcheck
			w.Write([]byte{'a', 'b', 'c'})
		}),
	)
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	rr := httptest.NewRecorder()
	mw.ServeHTTP(rr, req)

	assert.Equal(t, "bar", rr.Header().Get("x-foo"))
	assert.True(t, strings.HasSuffix(buf.String(), "\n"), "log output should end with a newline")
	assert.Contains(t, buf.String(), "202 GET http://example.com/foo")
	assert.True(t, durationRegex.MatchString(buf.String()), "log output should contain a duration in parentheses %q", buf.String())
	assert.Contains(t, buf.String(), "3 b")
}
