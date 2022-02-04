package util

import (
	"testing"
)

func TestBasicLeaderboard(t *testing.T) {
	leaderboard := &Leaderboard{
		counters: make(map[string]int),
		top:      []string{},
	}

	AddToLeaderboard(leaderboard, "a")
	AddToLeaderboard(leaderboard, "a")
	AddToLeaderboard(leaderboard, "b")
	AddToLeaderboard(leaderboard, "c")
	AddToLeaderboard(leaderboard, "d")

	if leaderboard.counters["a"] != 2 {
		t.Errorf("Expected a to have value 2, got %d", leaderboard.counters["a"])
	}

	if leaderboard.counters["b"] != 1 {
		t.Errorf("Expected b to have value 1, got %d", leaderboard.counters["b"])
	}

	if leaderboard.counters["c"] != 1 {
		t.Errorf("Expected c to have value 1, got %d", leaderboard.counters["c"])
	}

	if leaderboard.counters["d"] != 1 {
		t.Errorf("Expected d to have value 1, got %d", leaderboard.counters["d"])
	}

	if leaderboard.top[0] != "a" {
		t.Errorf("Expected top[0] to be a, got %s", leaderboard.top[0])
	}

	if leaderboard.totalCount != 5 {
		t.Errorf("Expected totalCount to be 5, got %d", leaderboard.totalCount)
	}
}

func TestFullLeaderboard(t *testing.T) {
	leaderboard := &Leaderboard{
		counters:    make(map[string]int),
		top:         []string{},
		numberOfTop: 3,
	}

	//note loading in weird order to cause misbehavior if possible
	AddToLeaderboard(leaderboard, "d")

	// 4 a's
	AddToLeaderboard(leaderboard, "a")
	AddToLeaderboard(leaderboard, "a")
	AddToLeaderboard(leaderboard, "a")
	AddToLeaderboard(leaderboard, "a")

	// 2 c's
	AddToLeaderboard(leaderboard, "c")
	AddToLeaderboard(leaderboard, "c")

	// 3 b's
	AddToLeaderboard(leaderboard, "b")
	AddToLeaderboard(leaderboard, "b")
	AddToLeaderboard(leaderboard, "b")

	if leaderboard.minimumTopValue != 2 {
		t.Errorf("Expected minimumTopValue to be 2, got %d", leaderboard.minimumTopValue)
	}

	//the top should be a, b, c in that order
	if leaderboard.top[0] != "a" {
		t.Errorf("Expected top[0] to be a, got %s", leaderboard.top[0])
	}

	if leaderboard.top[1] != "b" {
		t.Errorf("Expected top[1] to be b, got %s", leaderboard.top[1])
	}

	if leaderboard.top[2] != "c" {
		t.Errorf("Expected top[2] to be c, got %s", leaderboard.top[2])
	}

	//actual counter values
	if leaderboard.counters["a"] != 4 {
		t.Errorf("Expected a to have value 4, got %d", leaderboard.counters["a"])
	}

	if leaderboard.counters["b"] != 3 {
		t.Errorf("Expected b to have value 3, got %d", leaderboard.counters["b"])
	}

	if leaderboard.counters["c"] != 2 {
		t.Errorf("Expected c to have value 2, got %d", leaderboard.counters["c"])
	}

	if leaderboard.counters["d"] != 1 {
		t.Errorf("Expected d to have value 1, got %d", leaderboard.counters["d"])
	}
}

func TestManyEntriesLeaderboard(t *testing.T) {
	leaderboard := &Leaderboard{
		counters:    make(map[string]int),
		top:         []string{},
		numberOfTop: 3,
	}

	//note loading in weird order to cause misbehavior if possible
	AddToLeaderboard(leaderboard, "d")

	// 4 a's
	AddToLeaderboard(leaderboard, "a")
	AddToLeaderboard(leaderboard, "a")
	AddToLeaderboard(leaderboard, "a")
	AddToLeaderboard(leaderboard, "a")

	// 2 c's
	AddToLeaderboard(leaderboard, "c")
	AddToLeaderboard(leaderboard, "c")

	// 3 b's
	AddToLeaderboard(leaderboard, "b")
	AddToLeaderboard(leaderboard, "b")
	AddToLeaderboard(leaderboard, "b")

	//add a bunch of entries
	for i := 0; i < 1000; i++ {
		AddToLeaderboard(leaderboard, "b")
		AddToLeaderboard(leaderboard, "a")
		AddToLeaderboard(leaderboard, "a")
		AddToLeaderboard(leaderboard, "c")
	}

	//the top should be a, b, c in that order
	if leaderboard.top[0] != "a" {
		t.Errorf("Expected top[0] to be a, got %s", leaderboard.top[0])
	}

	if leaderboard.top[1] != "b" {
		t.Errorf("Expected top[1] to be b, got %s", leaderboard.top[1])
	}

	if leaderboard.top[2] != "c" {
		t.Errorf("Expected top[2] to be c, got %s", leaderboard.top[2])
	}
}

func TestBubbleDown(t *testing.T) {
	lb := &Leaderboard{
		top: []string{"a", "b", "c"},
	}

	bubbleDown(lb, 0, "a")
	if !Equal(lb.top, []string{"a", "b", "c"}) {
		t.Errorf("Expected top to be [a, b, c], got %v", lb.top)
	}

	bubbleDown(lb, 0, "a")
	if !Equal(lb.top, []string{"a", "b", "c"}) {
		t.Errorf("Expected top to be [a, b, c], got %v", lb.top)
	}

	bubbleDown(lb, 0, "d")
	if !Equal(lb.top, []string{"d", "a", "b"}) {
		t.Errorf("Expected top to be [d, a, b], got %v", lb.top)
	}

	bubbleDown(lb, 1, "i")
	if !Equal(lb.top, []string{"d", "i", "a"}) {
		t.Errorf("Expected top to be [d, i, a], got %v", lb.top)
	}

	bubbleDown(lb, 1, "i")
	if !Equal(lb.top, []string{"d", "i", "a"}) {
		t.Errorf("Expected top to be [d, i, a], got %v", lb.top)
	}

	bubbleDown(lb, 0, "i")
	if !Equal(lb.top, []string{"i", "d", "a"}) {
		t.Errorf("Expected top to be [i, d, a], got %v", lb.top)
	}
}

//a quick benchmark of adding to the leaderboard
func BenchmarkAddToLeaderboard(b *testing.B) {
	leaderboard := &Leaderboard{
		counters: map[string]int{"a": 10, "b": 20, "c": 30,
			"d": 40, "e": 50, "f": 60, "g": 70, "h": 80, "i": 90, "j": 100,
			"k": 110, "l": 120, "m": 130, "n": 140, "o": 150, "p": 160, "q": 170, "r": 180,
			"s": 190, "t": 200, "u": 210, "v": 220, "w": 230, "x": 240, "y": 250, "z": 260},
		top:         []string{"q", "r", "s", "t", "u", "v", "w", "x", "y", "z"},
		numberOfTop: 10,
	}
	AddToLeaderboard(leaderboard, "b")
	AddToLeaderboard(leaderboard, "c")
	AddToLeaderboard(leaderboard, "d") //make sure it doesn't get added to the top initially

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		AddToLeaderboard(leaderboard, "a")
	}
}

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
