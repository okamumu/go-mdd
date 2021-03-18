package bdd

import (
// 	"fmt"
)

type ZDDTerminalValue bool

type ZDD struct {
	*NodeManager
	*UniqueTable
	values    map[*Node]ZDDTerminalValue
	terminals map[ZDDTerminalValue]*Node
	one       *Node
	zero      *Node
	notZDD    *NotZDD
	unionZDD  *OrZDD
	insectZDD *AndZDD
	diffZDD   *XorZDD
	prodZDD   *AndZDD
	// plusZDD   *OrZDD
}

func NewZDD() *ZDD {
	d := &ZDD{
		NodeManager: NewNodeManager(0),
		UniqueTable: NewUniqueTable(),
		values:      make(map[*Node]ZDDTerminalValue),
		terminals:   make(map[ZDDTerminalValue]*Node),
	}
	d.notZDD = &NotZDD{
		UniqueTable: NewUniqueTable(),
		d:           d,
	}
	d.unionZDD = &OrZDD{
		UniqueTable: NewUniqueTable(),
		d:           d,
	}
	d.insectZDD = &AndZDD{
		UniqueTable: NewUniqueTable(),
		d:           d,
	}
	d.diffZDD = &XorZDD{
		UniqueTable: NewUniqueTable(),
		d:           d,
	}
	d.prodZDD = &AndZDD{
		UniqueTable: NewUniqueTable(),
		d:           d,
	}
	// d.plusZDD = &OrZDD{
	// 	UniqueTable: NewUniqueTable(),
	// 	d:           d,
	// }
	h := d.NewHeader("", nil)
	d.zero = d.Terminal(h, false)
	d.one = d.Terminal(h, true)
	return d
}

func (d *ZDD) CreateNode(xs []int) *Node {
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

func (d *ZDD) Not(f *Node) *Node {
	return d.UniApply(d.notZDD, f)
}

func (d *ZDD) Intersect(x ...*Node) *Node {
	result := x[0]
	for i := 1; i < len(x); i++ {
		result = d.BinApply(d.insectZDD, result, x[i])
	}
	return result
}

func (d *ZDD) Union(x ...*Node) *Node {
	result := x[0]
	for i := 1; i < len(x); i++ {
		result = d.BinApply(d.unionZDD, result, x[i])
	}
	return result
}

func (d *ZDD) Setdiff(x ...*Node) *Node {
	result := x[0]
	for i := 1; i < len(x); i++ {
		result = d.BinApply(d.diffZDD, result, x[i])
	}
	return result
}

func (d *ZDD) Product(x ...*Node) *Node {
	result := x[0]
	for i := 1; i < len(x); i++ {
		result = d.BinApply2(d.prodZDD, result, x[i])
	}
	return result
}

// func (d *ZDD) Plus(x ...*Node) *Node {
// 	result := x[0]
// 	for i := 1; i < len(x); i++ {
// 		result = d.BinApply2(d.plusZDD, result, x[i])
// 	}
// 	return result
// }

func (d *ZDD) Node(header *NodeHeader, nodes ...*Node) *Node {
	// reduction
	if nodes[1] == d.zero {
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

func (d *ZDD) Terminal(header *NodeHeader, x ZDDTerminalValue) *Node {
	if n, ok := d.terminals[x]; ok {
		return n
	} else {
		n = d.NewNode(header, nil)
		d.terminals[x] = n
		d.values[n] = x
		return n
	}
}

func (d *ZDD) GetValue(n *Node) ZDDTerminalValue {
	if v, ok := d.values[n]; ok {
		return v
	} else {
		panic("!")
	}
}

func (d *ZDD) UniApply(op UniOperator, f *Node) *Node {
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

func (d *ZDD) BinApply(op BinOperator, f, g *Node) *Node {
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
		node1 := d.BinApply(op, f.nodes[1], d.zero)
		result = d.Node(f.header, node0, node1)
	case f.header.level < g.header.level:
		node0 := d.BinApply(op, f, g.nodes[0])
		node1 := d.BinApply(op, d.zero, g.nodes[1])
		result = d.Node(g.header, node0, node1)
	case f.header.level == g.header.level:
		node0 := d.BinApply(op, f.nodes[0], g.nodes[0])
		node1 := d.BinApply(op, f.nodes[1], g.nodes[1])
		result = d.Node(f.header, node0, node1)
	}
	op.Set(key, result)
	return result
}

func (d *ZDD) BinApply2(op BinOperator, f, g *Node) *Node {
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
		node0 := d.BinApply2(op, f.nodes[0], g)
		node1 := d.BinApply2(op, f.nodes[1], g)
		result = d.Node(f.header, node0, node1)
	case f.header.level < g.header.level:
		node0 := d.BinApply2(op, f, g.nodes[0])
		node1 := d.BinApply2(op, f, g.nodes[1])
		result = d.Node(g.header, node0, node1)
	case f.header.level == g.header.level:
		node0 := d.BinApply2(op, f.nodes[0], g.nodes[0])
		node1 := d.BinApply2(op, f.nodes[1], g.nodes[1])
		result = d.Node(f.header, node0, node1)
	}
	op.Set(key, result)
	return result
}

type NotZDD struct {
	*UniqueTable
	d *ZDD
}

type OrZDD struct {
	*UniqueTable
	d *ZDD
}

type AndZDD struct {
	*UniqueTable
	d *ZDD
}

type XorZDD struct {
	*UniqueTable
	d *ZDD
}

func (op *NotZDD) checkTerminal(node *Node) (*Node, bool) {
	switch {
	case node.isTerminal():
		return op.d.Terminal(node.header, !op.d.GetValue(node)), true
	}
	return nil, false
}

func (op *OrZDD) checkTerminal(f, g *Node) (*Node, bool) {
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
		// 	return op.d.Terminal(g.header, false), true
	}
	return nil, false
}

func (op *AndZDD) checkTerminal(f, g *Node) (*Node, bool) {
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

func (op *XorZDD) checkTerminal(f, g *Node) (*Node, bool) {
	switch {
	case f.isTerminal() && g.isTerminal():
		if op.d.GetValue(f) == op.d.GetValue(g) {
			return op.d.Terminal(f.header, false), true
		} else {
			return op.d.Terminal(f.header, true), true
		}
	}
	return nil, false
}
