package linkedlist

import "testing"

func TestListOperations(t *testing.T) {
	list := New[int]()
	list.Append(2)
	list.Prepend(1)
	list.Append(4)
	if got := list.Len(); got != 3 {
		t.Fatalf("length = %d", got)
	}
	filtered := Filter(list, func(value int) bool { return value%2 == 0 })
	if got := filtered.Values(); len(got) != 2 || got[0] != 2 || got[1] != 4 {
		t.Fatalf("filtered values = %#v", got)
	}
	value, ok := list.PopFront()
	if !ok || value != 1 {
		t.Fatalf("pop = %d, %v", value, ok)
	}
}
