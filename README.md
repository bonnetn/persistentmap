# Persistent hash map


Just a simple persistent hash map implemented in go.

The implementation is inspired from Clojure's persistent hash maps.

## How to use

Just 
```go
var emptyMap persistentMap
map1  := emptyMap.Set(key1, value1) 
map12 := map1.Set(key2, value2) 

// At this point:
// Map is empty
// Map1 contains key1:value1
// Map12 contains key1:value1 AND key2:value2


value, ok := map1.Get(key1) // value1, true
value, ok = map1.Get(key2) // nil, false
```

## Complexity

With `n` being the elements count in the map.

Amortized time complexity for Get/Set: `O(log32(n))` but it is sometimes considered as `O(1)`

Shantanu Kumar wrote in his book `Clojure High Performance Programming`:

> [...] even though the complexity is O(log32(n)), only 2^32 hash codes can fit into the trie nodes. Hence, log32(2^32), which turns out to be 6.4 and is less than 7, is the worst-case complexity and can be considered near-constant time.

Space complexity: `O(n)`


