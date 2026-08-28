package linkedlist

import "sync"

type Node[T any] struct {
	Value T
	Next  *Node[T]
}

type List[T any] struct {
	mu    sync.RWMutex
	head  *Node[T]
	tail  *Node[T]
	count int
}

func New[T any]() *List[T] {
	return &List[T]{}
}

func (l *List[T]) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.count
}

func (l *List[T]) Append(value T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	node := &Node[T]{Value: value}
	if l.tail == nil {
		l.head = node
		l.tail = node
	} else {
		l.tail.Next = node
		l.tail = node
	}
	l.count++
}

func (l *List[T]) Prepend(value T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	node := &Node[T]{Value: value, Next: l.head}
	l.head = node
	if l.tail == nil {
		l.tail = node
	}
	l.count++
}

func (l *List[T]) PopFront() (T, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.head == nil {
		var zero T
		return zero, false
	}
	node := l.head
	l.head = node.Next
	node.Next = nil
	l.count--
	if l.head == nil {
		l.tail = nil
	}
	return node.Value, true
}

func (l *List[T]) Values() []T {
	l.mu.RLock()
	defer l.mu.RUnlock()
	values := make([]T, 0, l.count)
	for node := l.head; node != nil; node = node.Next {
		values = append(values, node.Value)
	}
	return values
}

func (l *List[T]) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.head = nil
	l.tail = nil
	l.count = 0
}

func (l *List[T]) Clone() *List[T] {
	copyList := New[T]()
	for _, value := range l.Values() {
		copyList.Append(value)
	}
	return copyList
}
