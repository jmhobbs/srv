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