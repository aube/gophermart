package memory

import "errors"

// OrdersQueueRepository ...
type OrdersQueueRepository struct {
	oq []int
}

// Enqueue ...
func (q *OrdersQueueRepository) Enqueue(item int) {
	q.oq = append(q.oq, item)
}

// Dequeue ...
func (q *OrdersQueueRepository) Dequeue() (int, error) {
	if q.IsEmpty() {
		return 0, errors.New("Queue is empty")
	}
	item := q.oq[0]
	q.oq = q.oq[1:]
	return item, nil
}

// IsEmpty ...
func (q *OrdersQueueRepository) IsEmpty() bool {
	return len(q.oq) == 0
}

// Size ...
func (q *OrdersQueueRepository) Size() int {
	return len(q.oq)
}
