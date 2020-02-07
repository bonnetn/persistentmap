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
