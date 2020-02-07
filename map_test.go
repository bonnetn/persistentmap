package persistentmap

import (
	"testing"
)

type IntKey int32

func (s IntKey) Equal(c Comparable) bool {
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

func (s CollisionIntKey) Equal(c Comparable) bool {
	i, ok := c.(CollisionIntKey)
	if !ok {
		return false
	}
	return i == s
}

func (s CollisionIntKey) Hash() int32 {
	return 1
}

func validateMap(t *testing.T, myMap persistentMap, want map[Key]*string) {
	for k, wantValue := range want {
		wantOK := false
		if wantValue != nil {
			wantOK = true
		}

		gotValue, gotOK := myMap.Get(k)
		if gotOK != wantOK {
			t.Errorf("map[%v] returned ok=%t instead of %t", k, gotOK, wantOK)
		}
		if wantValue == nil && gotValue != nil {
			t.Errorf("map[%v] returned a non-nil value", k)
		}
		if wantValue != nil && gotValue != *wantValue {
			t.Errorf("map[%v] returned %v instead of %v", k, gotValue, *wantValue)
		}
	}
}

func TestMap(t *testing.T) {
	tests := []struct {
		name string
		key1 Key
		key2 Key
		key3 Key
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
				map131    = map13.Set(tt.key1, value2)
			)

			validateMap(t, emptyMap, map[Key]*string{
				tt.key1: nil,
				tt.key2: nil,
				tt.key3: nil,
			})
			validateMap(t, map1, map[Key]*string{
				tt.key1: &value1,
				tt.key2: nil,
				tt.key3: nil,
			})
			validateMap(t, map123, map[Key]*string{
				tt.key1: &value1,
				tt.key2: &value2,
				tt.key3: &value3,
			})
			validateMap(t, map13, map[Key]*string{
				tt.key1: &value1,
				tt.key2: nil,
				tt.key3: &value3,
			})
			validateMap(t, map131, map[Key]*string{
				tt.key1: &value2,
				tt.key2: nil,
				tt.key3: &value3,
			})
		})
	}
}
