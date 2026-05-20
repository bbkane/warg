// Package set provides a generic set backed by a map.
package set

// Set is a generic set backed by map[T]struct{}. Not safe for concurrent use.
// Pass by pointer or use the pointer receiver methods.
type Set[T comparable] struct {
	data map[T]struct{}
}

// New creates an empty [Set].
func New[T comparable]() Set[T] {
	return Set[T]{data: make(map[T]struct{})}
}

// Add inserts a value into the set. No-op if already present.
func (s *Set[T]) Add(value T) {
	s.data[value] = struct{}{}
}

// AddAll inserts all given values into the set.
func (s *Set[T]) AddAll(values ...T) {
	for _, v := range values {
		s.Add(v)
	}
}

// Delete removes a value from the set. No-op if not present.
func (s *Set[T]) Delete(value T) {
	delete(s.data, value)
}

// Contains reports whether the set contains the given value.
func (s *Set[T]) Contains(value T) bool {
	_, exists := s.data[value]
	return exists
}
