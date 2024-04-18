// Copyright 2024 Zachary S Young.  All rights reserved.
// Use of this source code is governed by an MIT
// license that can be found in the LICENSE file.

// Tree prints the directory tree starting at PATH.
//
// Usage:
//
// tree [-h] [options] PATH
package main

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
	indent     = flag.String("indent", "  ", "what each level of identation should add before the object")
	dirPrefix  = flag.String("dirprefix", "+ ", "what directly precedes a directory-like object")
	filePrefix = flag.String("fileprefix", "- ", "what directly precedes a file-like object")
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: tree [-h] [options] PATH")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if len(flag.Args()) != 1 {
		usage()
	}

	if err := printTree(flag.Arg(0), os.Stdout); err != nil {
		errorExit(err.Error())
	}
}

func printTree(root string, w io.Writer) error {
	root = filepath.Clean(root)
	nRoot, err := getDepth(root)
	if err != nil {
		return err
	}

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		n, err := getDepth(path)
		if err != nil {
			return err
		}

		s := strings.Repeat(*indent, n-nRoot)

		op := filePrefix
		if info.IsDir() {
			op = dirPrefix
		}

		name := filepath.Base(path)

		fmt.Fprintf(w, "%s%s%s%s\n", *prefix, s, *op, name)

		return nil
	})
}

func getDepth(path string) (depth int, err error) {
	path, err = filepath.Abs(path)
	if err != nil {
		return
	}
	for path != "/" {
		depth++
		path = filepath.Dir(path)
	}
	return
}

func errorExit(msg string) {
	fmt.Fprintf(os.Stderr, "error: %s\n", msg)
	os.Exit(1)
}
