package chord

import (
    "math/big"
)

type Finger struct {
    Start *big.Int
    Node  *Node
}

func (n *Node) FixFingers() {
    for i := 0; i < fingerSize; i++ {
        successor, err := n.FindSuccessor(n.Finger[i].Start)
        if err != nil {
            continue
        }
        n.Finger[i].Node = successor
    }
}

func (n *Node) UpdateFingerTable(s *Node, i int) bool {
    if between(n.ID, s.ID, n.Finger[i].Node.ID, false) {
        n.Finger[i].Node = s
        p := n.Predecessor
        if p != nil {
            p.UpdateFingerTable(s, i)
        }
        return true
    }
    return false
}