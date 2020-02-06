package persistentmap

type Key interface {
	Comparable
	Hashable
}

type Value = interface{}

type Comparable interface {
	Equal(Comparable) bool
}

type Hashable interface {
	Hash() int32
}

type persistentMap struct {
	rootLeaf node
}

func (m persistentMap) Get(key Key) (interface{}, bool) {
	if m.rootLeaf == nil {
		return nil, false
	}
	return m.rootLeaf.find(key, key.Hash())
}

func (m persistentMap) Set(key Key, value interface{}) persistentMap {
	if m.rootLeaf == nil {
		return persistentMap{rootLeaf: &leafNode{
			baseNode: baseNode{level: 0},
			key:      key,
			hash:     key.Hash(),
			value:    value,
		}}
	}
	return persistentMap{
		rootLeaf: m.rootLeaf.assoc(key, key.Hash(), value),
	}
}
