package util

import (
	"database/sql"
	"testing"
	"time"
)

func TestTaskQueue(t *testing.T) {
	tq := NewTaskQueue(&sql.DB{}, func(db *sql.DB, item interface{}, queue *TaskQueue) TaskQueueReturn {
		//if we send a bool as an item, it won't wait to process the next item
		if v, ok := item.(bool); ok && v {
			return TaskQueueReturn{
				NoWait: true,
			}
		}
		return TaskQueueReturn{
			NoWait: false,
		}
	}, time.Millisecond*15, 3)

	//The first one will get run immediately so 4 need to be added to saturate queue
	for i := 0; i < 4; i++ {
		err := tq.Enqueue(false)
		if err != nil {
			t.Errorf("Should not have failed to enqueue")
		}
		time.Sleep(time.Millisecond)
	}

	err := tq.Enqueue(false)
	if err == nil {
		t.Errorf("Should have failed to enqueue once queue is full")
	}

	time.Sleep(time.Millisecond * 80)

	//the item that failed to enqueue doesn't count
	if tq.processedItems != 4 {
		t.Errorf("Should have processed 4 items by now")
	}

	//check no wait
	for i := 0; i < 10; i++ {
		err := tq.Enqueue(true)
		if err != nil {
			t.Errorf("Should not have failed to enqueue")
		}
		time.Sleep(time.Millisecond)
	}
	time.Sleep(time.Millisecond * 20) // this should be plenty given they aren't waiting in between
	if tq.processedItems != 14 {
		t.Errorf("Should have processed 14 items by now")
	}

}
