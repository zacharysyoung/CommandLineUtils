package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"testing"
)

var (
	x = strings.ReplaceAll

	pre  = func(s string) string { return x(x(x(s, "<CRLF>", "\r\n"), "<LF>", "\n"), "<CR>", "\r") }
	post = func(s string) string { return x(x(x(s, "\r\n", "<CRLF>"), "\n", "<LF>"), "\r", "<CR>") }
)

var testCases = []struct {
	in, want string
}{
	{
		"foo<LF>bar",
		"foo<CRLF>bar",
	},
	{
		"foo<LF>bar<LF>",
		"foo<CRLF>bar<CRLF>",
	},
	{
		"foo<CRLF>bar<LF>",
		"foo<CRLF>bar<CRLF>",
	},
	{
		"foo<CR>bar<LF>",
		"foo<CR>bar<CRLF>",
	},
	{
		"foo<CR>bar<CRLF>",
		"foo<CR>bar<CRLF>",
	},
}

func Test_replace(t *testing.T) {
	for _, runner := range []struct {
		f    func(io.Reader, io.Writer) error
		name string
	}{
		{replaceBytes, "replaceBytes"},
		{replaceStrings, "replaceStrings"},
	} {
		t.Run(runner.name, func(t *testing.T) {
			for _, tc := range testCases {
				buf := &bytes.Buffer{}
				err := runner.f(strings.NewReader(pre(tc.in)), buf)
				if err != nil {
					t.Fatalf("got non-nil err: %v", err)
				}

				if got := post(buf.String()); got != tc.want {
					t.Errorf("\nin   %s\ngot  %s\nwant %s", tc.in, got, tc.want)
				}
			}
		})
	}
}

func Benchmark_performance(b *testing.B) {
	testFile := path.Join(b.TempDir(), "test.txt")
	{
		f, err := os.Create(testFile)
		if err != nil {
			b.Fatal(err)
		}
		w := bufio.NewWriter(f)
		testSize := 100 * 1024 * 1024
		for i := range testSize {
			switch i % 10 {
			case 0:
				w.Write([]byte{'\n'})
			default:
				w.Write([]byte{'a'})
			}
		}
		if err := w.Flush(); err != nil {
			b.Fatal(err)
		}
		f.Close()
	}

	b.Run("replaceStrings", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			f, _ := os.Open(testFile)
			b.StartTimer()
			err := replaceStrings(f, io.Discard)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("replaceBytes", func(b *testing.B) {
		multiplier := 1
		for range 6 {
			size := 1024 * multiplier
			name := fmt.Sprintf("size_%dK", multiplier)
			b.Run(name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					f, _ := os.Open(testFile)
					b.StartTimer()
					err := replaceBytesSize(f, io.Discard, size)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
			multiplier *= 2
		}
	})
}
