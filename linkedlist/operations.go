package linkedlist

import "sort"

func Filter[T any](list *List[T], keep func(T) bool) *List[T] {
	result := New[T]()
	for _, value := range list.Values() {
		if keep(value) {
			result.Append(value)
		}
	}
	return result
}

func Find[T any](list *List[T], match func(T) bool) (T, bool) {
	for _, value := range list.Values() {
		if match(value) {
			return value, true
		}
	}
	var zero T
	return zero, false
}

func Partition[T any](list *List[T], keep func(T) bool) (*List[T], *List[T]) {
	yes := New[T]()
	no := New[T]()
	for _, value := range list.Values() {
		if keep(value) {
			yes.Append(value)
		} else {
			no.Append(value)
		}
	}
	return yes, no
}

func Sort[T any](list *List[T], less func(T, T) bool) *List[T] {
	values := list.Values()
	sort.SliceStable(values, func(i, j int) bool { return less(values[i], values[j]) })
	result := New[T]()
	for _, value := range values {
		result.Append(value)
	}
	return result
}

func Take[T any](list *List[T], limit int) *List[T] {
	result := New[T]()
	if limit <= 0 {
		return result
	}
	for index, value := range list.Values() {
		if index >= limit {
			break
		}
		result.Append(value)
	}
	return result
}

func Skip[T any](list *List[T], offset int) *List[T] {
	result := New[T]()
	if offset < 0 {
		offset = 0
	}
	for index, value := range list.Values() {
		if index >= offset {
			result.Append(value)
		}
	}
	return result
}

func AppendAll[T any](target *List[T], source *List[T]) {
	if target == nil || source == nil {
		return
	}
	for _, value := range source.Values() {
		target.Append(value)
	}
}
