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

// Removes and returns the front of the queue, if the queue is empty nil is
// returned
func (q *Queue[P]) Pop() (*P, bool) {
	var v *P = nil
	ok := false

	q.mutex.Lock()
	if q.front != nil {
		v = &q.front.val

		// Only one item so make sure the back field is cleared
		if q.front == q.back {
			q.back = nil
		}

		q.front = q.front.next
		ok = true
	}
	q.mutex.Unlock()

	return v, ok
}

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

func (q *Queue[P]) PushMultiple(items []P) {
	if len(items) == 0 {
		return
	}

	c := &queueItem[P]{
		val:  items[0],
		next: nil,
	}
	q.mutex.Lock()

	if q.back == nil {
		q.front = c
	} else {
		q.back.next = c
	}

	for _, i := range items[1:] {
		c.next = &queueItem[P]{
			val:  i,
			next: nil,
		}
		c = c.next
	}

	q.back = c
	q.mutex.Unlock()
}
