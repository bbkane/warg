// Package metadata provides a simple weakly-typed key/value store. Metadata must be set at creation time, and cannot be modified later.
//
// Intended use for warg:
//   - Attaching metadata when parsing, and then retrieving it later. Especially useful to mock dependencies.
//   - Attaching metadata Sections, Cmds and Flags for later retrieval. I intend to use this for TUI generation (i.e., you can set metadata on a flag that indicates the TUI should use a larger text field or a password field, etc)
package metadata

import "fmt"

// A simple type to hold key/value metadata pairs of type any

type Metadata struct {
	data map[any]any
}

func Empty() Metadata {
	return Metadata{
		// Safe to use a nil map because it's not possible to modify later
		data: nil,
	}
}

// New creates a new Metadata instance. kvs should be an even number of arguments, alternating key, value. Panics if an odd number of arguments is provided.
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

func (m *Metadata) Get(key any) (any, bool) {
	value, exists := m.data[key]
	return value, exists
}

func (m *Metadata) MustGet(key any) any {
	value, exists := m.Get(key)
	if !exists {
		panic("metadata: key does not exist: " + fmt.Sprint(key))
	}
	return value
}
