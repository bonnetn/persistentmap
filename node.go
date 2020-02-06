package persistentmap

import (
	"math"
	"math/bits"
)

const (
	chunkSize     = 5
	childrenCount = 1 << chunkSize
	mask          = childrenCount - 1
)

var maxDepth = int32(math.Ceil(32 / chunkSize))

type node interface {
	find(key Key, hash int32) (interface{}, bool)
	assoc(key Key, hash int32, value Value) node
}

//////////////////////////////////////////////////////

type baseNode struct {
	level int32
}

//////////////////////////////////////////////////////
type emptyNode struct{ baseNode }

func (l *emptyNode) assoc(key Key, hash int32, value Value) node {
	return &leafNode{
		baseNode: l.baseNode,
		key:      key,
		value:    value,
		hash:     hash,
	}
}

func (*emptyNode) find(Key, int32) (interface{}, bool) { return nil, false }

//////////////////////////////////////////////////////

type leafNode struct {
	baseNode
	key   Key
	hash  int32
	value Value
}

func (l *leafNode) assoc(key Key, hash int32, value Value) node {
	if l.level > maxDepth {
		return &arrayLeafNode{
			baseNode: l.baseNode,
			children: []*leafNode{
				{
					baseNode: baseNode{level: l.level},
					key:      key,
					hash:     hash,
					value:    value,
				},
				l,
			},
		}
	}
	node := &bitmapIndexedNode{
		baseNode: l.baseNode,
		children: nil,
		bitmap:   0,
	}
	return node.
		assoc(key, hash, value).
		assoc(l.key, l.hash, l.value)
}

func (l *leafNode) find(key Key, hash int32) (interface{}, bool) {
	if l.hash != hash {
		return nil, false
	}
	if !l.key.Equal(key) {
		return nil, false
	}
	return l.value, true
}

//////////////////////////////////////////////////////
type bitmapIndexedNode struct {
	baseNode
	children []node
	bitmap   int
}

func (b *bitmapIndexedNode) find(key Key, hash int32) (interface{}, bool) {
	bitPos := bitpos(hash, b.level)
	if b.bitmap&bitPos != 0 {
		index := computeIndex(b.bitmap, bitPos)
		return b.children[index].find(key, hash)
	} else {
		return nil, false
	}
}

func (b *bitmapIndexedNode) assoc(key Key, hash int32, value Value) node {
	bitPos := bitpos(hash, b.level)
	index := computeIndex(b.bitmap, bitPos)
	if b.bitmap&bitPos != 0 {
		// Already a node with the same mask value, assoc to it.
		children := make([]node, 0, len(b.children))
		for i := 0; i < len(b.children); i++ {
			child := b.children[i]
			if i == index {
				child = child.assoc(key, hash, value)
			}
			children = append(children, child)
		}
		return &bitmapIndexedNode{
			baseNode: baseNode{level: b.level},
			children: children,
			bitmap:   b.bitmap | bitPos,
		}
	} else {
		// No value in the map, adding one.
		children := make([]node, len(b.children)+1)
		children[index] = &leafNode{
			baseNode: baseNode{level: b.level + 1},
			key:      key,
			hash:     hash,
			value:    value,
		}
		for i := 0; i < len(b.children); i++ {
			if i >= index {
				children[i+1] = b.children[i]
			} else {
				children[i] = b.children[i]
			}
		}
		return &bitmapIndexedNode{
			baseNode: baseNode{level: b.level},
			children: children,
			bitmap:   b.bitmap | bitPos,
		}
	}
}

func maskFunc(hash int32, shift int32) int32 {
	return (hash << (shift * chunkSize)) & mask
}

func bitpos(hash int32, level int32) int {
	return 1 << maskFunc(hash, level)
}

func computeIndex(bitmap, bitPos int) int {
	return bits.OnesCount32(uint32(bitmap & (bitPos - 1)))
}

//////////////////////////////////////////////////////

type arrayLeafNode struct {
	baseNode
	children []*leafNode
}

func (l *arrayLeafNode) assoc(key Key, hash int32, value Value) node {
	children := make([]*leafNode, 0, len(l.children))
	children = append(children, &leafNode{
		baseNode: baseNode{
			level: l.level + 1,
		},
		key:   key,
		hash:  hash,
		value: value,
	})
	for _, v := range l.children {
		children = append(children, v)
	}

	return &arrayLeafNode{
		baseNode: l.baseNode,
		children: children,
	}
}

func (l *arrayLeafNode) find(key Key, hash int32) (interface{}, bool) {
	for _, v := range l.children {
		if v.key.Equal(key) {
			return v.value, true
		}
	}
	return nil, false
}

//////////////////////////////////////////////////////
var (
	_ node = &emptyNode{}
	_ node = &leafNode{}
	_ node = &arrayLeafNode{}
	_ node = &bitmapIndexedNode{}
)
