package priority

import "container/heap"

type Queue[V any] struct {
	items items[V]
}

func NewQueue[V any](size int) *Queue[V] {
	items := make(items[V], 0, size)
	heap.Init(&items)
	return &Queue[V]{
		items: items,
	}
}

func (q Queue[V]) Len() int {
	return q.items.Len()
}

func (q *Queue[V]) Push(v V, p float64) {
	heap.Push(&q.items, &item[V]{value: v, priority: p})
}

func (q *Queue[V]) Pop() V {
	return heap.Pop(&q.items).(*item[V]).value
}

type item[V any] struct {
	value    V
	priority float64
	index    int
}

type items[V any] []*item[V]

func (q items[V]) Len() int {
	return len(q)
}

func (q items[V]) Less(i int, j int) bool {
	return q[i].priority < q[j].priority
}

func (q *items[V]) Pop() any {
	old := *q
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*q = old[0 : n-1]
	return item
}

func (q *items[V]) Push(x any) {
	n := len(*q)
	item := x.(*item[V])
	item.index = n
	*q = append(*q, item)
}

func (q items[V]) Swap(i int, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}
