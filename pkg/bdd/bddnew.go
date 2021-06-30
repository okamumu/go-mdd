package bdd

import (
	"fmt"
)

type Id uint64
type Level uint64

type NodeHeader struct {
	id    Id
	level Level  // Level of node
	label string // Label of node
}

type Node struct {
	id     Id
	header *NodeHeader // info of Node
	nodes  [2]*Node    // the pointer for Nodes
}

// struct to identify a node
type nodeId struct {
	head Id
	low  Id
	high Id
}

type UniqueNodeTable struct {
	buf   nodeId
	table map[nodeId]*Node
}

func NewUniqueNodeTable() *UniqueNodeTable {
	return &UniqueNodeTable{
		buf: nodeId{
			head: 0,
			low:  0,
			high: 0,
		},
		table: make(map[nodeId]*Node),
	}
}

func (t *UniqueNodeTable) get(head *NodeHeader, low, high *Node) (*Node, bool) {
	t.buf.head, t.buf.low, t.buf.high = head.id, low.id, high.id
	v, ok := t.table[t.buf]
	return v, ok
}

func (t *UniqueNodeTable) set(head *NodeHeader, low, high, x *Node) {
	t.buf.head, t.buf.low, t.buf.high = head.id, low.id, high.id
	t.table[t.buf] = x
}

// struct for operation cache
type tupleNode struct {
	left  Id
	right Id
}

type OpeCache2 struct {
	buf   tupleNode
	table map[tupleNode]*Node
}

func NewOpeCache2() *OpeCache2 {
	return &OpeCache2{
		buf:   tupleNode{0, 0},
		table: make(map[tupleNode]*Node),
	}
}

func (c *OpeCache2) get(f, g *Node) (*Node, bool) {
	c.buf.left, c.buf.right = f.id, g.id
	v, ok := c.table[c.buf]
	return v, ok
}

func (c *OpeCache2) set(f, g, x *Node) {
	c.buf.left, c.buf.right = f.id, g.id
	c.table[c.buf] = x
}

type BDD struct {
	total    Id
	headers  []*NodeHeader
	utable   *UniqueNodeTable
	zero     *Node
	one      *Node
	cacheNot map[Id]*Node
	cacheOr  *OpeCache2
	cacheAnd *OpeCache2
	cacheXor *OpeCache2
}

func NewBDD() *BDD {
	h := &NodeHeader{
		id:    0,
		level: 0,
	}
	zero := &Node{
		id:     0,
		header: h,
	}
	one := &Node{
		id:     1,
		header: h,
	}
	return &BDD{
		total:    2,
		headers:  []*NodeHeader{h},
		utable:   NewUniqueNodeTable(),
		zero:     zero,
		one:      one,
		cacheNot: make(map[Id]*Node),
		cacheOr:  NewOpeCache2(),
		cacheAnd: NewOpeCache2(),
		cacheXor: NewOpeCache2(),
	}
}

func (b *BDD) Var(name string) *Node {
	h := &NodeHeader{
		id:    Id(len(b.headers)),
		level: Level(len(b.headers)),
		label: name,
	}
	b.headers = append(b.headers, h)
	return b.createNode(h, b.zero, b.one)
}

func (b *BDD) GetVarLabels() []string {
	result := make([]string, len(b.headers))
	for i, x := range b.headers {
		result[i] = x.label
	}
	return result
}

func (b *BDD) GetVars() []*Node {
	result := make([]*Node, len(b.headers))
	for i, h := range b.headers {
		result[i] = b.createNode(h, b.zero, b.one)
	}
	return result
}

func (b *BDD) Zero() *Node {
	return b.zero
}

func (b *BDD) One() *Node {
	return b.one
}

func (b *BDD) Not(f *Node) *Node {
	return b.not(f)
}

func (b *BDD) And(xs ...*Node) *Node {
	result := b.one
	for _, x := range xs {
		result = b.and(result, x)
	}
	return result
}

func (b *BDD) Or(xs ...*Node) *Node {
	result := b.zero
	for _, x := range xs {
		result = b.or(result, x)
	}
	return result
}

func (b *BDD) Xor(xs ...*Node) *Node {
	result := b.zero
	for _, x := range xs {
		result = b.xor(result, x)
	}
	return result
}

func (b *BDD) Imp(f, g *Node) *Node {
	return b.Or(b.Not(f), g)
}

func (b *BDD) Ite(f, g, h *Node) *Node {
	return b.Or(b.And(f, g), b.And(b.Not(f), h))
}

func (n *Node) Zero() *Node {
	return n.nodes[0]
}

func (n *Node) One() *Node {
	return n.nodes[1]
}

func (n *Node) Level() Level {
	return n.header.level - 1
}

func (n *Node) Label() string {
	return n.header.label
}

func (b *BDD) createNode(h *NodeHeader, low, high *Node) *Node {
	if low.id == high.id {
		return low
	}
	if node, ok := b.utable.get(h, low, high); ok {
		return node
	}
	node := &Node{
		id:     b.total,
		header: h,
		nodes:  [2]*Node{low, high},
	}
	b.total++
	b.utable.set(h, low, high, node)
	return node
}

func (b *BDD) not(f *Node) *Node {
	switch {
	case f == b.zero:
		return b.one
	case f == b.one:
		return b.zero
	}
	key := f.id
	if node, ok := b.cacheNot[key]; ok {
		return node
	}
	result := b.createNode(f.header, b.not(b.zero), b.not(b.one))
	b.cacheNot[key] = result
	return result
}

func (b *BDD) and(f, g *Node) *Node {
	switch {
	case f == b.zero || g == b.zero:
		return b.zero
	case f == b.one:
		return g
	case g == b.one:
		return f
	}
	if node, ok := b.cacheAnd.get(f, g); ok {
		return node
	}
	var result *Node
	switch {
	case f.header.level > g.header.level:
		result = b.createNode(f.header, b.and(f.Zero(), g), b.and(f.One(), g))
	case f.header.level < g.header.level:
		result = b.createNode(g.header, b.and(f, g.Zero()), b.and(f, g.One()))
	case f.header.level == g.header.level:
		result = b.createNode(f.header, b.and(f.Zero(), g.Zero()), b.and(f.One(), g.One()))
	}
	b.cacheAnd.set(f, g, result)
	return result
}

func (b *BDD) or(f, g *Node) *Node {
	switch {
	case f == b.one || g == b.one:
		return b.one
	case f == b.zero:
		return g
	case g == b.zero:
		return f
	}
	if node, ok := b.cacheOr.get(f, g); ok {
		return node
	}
	var result *Node
	switch {
	case f.header.level > g.header.level:
		result = b.createNode(f.header, b.or(f.Zero(), g), b.or(f.One(), g))
	case f.header.level < g.header.level:
		result = b.createNode(g.header, b.or(f, g.Zero()), b.or(f, g.One()))
	case f.header.level == g.header.level:
		result = b.createNode(f.header, b.or(f.Zero(), g.Zero()), b.or(f.One(), g.One()))
	}
	b.cacheOr.set(f, g, result)
	return result
}

func (b *BDD) xor(f, g *Node) *Node {
	switch {
	case f == b.one:
		return b.not(g)
	case g == b.one:
		return b.not(f)
	case f == b.zero:
		return g
	case g == b.zero:
		return f
	}
	if node, ok := b.cacheXor.get(f, g); ok {
		return node
	}
	var result *Node
	switch {
	case f.header.level > g.header.level:
		result = b.createNode(f.header, b.xor(f.Zero(), g), b.xor(f.One(), g))
	case f.header.level < g.header.level:
		result = b.createNode(g.header, b.xor(f, g.Zero()), b.xor(f, g.One()))
	case f.header.level == g.header.level:
		result = b.createNode(f.header, b.xor(f.Zero(), g.Zero()), b.xor(f.One(), g.One()))
	}
	b.cacheXor.set(f, g, result)
	return result
}

// todot

func (b *BDD) ToDot(f *Node) {
	fmt.Println("digraph { layout=dot; overlap=false; splines=true; node [fontsize=10];")
	visited := make(map[Id]struct{})
	b.todot(f, visited)
	fmt.Println("}")
}

func (b *BDD) todot(f *Node, visited map[Id]struct{}) {
	key := f.id
	if _, ok := visited[key]; ok {
		return
	}
	switch {
	case f == b.zero:
		fmt.Printf(" n%p [shape = square, label = \"0\"];\n", f)
	case f == b.one:
		fmt.Printf(" n%p [shape = square, label = \"1\"];\n", f)
	default:
		fmt.Printf(" n%p [shape = circle, label = \"%s\"];\n", f, f.header.label)
		for i, x := range f.nodes {
			b.todot(x, visited)
			fmt.Printf(" n%p -> n%p [label = \"%d\"];\n", f, x, i)
		}
	}
	visited[key] = struct{}{}
}
