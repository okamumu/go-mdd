package bdd

import (
	"fmt"
	"sort"
	"testing"
)

// utils

func NewBDDWithLabels(labels []int) (*BDD, map[int]*Node) {
	d := NewBDD()
	for _, i := range labels {
		d.NewHeader(fmt.Sprint(i), []DomainInt{0, 1})
	}
	nodes := make(map[int]*Node)
	for _, i := range labels {
		// make create vector
		x := []int{1}
		for _, j := range labels {
			if i != j {
				x = append(x, -1)
			} else {
				x = append(x, 1)
			}
		}
		nodes[i] = d.CreateNode(x)
	}
	return d, nodes
}

func (d *BDD) todot(x *Node) {
	x.ToDot(func(x *Node) string {
		if x == d.one {
			return "1"
		} else {
			return "0"
		}
	})
}

func (d *BDD) getmcs(x *Node) [][]string {
	switch {
	case x == d.one:
		return [][]string{[]string{}}
	case x == d.zero:
		return [][]string{}
	default:
		result0 := d.getmcs(x.nodes[0])
		result1 := d.getmcs(x.nodes[1])
		for i, _ := range result1 {
			result1[i] = append(result1[i], x.header.label)
		}
		ans := append(result0, result1...)
		sort.Slice(ans, func(i, j int) bool {
			if len(ans[i]) < len(ans[j]) {
				return true
			} else {
				return false
			}
		})
		s := [][]string{}
		for i := 0; i < len(ans); i++ {
			if isIncludeList(ans[i], s) == false {
				s = append(s, ans[i])
			}
		}
		return s
	}
}

func TestBDD1(t *testing.T) {
	f := NewBDD()
	one := f.terminals[true]
	zero := f.terminals[false]
	fmt.Println(one)
	fmt.Println(zero)
	h1 := f.NewHeader("x", []DomainInt{0, 1})
	h2 := f.NewHeader("y", []DomainInt{0, 1})
	fmt.Println(h1, h2)
	x := f.Node(h2, f.Node(h1, one, zero), f.Node(h1, one, zero))
	y := f.Node(h2, f.Node(h1, zero, zero), f.Node(h1, one, one))
	z := f.And(x, y)
	fmt.Println(z)
	x.ToDot(func(x *Node) string {
		if x == one {
			return "1"
		} else {
			return "0"
		}
	})
	y.ToDot(func(x *Node) string {
		if x == one {
			return "1"
		} else {
			return "0"
		}
	})
	z.ToDot(func(x *Node) string {
		if x == one {
			return "1"
		} else {
			return "0"
		}
	})
}

func TestCreateBDD1(t *testing.T) {
	d := NewBDD()
	d.NewHeader("x1", []DomainInt{0, 1})
	d.NewHeader("x2", []DomainInt{0, 1})
	d.NewHeader("x3", []DomainInt{0, 1})
	x := d.CreateNode([]int{1, -1, 1, -1})
	x.ToDot(func(x *Node) string {
		if x == d.one {
			return "1"
		} else {
			return "0"
		}
	})
}

func TestBDDFT1(t *testing.T) {
	labels, result := makeTestData()
	d, vars := NewBDDWithLabels(labels)
	fts := []*Node{}
	for _, xs := range result {
		ans := vars[xs[0]]
		for i := 1; i < len(xs); i++ {
			ans = d.And(ans, vars[xs[i]])
		}
		fts = append(fts, ans)
	}
	ft := d.Or(fts...)
	fmt.Println(ft)
	mcs := d.getmcs(ft)
	fmt.Println(mcs)
}

func TestBDDFT2(t *testing.T) {
	labels := []int{1, 2, 3, 4, 5, 6, 7}
	d, vars := NewBDDWithLabels(labels)
	xs := []*Node{
		d.And(vars[2]),
		d.And(vars[4], vars[6], vars[7]),
		d.And(vars[3], vars[5], vars[2]),
		d.And(vars[1], vars[7], vars[3]),
		d.And(vars[5], vars[7], vars[3]),
		d.And(vars[1], vars[2], vars[6]),
	}
	ft := d.Or(xs...)
	for _, x := range xs {
		d.todot(x)
	}
	d.todot(ft)
	result := d.getmcs(ft)
	fmt.Println(result)
}
