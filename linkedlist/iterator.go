package linkedlist

type Iterator[T any] struct {
	next *Node[T]
}

func (l *List[T]) Iterator() *Iterator[T] {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return &Iterator[T]{next: l.head}
}

func (it *Iterator[T]) Next() (T, bool) {
	if it.next == nil {
		var zero T
		return zero, false
	}
	value := it.next.Value
	it.next = it.next.Next
	return value, true
}

func (it *Iterator[T]) Drain() []T {
	values := make([]T, 0)
	for {
		value, ok := it.Next()
		if !ok {
			return values
		}
		values = append(values, value)
	}
}

func Map[T any, R any](list *List[T], fn func(T) R) *List[R] {
	result := New[R]()
	for _, value := range list.Values() {
		result.Append(fn(value))
	}
	return result
}

func Reduce[T any, R any](list *List[T], seed R, fn func(R, T) R) R {
	result := seed
	for _, value := range list.Values() {
		result = fn(result, value)
	}
	return result
}
