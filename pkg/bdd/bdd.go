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
	one       *Node
	zero      *Node
	notBDD    *NotBDD
	orBDD     *OrBDD
	andBDD    *AndBDD
	xorBDD    *XorBDD
}

func NewBDD() *BDD {
	d := &BDD{
		NodeManager: NewNodeManager(0),
		UniqueTable: NewUniqueTable(),
		values:      make(map[*Node]BDDTerminalValue),
		terminals:   make(map[BDDTerminalValue]*Node),
	}
	d.notBDD = &NotBDD{
		UniqueTable: NewUniqueTable(),
		d:           d,
	}
	d.orBDD = &OrBDD{
		UniqueTable: NewUniqueTable(),
		d:           d,
	}
	d.andBDD = &AndBDD{
		UniqueTable: NewUniqueTable(),
		d:           d,
	}
	d.xorBDD = &XorBDD{
		UniqueTable: NewUniqueTable(),
		d:           d,
	}
	h := d.NewHeader("", nil)
	d.zero = d.Terminal(h, false)
	d.one = d.Terminal(h, true)
	return d
}

func (d *BDD) CreateNode(xs []int) *Node {
	var zero, one *Node
	if xs[0] == 0 {
		zero, one = d.one, d.zero
	} else {
		zero, one = d.zero, d.one
	}

	for i := 1; i < len(xs); i++ {
		switch xs[i] {
		case 0:
			zero, one = d.Node(d.headers[i], zero, zero), d.Node(d.headers[i], one, zero)
		case 1:
			zero, one = d.Node(d.headers[i], zero, zero), d.Node(d.headers[i], zero, one)
		default:
			zero, one = d.Node(d.headers[i], zero, zero), d.Node(d.headers[i], one, one)
		}
	}

	return one
}

func (d *BDD) Not(f *Node) *Node {
	return d.UniApply(d.notBDD, f)
}

func (d *BDD) And(x ...*Node) *Node {
	result := x[0]
	for i := 1; i < len(x); i++ {
		result = d.BinApply(d.andBDD, result, x[i])
	}
	return result
}

func (d *BDD) Or(x ...*Node) *Node {
	result := x[0]
	for i := 1; i < len(x); i++ {
		result = d.BinApply(d.orBDD, result, x[i])
	}
	return result
}

func (d *BDD) Xor(x ...*Node) *Node {
	result := x[0]
	for i := 1; i < len(x); i++ {
		result = d.BinApply(d.xorBDD, result, x[i])
	}
	return result
}

func (d *BDD) Ite(f, g, h *Node) *Node {
	return d.Or(d.And(f, g), d.And(d.Not(f), h))
}

func (d *BDD) Node(header *NodeHeader, nodes ...*Node) *Node {
	// reduction
	if nodes[0] == nodes[1] {
		return nodes[0]
	}
	// sharing
	key := d.GenNodeKey(header.level, nodes...)
	if node, ok := d.Get(key); ok {
		return node
	}
	node := d.NewNode(header, nodes)
	d.Set(key, node)
	return node
}

func (d *BDD) Terminal(header *NodeHeader, x BDDTerminalValue) *Node {
	if n, ok := d.terminals[x]; ok {
		return n
	} else {
		n = d.NewNode(header, nil)
		d.terminals[x] = n
		d.values[n] = x
		return n
	}
}

func (d *BDD) GetValue(n *Node) BDDTerminalValue {
	if v, ok := d.values[n]; ok {
		return v
	} else {
		panic("Cannot find a value for a given node.")
	}
}

func (d *BDD) UniApply(op UniOperator, f *Node) *Node {
	if node, ok := op.checkTerminal(f); ok {
		return node
	}
	key := op.GenKey(f)
	if node, ok := op.Get(key); ok {
		return node
	}
	node0 := d.UniApply(op, f.nodes[0])
	node1 := d.UniApply(op, f.nodes[1])
	result := d.Node(f.header, node0, node1)
	op.Set(key, result)
	return result
}

func (d *BDD) BinApply(op BinOperator, f, g *Node) *Node {
	if node, ok := op.checkTerminal(f, g); ok {
		return node
	}
	key := op.GenKey(f, g)
	if node, ok := op.Get(key); ok {
		return node
	}

	var result *Node
	switch {
	case f.header.level > g.header.level:
		node0 := d.BinApply(op, f.nodes[0], g)
		node1 := d.BinApply(op, f.nodes[1], g)
		result = d.Node(f.header, node0, node1)
	case f.header.level < g.header.level:
		node0 := d.BinApply(op, f, g.nodes[0])
		node1 := d.BinApply(op, f, g.nodes[1])
		result = d.Node(g.header, node0, node1)
	case f.header.level == g.header.level:
		node0 := d.BinApply(op, f.nodes[0], g.nodes[0])
		node1 := d.BinApply(op, f.nodes[1], g.nodes[1])
		result = d.Node(f.header, node0, node1)
	}
	op.Set(key, result)
	return result
}

type NotBDD struct {
	*UniqueTable
	d *BDD
}

type OrBDD struct {
	*UniqueTable
	d *BDD
}

type AndBDD struct {
	*UniqueTable
	d *BDD
}

type XorBDD struct {
	*UniqueTable
	d *BDD
}

func (op *NotBDD) checkTerminal(node *Node) (*Node, bool) {
	switch {
	case node.isTerminal():
		return op.d.Terminal(node.header, !op.d.GetValue(node)), true
	}
	return nil, false
}

func (op *OrBDD) checkTerminal(f, g *Node) (*Node, bool) {
	switch {
	case f.isTerminal() && g.isTerminal():
		return op.d.Terminal(f.header, op.d.GetValue(f) || op.d.GetValue(g)), true
		// case f.isTerminal() && op.d.GetValue(f) == false:
		// 	return g, true
		// case g.isTerminal() && op.d.GetValue(g) == false:
		// 	return f, true
		// case f.isTerminal() && op.d.GetValue(f) == true:
		// 	return op.d.Terminal(f.header, true), true
		// case g.isTerminal() && op.d.GetValue(g) == true:
		// 	return op.d.Terminal(g.header, true), true
	}
	return nil, false
}

func (op *AndBDD) checkTerminal(f, g *Node) (*Node, bool) {
	switch {
	case f.isTerminal() && g.isTerminal():
		return op.d.Terminal(f.header, op.d.GetValue(f) && op.d.GetValue(g)), true
		// case f.isTerminal() && op.d.GetValue(f) == false:
		// 	return op.d.Terminal(f.header, false), true
		// case g.isTerminal() && op.d.GetValue(g) == false:
		// 	return op.d.Terminal(g.header, false), true
		// case f.isTerminal() && op.d.GetValue(f) == true:
		// 	return g, true
		// case g.isTerminal() && op.d.GetValue(g) == true:
		// 	return f, true
	}
	return nil, false
}

func (op *XorBDD) checkTerminal(f, g *Node) (*Node, bool) {
	switch {
	case f.isTerminal() && g.isTerminal():
		if op.d.GetValue(f) == op.d.GetValue(g) {
			return op.d.Terminal(f.header, false), true
		} else {
			return op.d.Terminal(f.header, true), true
		}
		// case f.isTerminal() && op.d.GetValue(f) == false:
		// 	return g, true
		// case g.isTerminal() && op.d.GetValue(g) == false:
		// 	return f, true
		// case f.isTerminal() && op.d.GetValue(f) == true:
		// 	return op.d.UniApply(op.d.notBDD, g), true
		// case g.isTerminal() && op.d.GetValue(g) == true:
		// 	return op.d.UniApply(op.d.notBDD, f), true
	}
	return nil, false
}
