package util

import (
	"sync"
	"sync/atomic"
)

//TODO eventually decay values over time
//Leaderboard keeps track of the top values in a list
//thread safe for the top calculations and the total count, not for the counters since that matters less
type Leaderboard struct {
	counters   map[string]int
	Top        []string
	TotalCount int64

	//internal help stuff
	minimumTopValue int
	topLock         *sync.Mutex //lock when updating top list

	//config
	NumberOfTop int
}

func AddToLeaderboard(leaderboard *Leaderboard, key string) {
	//make sure counters exists
	if leaderboard.counters == nil {
		leaderboard.counters = make(map[string]int)
	}

	if leaderboard.topLock == nil {
		leaderboard.topLock = &sync.Mutex{}
	}

	//actually add to the counter
	leaderboard.counters[key]++
	atomic.AddInt64(&leaderboard.TotalCount, 1) // add to total count w/ thread safety

	//always have a default number of top entries
	if leaderboard.NumberOfTop == 0 {
		leaderboard.NumberOfTop = 10
	}

	//move towards the max number of top entries if not already there
	if len(leaderboard.Top) < leaderboard.NumberOfTop {
		var alreadyInTop = false
		for _, topKey := range leaderboard.Top {
			if topKey == key {
				alreadyInTop = true
				break
			}
		}
		if !alreadyInTop {
			leaderboard.Top = append(leaderboard.Top, key)
		}
	}

	if leaderboard.counters[key] >= leaderboard.minimumTopValue {
		figureOutNewTop(leaderboard, key)

		// update min value if necessary
		if leaderboard.minimumTopValue == 0 && len(leaderboard.Top) > 0 {
			var lastIndex = len(leaderboard.Top) - 1
			leaderboard.minimumTopValue = leaderboard.counters[leaderboard.Top[lastIndex]]
		}
	}
}

// loop through the top list and see if new key is greater than any of them (starting at smallest)
// if it is, then update the list and the min value
func figureOutNewTop(leaderboard *Leaderboard, key string) {
	// make sure we have lock so we are the only one updating the top list at a time
	leaderboard.topLock.Lock()
	defer leaderboard.topLock.Unlock()

	var lastIndex = len(leaderboard.Top) - 1

	for i := 0; i <= lastIndex; i++ {
		if leaderboard.counters[key] >= leaderboard.counters[leaderboard.Top[i]] {
			bubbleDown(leaderboard, i, key)
			leaderboard.minimumTopValue = leaderboard.counters[leaderboard.Top[lastIndex]]
			return
		}
	}
}

// shifts everything down from position and inserts key at position
// overwrites duplicate if it exists
func bubbleDown(leaderboard *Leaderboard, pos int, key string) {
	var lastIndex = len(leaderboard.Top) - 1

	// handle case where key is already in the right place
	if leaderboard.Top[pos] == key {
		return
	}

	//handle case where key is already in list
	for i := pos; i <= lastIndex; i++ {
		if leaderboard.Top[i] == key {
			for j := i; j >= pos; j-- {
				if j == 0 {
					continue
				}
				leaderboard.Top[j] = leaderboard.Top[j-1]
			}
			leaderboard.Top[pos] = key
			return
		}
	}

	//handle normal case
	for i := lastIndex; i >= pos; i-- {
		if i == 0 {
			continue
		}
		leaderboard.Top[i] = leaderboard.Top[i-1]
	}

	leaderboard.Top[pos] = key
}
