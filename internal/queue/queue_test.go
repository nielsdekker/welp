package queue_test

import (
	"testing"

	"github.com/nielsdekker/welp/internal/_tests/asserts"
	"github.com/nielsdekker/welp/internal/queue"
)

func Test_Queue(t *testing.T) {
	q := queue.New[int]()

	asserts.Eq(t, q.Pop(), nil)

	q.Push(1)
	asserts.Eq(t, *q.Pop(), 1)
	asserts.Eq(t, q.Pop(), nil)

	q.Push(2)
	q.Push(3)
	q.Push(4)
	asserts.Eq(t, *q.Pop(), 2)
	asserts.Eq(t, *q.Pop(), 3)
	asserts.Eq(t, *q.Pop(), 4)
	asserts.Eq(t, q.Pop(), nil)
}
