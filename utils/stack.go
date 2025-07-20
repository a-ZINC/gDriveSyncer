package utils

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
	if s.IsEmpty() {
		var zero T
		return zero, false
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, true
}
func (s *Stack[T]) Peek() (*T, bool) {
	if s.IsEmpty() {
		var zero *T
		return zero, false
	}
	item := s.items[len(s.items)-1]
	return &item, true
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}