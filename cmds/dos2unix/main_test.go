package main

import (
	"bytes"
	"strings"
	"testing"
)

func Test_run(t *testing.T) {
	var (
		x = strings.ReplaceAll

		pre  = func(s string) string { return x(x(s, "<CRLF>", "\r\n"), "<LF>", "\n") }
		post = func(s string) string { return x(x(s, "\r\n", "<CRLF>"), "\n", "<LF>") }
	)
	for _, tc := range []struct {
		in, want string
	}{
		{
			"foo<CRLF>bar",
			"foo<LF>bar",
		},
		{
			"foo<LF>bar<CRLF>",
			"foo<LF>bar<LF>",
		},
	} {
		buf := &bytes.Buffer{}
		err := run(strings.NewReader(pre(tc.in)), buf)
		if err != nil {
			t.Fatalf("got non-nil err: %v", err)
		}

		if got := post(buf.String()); got != tc.want {
			t.Errorf("\nin   %s\ngot  %s\nwant %s", tc.in, got, tc.want)
		}

	}
}
