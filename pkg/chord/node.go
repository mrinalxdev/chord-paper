package chord

import (
    "crypto/sha1"
    "fmt"
    "math/big"
    "sync"
    // "time"
)

const (
    m            = 160 // SHA1 produces 160-bit identifiers
    fingerSize   = m
    successorNum = 3
)

type Node struct {
    ID         *big.Int
    Address    string
    Successors []*Node
    Finger     []*Finger
    Predecessor *Node
    Storage    Storage
    mutex      sync.RWMutex
}

type Storage interface {
    Put(key []byte, value []byte) error
    Get(key []byte) ([]byte, error)
    Delete(key []byte) error
}

func NewNode(address string, storage Storage) *Node {
    id := generateID(address)
    node := &Node{
        ID:         id,
        Address:    address,
        Successors: make([]*Node, successorNum),
        Finger:     make([]*Finger, fingerSize),
        Storage:    storage,
    }
    
    for i := 0; i < fingerSize; i++ {
        node.Finger[i] = &Finger{
            Start: calculateStart(node.ID, i),
        }
    }
    
    return node
}

func (n *Node) Join(introducerNode *Node) error {
    if introducerNode == nil {
        // First node in the network
        for i := 0; i < successorNum; i++ {
            n.Successors[i] = n
        }
        n.Predecessor = n
        return nil
    }

    successor, err := introducerNode.FindSuccessor(n.ID)
    if err != nil {
        return fmt.Errorf("failed to find successor: %v", err)
    }

    n.Successors[0] = successor
    return n.Stabilize()
}

func (n *Node) FindSuccessor(id *big.Int) (*Node, error) {
    if id.Cmp(n.ID) == 0 {
        return n, nil
    }

    successor := n.GetSuccessor()
    if between(n.ID, id, successor.ID, true) {
        return successor, nil
    }

    node := n.closestPrecedingNode(id)
    if node.ID.Cmp(n.ID) == 0 {
        return successor, nil
    }
    
    return node.FindSuccessor(id)
}

func (n *Node) Stabilize() error {
    successor := n.GetSuccessor()
    if successor == nil {
        return fmt.Errorf("no successor found")
    }

    x := successor.Predecessor
    if x != nil && between(n.ID, x.ID, successor.ID, false) {
        n.Successors[0] = x
        successor = x
    }

    successor.Notify(n)
    return nil
}

func (n *Node) Notify(node *Node) {
    n.mutex.Lock()
    defer n.mutex.Unlock()

    if n.Predecessor == nil || between(n.Predecessor.ID, node.ID, n.ID, false) {
        n.Predecessor = node
    }
}

func (n *Node) GetSuccessor() *Node {
    n.mutex.RLock()
    defer n.mutex.RUnlock()
    return n.Successors[0]
}

func (n *Node) closestPrecedingNode(id *big.Int) *Node {
    for i := fingerSize - 1; i >= 0; i-- {
        finger := n.Finger[i].Node
        if finger != nil && between(n.ID, finger.ID, id, false) {
            return finger
        }
    }
    return n
}

func generateID(address string) *big.Int {
    h := sha1.New()
    h.Write([]byte(address))
    return new(big.Int).SetBytes(h.Sum(nil))
}

func calculateStart(id *big.Int, i int) *big.Int {
    two := big.NewInt(2)
    start := new(big.Int).Exp(two, big.NewInt(int64(i)), nil)
    start.Add(id, start)
    max := new(big.Int).Exp(two, big.NewInt(m), nil)
    return start.Mod(start, max)
}