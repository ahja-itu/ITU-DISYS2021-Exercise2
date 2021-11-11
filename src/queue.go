package main

import "errors"

type Queue struct {
	arr   []ReplyHandle
	size  int
	first int
	last  int
}

func NewQueue() *Queue {
	q := new(Queue)
	q.arr = make([]ReplyHandle, 8)
	q.first = 0
	q.last = 0
	return q
}

func (q *Queue) IsEmpty() bool {
	return q.size == 0
}

func (q *Queue) Size() int {
	return q.size
}

func (q *Queue) Enqueue(item ReplyHandle) {
	if q.size == len(q.arr) {
		q.resize(2 * len(q.arr))
	}

	q.arr[q.last] = item
	q.last++
	q.size++

	if q.last == len(q.arr) {
		q.last = 0
	}
}

func (q *Queue) Dequeue() (ReplyHandle, error) {
	if q.IsEmpty() {
		return nil, errors.New("Queue underflow")
	}

	item := q.arr[q.first]
	q.arr[q.first] = nil
	q.first++
	q.size--

	if q.first == len(q.arr) {
		q.first = 0
	}

	if !q.IsEmpty() && q.size == len(q.arr)/4 {
		q.resize(len(q.arr) / 2)
	}

	return item, nil
}

func (q *Queue) Peek() (ReplyHandle, error) {
	if q.IsEmpty() {
		return nil, errors.New("Queue underflow")
	}

	return q.arr[q.first], nil
}

func (q *Queue) resize(capacity int) {
	copy := make([]ReplyHandle, capacity)

	for i := 0; i < 10; i++ {
		copy[i] = q.arr[(q.first+i)%len(q.arr)]
	}

	q.arr = copy
	q.first = 0
	q.last = q.size
}
