package persistent

import (
	"github.com/bonnetn/persistentmap"
		"github.com/bonnetn/persistentmap/internal/node"
)

type persistentMap struct {
	rootLeaf node.Node
}

func (m persistentMap) Get(key persistentmap.Key) (interface{}, bool) {
	if m.rootLeaf == nil {
		return nil, false
	}
	return m.rootLeaf.Find(key, key.Hash())
}

func (m persistentMap) Set(key persistentmap.Key, value interface{}) persistentMap {
	if m.rootLeaf == nil {
		return persistentMap{rootLeaf: node.FromKeyValue(key,value)}
	}
	return persistentMap{
		rootLeaf: m.rootLeaf.Assoc(key, key.Hash(), value),
	}
}
