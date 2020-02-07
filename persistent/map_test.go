package persistent

import (
	"testing"
	"github.com/bonnetn/persistentmap"
)

type IntKey int32

func (s IntKey) Equal(c persistentmap.Comparable) bool {
	i, ok := c.(IntKey)
	if !ok {
		return false
	}
	return i == s
}

func (s IntKey) Hash() int32 {
	return int32(s)
}

type CollisionIntKey int32

func (s CollisionIntKey) Equal(c persistentmap.Comparable) bool {
	i, ok := c.(CollisionIntKey)
	if !ok {
		return false
	}
	return i == s
}

func (s CollisionIntKey) Hash() int32 {
	return 1
}

func validateMap(t *testing.T, mapName string, myMap persistentMap, want map[persistentmap.Key]*string) {
	for k, wantValue := range want {
		wantOK := false
		if wantValue != nil {
			wantOK = true
		}

		gotValue, gotOK := myMap.Get(k)
		if gotOK != wantOK {
			t.Errorf("%s[%v] returned ok=%t instead of %t", mapName, k, gotOK, wantOK)
		}
		if wantValue == nil && gotValue != nil {
			t.Errorf("%s[%v] returned a non-nil value", mapName, k)
		}
		if wantValue != nil && gotValue != *wantValue {
			t.Errorf("%s[%v] returned %v instead of %v", mapName, k, gotValue, *wantValue)
		}
	}
}

func TestMap(t *testing.T) {
	tests := []struct {
		name string
		key1 persistentmap.Key
		key2 persistentmap.Key
		key3 persistentmap.Key
	}{
		{
			name: "int key",
			key1: IntKey(1),
			key2: IntKey(2),
			key3: IntKey(3),
		},
		{
			name: "key collision",
			key1: CollisionIntKey(1),
			key2: CollisionIntKey(2),
			key3: CollisionIntKey(3),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				value1 = "value for key1"
				value2 = "value for key2"
				value3 = "value for key3"
			)
			var (
				emptyMap persistentMap
				map1     = emptyMap.Set(tt.key1, value1)
				map12    = map1.Set(tt.key2, value2)
				map123   = map12.Set(tt.key3, value3)
				map13    = map1.Set(tt.key3, value3)
				map131   = map13.Set(tt.key1, value2)
			)

			validateMap(t, "emptyMap", emptyMap, map[persistentmap.Key]*string{
				tt.key1: nil,
				tt.key2: nil,
				tt.key3: nil,
			})
			validateMap(t, "mapWithKey1", map1, map[persistentmap.Key]*string{
				tt.key1: &value1,
				tt.key2: nil,
				tt.key3: nil,
			})
			validateMap(t, "mapWithKey1And2", map12, map[persistentmap.Key]*string{
				tt.key1: &value1,
				tt.key2: &value2,
				tt.key3: nil,
			})
			validateMap(t, "mapWithKey1And2And3", map123, map[persistentmap.Key]*string{
				tt.key1: &value1,
				tt.key2: &value2,
				tt.key3: &value3,
			})
			validateMap(t, "mapWithKey1And3", map13, map[persistentmap.Key]*string{
				tt.key1: &value1,
				tt.key2: nil,
				tt.key3: &value3,
			})
			validateMap(t, "mapWithKey1And2And1", map131, map[persistentmap.Key]*string{
				tt.key1: &value2,
				tt.key2: nil,
				tt.key3: &value3,
			})
		})
	}
}
