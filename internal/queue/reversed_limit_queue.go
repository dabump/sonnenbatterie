package queue

type ReversedLimitedQueue struct {
	queue    []interface{}
	size     int
	capacity int
}

func NewReversedLimitedQueue(capacity int) *ReversedLimitedQueue {
	return &ReversedLimitedQueue{
		queue:    make([]interface{}, 0, capacity),
		size:     0,
		capacity: capacity,
	}
}

func (q *ReversedLimitedQueue) IsEmpty() bool {
	return q.size == 0
}

func (q *ReversedLimitedQueue) IsFull() bool {
	return q.size == q.capacity
}

func (q *ReversedLimitedQueue) Queue() []interface{} {
	return q.queue
}

func (q *ReversedLimitedQueue) Enqueue(data interface{}) {
	if q.IsFull() {
		q.Dequeue()
	}
	q.queue = append([]interface{}{data}, q.queue...)
	q.size++
}

func (q *ReversedLimitedQueue) Dequeue() interface{} {
	if q.IsEmpty() {
		return nil
	}
	data := q.queue[q.size-1]
	q.queue = q.queue[:q.size-1]
	q.size--
	return data
}
