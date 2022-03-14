package dispatch

import (
	"errors"
	"fmt"
	"go_im/pkg/murmur"
	"strconv"
)

const (
	duplicateVirtual = 100 // 1_000_000
	seed             = 0xabcd1234
)

type Uint32 interface {
	Val() uint32
}

type Node struct {
	val     string
	hash    uint32
	virtual bool
	real    *Node
}

func (n *Node) Val() uint32 {
	return n.hash
}
func (n *Node) String() string {
	return strconv.FormatInt(int64(n.hash), 10)
}

type Nodes struct {
	nd      Node
	virtual []Node
	hit     int64
}

func (n *Nodes) appendVirtual(node Node) {
	n.virtual = append(n.virtual, node)
}

type ConsistentHash struct {
	nodes   []Node
	nodeMap map[string]*Nodes
}

func NewConsistentHash() *ConsistentHash {
	hash := &ConsistentHash{
		nodes:   []Node{},
		nodeMap: map[string]*Nodes{},
	}
	return hash
}

func (c *ConsistentHash) Remove(id string) error {
	nodes, ok := c.nodeMap[id]
	if !ok {
		return errors.New("node does not exist, id:" + id)
	}
	for _, vNd := range nodes.virtual {
		ndIndex, exist := c.findIndex(vNd.hash)
		if exist {
			ndIndex--
		} else {
			return errors.New("virtual node does not exist, id:" + vNd.val)
		}
		nd := c.nodes[ndIndex]
		if nd.hash != vNd.hash {
			return errors.New("could not find virtual node, id:" + vNd.val)
		} else {
			c.removeIndex(ndIndex)
		}
	}
	index, exist := c.findIndex(nodes.nd.hash)
	if !exist {
		return errors.New("real node not fund")
	}
	index--
	c.removeIndex(index)
	delete(c.nodeMap, id)
	return nil
}

func (c *ConsistentHash) Get(data string) *Node {
	hash := murmur.Hash([]byte(data), seed)
	index, _ := c.findIndex(hash)
	if index == len(c.nodes) {
		index = len(c.nodes) - 1
	}
	n := c.nodes[index]
	if n.virtual {
		return n.real
	}
	return &n
}

func (c *ConsistentHash) addVirtual(real *Node, duplicate int) {
	for i := 0; i < duplicate; i++ {
		vNodeID := fmt.Sprintf("%s_#%d", real.val, i)
		hash := murmur.Hash([]byte(vNodeID), seed)
		vNode := Node{
			val:     vNodeID,
			hash:    hash,
			virtual: true,
			real:    real,
		}
		c.addNode(vNode)
		nds := c.nodeMap[real.val]
		nds.appendVirtual(vNode)
	}
}

func (c *ConsistentHash) addNode(nd Node) {
	index, _ := c.findIndex(nd.hash)
	p1 := c.nodes[:index]
	p2 := c.nodes[index:]
	n := make([]Node, len(p1))
	copy(n, p1)
	n = append(n, nd)
	for _, i := range p2 {
		n = append(n, i)
	}
	c.nodes = n
}

func (c *ConsistentHash) Add(id string) {
	_, ok := c.nodeMap[id]
	if ok {
		// exist
	}
	hash := murmur.Hash([]byte(id), seed)
	nd := Node{
		val:     id,
		hash:    hash,
		virtual: false,
		real:    nil,
	}
	c.nodeMap[id] = &Nodes{
		nd:      nd,
		virtual: []Node{},
	}
	c.addNode(nd)
	c.addVirtual(&nd, duplicateVirtual)
}

func (c *ConsistentHash) removeHash(hash uint32) {

}

func (c *ConsistentHash) removeIndex(index int) {
	if index == len(c.nodes)-1 {
		c.nodes = c.nodes[:len(c.nodes)-1]
		return
	}
	p2 := c.nodes[index+1:]
	c.nodes = c.nodes[:index]
	for _, n := range p2 {
		c.nodes = append(c.nodes, n)
	}
}

func (c *ConsistentHash) findIndex(s uint32) (int, bool) {
	left := 0
	right := len(c.nodes)
	exist := false
LOOP:
	if left < right {
		middle := (left + right) / 2
		hash := c.nodes[middle].hash
		if hash < s {
			left = middle + 1
		} else if hash == s {
			left = middle + 1
			exist = true
		} else {
			right = middle
		}
		goto LOOP
	}
	return left, exist
}
