package node

import (
	"github.com/bonnetn/persistentmap"
	"github.com/bonnetn/persistentmap/internal"
	"math"
)

var maxTrieLevel = int32(math.Ceil(32 / bitmapSize))

type leafNode struct {
	internal.NodeMetadata

	key   persistentmap.Key
	hash  int32
	value persistentmap.Value
}

func (l *leafNode) Assoc(key persistentmap.Key, hash int32, value persistentmap.Value) Node {
	// If the key is the same, replace this node.
	if (l.hash == hash && l.key.Equal(key)) || l.key == nil {
		return &leafNode{
			NodeMetadata: l.NodeMetadata,

			key:   key,
			hash:  hash,
			value: value,
		}
	}

	if l.Level() <= maxTrieLevel {
		// If we can store it in the trie, use a bitmap indexed node.
		return (&bitmapIndexedNode{NodeMetadata: l.NodeMetadata,}).
			Assoc(l.key, l.hash, l.value).
			Assoc(key, hash, value)
	}

	// If the trie max level is exceeded, replace by a simple array of nodes.
	return &arrayLeafNode{
		NodeMetadata: l.NodeMetadata,
		children: []*leafNode{
			{
				NodeMetadata: l.Child(),
				key:          key,
				hash:         hash,
				value:        value,
			},
			l,
		},
	}
}

func (l *leafNode) Find(key persistentmap.Key, hash int32) (interface{}, bool) {
	if l.key == nil {
		return nil, false
	}

	if l.hash != hash {
		return nil, false
	}
	if !l.key.Equal(key) {
		return nil, false
	}
	return l.value, true
}
