package node

import (
	"github.com/bonnetn/persistentmap/internal"
	"github.com/bonnetn/persistentmap"
	"math/bits"
)

const (
	bitmapSize    = 5
	childrenCount = 1 << bitmapSize
	mask          = childrenCount - 1
)

type bitmapIndexedNode struct {
	internal.NodeMetadata
	children []Node
	bitmap   int
}

var _ Node = &bitmapIndexedNode{}

func (b *bitmapIndexedNode) Find(key persistentmap.Key, hash int32) (interface{}, bool) {
	bitPos := bitpos(hash, b.Level())
	if b.bitmap&bitPos != 0 {
		index := computeIndex(b.bitmap, bitPos)
		return b.children[index].Find(key, hash)
	} else {
		return nil, false
	}
}

func (b *bitmapIndexedNode) Assoc(key persistentmap.Key, hash int32, value persistentmap.Value) Node {
	bitPos := bitpos(hash, b.Level())
	index := computeIndex(b.bitmap, bitPos)
	if b.bitmap&bitPos != 0 {
		// Already a Node with the same mask value, Assoc to it.
		children := make([]Node, 0, len(b.children))
		for i := 0; i < len(b.children); i++ {
			child := b.children[i]
			if i == index {
				child = child.Assoc(key, hash, value)
			}
			children = append(children, child)
		}
		return &bitmapIndexedNode{
			NodeMetadata: b.NodeMetadata,
			children:     children,
			bitmap:       b.bitmap | bitPos,
		}
	} else {
		// No value in the map, adding one.
		children := make([]Node, len(b.children)+1)
		children[index] = &leafNode{
			NodeMetadata: b.Child(),
			key:          key,
			hash:         hash,
			value:        value,
		}
		for i := 0; i < len(b.children); i++ {
			if i >= index {
				children[i+1] = b.children[i]
			} else {
				children[i] = b.children[i]
			}
		}
		return &bitmapIndexedNode{
			NodeMetadata: b.NodeMetadata,
			children:     children,
			bitmap:       b.bitmap | bitPos,
		}
	}
}

func maskFunc(hash int32, shift int32) int32 {
	return (hash << (shift * bitmapSize)) & mask
}

func bitpos(hash int32, level int32) int {
	return 1 << maskFunc(hash, level)
}

func computeIndex(bitmap, bitPos int) int {
	return bits.OnesCount32(uint32(bitmap & (bitPos - 1)))
}
