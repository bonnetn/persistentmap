package node

import (
	"github.com/bonnetn/persistentmap"
)

type Node interface {
	Find(key persistentmap.Key, hash int32) (interface{}, bool)
	Assoc(key persistentmap.Key, hash int32, value persistentmap.Value) Node
}

func FromKeyValue(key persistentmap.Key, value persistentmap.Value) Node {
	return &leafNode{
		key:   key,
		hash:  key.Hash(),
		value: value,
	}
}
