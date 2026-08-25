package linkedlist

func (l *List[T]) Reverse() {
	l.mu.Lock()
	defer l.mu.Unlock()
	var previous *Node[T]
	current := l.head
	l.tail = l.head
	for current != nil {
		next := current.Next
		current.Next = previous
		previous = current
		current = next
	}
	l.head = previous
}

func (l *List[T]) FindIndex(match func(T) bool) int {
	index := 0
	for _, value := range l.Values() {
		if match(value) {
			return index
		}
		index++
	}
	return -1
}

func (l *List[T]) InsertAt(index int, value T) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if index < 0 || index > l.count {
		return false
	}
	if index == 0 {
		node := &Node[T]{Value: value, Next: l.head}
		l.head = node
		if l.tail == nil {
			l.tail = node
		}
		l.count++
		return true
	}
	previous := l.head
	for position := 1; position < index; position++ {
		previous = previous.Next
	}
	node := &Node[T]{Value: value, Next: previous.Next}
	previous.Next = node
	if node.Next == nil {
		l.tail = node
	}
	l.count++
	return true
}

func (l *List[T]) RemoveIf(match func(T) bool) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	removed := 0
	for l.head != nil && match(l.head.Value) {
		l.head = l.head.Next
		removed++
		l.count--
	}
	if l.head == nil {
		l.tail = nil
		return removed
	}
	previous := l.head
	for current := previous.Next; current != nil; {
		if match(current.Value) {
			previous.Next = current.Next
			if current == l.tail {
				l.tail = previous
			}
			removed++
			l.count--
			current = previous.Next
			continue
		}
		previous = current
		current = current.Next
	}
	return removed
}

func Zip[A any, B any](left *List[A], right *List[B]) *List[struct {
	Left  A
	Right B
}] {
	result := New[struct {
		Left  A
		Right B
	}]()
	leftValues, rightValues := left.Values(), right.Values()
	limit := len(leftValues)
	if len(rightValues) < limit {
		limit = len(rightValues)
	}
	for index := 0; index < limit; index++ {
		result.Append(struct {
			Left  A
			Right B
		}{Left: leftValues[index], Right: rightValues[index]})
	}
	return result
}
