package blocking_queue

// TODO: unit testing
import (
	"errors"
	"sync"
	"time"
)

var ErrQueueClosed = errors.New("queue is closed")

type BlockingQueue struct {
	lock    sync.Mutex // lock is used to synchronize access to the queue.
	cond    *sync.Cond // cond is a conditional variable for queue synchronization.
	data    TaskQueue  // data is the underlying task queue.
	closed  bool       // closed indicates whether the queue is closed (used for graceful shutdown).
	waiting int        // waiting counts the number of consumers waiting for tasks (used for optimization).
}

func NewBlockingQueue() *BlockingQueue {
	bq := &BlockingQueue{}
	bq.cond = sync.NewCond(&bq.lock)
	bq.data = make(TaskQueue, 0)
	return bq
}

func (q *BlockingQueue) Put(t *Task) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.closed {
		return ErrQueueClosed
	}

	// Set StartTime and CreatedAt when adding a task
	t.StartTime = time.Now()
	t.CreatedAt = t.StartTime
	q.data = append(q.data, t)

	// Notify a waiting consumer
	if q.waiting > 0 {
		q.cond.Signal()
	}
	return nil
}

func (q *BlockingQueue) Get() (*Task, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	// Wait until a task is available
	for len(q.data) == 0 && !q.closed {
		q.waiting++
		q.cond.Wait()
		q.waiting--
	}

	// Check if the queue is closed
	if q.closed && len(q.data) == 0 {
		return nil, ErrQueueClosed
	}

	// Get the task from the front of the queue
	t := q.data[0]
	q.data = q.data[1:]

	t.EndTime = time.Now()
	t.UpdatedAt = t.EndTime

	return t, nil
}

func (q *BlockingQueue) Size() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return len(q.data)
}

func (q *BlockingQueue) Close() {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.closed = true
	q.cond.Broadcast()
}
