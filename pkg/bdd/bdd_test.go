package bdd

import (
	"fmt"
	"testing"
)

func TestBDD1(t *testing.T) {
	b := NewBDD()
	x := b.Var("x")
	y := b.Var("y")
	z := b.And(x, y)
	fmt.Println(z)
	b.ToDot(x)
	b.ToDot(y)
	b.ToDot(z)
}

func TestBDD2(t *testing.T) {
	b := NewBDD()
	x := b.Var("x")
	y := b.Var("y")
	z := b.Or(x, y)
	fmt.Println(z)
	b.ToDot(x)
	b.ToDot(y)
	b.ToDot(z)
}

func TestBDD3(t *testing.T) {
	b := NewBDD()
	x := b.Var("x")
	y := b.Var("y")
	z := b.Xor(x, y)
	fmt.Println(z)
	b.ToDot(x)
	b.ToDot(y)
	b.ToDot(z)
}

func TestBDD01(t *testing.T) {
	b := NewBDD0()
	x := b.Var("x")
	y := b.Var("y")
	z := b.And(x, y)
	fmt.Println(z)
	b.ToDot(x)
	b.ToDot(y)
	b.ToDot(z)
}

func TestBDD02(t *testing.T) {
	b := NewBDD0()
	x := b.Var("x")
	y := b.Var("y")
	z := b.Or(x, y)
	fmt.Println(z)
	b.ToDot(x)
	b.ToDot(y)
	b.ToDot(z)
}

func TestBDD03(t *testing.T) {
	b := NewBDD0()
	x := b.Var("x")
	y := b.Var("y")
	z := b.Xor(x, y)
	fmt.Println(z)
	b.ToDot(x)
	b.ToDot(y)
	b.ToDot(z)
}
