package bdd

import (
	"fmt"
	"testing"
)

// utils

func (d *ZDD) todot(x *Node) {
	x.ToDot(func(x *Node) string {
		if x == d.one {
			return "1"
		} else {
			return "0"
		}
	})
}

func NewZDDWithLabels(labels []int) (*ZDD, map[int]*Node) {
	d := NewZDD()
	for _, i := range labels {
		d.NewHeader(fmt.Sprint(i), []DomainInt{0, 1})
	}
	nodes := make(map[int]*Node)
	for _, i := range labels {
		// make create vector
		x := []int{1}
		for _, j := range labels {
			if i != j {
				x = append(x, 0)
			} else {
				x = append(x, 1)
			}
		}
		nodes[i] = d.CreateNode(x)
	}
	return d, nodes
}

// tests

func TestZDD1(t *testing.T) {
	f := NewZDD()
	one := f.terminals[true]
	zero := f.terminals[false]
	fmt.Println(one)
	fmt.Println(zero)
	h1 := f.NewHeader("x", []DomainInt{0, 1})
	h2 := f.NewHeader("y", []DomainInt{0, 1})
	fmt.Println(h1, h2)
	x := f.Node(h2, f.Node(h1, one, zero), f.Node(h1, one, zero))
	y := f.Node(h2, f.Node(h1, zero, zero), f.Node(h1, one, one))
	z := f.Intersect(x, y)
	fmt.Println(z)
	f.todot(x)
	f.todot(y)
	f.todot(z)
}

func TestCreateZDD1(t *testing.T) {
	d := NewZDD()
	d.NewHeader("x1", []DomainInt{0, 1})
	d.NewHeader("x2", []DomainInt{0, 1})
	d.NewHeader("x3", []DomainInt{0, 1})
	x := d.CreateNode([]int{1, 0, 1, 0})
	d.todot(x)
}

func TestZDDFT2(t *testing.T) {
	labels := []int{1, 2, 3, 4, 5, 6, 7}
	d, vars := NewZDDWithLabels(labels)
	ft := d.Union(d.Product(vars[2]), d.Product(vars[4], vars[6], vars[7]), d.Product(vars[3], vars[5], vars[2]),
		d.Product(vars[1], vars[7], vars[3]), d.Product(vars[5], vars[7], vars[3]), d.Product(vars[1], vars[2], vars[6]))
	d.todot(ft)
	result := d.getmcs(ft)
	fmt.Println(result)
}
