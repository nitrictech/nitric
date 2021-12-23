package utils

import "strings"

// slashSplitter - used to split strings, with the same output regardless of leading or trailing slashes
// e.g - strings.FieldsFunc("/one/two/three/", f) == strings.FieldsFunc("/one/two/three", f) == strings.FieldsFunc("one/two/three", f) == ["one" "two" "three"]
func slashSplitter(c rune) bool {
	return c == '/'
}

// SplitPath - splits a path into its component parts, ignoring leading or trailing slashes.
// e.g - SplitPath("/one/two/three/") == SplitPath("/one/two/three") == SplitPath("one/two/three") == ["one" "two" "three"]
func SplitPath(p string) []string {
	return strings.FieldsFunc(p, slashSplitter)
}
