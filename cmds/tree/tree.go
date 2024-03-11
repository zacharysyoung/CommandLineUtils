package main

// Print directory tree from first arg.
//
// https://gist.github.com/zacharysyoung/64b6593f7d0314d0eb29bbc9ef121f1e

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	prefix     = flag.String("prefix", "", "what each line should begin with")
	indent     = flag.String("indent", "  ", "what each level of identation should add before the file name")
	dirPrefix  = flag.String("dirprefix", "+ ", "what directly precedes a directory object")
	filePrefix = flag.String("fileprefix", "- ", "what directly precedes a file object")
)

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Fprintln(os.Stderr, "usage: tree [-h] [-prefix] [-indent] PATH")
		os.Exit(1)
	}

	if err := printTree(flag.Arg(0), os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
	}
}

func printTree(root string, w io.Writer) error {
	root = filepath.Clean(root)
	rootDepth := getDepth(root)

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		objPrefix := filePrefix
		if info.IsDir() {
			objPrefix = dirPrefix
		}

		n := getDepth(path)
		nIndent := strings.Repeat(*indent, n-rootDepth)

		fname := filepath.Base(path)
		fmt.Fprintf(w, "%s%s%s%s\n", *prefix, nIndent, *objPrefix, fname)

		return nil
	})
}

func getDepth(path string) (depth int) {
	var err error
	if path, err = filepath.Abs(path); err != nil {
		panic(err)
	}
	for path != "/" {
		depth++
		path = filepath.Dir(path)
	}
	return
}
