package bdd

import (
	"fmt"
)

type Node0 struct {
	id      Id
	header  *NodeHeader // info of Node
	nodes   [2]*Node0   // the pointer for Nodes
	notedge bool        // indicator whether 1-edge is not-edge or not
}

func (n *Node0) hasNotEdge() bool {
	return n.notedge
}

type Node01 struct {
	*Node0
	Neg bool
}

// struct to identify a node
type nodeId0 struct {
	h   Id
	f   Id
	g   Id
	neg bool
}

type UniqueNodeTable0 struct {
	buf   nodeId0
	table map[nodeId]*Node01
}

func NewUniqueNodeTable0() *UniqueNodeTable0 {
	return &UniqueNodeTable0{
		buf:   nodeId0{0, 0, 0, false},
		table: make(map[nodeId0]*Node01),
	}
}

func (t *UniqueNodeTable0) get(head *NodeHeader, low, high *Node) (*Node, bool) {
	t.buf.head, t.buf.low, t.buf.high = head.id, low.id, high.id
	v, ok := t.table[t.buf]
	return v, ok
}

func (t *UniqueNodeTable) set(head *NodeHeader, low, high, x *Node) {
	t.buf.head, t.buf.low, t.buf.high = head.id, low.id, high.id
	t.table[t.buf] = x
}

// struct for operation cache
type tupleNode0 struct {
	f    Id
	g    Id
	fneg bool
	gneg bool
}

type BDD0 struct {
	total    Id
	headers  []*NodeHeader
	nodes    map[nodeId0]*Node0
	Zero     *Node0
	cacheOr  map[tupleNode0]*Node01
	cacheAnd map[tupleNode0]*Node01
	cacheXor map[tupleNode0]*Node01
	key1     nodeId0
	key2     tupleNode0
}

func NewBDD0() *BDD0 {
	h := &NodeHeader{
		id:    0,
		level: 0,
	}
	zero := &Node0{
		id:     0,
		header: h,
	}
	return &BDD0{
		total:    1,
		headers:  []*NodeHeader{h},
		nodes:    make(map[nodeId0]*Node0),
		Zero:     zero,
		cacheOr:  make(map[tupleNode0]*Node01),
		cacheAnd: make(map[tupleNode0]*Node01),
		cacheXor: make(map[tupleNode0]*Node01),
		key1:     nodeId0{},
		key2:     tupleNode0{},
	}
}

func (b *BDD0) Var(name string) *Node01 {
	h := &NodeHeader{
		id:    Id(len(b.headers)),
		level: Level(len(b.headers)),
		label: name,
	}
	b.headers = append(b.headers, h)
	result, resneg := b.createNode(h, b.Zero, false, b.Zero, true)
	return &Node01{result, resneg}
}

func (b *BDD0) GetVarLabels() []string {
	result := make([]string, len(b.headers))
	for i, x := range b.headers {
		result[i] = x.label
	}
	return result
}

func (b *BDD0) GetVars() []*Node01 {
	result := make([]*Node01, len(b.headers))
	for i, h := range b.headers {
		r0, r1 := b.createNode(h, b.Zero, false, b.Zero, true)
		result[i] = &Node01{r0, r1}
	}
	return result
}

func (b *BDD0) Not(f *Node01) *Node01 {
	result, resneg := b.not(f.Node0, f.Neg)
	return &Node01{result, resneg}
}

func (b *BDD0) And(xs ...*Node01) *Node01 {
	result, resneg := b.Zero, true
	for _, x := range xs {
		result, resneg = b.and(result, resneg, x.Node0, x.Neg)
	}
	return &Node01{result, resneg}
}

func (b *BDD0) Or(xs ...*Node01) *Node01 {
	result, resneg := b.Zero, false
	for _, x := range xs {
		result, resneg = b.or(result, resneg, x.Node0, x.Neg)
	}
	return &Node01{result, resneg}
}

func (b *BDD0) Xor(xs ...*Node01) *Node01 {
	result, resneg := b.Zero, false
	for _, x := range xs {
		result, resneg = b.xor(result, resneg, x.Node0, x.Neg)
	}
	return &Node01{result, resneg}
}

func (b *BDD0) Imp(f, g *Node01) *Node01 {
	return b.Or(b.Not(f), g)
}

func (b *BDD0) Ite(f, g, h *Node01) *Node01 {
	return b.Or(b.And(f, g), b.And(b.Not(f), h))
}

func (n *Node0) Zero(neg bool) (*Node0, bool) {
	return n.nodes[0], neg
}

func (n *Node0) One(neg bool) (*Node0, bool) {
	return n.nodes[1], (neg && !n.hasNotEdge()) || (!neg && n.hasNotEdge())
}

func (n *Node0) Level() Level {
	return n.header.level - 1
}

func (n *Node0) Label() string {
	return n.header.label
}

func (b *BDD0) createNode(h *NodeHeader, low *Node0, lneg bool, high *Node0, hneg bool) (*Node0, bool) {
	if low == high && lneg == hneg {
		return low, lneg
	}
	hneg = lneg != hneg
	b.key1.h, b.key1.f, b.key1.g, b.key1.neg = h.id, low.id, high.id, hneg
	if node, ok := b.nodes[b.key1]; ok {
		return node, lneg
	}
	node := &Node0{
		id:      b.total,
		header:  h,
		nodes:   [2]*Node0{low, high},
		notedge: hneg,
	}
	b.total++
	b.nodes[b.key1] = node
	return node, lneg
}

func (b *BDD0) not(f *Node0, neg bool) (*Node0, bool) {
	return f, !neg
}

func (b *BDD0) and(f *Node0, fneg bool, g *Node0, gneg bool) (*Node0, bool) {
	switch {
	case (f == b.Zero && fneg == false) || (g == b.Zero && gneg == false):
		return b.Zero, false
	case f == b.Zero && fneg == true:
		return g, gneg
	case g == b.Zero && gneg == true:
		return f, fneg
	}
	b.key2.f, b.key2.g, b.key2.fneg, b.key2.gneg = f.id, g.id, fneg, gneg
	if node, ok := b.cacheAnd[b.key2]; ok {
		return node.Node0, node.Neg
	}
	var result *Node0
	var resneg bool
	switch {
	case f.header.level > g.header.level:
		n00, b00 := f.Zero(fneg)
		n01, b01 := b.and(n00, b00, g, gneg)
		n10, b10 := f.One(fneg)
		n11, b11 := b.and(n10, b10, g, gneg)
		result, resneg = b.createNode(f.header, n01, b01, n11, b11)
	case f.header.level < g.header.level:
		n00, b00 := g.Zero(gneg)
		n01, b01 := b.and(f, fneg, n00, b00)
		n10, b10 := g.One(gneg)
		n11, b11 := b.and(f, fneg, n10, b10)
		result, resneg = b.createNode(g.header, n01, b01, n11, b11)
	case f.header.level == g.header.level:
		fn00, fb00 := f.Zero(fneg)
		gn00, gb00 := g.Zero(gneg)
		n01, b01 := b.and(fn00, fb00, gn00, gb00)
		fn10, fb10 := f.One(fneg)
		gn10, gb10 := g.One(gneg)
		n11, b11 := b.and(fn10, fb10, gn10, gb10)
		result, resneg = b.createNode(f.header, n01, b01, n11, b11)
	}
	b.key2.f, b.key2.g, b.key2.fneg, b.key2.gneg = f.id, g.id, fneg, gneg
	b.cacheAnd[b.key2] = &Node01{result, resneg}
	return result, resneg
}

func (b *BDD0) or(f *Node0, fneg bool, g *Node0, gneg bool) (*Node0, bool) {
	switch {
	case (f == b.Zero && fneg == true) || (g == b.Zero && gneg == true):
		return b.Zero, true
	case f == b.Zero && fneg == false:
		return g, gneg
	case g == b.Zero && gneg == false:
		return f, fneg
	}
	b.key2.f, b.key2.g, b.key2.fneg, b.key2.gneg = f.id, g.id, fneg, gneg
	if node, ok := b.cacheOr[b.key2]; ok {
		return node.Node0, node.Neg
	}
	var result *Node0
	var resneg bool
	switch {
	case f.header.level > g.header.level:
		n00, b00 := f.Zero(fneg)
		n01, b01 := b.or(n00, b00, g, gneg)
		n10, b10 := f.One(fneg)
		n11, b11 := b.or(n10, b10, g, gneg)
		result, resneg = b.createNode(f.header, n01, b01, n11, b11)
	case f.header.level < g.header.level:
		n00, b00 := g.Zero(gneg)
		n01, b01 := b.or(f, fneg, n00, b00)
		n10, b10 := g.One(gneg)
		n11, b11 := b.or(f, fneg, n10, b10)
		result, resneg = b.createNode(g.header, n01, b01, n11, b11)
	case f.header.level == g.header.level:
		fn00, fb00 := f.Zero(fneg)
		gn00, gb00 := g.Zero(gneg)
		n01, b01 := b.or(fn00, fb00, gn00, gb00)
		fn10, fb10 := f.One(fneg)
		gn10, gb10 := g.One(gneg)
		n11, b11 := b.or(fn10, fb10, gn10, gb10)
		result, resneg = b.createNode(f.header, n01, b01, n11, b11)
	}
	b.key2.f, b.key2.g, b.key2.fneg, b.key2.gneg = f.id, g.id, fneg, gneg
	b.cacheOr[b.key2] = &Node01{result, resneg}
	return result, resneg
}

func (b *BDD0) xor(f *Node0, fneg bool, g *Node0, gneg bool) (*Node0, bool) {
	switch {
	case f == b.Zero && fneg == true:
		return g, !gneg
	case f == b.Zero && fneg == false:
		return g, gneg
	case g == b.Zero && fneg == true:
		return f, !fneg
	case g == b.Zero && gneg == false:
		return f, fneg
	}
	b.key2.f, b.key2.g, b.key2.fneg, b.key2.gneg = f.id, g.id, fneg, gneg
	if node, ok := b.cacheXor[b.key2]; ok {
		return node.Node0, node.Neg
	}
	var result *Node0
	var resneg bool
	switch {
	case f.header.level > g.header.level:
		n00, b00 := f.Zero(fneg)
		n01, b01 := b.xor(n00, b00, g, gneg)
		n10, b10 := f.One(fneg)
		n11, b11 := b.xor(n10, b10, g, gneg)
		result, resneg = b.createNode(f.header, n01, b01, n11, b11)
	case f.header.level < g.header.level:
		n00, b00 := g.Zero(gneg)
		n01, b01 := b.xor(f, fneg, n00, b00)
		n10, b10 := g.One(gneg)
		n11, b11 := b.xor(f, fneg, n10, b10)
		result, resneg = b.createNode(g.header, n01, b01, n11, b11)
	case f.header.level == g.header.level:
		fn00, fb00 := f.Zero(fneg)
		gn00, gb00 := g.Zero(gneg)
		n01, b01 := b.xor(fn00, fb00, gn00, gb00)
		fn10, fb10 := f.One(fneg)
		gn10, gb10 := g.One(gneg)
		n11, b11 := b.xor(fn10, fb10, gn10, gb10)
		result, resneg = b.createNode(f.header, n01, b01, n11, b11)
	}
	b.key2.f, b.key2.g, b.key2.fneg, b.key2.gneg = f.id, g.id, fneg, gneg
	b.cacheXor[b.key2] = &Node01{result, resneg}
	return result, resneg
}

// todot

func (b *BDD0) ToDot(f *Node01) {
	fmt.Println("digraph { layout=dot; overlap=false; splines=true; node [fontsize=10];")
	visited := make(map[*Node0]struct{})
	b.todot(f.Node0, visited)
	fmt.Println("}")
}

func (b *BDD0) todot(f *Node0, visited map[*Node0]struct{}) {
	if _, ok := visited[f]; ok {
		return
	}
	switch {
	case f == b.Zero:
		fmt.Printf(" n%p [shape = square, label = \"0\"];\n", f)
	default:
		fmt.Printf(" n%p [shape = circle, label = \"%s\"];\n", f, f.header.label)
		b.todot(f.nodes[0], visited)
		fmt.Printf(" n%p -> n%p [label = \"0\"];\n", f, f.nodes[0])
		b.todot(f.nodes[1], visited)
		if f.notedge == false {
			fmt.Printf(" n%p -> n%p [label = \"1\"];\n", f, f.nodes[1])
		} else {
			fmt.Printf(" n%p -> n%p [label = \"1\", arrowhead = odot];\n", f, f.nodes[1])
		}
	}
	visited[f] = struct{}{}
}
