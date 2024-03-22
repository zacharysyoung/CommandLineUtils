package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"zacharysyoung/CLUtils/pkg/temptree"
)

// TestGetDepthRelative tests the depth of the relative paths
// (not rooted at "/") by comparing the depth of the test paths
// to the depth of the cwd.
func TestGetDepthRelative(t *testing.T) {
	cwdDepth := getDepth(".")

	for _, tc := range []struct {
		path string
		want int
	}{
		{".", 0},
		{"foo", 1},
		{"./foo", 1},
		{"foo/bar", 2},
		{"./foo/bar", 2},
	} {
		if got := getDepth(tc.path); got-cwdDepth != tc.want {
			t.Errorf("getDepth(%s) = %d; want %d", tc.path, got, tc.want)
		}
	}
}
func TestGetDepthAbs(t *testing.T) {
	for _, tc := range []struct {
		path string
		want int
	}{
		{"/", 0},
		{"/foo", 1},
		{"/foo/bar", 2},
	} {
		if got := getDepth(tc.path); got != tc.want {
			t.Errorf("getDepth(%s) = %d; want %d", tc.path, got, tc.want)
		}
	}
}

var (
	D = temptree.D
	F = temptree.F
)

func popFirstLine(s string) (first, rem string) {
	lines := strings.Split(s, "\n")
	first = lines[0]
	rem = strings.Join(lines[1:], "\n")
	return
}

func TestPrint(t *testing.T) {
	testCases := []struct {
		files []temptree.File
		want  string
	}{
		{
			files: []temptree.File{F("foo"), F("bar"), F("baz")},
			want: `+ Some temp dir
  - bar
  - baz
  - foo
`,
		},
		{
			files: []temptree.File{D("foo", F("bar"), F("baz"))},
			want: `+ Some temp dir
  + foo
    - bar
    - baz
`,
		},
	}

	buf := &bytes.Buffer{}
	for _, tc := range testCases {
		tree, tempPath, err := temptree.NewTree(tc.files...)
		if err != nil {
			t.Fatal(err)
		}

		buf.Reset()
		printTree(filepath.Join(tempPath, "."), buf)
		got := buf.String()

		got1, gotRem := popFirstLine(got)

		if got1[:2] != "+ " {
			t.Errorf("tree starts with %q; want \"+ ...\"", got1)
		}

		_, want := popFirstLine(tc.want)
		if gotRem != want {
			t.Errorf("printTree()\n  got %s\n want %s", got, tc.want)
		}

		if err = tree.Remove(); err != nil {
			t.Fatal(err)
		}
	}
}
