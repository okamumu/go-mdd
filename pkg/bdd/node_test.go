package bdd

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestNode1(t *testing.T) {
	m := NewNodeManager(0)
	m.NewHeader("x", 1, []DomainInt{0, 1, 2})
	m.NewHeader("y", 2, []DomainInt{0, 1, 2})
	fmt.Println(m)
}

func TestNode2(t *testing.T) {
	m := NewNodeManager(0)
	h1 := m.NewHeader("x", 1, []DomainInt{0, 1, 2})
	h2 := m.NewHeader("y", 2, []DomainInt{0, 1, 2})
	fmt.Println(m)
	n1 := m.NewNode(h1, nil)
	fmt.Println(n1)
	n2 := m.NewNode(h1, nil)
	fmt.Println(n2)
	n3 := m.NewNode(h2, []*Node{n1, n2})
	fmt.Println(n3)
}

func TestUniqueTable1(t *testing.T) {
	table := NewUniqueTable()
	m := NewNodeManager(0)
	h1 := m.NewHeader("x", 1, []DomainInt{0, 1, 2})
	h2 := m.NewHeader("y", 2, []DomainInt{0, 1, 2})
	fmt.Println(m)
	n1 := m.NewNode(h1, nil)
	fmt.Println(n1)
	n2 := m.NewNode(h1, nil)
	fmt.Println(n2)
	n3 := m.NewNode(h2, []*Node{n1, n2})
	fmt.Println(n3)
	key := table.GenKey([]*Node{n1, n2}...)
	table.Set(key, n3)
	fmt.Println(n3)
	fmt.Println(table)
	fmt.Println(table.Get(table.GenKey(n1, n2)))
	fmt.Println(unsafe.Sizeof(*n3))
}

func TestNode3(t *testing.T) {
	m := NewNodeManager(0)
	h1 := m.NewHeader("x", 0, []DomainInt{0, 1, 2})
	h2 := m.NewHeader("y", 2, []DomainInt{0, 1, 2})
	fmt.Println(m)
	n1 := m.NewNode(h1, nil)
	fmt.Println(n1)
	n2 := m.NewNode(h1, nil)
	fmt.Println(n2)
	n3 := m.NewNode(h2, []*Node{n1, n2})
	fmt.Println(n3)
	n3.ToDot(func(x *Node) string {
		if x == n1 {
			return "true"
		} else {
			return "false"
		}
	})
}
