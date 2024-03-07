// Package tempfs provides a simple minilanguage for creating
// temporary file trees.
//
// To create a tree that contains the files f1 and f6, the empty dir d2,
// and the dir d3, with dir d4 below that, and the empty dir d5 below that:
//
//	/f1
//	/d2/
//	/d3/d4/d5/
//	/f6
package tempfs

import (
	"fmt"
	"slices"
	"strings"
)

// node represents a file or directory in the tree.
type node struct {
	Name     string
	IsDir    bool
	Children []*node
}

func N(name string, children []*node) *node {
	return &node{
		Name:     name,
		IsDir:    children != nil,
		Children: children,
	}
}
func Ns(nodes ...*node) []*node {
	if nodes == nil {
		return []*node{}
	}
	return nodes
}

// Parse parses tree into a Node.
func Parse(tree string) (*node, error) {
	lines := strings.Split(tree, "\n")
	lines = slices.DeleteFunc(lines, isEmpty)
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty tree")
	}
	name, dir, depth := parse(lines[0])
	if !dir {
		return nil, fmt.Errorf("root %q must be directory", lines[0])
	}

	root := N(name, Ns())
	m := map[int]*node{
		depth: root,
	}

	var node *node
	for i := 1; i < len(lines); i++ {
		name, dir, n := parse(lines[i])
		if dir {
			node = N(name, Ns())
			m[n] = node
		} else {
			node = N(name, nil)
		}
		parent := m[n-1]
		parent.Children = append(parent.Children, node)
	}

	return root, nil
}

func isEmpty(x string) bool { return strings.TrimSpace(x) == "" }

// parse parses a single line returning its name, directory status
// and depth in the tree (by counting indentation).
func parse(line string) (name string, isDir bool, depth int) {
	i := strings.Index(line, "-")
	if i < 0 {
		i = strings.Index(line, "+")
		isDir = true
	}
	depth = i
	name = strings.TrimSpace(line[i+2:])
	return
}

// print a compact, diagnostic string for n.
func (n *node) print() string {
	if !n.IsDir {
		return n.Name
	}

	children, sep := "", ""
	for _, c := range n.Children {
		children += sep + c.print()
		sep = " "
	}

	return fmt.Sprintf("%s[%s]", n.Name, children)
}
