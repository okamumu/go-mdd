package bdd

import (
	"fmt"
	"testing"
)

type BDDNode struct {
	*Node
	edge bool
}

func TestNode1(t *testing.T) {
	m := NewNodeManager()
	m.NewHeader("x")
	m.NewHeader("y")
	fmt.Println(m)
}

func TestNode2(t *testing.T) {
	m := NewNodeManager()
	h1 := m.NewHeader("x")
	h2 := m.NewHeader("y")
	fmt.Println(m)
	n1 := m.NewNode(h1, nil)
	fmt.Println(n1)
	n2 := m.NewNode(h1, nil)
	fmt.Println(n2)
	n3 := m.NewNode(h2, n1, n2)
	fmt.Println(n3)
}

func TestUniqueTable1(t *testing.T) {
	table := NewUniqueTable()
	m := NewNodeManager()
	h1 := m.NewHeader("x")
	h2 := m.NewHeader("y")
	fmt.Println(m)
	n1 := m.NewNode(h1, nil)
	fmt.Println(n1)
	n2 := m.NewNode(h1, nil)
	fmt.Println(n2)
	n3 := m.NewNode(h2, n1, n2)
	fmt.Println(n3)
	key := table.GenKey(n1, n2)
	table.Set(key, n3)
	fmt.Println(n3)
	fmt.Println(table)
	fmt.Println(table.Get(table.GenKey(n1, n2)))
}

func TestNode3(t *testing.T) {
	m := NewNodeManager()
	h1 := m.NewHeader("x")
	h2 := m.NewHeader("y")
	fmt.Println(m)
	n1 := m.NewNode(h1, nil)
	fmt.Println(n1)
	n2 := m.NewNode(h1, nil)
	fmt.Println(n2)
	n3 := m.NewNode(h2, n1, n2)
	fmt.Println(n3)
	n3.ToDot(func(x *Node) string {
		if x == n1 {
			return "true"
		} else {
			return "false"
		}
	})
}

func TestNode4(t *testing.T) {
	m := NewNodeManager()
	h1 := m.NewHeader("x")
	h2 := m.NewHeader("y")
	fmt.Println(m)
	n1 := m.NewNode(h1, nil)
	fmt.Println(n1)
	n2 := m.NewNode(h1, nil)
	fmt.Println(n2)
	n3 := BDDNode{
		Node: m.NewNode(h2, n1, n2),
		edge: true,
	}
	table := NewUniqueTable()
	key := table.GenKey(h2, n1, n2)
	table.Set(key, n3.Node)
	fmt.Println(n3)
	n3.ToDot(func(x *Node) string {
		if x == n1 {
			return "true"
		} else {
			return "false"
		}
	})
}
