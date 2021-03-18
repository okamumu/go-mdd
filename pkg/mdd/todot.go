package bdd

import (
	"fmt"
)

// todot

func (f *Node) ToDot(conv func(x *Node) string) {
	fmt.Println("digraph { layout=dot; overlap=false; splines=true; node [fontsize=10];")
	visited := make(map[*Node]struct{})
	f.todot(visited, conv)
	fmt.Println("}")
}

func (f *Node) todot(visited map[*Node]struct{}, conv func(x *Node) string) {
	if _, ok := visited[f]; ok {
		return
	}
	if f.IsTerminal() {
		fmt.Printf(" n%d [shape = square, label = \"%s\"];\n", f.id, conv(f))
	} else {
		fmt.Printf(" n%d [shape = circle, label = \"%s\"];\n", f.id, f.header.label)
		for i, x := range f.nodes {
			x.todot(visited, conv)
			fmt.Printf(" n%d -> n%d [label = \"%d\"];\n", f.id, x.id, i)
		}
	}
	visited[f] = struct{}{}
}
