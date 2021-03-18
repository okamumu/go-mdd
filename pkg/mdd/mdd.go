package bdd

// import (
// // 	"fmt"
// )

// type MtMDDIntTerminalValue int

// type MtMDDInt struct {
// 	*NodeManager
// 	*UniqueTable
// 	values    map[*Node]MtMDDIntTerminalValue
// 	terminals map[MtMDDIntTerminalValue]*Node
// }

// func (b *MtMDDInt) Node(header *NodeHeader, nodes ...*Node) *Node {
// 	key := b.GenKey(nodes...)
// 	if node, ok := b.Get(key); ok {
// 		return node
// 	}
// 	node := b.NewNode(header, nodes)
// 	b.Set(key, node)
// 	return node
// }

// func (b *MtMDDInt) Terminal(header *NodeHeader, x MtMDDIntTerminalValue) *Node {
// 	if n, ok := b.terminals[x]; ok {
// 		return n
// 	} else {
// 		n = b.NewNode(header, nil)
// 		b.terminals[x] = n
// 		return n
// 	}
// }

// func (b *MtMDDInt) GetValue(n *Node) MtMDDIntTerminalValue {
// 	if v, ok := b.values[n]; ok {
// 		return v
// 	} else {
// 		panic("!")
// 	}
// }

// type PlusMtMDDInt struct {
// 	*UniqueTable
// 	f *MtMDDInt
// }

// func (op *PlusMtMDDInt) checkTerminal(left, right *Node) (*Node, bool) {
// 	switch {
// 	case left.isTerminal() && right.isTerminal():
// 		return op.f.Terminal(left.header, op.f.GetValue(left)+op.f.GetValue(right)), true
// 	case left.isTerminal() && op.f.GetValue(left) == 0:
// 		return right, true
// 	case right.isTerminal() && op.f.GetValue(right) == 0:
// 		return left, true
// 	}
// 	return nil, false
// }

// type MultipleMtMDDInt struct {
// 	*UniqueTable
// 	f *MtMDDInt
// }

// func (op *MultipleMtMDDInt) checkTerminal(left, right *Node) (*Node, bool) {
// 	switch {
// 	case left.isTerminal() && op.f.GetValue(left) == 0:
// 		return op.f.Terminal(left.header, 0), true
// 	case right.isTerminal() && op.f.GetValue(right) == 0:
// 		return op.f.Terminal(right.header, 0), true
// 	case left.isTerminal() && op.f.GetValue(left) == 1:
// 		return right, true
// 	case right.isTerminal() && op.f.GetValue(right) == 1:
// 		return left, true
// 	case left.isTerminal() && right.isTerminal():
// 		return op.f.Terminal(left.header, op.f.GetValue(left)*op.f.GetValue(right)), true
// 	}
// 	return nil, false
// }
