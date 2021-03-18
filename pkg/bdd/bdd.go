package bdd

import (
	"fmt"
)

type Level uint64

type NodeHeader struct {
	level Level  // Level of node
	label string // Label of node
}

type Node struct {
	header *NodeHeader // info of Node
	nodes  [2]*Node    // the pointer for Nodes
}

// struct to identify a node
type nodeId struct {
	head *NodeHeader
	low  *Node
	high *Node
}

// struct for operation cache
type tupleNode struct {
	left  *Node
	right *Node
}

type BDD struct {
	headers  []*NodeHeader
	nodes    map[nodeId]*Node
	zero     *Node
	one      *Node
	cacheNot map[*Node]*Node
	cacheOr  map[tupleNode]*Node
	cacheAnd map[tupleNode]*Node
	cacheXor map[tupleNode]*Node
}

func NewBDD() *BDD {
	h := &NodeHeader{
		level: 0,
	}
	zero := &Node{
		header: h,
	}
	one := &Node{
		header: h,
	}
	return &BDD{
		headers:  []*NodeHeader{h},
		nodes:    make(map[nodeId]*Node),
		zero:     zero,
		one:      one,
		cacheNot: make(map[*Node]*Node),
		cacheOr:  make(map[tupleNode]*Node),
		cacheAnd: make(map[tupleNode]*Node),
		cacheXor: make(map[tupleNode]*Node),
	}
}

func (b *BDD) Var(name string) *Node {
	h := &NodeHeader{
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
	if low == high {
		return low
	}
	key := nodeId{h, low, high}
	if node, ok := b.nodes[key]; ok {
		return node
	}
	node := &Node{
		header: h,
		nodes:  [2]*Node{low, high},
	}
	b.nodes[key] = node
	return node
}

func (b *BDD) not(f *Node) *Node {
	switch {
	case f == b.zero:
		return b.one
	case f == b.one:
		return b.zero
	}
	if node, ok := b.cacheNot[f]; ok {
		return node
	}
	result := b.createNode(f.header, b.not(f.Zero()), b.not(f.One()))
	b.cacheNot[f] = result
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
	key := tupleNode{f, g}
	if node, ok := b.cacheAnd[key]; ok {
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
	b.cacheAnd[key] = result
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
	key := tupleNode{f, g}
	if node, ok := b.cacheOr[key]; ok {
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
	b.cacheOr[key] = result
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
	key := tupleNode{f, g}
	if node, ok := b.cacheXor[key]; ok {
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
	b.cacheXor[key] = result
	return result
}

// todot

func (b *BDD) ToDot(f *Node) {
	fmt.Println("digraph { layout=dot; overlap=false; splines=true; node [fontsize=10];")
	visited := make(map[*Node]struct{})
	b.todot(f, visited)
	fmt.Println("}")
}

func (b *BDD) todot(f *Node, visited map[*Node]struct{}) {
	if _, ok := visited[f]; ok {
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
	visited[f] = struct{}{}
}
