package bdd

import (
	"fmt"
	"testing"
)

func TestBDD1(t *testing.T) {
	f := NewBDD()
	one := f.terminals[true]
	zero := f.terminals[false]
	fmt.Println(one)
	fmt.Println(zero)
	h1 := f.NewHeader("x", 1, []DomainInt{0, 1})
	h2 := f.NewHeader("y", 2, []DomainInt{0, 1})
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

// func TestDDVar(t *testing.T) {
// 	forest := NewForest()
// 	zero := forest.MakeDDVal(0)
// 	fmt.Println(zero)
// 	one := forest.MakeDDVal(1)
// 	fmt.Println(one)
// 	n1 := forest.MakeDDVar(0, zero, one)
// 	fmt.Println(n1)
// 	n2 := forest.MakeDDVar(0, one, zero)
// 	fmt.Println(n2)
// 	n3 := forest.MakeDDVar(0, zero, one)
// 	fmt.Println(n3)
// 	a := fmt.Sprintf("%p\n", n1)
// 	b := fmt.Sprintf("%p\n", n3)
// 	if a != b {
// 		t.Error("fail")
// 	}
// }

// func TestApply1(t *testing.T) {
// 	forest := NewForest()
// 	forest.AddOperator("plus", func(f *Forest, r *DDVal, l *DDVal) *DDVal {
// 		if int(l.val) < int(r.val) {
// 			return f.MakeDDVal(l.val)
// 		} else {
// 			return f.MakeDDVal(r.val)
// 		}
// 	})
// 	zero := forest.MakeDDVal(0)
// 	fmt.Println(zero)
// 	one := forest.MakeDDVal(1)
// 	fmt.Println(one)
// 	n1 := forest.MakeDDVar(0, zero, one)
// 	fmt.Println(n1)
// 	n2 := forest.MakeDDVar(0, one, zero)
// 	fmt.Println(n2)
// 	fmt.Println(forest.Apply("plus", n1, n2))
// }
