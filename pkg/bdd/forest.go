package bdd

// import (
// // 	"fmt"
// )

type UniOperator interface {
	UniqueTableInterface
	checkTerminal(node *Node) (*Node, bool)
}

type BinOperator interface {
	UniqueTableInterface
	checkTerminal(left, right *Node) (*Node, bool)
}

// type Forest interface {
// 	Node(header *NodeHeader, nodes ...*Node) *Node
// }

// func UniApply(f Forest, op UniOperator, node *Node) *Node {
// 	if node, ok := op.checkTerminal(node); ok {
// 		return node
// 	}
// 	key := op.GenKey(node)
// 	if node, ok := op.Get(key); ok {
// 		return node
// 	}

// 	n := len(node.header.domain)
// 	nodes := make([]*Node, n, n)
// 	for i := 0; i < n; i++ {
// 		nodes[i] = UniApply(f, op, node.nodes[i])
// 	}
// 	result := f.Node(node.header, nodes...)
// 	op.Set(key, result)
// 	return result
// }

// func BinApply(f Forest, op BinOperator, left, right *Node) *Node {
// 	if node, ok := op.checkTerminal(left, right); ok {
// 		return node
// 	}
// 	key := op.GenKey(left, right)
// 	if node, ok := op.Get(key); ok {
// 		return node
// 	}

// 	var result *Node
// 	switch {
// 	case left.header.level > right.header.level:
// 		n := len(left.header.domain)
// 		nodes := make([]*Node, n, n)
// 		for i := 0; i < n; i++ {
// 			nodes[i] = BinApply(f, op, left.nodes[i], right)
// 		}
// 		result = f.Node(left.header, nodes...)
// 	case left.header.level < right.header.level:
// 		n := len(right.header.domain)
// 		nodes := make([]*Node, n, n)
// 		for i := 0; i < n; i++ {
// 			nodes[i] = BinApply(f, op, left, right.nodes[i])
// 		}
// 		result = f.Node(right.header, nodes...)
// 	case left.header.level == right.header.level:
// 		n := len(left.header.domain)
// 		nodes := make([]*Node, n, n)
// 		for i := 0; i < n; i++ {
// 			nodes[i] = BinApply(f, op, left.nodes[i], right.nodes[i])
// 		}
// 		result = f.Node(left.header, nodes...)
// 	}
// 	op.Set(key, result)
// 	return result
// }
