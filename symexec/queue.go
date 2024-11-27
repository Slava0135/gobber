package symexec

import "math/rand"

type Queue interface {
	push(state State)
	pop() State
	empty() bool
}

type RandomQueue struct {
	q []State
}

type BFSQueue struct {
	q []State
}

// stack, actualy
type DFSQueue struct {
	q []State
}

// order doesn't matter
func (q *RandomQueue) push(state State) {
	q.q = append(q.q, state)
}

func (q *RandomQueue) pop() State {
	index := rand.Intn(len(q.q))
	next := q.q[index]
	q.q = append(q.q[:index], q.q[index+1:]...)
	return next
}

func (q *RandomQueue) empty() bool {
	return len(q.q) <= 0
}

func (q *BFSQueue) push(state State) {
	q.q = append(q.q, state)
}

// pop oldest - first in queue
func (q *BFSQueue) pop() State {
	next := q.q[0]
	q.q = q.q[1:]
	return next
}

func (q *BFSQueue) empty() bool {
	return len(q.q) <= 0
}

func (q *DFSQueue) push(state State) {
	q.q = append(q.q, state)
}

// pop newest - last in queue (first in stack)
func (q *DFSQueue) pop() State {
	last := len(q.q) - 1
	next := q.q[last]
	q.q = q.q[:last]
	return next
}

func (q *DFSQueue) empty() bool {
	return len(q.q) <= 0
}
