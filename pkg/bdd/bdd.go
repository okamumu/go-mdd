package bdd

import (
// 	"fmt"
)

type BDDTerminalValue bool

type BDD struct {
	*NodeManager
	*UniqueTable
	values    map[*Node]BDDTerminalValue
	terminals map[BDDTerminalValue]*Node
	notbdd    *NotBDD
	orbdd     *OrBDD
	andbdd    *AndBDD
}

type NotBDD struct {
	*UniqueTable
	f *BDD
}

type OrBDD struct {
	*UniqueTable
	f *BDD
}

type AndBDD struct {
	*UniqueTable
	f *BDD
}

func NewBDD() *BDD {
	f := &BDD{
		NodeManager: NewNodeManager(0),
		UniqueTable: NewUniqueTable(),
		values:      make(map[*Node]BDDTerminalValue),
		terminals:   make(map[BDDTerminalValue]*Node),
	}
	f.notbdd = &NotBDD{
		UniqueTable: NewUniqueTable(),
		f:           f,
	}
	f.orbdd = &OrBDD{
		UniqueTable: NewUniqueTable(),
		f:           f,
	}
	f.andbdd = &AndBDD{
		UniqueTable: NewUniqueTable(),
		f:           f,
	}
	h := f.NewHeader("", 0, nil)
	f.Terminal(h, true)
	f.Terminal(h, false)
	return f
}

func (b *BDD) Not(f *Node) *Node {
	return UniApply(b, b.notbdd, f)
}

func (b *BDD) And(f, g *Node) *Node {
	return BinApply(b, b.andbdd, f, g)
}

func (b *BDD) Or(f, g *Node) *Node {
	return BinApply(b, b.andbdd, f, g)
}

func (b *BDD) Ite(f, g, h *Node) *Node {
	return b.Or(b.And(f, g), b.And(b.Not(f), h))
}

func (b *BDD) Show(s []byte, f *Node) {

}

func (b *BDD) Node(header *NodeHeader, nodes ...*Node) *Node {
	key := b.GenKey(nodes...)
	if node, ok := b.Get(key); ok {
		return node
	}
	node := b.NewNode(header, nodes)
	b.Set(key, node)
	return node
}

func (b *BDD) Terminal(header *NodeHeader, x BDDTerminalValue) *Node {
	if n, ok := b.terminals[x]; ok {
		return n
	} else {
		n = b.NewNode(header, nil)
		b.terminals[x] = n
		b.values[n] = x
		return n
	}
}

func (b *BDD) GetValue(n *Node) BDDTerminalValue {
	if v, ok := b.values[n]; ok {
		return v
	} else {
		panic("!")
	}
}

func (op *NotBDD) checkTerminal(node *Node) (*Node, bool) {
	switch {
	case node.isTerminal():
		return op.f.Terminal(node.header, !op.f.GetValue(node)), true
	}
	return nil, false
}

func (op *OrBDD) checkTerminal(left, right *Node) (*Node, bool) {
	switch {
	case left.isTerminal() && right.isTerminal():
		return op.f.Terminal(left.header, op.f.GetValue(left) || op.f.GetValue(right)), true
	case left.isTerminal() && op.f.GetValue(left) == false:
		return right, true
	case right.isTerminal() && op.f.GetValue(right) == false:
		return left, true
	case left.isTerminal() && op.f.GetValue(left) == true:
		return op.f.Terminal(left.header, true), true
	case right.isTerminal() && op.f.GetValue(right) == true:
		return op.f.Terminal(right.header, false), true
	}
	return nil, false
}

func (op *AndBDD) checkTerminal(left, right *Node) (*Node, bool) {
	switch {
	case left.isTerminal() && op.f.GetValue(left) == false:
		return op.f.Terminal(left.header, false), true
	case right.isTerminal() && op.f.GetValue(right) == false:
		return op.f.Terminal(right.header, false), true
	case left.isTerminal() && op.f.GetValue(left) == true:
		return right, true
	case right.isTerminal() && op.f.GetValue(right) == true:
		return left, true
	case left.isTerminal() && right.isTerminal():
		return op.f.Terminal(left.header, op.f.GetValue(left) && op.f.GetValue(right)), true
	}
	return nil, false
}
