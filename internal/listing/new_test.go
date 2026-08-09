package listing_test

import (
	"testing"

	"github.com/jmhobbs/srv/internal/listing"
)

func TestFindParentPath(t *testing.T) {
	tests := map[string]string{
		"/":               "",
		"/sub/directory":  "/sub",
		"/sub":            "/",
		"/sub/":           "/",
		"/sub/directory/": "/sub",
	}

	for input, expected := range tests {
		actual := listing.FindParentPath(input)
		if actual != expected {
			t.Errorf("output did not match expected for %q\n  actual: %q\nexpected: %q", input, actual, expected)
		}
	}
}
