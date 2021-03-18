package bdd

import (
	"strconv"
)

type Level uint64
type Id uint64

type NodeInterface interface {
	GetId() Id
}

// The structure of NodeHeader is shared with Node
type NodeHeader struct {
	id    Id     // Id
	level Level  // Level of node
	label string // Label of node
}

func (n *NodeHeader) GetId() Id {
	return n.id
}

type Node struct {
	header *NodeHeader // info of Node
	id     Id          // Id
	nodes  []*Node     // the pointer for Nodes
}

func (n *Node) IsTerminal() bool {
	return n.header.level == 0
}

func (n *Node) GetId() Id {
	return n.id
}

type NodeManager struct {
	nodeid  Id
	headers []*NodeHeader
}

func NewNodeManager() *NodeManager {
	return &NodeManager{
		nodeid:  Id(0),
		headers: make([]*NodeHeader, 0),
	}
}

func (m *NodeManager) NewHeader(label string) *NodeHeader {
	h := &NodeHeader{
		id:    Id(len(m.headers)),
		level: Level(len(m.headers)),
		label: label,
	}
	m.headers = append(m.headers, h)
	return h
}

func (m *NodeManager) NewNode(h *NodeHeader, nodes ...*Node) *Node {
	node := &Node{
		header: h,
		id:     m.nodeid,
		nodes:  nodes,
	}
	m.nodeid++
	return node
}

type UniqueTable struct {
	buf  []byte
	data map[string]NodeInterface
}

func NewUniqueTable() *UniqueTable {
	return &UniqueTable{
		buf:  make([]byte, 0, 5),
		data: make(map[string]NodeInterface),
	}
}

func (t *UniqueTable) GenKey(xs ...NodeInterface) string {
	t.buf = t.buf[:0]
	for _, x := range xs {
		t.buf = append(t.buf, strconv.FormatUint(uint64(x.GetId()), 16)...)
		t.buf = append(t.buf, ',')
	}
	return string(t.buf)
}

func (t *UniqueTable) Get(key string) (NodeInterface, bool) {
	if node, ok := t.data[key]; ok {
		return node, true
	} else {
		return nil, false
	}
}

func (t *UniqueTable) Set(key string, node NodeInterface) {
	t.data[key] = node
}
