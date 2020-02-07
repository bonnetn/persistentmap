package internal

import (
	"testing"
)

func TestNodeMetadata(t *testing.T) {
	tests := []struct {
		name string
		node NodeMetadata
		want int32
	}{
		{
			name: "zero node should have level 0",
			node: NodeMetadata{},
			want: 0,
		},
		{
			name: "first child should have level 1",
			node: NodeMetadata{}.Child(),
			want: 1,
		},
		{
			name: "5th child should have level 5",
			node: NodeMetadata{}.
				Child().
				Child().
				Child().
				Child().
				Child() ,
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.node.Level(); got != tt.want {
				t.Errorf("Level() = %v, want %v", got, tt.want)
			}
		})
	}
}
