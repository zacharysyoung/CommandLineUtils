// Package tree provides a simple minilanguage to create temporary
// file trees (usually for testing).
package tree

import (
	"fmt"
	"slices"
	"strings"
)

// Node represents a file or directory in a tree.
type Node struct {
	Name     string
	IsDir    bool
	Children []*Node
}

func isEmpty(x string) bool { return strings.TrimSpace(x) == "" }

// Parse parses tree into a Node.
func Parse(tree string) (*Node, error) {
	lines := strings.Split(tree, "\n")
	lines = slices.DeleteFunc(lines, isEmpty)
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty tree")
	}
	name, isDir, depth := parse(lines[0])
	if !isDir {
		return nil, fmt.Errorf("root %q must be directory", lines[0])
	}

	root := &Node{Name: name, IsDir: isDir}

	var m = map[int]*Node{
		depth: root,
	}
	for i := 1; i < len(lines); i++ {
		x, dir, n := parse(lines[i])
		node := &Node{Name: x, IsDir: dir}
		if dir {
			m[n] = node
		}
		parent := m[n-1]
		parent.Children = append(parent.Children, node)
	}

	return root, nil
}

// parse parses a single line
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

func (n *Node) String() string {
	if !n.IsDir {
		return n.Name
	}

	children, sep := "", ""
	for _, c := range n.Children {
		children += sep + c.String()
		sep = " "
	}

	return fmt.Sprintf("%s[%s]", n.Name, children)
}
