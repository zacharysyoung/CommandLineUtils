// Package temptree provides a simple minilanguage for creating
// temporary file trees.
//
// To create a tree that contains the files f1 and f6, the empty dir d2,
// and the dir d3, with dir d4 below that, and the empty dir d5 below that:
//
//	/f1
//	/d2/
//	/d3/d4/d5/
//	/f6
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

func NewTree(files ...File) *Tree {
	return &Tree{files: files}
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
func (t *Tree) MakeTemp() (string, error) {
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

// print a compact, diagnostic string for n.
func (f *File) print() string {
	if f.children == nil {
		return f.name
	}

	children, sep := "", ""
	for _, cf := range f.children {
		children += sep + cf.print()
		sep = " "
	}

	return fmt.Sprintf("%s[%s]", f.name, children)
}
