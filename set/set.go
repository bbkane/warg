package set

// Set is a generic set implementation using a map[T]struct{} as the backing store.
type Set[T comparable] struct {
	data map[T]struct{}
}

// New creates and returns a new empty set. Sets use maps, so this type is not "copy by value".
func New[T comparable]() Set[T] {
	return Set[T]{data: make(map[T]struct{})}
}

// Add inserts an element into the set.
func (s *Set[T]) Add(value T) {
	s.data[value] = struct{}{}
}

// Delete removes an element from the set.
func (s *Set[T]) Delete(value T) {
	delete(s.data, value)
}

// Contains checks if the set contains an element.
func (s *Set[T]) Contains(value T) bool {
	_, exists := s.data[value]
	return exists
}
