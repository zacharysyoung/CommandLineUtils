package tree

import (
	"bytes"
	"reflect"
	"testing"
)

func TestLoop(t *testing.T) {

	testCases := []struct {
		tree string
		node *Node
	}{
		{tree: `+ foo`,
			node: &Node{Name: "foo", IsDir: true},
		},
		{tree: `
				+ bar
					+ d1
						- f2
					+ d3
						- f4
				`,
			node: &Node{
				Name: "bar", IsDir: true, Children: []*Node{
					{Name: "d1", IsDir: true, Children: []*Node{
						{Name: "f2"},
					}},
					{Name: "d3", IsDir: true, Children: []*Node{
						{Name: "f4"},
					}},
				},
			},
		},
		{tree: `
				+ baz
					+ d1
						- f2
						- f3
					+ d4
						+ d5
							- f6
					- f7
				`,
			node: &Node{
				Name: "baz", IsDir: true, Children: []*Node{
					{Name: "d1", IsDir: true, Children: []*Node{
						{Name: "f2"},
						{Name: "f3"},
					}},
					{Name: "d4", IsDir: true, Children: []*Node{
						{Name: "d5", IsDir: true, Children: []*Node{
							{Name: "f6"},
						}},
					}},
					{Name: "f7"},
				},
			},
		},
	}
	buf := &bytes.Buffer{}
	for _, tc := range testCases {
		got, err := Parse(tc.tree)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, tc.node) {
			buf.Reset()
			got.print(buf, "@")
			gotS := buf.String()
			buf.Reset()
			tc.node.print(buf, "@")
			wantS := buf.String()
			t.Errorf("Parse(%s)\n  got %v\n want %v", tc.tree, gotS, wantS)
		}
	}
}

func TestPrint(t *testing.T) {
	buf := &bytes.Buffer{}
	node := &Node{Name: "root", IsDir: true, Children: []*Node{
		{Name: "f1"},
		{Name: "d2", IsDir: true, Children: []*Node{
			{Name: "f3"},
		}},
	}}
	want := `(d) root
  (f) f1
  (d) d2
    (f) f3
`
	node.print(buf, "")
	got := buf.String()
	if got != want {
		t.Errorf("\n  got %q;\n want %q", got, want)
	}
}
