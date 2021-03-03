package bdd

import (
	"fmt"
	"strconv"
)

type Level uint
type NodeId int64
type DomainInt int

// The structure of NodeHeader is shared with Node
type NodeHeader struct {
	// 	forest Forest // Forest
	level  Level       // Level of node
	label  string      // Label of node
	domain []DomainInt // Domain
}

type Node struct {
	header *NodeHeader // info of Node
	id     NodeId      // Id
	nodes  []*Node     // the pointer for Nodes
}

func (n *Node) isTerminal() bool {
	return n.header.level == 0
}

type NodeManager struct {
	id      NodeId
	headers map[string]*NodeHeader
}

func NewNodeManager(start_id NodeId) *NodeManager {
	return &NodeManager{
		id:      start_id,
		headers: make(map[string]*NodeHeader),
	}
}

func (ng *NodeManager) NewHeader(label string, level Level, domain []DomainInt) *NodeHeader {
	h := &NodeHeader{
		level:  level,
		label:  label,
		domain: domain,
	}
	ng.headers[label] = h
	return h
}

func (ng *NodeManager) NewNode(h *NodeHeader, nodes []*Node) *Node {
	node := &Node{
		header: h,
		id:     ng.id,
		nodes:  nodes,
	}
	ng.id++
	return node
}

// UniqueTable

type UniqueTableInterface interface {
	GenKey(nodes ...*Node) string
	Get(key string) (*Node, bool)
	Set(key string, node *Node)
}

type UniqueTable struct {
	key  []byte
	data map[string]*Node
}

func NewUniqueTable() *UniqueTable {
	return &UniqueTable{
		key:  make([]byte, 0, 5),
		data: make(map[string]*Node),
	}
}

func (t *UniqueTable) GenKey(nodes ...*Node) string {
	t.key = t.key[:0]
	for _, x := range nodes {
		t.key = append(t.key, strconv.Itoa(int(x.id))...)
		t.key = append(t.key, ',')
	}
	return string(t.key)
}

func (t *UniqueTable) Get(key string) (*Node, bool) {
	if node, ok := t.data[key]; ok {
		return node, true
	} else {
		return nil, false
	}
}

func (t *UniqueTable) Set(key string, node *Node) {
	t.data[key] = node
}

// todot

func (f *Node) ToDot(conv func(x *Node) string) {
	visited := make(map[*Node]struct{})
	fmt.Println("digraph { layout=dot; overlap=false; splines=true; node [fontsize=10];")
	f.todot(visited, conv)
	fmt.Println("}")
}

func (f *Node) todot(visited map[*Node]struct{}, conv func(x *Node) string) {
	if _, ok := visited[f]; ok {
		return
	}
	if f.isTerminal() {
		fmt.Printf(" n%d [shape = square, label = \"%s\"];\n", f.id, conv(f))
	} else {
		fmt.Printf(" n%d [shape = circle, label = \"%s\"];\n", f.id, f.header.label)
		for i, x := range f.nodes {
			x.todot(visited, conv)
			fmt.Printf(" n%d -> n%d [label = \"%d\"];\n", f.id, x.id, i)
		}
	}
	visited[f] = struct{}{}
}
