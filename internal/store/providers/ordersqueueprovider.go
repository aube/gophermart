package providers

type OrdersQueueRepositoryProvider interface {
	Enqueue(item int)
	Dequeue() (int, error)
	IsEmpty() bool
	Size() int
}
