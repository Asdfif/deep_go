package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type CircularQueue struct {
	values []int
	start  int
	finish int
}

func NewCircularQueue(size int) CircularQueue {
	return CircularQueue{values: make([]int, size), start: -1, finish: -1}
}

func (q *CircularQueue) RemoveEdges() {
	q.start = -1
	q.finish = -1
}

func (q *CircularQueue) Push(value int) bool {
	if q.Full() {
		return false
	}
	if q.Empty() {
		q.finish = 0
		q.start = 0
		q.values[q.finish] = value
		return true
	}

	q.finish = q.Next(q.finish)
	q.values[q.finish] = value
	return true
}

func (q *CircularQueue) Last() bool {
	return q.start == q.finish
}

func (q *CircularQueue) Pop() bool {
	if q.Empty() {
		return false
	}

	q.values[q.start] = 0

	if q.Last() {
		q.RemoveEdges()
	} else {
		q.start = q.Next(q.start)
	}

	return true
}

func (q *CircularQueue) Next(idx int) int {
	return (idx + 1) % len(q.values)
}

func (q *CircularQueue) Front() int {
	if q.Empty() {
		return -1
	}

	front := q.values[q.start]
	if front == 0 {
		return -1
	}

	return front
}

func (q *CircularQueue) Back() int {
	if q.Empty() {
		return -1
	}

	back := q.values[q.finish]

	return back
}

func (q *CircularQueue) Empty() bool {
	return q.start == -1 && q.finish == -1
}

func (q *CircularQueue) Full() bool {
	return !q.Empty() && q.Next(q.finish) == q.start
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue(queueSize)

	fmt.Printf("INIT %+v\n", queue)
	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	fmt.Printf("PUSH 1 %+v\n", queue)

	assert.True(t, queue.Push(2))
	fmt.Printf("PUSH 2 %+v\n", queue)
	assert.True(t, queue.Push(3))
	fmt.Printf("PUSH 3 %+v\n", queue)
	assert.False(t, queue.Push(4))
	fmt.Printf("PUSH 4 %+v\n", queue)

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	fmt.Printf("POP %+v\n", queue)

	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))
	fmt.Printf("PUSH 4 %+v\n", queue)

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())
	fmt.Printf("POP %+v\n", queue)
	assert.True(t, queue.Pop())
	fmt.Printf("POP %+v\n", queue)
	assert.True(t, queue.Pop())
	fmt.Printf("POP %+v\n", queue)
	assert.False(t, queue.Pop())
	fmt.Printf("POP %+v\n", queue)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}
