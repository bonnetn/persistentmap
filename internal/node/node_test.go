package node

import (
	"testing"
	)

func TestFromKeyValue(t *testing.T) {
	k := key{1, 2}
	node := FromKeyValue(k, 3)
	leaf, ok := node.(*leafNode)
	if !ok {
		t.Error("should create a leaf node")
	}

	got, ok := leaf.key.(key)
	if !ok {
		t.Error("leaf.key is of the wrong type")
	}
	if got != k {
		t.Errorf("leaf.key is %v, should be %v", got, k)
	}
	if leaf.hash != k.Hash() {
		t.Errorf("leaf.hash is %v, should be %v", leaf.hash, k.Hash())
	}
	if leaf.value != 3 {
		t.Errorf("leaf.value is %v, should be %v", leaf.value, 3)
	}
}
