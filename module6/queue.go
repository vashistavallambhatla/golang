package main

import (
	"fmt"
	"errors"
)

type Queue[T any] struct {
	elements []T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

func (q *Queue[T]) Enqueue(ele T) {
	q.elements = append(q.elements, ele)
}

func (q *Queue[T]) Dequeue() (T,error) {
	if len(q.elements) == 0 {
		var zeroValue T
		return zeroValue,errors.New("empty Queue")
	}

	res := q.elements[0]
	q.elements = q.elements[1:]
	return res,nil
}

func (q *Queue[T]) Peek() (T,error) {
	if len(q.elements) == 0 {
		var zeroValue T
		return zeroValue,errors.New("Empty Queue")
	}

	return q.elements[0],nil
}

func main() {
	queue := NewQueue[int]()

	queue.Enqueue(1)
	queue.Enqueue(2)
	ele , _ := queue.Peek()
	fmt.Printf("%d\n",ele)
	queue.Dequeue()
	ele2 , _ := queue.Peek()
	fmt.Printf("%d\n",ele2)
}



