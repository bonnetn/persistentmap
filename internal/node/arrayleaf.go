package node

import (
	"github.com/bonnetn/persistentmap/internal"
	"github.com/bonnetn/persistentmap"
)

type arrayLeafNode struct {
	internal.NodeMetadata

	children []*leafNode
}

var _ Node = &arrayLeafNode{}

func (l *arrayLeafNode) Assoc(key persistentmap.Key, hash int32, value persistentmap.Value) Node {
	if _, ok := l.Find(key, hash); ok {
		children := make([]*leafNode, 0, len(l.children))
		for _, v := range l.children {
			child := v
			if v.hash == hash && v.key.Equal(key) {
				child = &leafNode{
					NodeMetadata: v.NodeMetadata,
					key:          key,
					hash:         hash,
					value:        value,
				}
			}
			children = append(children, child)
		}
		return &arrayLeafNode{
			NodeMetadata: l.NodeMetadata,
			children:     children,
		}
	}
	children := make([]*leafNode, 0, len(l.children))
	for _, v := range l.children {
		children = append(children, v)
	}
	children = append(children, &leafNode{
		NodeMetadata: l.Child(),
		key:          key,
		hash:         hash,
		value:        value,
	})

	return &arrayLeafNode{
		NodeMetadata: l.NodeMetadata,
		children:     children,
	}
}

func (l *arrayLeafNode) Find(key persistentmap.Key, hash int32) (interface{}, bool) {
	for _, v := range l.children {
		if v.key.Equal(key) {
			return v.value, true
		}
	}
	return nil, false
}
