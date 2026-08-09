package accesslogs

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func Setup(handler http.Handler, accessLogs string) (http.Handler, error) {
	if accessLogs != "" {
		var sink io.Writer
		if accessLogs == "-" {
			sink = os.Stdout
		} else {
			f, err := os.OpenFile(accessLogs, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return nil, fmt.Errorf("error opening access log: %v", err)
			}
			//nolint:errcheck
			defer f.Close()
			sink = f
		}
		return LoggingMiddleware(sink, handler), nil
	}
	return handler, nil
}
