package util

import (
	"sync"
	"sync/atomic"
	"time"
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
	Decay       time.Duration
	DecayAmount int
}

func SetUpLeaderboard(leaderboard *Leaderboard) chan struct{} {
	//make sure counters exists
	if leaderboard.counters == nil {
		leaderboard.counters = make(map[string]int)
	}

	if leaderboard.topLock == nil {
		leaderboard.topLock = &sync.Mutex{}
	}

	//always have a default number of top entries
	if leaderboard.NumberOfTop == 0 {
		leaderboard.NumberOfTop = 10
	}

	if leaderboard.Decay == 0 {
		leaderboard.Decay = time.Hour
	}

	if leaderboard.DecayAmount == 0 {
		leaderboard.DecayAmount = 1
	}

	//start the decayer every Decay interval
	//this makes sure that "trending" is actually trending
	ticker := time.NewTicker(leaderboard.Decay)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				DecayLeaderboard(leaderboard)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return quit
}

func AddToLeaderboard(leaderboard *Leaderboard, key string) {
	//actually add to the counter
	leaderboard.counters[key]++
	atomic.AddInt64(&leaderboard.TotalCount, 1) // add to total count w/ thread safety

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

func DecayLeaderboard(leaderboard *Leaderboard) {
	leaderboard.topLock.Lock()
	defer leaderboard.topLock.Unlock()

	// decay all counters by set amount
	for key := range leaderboard.counters {
		leaderboard.counters[key] -= leaderboard.DecayAmount

		// if we are below the minimum value, then remove from counter list
		if leaderboard.counters[key] <= 0 {
			// if in top list, then remove from top list
			if contains(leaderboard.Top, key) {
				leaderboard.Top = remove(leaderboard.Top, key)
			}
			delete(leaderboard.counters, key) // remove from counter list
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

// func contains
func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// func remove
func remove(slice []string, val string) []string {
	for i, item := range slice {
		if item == val {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
