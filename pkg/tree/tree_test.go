package tree

import (
	"reflect"
	"testing"
)

type c []*Node // shorthand for children

func TestLoop(t *testing.T) {
	testCases := []struct {
		tree string
		node *Node
	}{
		{tree: `+ foo`,
			node: &Node{
				"foo", true, nil,
			},
		},
		{tree: `
				+ bar
					+ d1
						- f2
					+ d3
						- f4
				`,
			node: &Node{
				"bar", true, c{
					{"d1", true, c{
						{"f2", false, nil},
					}},
					{"d3", true, c{
						{"f4", false, nil},
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
							+ d6
						- f7
					- f8
				`,
			node: &Node{
				"baz", true, c{
					{"d1", true, c{
						{"f2", false, nil},
						{"f3", false, nil},
					}},
					{"d4", true, c{
						{"d5", true, c{
							{"d6", true, nil},
						}},
						{"f7", false, nil},
					}},
					{"f8", false, nil},
				},
			},
		},
	}
	for _, tc := range testCases {
		got, err := Parse(tc.tree)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, tc.node) {
			t.Errorf("Parse(%s)\n  got %s\n want %s", tc.tree, got, tc.node)
		}
	}
}

func TestString(t *testing.T) {
	for _, tc := range []struct {
		n    *Node
		want string
	}{
		{
			n: &Node{"root", true, c{
				{"f1", false, nil},
			}},
			want: "root[f1]",
		},
		{
			n: &Node{"root", true, c{
				{"f1", false, nil},
				{"d2", true, c{
					{"d3", true, nil},
				}},
				{"f4", false, nil},
			}},
			want: "root[f1 d2[d3[]] f4]",
		},
	} {
		if got := tc.n.String(); got != tc.want {
			t.Errorf("%v.String() = %s; want %s", tc.n, got, tc.want)
		}
	}
}
