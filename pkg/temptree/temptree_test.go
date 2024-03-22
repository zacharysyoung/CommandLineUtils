package temptree

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestRemove(t *testing.T) {
	tree, tempPath, err := NewTree(F("foo"), D("bar", F("baz")))
	if err != nil {
		t.Fatal(err)
	}

	if err := tree.Remove(); err != nil {
		t.Fatal("could not cleanly remove temp tree:", err)
	}

	if _, err := os.Lstat(tempPath); os.IsExist(err) {
		t.Errorf("found temp dir %s; it should have been removed", tempPath)
	}
}

func TestNewTree(t *testing.T) {
	files := []File{
		F("f1"),
		D("d2",
			F("f3"),
			D("d4")), // empty dir
	}
	wants := map[string]bool{
		"f1":    false,
		"d2":    true,
		"d2/f3": false,
		"d2/d4": true,
	}

	tree, tempPath, err := NewTree(files...)
	if err != nil {
		t.Fatal(err)
	}

	// Assert files were added to temp folder
	for wPath, wIsDir := range wants {
		path := filepath.Join(tempPath, wPath)

		got, err := os.Lstat(path)
		if err != nil {
			t.Fatal(err)
		}

		if isDir := got.IsDir(); isDir != wIsDir {
			t.Errorf("%s.IsDir() = %t; want %t", path, isDir, wIsDir)
		}
	}

	// Assert only paths from files were added to temp folder
	filepath.WalkDir(tempPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			t.Fatal(err)
		}
		if path == tempPath {
			return nil
		}

		if path, err = filepath.Rel(tempPath, path); err != nil {
			t.Fatal(err)
		}
		if _, ok := wants[path]; !ok {
			t.Errorf("found path %s in temp folder, but not in wants", path)
		}
		return nil
	})

	if err := tree.Remove(); err != nil {
		t.Fatal(err)
	}
}

func TestString(t *testing.T) {
	for _, tc := range []struct {
		f    File
		want string
	}{
		{
			f: D("root",
				F("f1")),
			want: "root[f1]",
		},
		{
			f: D("root",
				F("f1"),
				D("d2",
					D("d3")),
				F("f4"),
			),
			want: "root[f1 d2[d3[]] f4]",
		},
	} {
		if got := tc.f.String(); got != tc.want {
			t.Errorf("got %s; want %s", got, tc.want)
		}
	}
}

func TestDebug(t *testing.T) {

	for _, tc := range []struct {
		files []File
		want  string
	}{
		{
			files: []File{F("f1"), F("f2"), F("f3")},
			want:  `NewTree(F("f1"), F("f2"), F("f3"))`,
		},

		{
			files: []File{F("f1"), D("d2", F("f3"), D("d4"))},
			want:  `NewTree(F("f1"), D("d2", F("f3"), D("d4")))`,
		},
	} {
		tree, _, _ := NewTree(tc.files...)
		if got := tree.Debug(); got != tc.want {
			t.Errorf("\n%v.Debug()\n  got %s\n want %s", tree.files, got, tc.want)
		}
		tree.Remove()
	}
}
