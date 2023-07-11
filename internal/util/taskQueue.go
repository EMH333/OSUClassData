package util

import (
	"context"
	"database/sql"
	"errors"
	"github.com/enriquebris/goconcurrentqueue"
	"time"
)

type TaskQueue struct {
	internalQueue *goconcurrentqueue.FixedFIFO
	taskRunning   bool
	//stateMutex    sync.Mutex
	task func(db *sql.DB, item interface{}, queue *TaskQueue) TaskQueueReturn
	db   *sql.DB

	WaitDuration time.Duration

	//stats
	processedItems uint64
}

type TaskQueueReturn struct {
	// if true, queue won't wait before processing next item
	NoWait bool
}

type TaskQueueStats struct {
	CurrentQueue        int
	MaxSize             int
	TotalItemsProcessed uint64
}

type TaskQueueImp interface {
	Enqueue(item interface{}) error
	GetStats() TaskQueueStats
	runQueue()
}

func NewTaskQueue(db *sql.DB, task func(db *sql.DB, item interface{}, queue *TaskQueue) TaskQueueReturn, waitDuration time.Duration, maxSize int) *TaskQueue {
	return &TaskQueue{
		internalQueue: goconcurrentqueue.NewFixedFIFO(maxSize),
		taskRunning:   false,
		//stateMutex: sync.Mutex{},
		task:           task,
		db:             db,
		WaitDuration:   waitDuration,
		processedItems: 0,
	}
}

func (q *TaskQueue) Enqueue(item interface{}) error {
	if q.internalQueue.GetLen() == q.internalQueue.GetCap() {
		return errors.New("queue full")
	}

	err := q.internalQueue.Enqueue(item)
	if err != nil {
		return err
	}

	if !q.taskRunning {
		go q.runQueue()
		q.taskRunning = true
	}

	return nil
}

func (q *TaskQueue) GetStats() TaskQueueStats {
	return TaskQueueStats{
		CurrentQueue:        q.internalQueue.GetLen(),
		MaxSize:             q.internalQueue.GetCap(),
		TotalItemsProcessed: q.processedItems,
	}
}

func (q *TaskQueue) runQueue() {
	// if this func exits, make sure we know to start a new one
	defer func() {
		q.taskRunning = false
	}()
	for {
		ctx, cancel := context.WithTimeout(context.Background(), q.WaitDuration*2)
		element, err := q.internalQueue.DequeueOrWaitForNextElementContext(ctx)
		if err != nil {
			q.taskRunning = false
			cancel()
			return // return if we don't have any more elements to process
		}
		cancel()

		// run the actual task
		ret := q.task(q.db, element, q)
		q.processedItems++

		if !ret.NoWait {
			time.Sleep(q.WaitDuration)
		}
	}
}
