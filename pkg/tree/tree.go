// Package tree provides a simple minilanguage to create temporary
// file trees (usually for testing).
package tree

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

// Node represents a file or directory in a tree.
type Node struct {
	Name     string
	IsDir    bool
	Children []*Node
}

func emptyLine(x string) bool { return strings.TrimSpace(x) == "" }

// Parse parses tree into a Node.  For errors, it returns either
// EmptyTree or NonDirRoot.
func Parse(tree string) (*Node, error) {
	lines := strings.Split(tree, "\n")
	lines = slices.DeleteFunc(lines, emptyLine)
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

func (n *Node) print(w io.Writer, indent string) {
	prefix := "(f)"
	if n.IsDir {
		prefix = "(d)"
	}
	fmt.Fprintf(w, "%s%s %s\n", indent, prefix, n.Name)
	if len(n.Children) == 0 {
		return
	}
	for _, x := range n.Children {
		x.print(w, indent+"  ")
	}
}

func parse(line string) (name string, isDir bool, depth int) {
	i := strings.Index(line, "-")
	if i < 0 {
		i = strings.Index(line, "+")
		isDir = true
	}
	depth = i
	name = line[i+2:]
	return
}
