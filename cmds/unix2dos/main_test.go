package main

import (
	"bytes"
	"strings"
	"testing"
)

func Test_run(t *testing.T) {
	for _, tc := range []struct {
		in, want string
	}{
		{
			"foo\nbar",
			"foo\r\nbar",
		},
		{
			"foo\nbar\n",
			"foo\r\nbar\r\n",
		},
		{
			"foo\r\nbar\n",
			"foo\r\nbar\r\n",
		},
	} {
		buf := &bytes.Buffer{}
		err := run(strings.NewReader(tc.in), buf)
		if err != nil {
			t.Fatalf("got non-nil err: %v", err)
		}

		if got := buf.String(); got != tc.want {
			t.Errorf("\nin   %q\ngot  %q\nwant %q", tc.in, got, tc.want)
		}

	}

}
