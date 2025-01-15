package dag

import (
	"sync"
)

type key string

type Node struct {
	Parents  map[key]int8
	Children map[key]int8
	Weight   float32
}

type DAG struct {
	nodes map[key]*Node
	lock  *sync.RWMutex
}

func NewDAG() *DAG {
	g := new(DAG)
	g.nodes = make(map[key]*Node, 0)
	g.lock = new(sync.RWMutex)

	return g
}

func (g *DAG) Node(k string) *Node {
	g.lock.RLock()
	node := g.nodes[key(k)]
	g.lock.RUnlock()

	return node
}

func (g *DAG) Nodes() map[key]*Node {
	defer g.lock.RUnlock()
	g.lock.RLock()
	nodes := g.nodes

	return nodes
}

func (g *DAG) UpsertNode(k, parent string, weight ...float32) {
	// Upsert node.
	g.lock.Lock()
	if g.nodes[key(k)] == nil {
		g.nodes[key(k)] = new(Node)
	}

	if g.nodes[key(k)].Parents == nil {
		g.nodes[key(k)].Parents = make(map[key]int8, 0)
	}

	g.nodes[key(k)].Parents[key(parent)]++
	if len(weight) > 0 {
		g.nodes[key(k)].Weight += weight[0]
	}
	g.lock.Unlock()

	// Update parent's children.
	if g.Node(parent) == nil {
		g.lock.Lock()
		g.nodes[key(parent)] = new(Node)
		g.lock.Unlock()
	}
	if g.Node(parent).Children == nil {
		g.lock.Lock()
		g.nodes[key(parent)].Children = make(map[key]int8, 0)
		g.lock.Unlock()
	}

	g.lock.Lock()
	g.nodes[key(parent)].Children[key(k)]++
	g.lock.Unlock()
}
