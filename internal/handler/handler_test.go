package handler_test

import (
	"testing"

	"github.com/jmhobbs/srv/internal/handler"
)

func TestDetectContentType(t *testing.T) {
	tests := map[string]*string{
		"vendor.min.css":                 new("text/css"),
		"path/isnt/important.9284af2.js": new("application/javascript"),
		"you-dont-know-me.txt":           nil,
	}

	for input, expected := range tests {
		actual := handler.DetectContentType(input)
		if actual == nil && expected != nil {
			t.Errorf("output did not match expected for %q\n  actual: %v\nexpected: %q", input, actual, *expected)
			continue
		}
		if actual != nil && expected == nil {
			t.Errorf("output did not match expected for %q\n  actual: %q\nexpected: %v", input, *actual, expected)
			continue
		}
		if actual == nil && expected == nil {
			continue
		}
		if *actual != *expected {
			t.Errorf("output did not match expected for %q\n  actual: %q\nexpected: %q", input, *actual, *expected)
		}
	}
}
