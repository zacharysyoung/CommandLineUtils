package main

// List path, and optionally files, or pattern-matched files, in
// each path component.

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var files = flag.Bool("f", false, "list files inside dirs")
var fPattern = flag.String("re", "", "regex[] pattern of files to match; turns on -f")

func usage() {
	fmt.Fprintln(os.Stderr, `usage: lspath [-f | -re]

Parses $PATH env var and prints the directories, optionally printing files in those directories.`)
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	var (
		reFpat *regexp.Regexp
		err    error
	)
	if *fPattern != "" {
		*files = true
		reFpat, err = regexp.Compile(*fPattern)
		if err != nil {
			fmt.Fprintln(os.Stderr, "couldn't compile -fpat, bad regexp")
			os.Exit(1)
		}
	}

	dirs := strings.Split(os.Getenv("PATH"), ":")

	toPrint := make([]string, 0)
	for _, d := range dirs {
		if ignore(d) {
			continue
		}

		toPrint = toPrint[0:0]

		if *files {
			dirEntries, err := os.ReadDir(d)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error: ", err)
			}

			for _, f := range dirEntries {
				if reFpat != nil &&
					!reFpat.Match([]byte(f.Name())) {
					continue
				}
				toPrint = append(toPrint, f.Name())
			}
		}

		switch {
		case *files:
			if len(toPrint) > 0 {
				fmt.Println(d)
				for _, x := range toPrint {
					fmt.Println("  ", x)
				}
			}
		default:
			fmt.Println(d)
		}
	}
}

// ignore ignores dirs I don't actually care about
func ignore(d string) bool {
	switch {
	default:
		return false

	// https://apple.stackexchange.com/q/458277/189634
	case len(d) >= 18 && d[:18] == "/var/run/com.apple":
		return true
	}
}
