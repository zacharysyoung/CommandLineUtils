// Package temptree creates temporary, ad-hoc file trees (for
// testing) and can dispose of them when done.
//
// Use the F() and D() funcs to create files and directories.
// Both take a name argument. D can also take any number of
// F, or none for an empty directory.
//
// To create the tree:
//
//	tempPath/root
//	tempPath/root/foo
//	tempPath/root/bar
//	tempPath/root/baz
//
// run:
//
//	t := NewTree(D("root", F("foo"), F("bar"), D("baz")))
//	p := t.MakeTemp()
//
// p=tempPath, foo and bar are files, and baz is an empty directory.
//
// Call t.Remove() to remove the tree on disk.
package temptree

import (
	"fmt"
	"os"
	"path/filepath"
)

type Tree struct {
	files []File

	tempPath string
}

func NewTree(files ...File) (tree *Tree, tempPath string, err error) {
	tree = &Tree{files: files}
	tempPath, err = tree.makeTemp()
	return
}

func (t *Tree) Debug() string {
	var print func(File) string
	print = func(f File) string {
		if f.children == nil {
			return `F("` + f.name + `")`
		}

		s := ""
		for _, cf := range f.children {
			s += ", " + print(cf)
		}

		return `D("` + f.name + `"` + s + `)`
	}

	s, sep := "", ""
	for _, f := range t.files {
		s += sep + print(f)
		sep = ", "
	}

	return fmt.Sprintf("NewTree(%s)", s)
}

// Make creates a temp directory then recursively adds the files
// in t underneath it.
func (t *Tree) makeTemp() (string, error) {
	tempPath, err := os.MkdirTemp("", "")
	if err != nil {
		return "", err
	}
	t.tempPath = tempPath

	var walk func(File, string) error
	walk = func(n File, path string) error {
		path = filepath.Join(path, n.name)

		if n.children != nil {
			if err := os.Mkdir(path, 0700); err != nil {
				return err
			}
		} else {
			if err := touch(path); err != nil {
				return err
			}
		}

		for _, nc := range n.children {
			if err := walk(nc, path); err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range t.files {
		if err := walk(f, tempPath); err != nil {
			return tempPath, err
		}
	}

	return tempPath, nil
}

// touch creates an empty file.
func touch(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	return f.Close()
}

func (t *Tree) Remove() error { return os.RemoveAll(t.tempPath) }

// File represents a file or directory in the tree.
type File struct {
	name     string
	children []File
}

// F wraps name in a File
func F(name string) File {
	return File{name, nil}
}

// D creates a directory File with name, and children files.  If files
// is nil, the directory is empty.
func D(name string, files ...File) File {
	if files == nil {
		files = []File{}
	}
	return File{name, files}
}

func (f File) String() string {
	return f.print()
}

// print a compact, diagnostic string for f.
func (f File) print() string {
	if f.children == nil {
		return f.name
	}

	s, sep := "", ""
	for _, x := range f.children {
		s += sep + x.print()
		sep = " "
	}

	return fmt.Sprintf("%s[%s]", f.name, s)
}
