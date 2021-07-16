package main

import "testing"

func TestFindParentPath(t *testing.T) {
	tests := map[string]string {
		"/": "",
		"/sub/directory": "/sub",
		"/sub": "/",
		"/sub/": "/",
		"/sub/directory/": "/sub",
	}

	for input, expected := range tests {
		actual := findParentPath(input)
		if actual != expected {
			t.Errorf("output did not match expected for %q\n  actual: %q\nexpected: %q", input, actual, expected)
		}
	}
}

func TestDetectContentType(t *testing.T) {
	tests := map[string]*string {
		"vendor.min.css": strPtr("text/css"),
		"path/isnt/important.9284af2.js": strPtr("application/javascript"),
		"you-dont-know-me.txt": nil,
	}

	for input, expected := range tests {
		actual := detectContentType(input)
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

func strPtr(str string) *string {
	return &str
}