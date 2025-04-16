package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"
)

func usage() {
	fmt.Fprintln(os.Stderr, `usage: unix2dos [-v] [file]

Reads file (or stdin) and converts all line feeds (LF) to carriage return line feeds.`)
	flag.PrintDefaults()
	os.Exit(2)
}

var versionFlag = flag.Bool("v", false, "print version/build info")

func main() {
	flag.Usage = usage
	flag.Parse()

	if *versionFlag {
		fmt.Fprintln(os.Stderr, version())
		os.Exit(2)
	}

	var (
		in  io.Reader
		err error
	)
	tail := flag.Args()
	switch len(tail) {
	case 0:
		in = os.Stdin
	case 1:
		in, err = os.Open(tail[0])
		if err != nil {
			errorOut(fmt.Sprintf("could not read from specified file: %v", err))
		}
		defer in.(*os.File).Close()
	default:
		badArgs(fmt.Sprintf("got %d files: %s; can only read from Stdin or a single file", len(tail), strings.Join(tail, ", ")))
	}

	if err = run(in, os.Stdout); err != nil {
		errorOut(err.Error())
	}
}

func run(in io.Reader, out io.Writer) error {
	b, err := io.ReadAll(in)
	if err != nil {
		return err
	}
	s := string(b)

	const (
		crlf = "\r\n"
		lf   = "\n"
	)
	crlfSegments := strings.Split(s, crlf)
	for i := range crlfSegments {
		crlfSegments[i] = strings.ReplaceAll(crlfSegments[i], lf, crlf)
	}

	_, err = io.WriteString(out, strings.Join(crlfSegments, crlf))
	if err != nil {
		return err
	}

	return nil
}

func version() string {
	s := "unix2dos"
	if bi, ok := debug.ReadBuildInfo(); ok {
		for _, x := range bi.Settings {
			if x.Key == "vcs.revision" {
				s += ":" + x.Value[:7] // short hash
				break
			}
		}
		s += ":" + bi.GoVersion
	}
	return s
}

func badArgs(s string) {
	fmt.Fprintf(os.Stderr, "error: %s\n", s)
	os.Exit(2)
}

func errorOut(s string) {
	fmt.Fprintf(os.Stderr, "error: %s\n", s)
	os.Exit(1)
}
