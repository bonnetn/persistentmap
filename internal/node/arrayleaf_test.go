package node

import "testing"

func TestArrayLeafAssoc(t *testing.T) {
	testMultipleNode(t, &arrayLeafNode{})
}

func testMultipleNode(t *testing.T, zeroNode Node) {
	var (
		key1 = key{1, 1}
		key2 = key{2, 2}
	)
	t.Run("should create a new array on every assoc", func(t *testing.T) {
		n1 := zeroNode.
			Assoc(key1, key1.Hash(), "test1")

		n2 := n1.Assoc(key2, key2.Hash(), "test2")

		_, ok := n1.Find(key1, key1.Hash())
		if !ok {
			t.Error("n1 should contain key1")
		}
		_, ok = n1.Find(key2, key2.Hash())
		if ok {
			t.Error("n1 should not contain key2")
		}
		_, ok = n2.Find(key1, key1.Hash())
		if !ok {
			t.Error("n2 should contain key1")
		}
		_, ok = n2.Find(key2, key2.Hash())
		if !ok {
			t.Error("n2 should contain key2")
		}

	})

	t.Run("should create a new array on every assoc invert order", func(t *testing.T) {
		n1 := zeroNode.
			Assoc(key2, key2.Hash(), "test2")

		n2 := n1.Assoc(key1, key1.Hash(), "test1")

		_, ok := n1.Find(key1, key1.Hash())
		if ok {
			t.Error("n1 should not contain key1")
		}
		_, ok = n1.Find(key2, key2.Hash())
		if !ok {
			t.Error("n1 should contain key2")
		}
		_, ok = n2.Find(key1, key1.Hash())
		if !ok {
			t.Error("n2 should contain key1")
		}
		_, ok = n2.Find(key2, key2.Hash())
		if !ok {
			t.Error("n2 should contain key2")
		}

	})
	t.Run("should be able to update a value for a given key", func(t *testing.T) {
		n1 := zeroNode.
			Assoc(key1, key1.Hash(), "test1")
		n2 := n1.Assoc(key1, key1.Hash(), "test2")

		value, _ := n1.Find(key1, key1.Hash())
		if value != "test1" {
			t.Error("n1 should contain 'test1'")
		}
		value, _ = n2.Find(key1, key1.Hash())
		if value != "test2" {
			t.Error("n2 should contain 'test1'")
		}
	})
}
