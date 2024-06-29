package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	bq := NewBlockingQueue(3)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		bq.Add(1)
		bq.Add(2)
		bq.Add(3)
		bq.Add(4)
		bq.Add(5)
		bq.Add(6)
		bq.Add(7)
		bq.Add(8)
		bq.Add(9)
		bq.Add(10)
	}()

	go func() {
		defer wg.Done()
		time.Sleep(5 * time.Second)
		for range 10 {
			fmt.Printf("Got %v from queue\n", bq.Get())
			time.Sleep(1 * time.Second)
		}
	}()

	wg.Wait()
}

type BlockingQueue struct {
	m        *sync.Mutex
	c        *sync.Cond
	buf      []interface{}
	capacity int
}

func NewBlockingQueue(capacity int) *BlockingQueue {
	m := sync.Mutex{}
	c := sync.NewCond(&m)
	return &BlockingQueue{
		m:        &m,
		capacity: capacity,
		c:        c,
		buf:      make([]interface{}, 0, capacity),
	}
}

func (q *BlockingQueue) Add(element interface{}) {
	q.c.L.Lock()
	defer q.c.L.Unlock()
	for q.capacity == len(q.buf) {
		q.c.Wait()
	}

	fmt.Printf("Added %v to the queue\n", element)
	q.buf = append(q.buf, element)
	fmt.Printf("Queue after Add operation %+v\n\n", q.buf)
	q.c.Broadcast()
}

func (q *BlockingQueue) Get() interface{} {
	q.c.L.Lock()
	defer q.c.L.Unlock()

	for len(q.buf) == 0 {
		q.c.Wait()
	}
	response := q.buf[0]
	q.buf = q.buf[1:]
	q.c.Broadcast()
	fmt.Printf("Returning %v from the queue\n", response)
	fmt.Printf("Queue after Get operation %+v\n\n", q.buf)
	return response
}
