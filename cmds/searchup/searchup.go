// Search up the directory tree for occurrences of NAME
// and print them.  Exits with status 1 if no occurrences
// were found.
//
// Use the -first flag to stop the search after the first
// occurrence of NAME.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var first = flag.Bool("first", false, "print first occurrence of NAME and stop")

func usage() {
	fmt.Fprintln(os.Stderr, "usage: searchup [-h] [-first] NAME")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if len(flag.Args()) != 1 {
		usage()
	}
	name := flag.Arg(0)

	path, err := os.Getwd()
	if err != nil {
		errorExit("could not get working directory", err)
	}

	found, err := searchUp(path, name, *first)
	if err != nil {
		errorExit("", err)
	}
	if found == nil {
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, strings.Join(found, "\n"))
}

// searchUp starts at path looking for name as it moves
// up the file tree.  Stops at the first occurrence if
// first is true.
func searchUp(path, name string, first bool) ([]string, error) {
	found := make([]string, 0)
	for {
		ok, err := search(path, name)
		if err != nil {
			return nil, err
		}
		if ok {
			found = append(found, filepath.Join(path, name))
			if first {
				break
			}
		}

		if path == "/" {
			break
		}
		path = filepath.Dir(path)
	}
	return found, nil
}

// search returns true if name matches an entry in path.
func search(path string, name string) (bool, error) {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.Name() == name {
			return true, nil
		}
	}

	return false, nil
}

func errorExit(msg string, err error) {
	switch {
	case msg != "" && err != nil:
		msg = fmt.Sprintf("%s: %v", msg, err)
	case err != nil:
		msg = err.Error()
	}
	fmt.Fprintf(os.Stderr, "error: %s\n", msg)
	os.Exit(1)
}
