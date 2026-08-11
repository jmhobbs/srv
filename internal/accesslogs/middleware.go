package accesslogs

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jmhobbs/srv/internal/humanize"
)

func LoggingMiddleware(logger io.Writer, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		w2 := &loggingResponseWriter{w, 0, 200}
		next.ServeHTTP(w2, r)
		elapsed := time.Since(start)
		fmt.Fprintf(logger, "%d %v %v (%v) %s\n", w2.code, r.Method, r.URL, elapsed, w2.Written())
	})
}

type loggingResponseWriter struct {
	wrapped http.ResponseWriter
	written int
	code    int
}

func (s *loggingResponseWriter) Header() http.Header {
	return s.wrapped.Header()
}

func (s *loggingResponseWriter) Write(buf []byte) (int, error) {
	written, err := s.wrapped.Write(buf)
	s.written += written
	return written, err
}

func (s *loggingResponseWriter) WriteHeader(code int) {
	s.code = code
	s.wrapped.WriteHeader(code)
}

func (s *loggingResponseWriter) Written() string {
	return humanize.Bytes(int64(s.written))
}
