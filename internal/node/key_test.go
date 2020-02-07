package node

import "github.com/bonnetn/persistentmap"

type key struct{ value, hash int32 }

var _ persistentmap.Key = key{}

func (k key) Equal(c persistentmap.Comparable) bool {
	if got, ok := c.(key); ok && got.value == k.value {
		return true
	}
	return false
}

func (k key) Hash() int32 { return k.hash }
