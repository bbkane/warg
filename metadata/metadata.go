// Package metadata provides a simple immutable key/value store for attaching
// arbitrary data during parsing. Common uses include injecting test mocks
// and storing UI hints for flags (planned).
package metadata

import "fmt"

// Metadata is an immutable key/value store. Keys and values can be any type.
// Create with [New] or [Empty]. Values cannot be modified after creation.
type Metadata struct {
	data map[any]any
}

// Empty returns a [Metadata] with no entries.
func Empty() Metadata {
	return Metadata{
		// Safe to use a nil map because it's not possible to modify later
		data: nil,
	}
}

// New creates a [Metadata] from alternating key/value pairs.
// Panics if an odd number of arguments is provided.
func New(kvs ...any) Metadata {
	data := make(map[any]any)
	if len(kvs)%2 != 0 {
		panic("metadata.New: odd number of arguments")
	}
	for i := 0; i < len(kvs); i += 2 {
		key := kvs[i]
		value := kvs[i+1]
		data[key] = value
	}
	return Metadata{
		data: data,
	}
}

// Get retrieves a value by key. Returns (value, true) if found, (nil, false) otherwise.
func (m *Metadata) Get(key any) (any, bool) {
	value, exists := m.data[key]
	return value, exists
}

// MustGet retrieves a value by key, panicking if the key does not exist.
func (m *Metadata) MustGet(key any) any {
	value, exists := m.Get(key)
	if !exists {
		panic("metadata: key does not exist: " + fmt.Sprint(key))
	}
	return value
}
