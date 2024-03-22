package main

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"zacharysyoung/CLUtils/pkg/temptree"
)

var d, f = temptree.D, temptree.F

func TestSearch(t *testing.T) {
	tree, prefix := newTree(
		d("foo",
			f("a"),
			d("bar",
				f("a"),
				f("b"),
				d("baz",
					f("a"),
					f("c"),
				),
			),
		), t)

	for _, tc := range []struct {
		path, name string
		want       bool
	}{
		{"foo", "a", true},
		{"foo", "x", false},
		{"foo/bar", "a", true},
		{"foo/bar", "b", true},
		{"foo/bar", "x", false},
		{"foo/bar/baz", "a", true},
		{"foo/bar/baz", "c", true},
		{"foo/bar/baz", "x", false},
	} {
		got, err := search(join(prefix, tc.path), tc.name)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.want {
			t.Errorf("search(%s, %s) = %t; want %t", tc.path, tc.name, got, tc.want)
		}
	}

	_, err := search(join(prefix, "fooz"), "a")
	if err == nil {
		t.Errorf("search with bad path didn't error")
	}

	removeTree(tree, t)
}

func TestSearchUp(t *testing.T) {
	type testCase struct {
		start, name string // search up from start, looking for name
		first       bool   // stop after first occurence
		want        []string
	}

	tree, prefix := newTree(
		d("foo",
			f("a"),
			d("bar",
				f("a"),
				f("b"),
				d("baz",
					f("a"),
					f("c")))), t)

	testCases := []testCase{
		{"foo", "a", false, []string{
			"/foo/a"}},
		{"foo", "x", false, []string{}},
		{"foo/bar", "a", false, []string{
			"/foo/bar/a",
			"/foo/a"}},
		{"foo/bar", "a", true, []string{
			"/foo/bar/a"}},
		{"foo/bar/baz", "a", false, []string{
			"/foo/bar/baz/a",
			"/foo/bar/a",
			"/foo/a"}},
		{"foo/bar/baz", "a", true, []string{
			"/foo/bar/baz/a"}},
	}

	for _, tc := range testCases {
		got, err := searchUp(join(prefix, tc.start), tc.name, tc.first)
		if err != nil {
			t.Fatal(err)
		}
		got = trim(got, prefix)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("searchUp(%s, %s, %t)\n  got %v\n want %v", tc.start, tc.name, tc.first, got, tc.want)
		}
	}

	_, err := search(join(prefix, "fooz"), "a")
	if err == nil {
		t.Errorf("search with bad path didn't error")
	}

	removeTree(tree, t)
}

func join(prefix, path string) string {
	return filepath.Join(prefix, path)
}

func trim(s []string, prefix string) []string {
	for i, x := range s {
		s[i] = strings.TrimPrefix(x, prefix)
	}
	return s
}

func newTree(file temptree.File, t *testing.T) (tree *temptree.Tree, prefix string) {
	tree, tempPath, err := temptree.NewTree(file)
	if err != nil {
		t.Fatal(err)
	}
	return tree, tempPath
}

func removeTree(tree *temptree.Tree, t *testing.T) {
	if err := tree.Remove(); err != nil {
		t.Fatal(err)
	}
}
