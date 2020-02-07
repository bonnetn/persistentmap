package node

import "testing"

func TestBitmapIndexedNode(t *testing.T) {
	testMultipleNode(t, &bitmapIndexedNode{})

}
