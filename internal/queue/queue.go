package queue

import "sync"

type Queue[P any] struct {
	mutex sync.Mutex
	front *queueItem[P]
	back  *queueItem[P]
}

type queueItem[P any] struct {
	val  P
	next *queueItem[P]
}

func New[P any]() *Queue[P] {
	return &Queue[P]{
		mutex: sync.Mutex{},
		front: nil,
		back:  nil,
	}
}

// Locks the queue and returns the front of the queue. The front of the queue is
// updated to the next value, if the queue is empty nil is returned.
func (q *Queue[P]) Pop() (*P, bool) {
	var v *P = nil
	ok := false

	q.mutex.Lock()
	if q.front != nil {
		v = &q.front.val

		// Only one item so make sure the back field is also cleared
		if q.front == q.back {
			q.back = nil
		}

		q.front = q.front.next
		ok = true
	}
	q.mutex.Unlock()

	return v, ok
}

// Locks the queue and adds the passed item
func (q *Queue[P]) Push(item P) {
	i := &queueItem[P]{
		val:  item,
		next: nil,
	}

	q.mutex.Lock()
	// The very first item
	if q.back == nil {
		q.front = i
	} else {
		q.back.next = i
	}
	q.back = i
	q.mutex.Unlock()
}
