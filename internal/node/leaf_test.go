package node

import (
	"testing"
	"github.com/bonnetn/persistentmap/internal"
	"github.com/bonnetn/persistentmap"
)

// TestFindZeroLeafNode tests the Find function.
func TestFindZeroLeafNode(t *testing.T) {
	var (
		key1               = key{1, 1}
		keyWithSameHashAs1 = key{-1, 1}
		key2               = key{2, 2}
		value1             = "test1"
	)
	tests := []struct {
		name      string
		node      Node
		key       persistentmap.Key
		wantValue persistentmap.Value
		wantOK    persistentmap.Value
	}{
		{
			name:      "zeroLeafNode[key1]",
			node:      &leafNode{},
			key:       key1,
			wantValue: nil,
			wantOK:    false,
		},
		{
			name:      "nodeWithKey1[key1]",
			node:      (&leafNode{}).Assoc(key1, key1.Hash(), value1),
			key:       key1,
			wantValue: value1,
			wantOK:    true,
		},
		{
			name:      "nodeWithKey1[keyWithSameHashThanKey1]",
			node:      (&leafNode{}).Assoc(key1, key1.Hash(), value1),
			key:       keyWithSameHashAs1,
			wantValue: nil,
			wantOK:    false,
		},
		{
			name:      "nodeWithKey1[key2]",
			node:      (&leafNode{}).Assoc(key1, key1.Hash(), value1),
			key:       key2,
			wantValue: nil,
			wantOK:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("ok", func(t *testing.T) {
				_, got := tt.node.Find(tt.key, tt.key.Hash())
				if got != tt.wantOK {
					t.Errorf("node[%v] should return ok=%v, returned %v instead", tt.key, tt.wantOK, got)
				}
			})
			t.Run("value", func(t *testing.T) {
				got, _ := tt.node.Find(tt.key, tt.key.Hash())
				if got != tt.wantValue {
					t.Errorf("node[%v] should return value=%v, returned %v instead", tt.key, tt.wantValue, got)
				}
			})
		})
	}
}

// TestAssocLeafNode tests the assoc function.
func TestAssocLeafNode(t *testing.T) {
	var (
		key1         = key{1, 1}
		key2         = key{2, 2}
		zeroLeafNode = leafNode{}
	)

	t.Run("assigning to zero leaf node should create a NEW instance of a leaf node", func(t *testing.T) {
		n := zeroLeafNode.Assoc(key1, key1.Hash(), "test")
		if n == &zeroLeafNode {
			t.Errorf(t.Name())
		}
	})

	t.Run("assigning to leaf node should create a NEW instance of a leaf node", func(t *testing.T) {
		n1 := zeroLeafNode.Assoc(key1, key1.Hash(), "test")
		n2 := n1.Assoc(key1, key1.Hash(), "test2")
		if n2 == n1 {
			t.Errorf(t.Name())
		}
	})

	t.Run("assigning value to an zero leaf node should result in node of type leafNode", func(t *testing.T) {
		n := zeroLeafNode.Assoc(key1, key1.Hash(), "test")
		_, ok := n.(*leafNode)
		if !ok {
			t.Error(t.Name())
		}
	})

	t.Run("assigning a value twice to leaf node with the same key1 should result in a node of type leafNode", func(t *testing.T) {
		n := zeroLeafNode.
			Assoc(key1, key1.Hash(), "test1").
			Assoc(key1, key1.Hash(), "test2")
		_, ok := n.(*leafNode)
		if !ok {
			t.Error(t.Name())
		}
	})
	t.Run("assigning two values to leaf node with level=0 should result in a bitmap indexed node", func(t *testing.T) {
		n := zeroLeafNode.
			Assoc(key1, key1.Hash(), "test1").
			Assoc(key2, key1.Hash(), "test2")
		_, ok := n.(*bitmapIndexedNode)
		if !ok {
			t.Error(t.Name())
		}
	})
	t.Run("assigning two values to leaf node with level>maxTrieLevel should result in an array node", func(t *testing.T) {
		var nodeMetadata internal.NodeMetadata
		for i := 0; i < int(maxTrieLevel)+1; i++ {
			nodeMetadata = nodeMetadata.Child()
		}

		n := (&leafNode{NodeMetadata: nodeMetadata}).
			Assoc(key1, key1.Hash(), "test1").
			Assoc(key2, key1.Hash(), "test2")
		_, ok := n.(*arrayLeafNode)
		if !ok {
			t.Error(t.Name())
		}
	})
}
