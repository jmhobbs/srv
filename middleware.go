package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func LoggingMiddleware(logger io.Writer, next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
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
	code int
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
	if(s.written > 1073741824) {
		return fmt.Sprintf("%0.2f gb", float64(s.written/1073741824))
	}	
	if(s.written > 1048576) {
		return fmt.Sprintf("%0.2f mb", float64(s.written/1048576))
	}
	if(s.written > 1024) {
		return fmt.Sprintf("%0.2f kb", float64(s.written/1024))
	}
	return strconv.Itoa(s.written) + " b"
}