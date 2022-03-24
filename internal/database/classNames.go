package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/enriquebris/goconcurrentqueue"
)

const classNameWaitTime = 30

var classNameQueue = goconcurrentqueue.NewFixedFIFO(20)
var classNameTaskRunning = false

func UpdateClassName(db *sql.DB, class string) {
	err := classNameQueue.Enqueue(class)
	// if there is an error then the queue is probably full so we don't add another item
	if err != nil {
		return
	}
	//TODO put class into queue
	//if there is a class in the queue && the thread hasn't been started yet, then start it
	//within the thread, every 30 seconds, check if there is a class in the queue:
	//pull a class from the queue
	//if it hasn't been queried, query the OSU API for the class name
	//if it has, then pull from database
	//normalize it
	//update the class name in the database

	//can use https://github.com/enriquebris/goconcurrentqueue

	if !classNameTaskRunning {
		go classNameTask()
	}
}

func classNameTask() {
	classNameTaskRunning = true
	for classNameTaskRunning {
		//context with a minute timeout
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

		element, err := classNameQueue.DequeueOrWaitForNextElementContext(ctx)
		if err != nil {
			classNameTaskRunning = false
			cancel()
			return // return if we don't have any more elements to process
		}

		refreshClassName(element.(string))

		// wait so we don't overload API
		time.Sleep(classNameWaitTime * time.Second)

		cancel()
	}
	classNameTaskRunning = false
}

func refreshClassName(class string) {
	//TODO pull from API
	//TODO update database
}
