package tempfs

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		tree string
		node *node
	}{
		{
			tree: `+ foo`,
			node: N("foo", Ns()),
		},
		{
			tree: `
				+ bar
					+ d1
						- f2
					+ d3
						- f4
				`,
			node: N(
				"bar", Ns(
					N("d1", Ns(
						N("f2", nil),
					)),
					N("d3", Ns(
						N("f4", nil),
					)),
				),
			),
		},
		{
			tree: `
				+ baz
					+ d1
						- f2
						- f3
					+ d4
						+ d5
							+ d6
						- f7
					- f8
				`,
			node: N(
				"baz", Ns(
					N("d1", Ns(
						N("f2", nil),
						N("f3", nil),
					)),
					N("d4", Ns(
						N("d5", Ns(
							N("d6", Ns()),
						)),
						N("f7", nil),
					)),
					N("f8", nil),
				),
			),
		},
	}
	for _, tc := range testCases {
		got, err := Parse(tc.tree)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, tc.node) {
			t.Errorf("Parse(%s)\n  got %+v\n want %+v", tc.tree, got, tc.node)
		}
	}
}

func TestPrint(t *testing.T) {
	for _, tc := range []struct {
		n    *node
		want string
	}{
		{
			n: N("root", Ns(
				N("f1", nil),
			)),
			want: "root[f1]",
		},
		{
			n: N("root", Ns(
				N("f1", nil),
				N("d2", Ns(
					N("d3", Ns()),
				)),
				N("f4", nil),
			)),
			want: "root[f1 d2[d3[]] f4]",
		},
	} {
		if got := tc.n.print(); got != tc.want {
			t.Errorf("%v.print() = %s; want %s", tc.n, got, tc.want)
		}
	}
}
