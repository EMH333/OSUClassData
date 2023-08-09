package util

import (
	"testing"
)

func TestBasicLeaderboard(t *testing.T) {
	leaderboard := &Leaderboard{
		counters: make(map[string]int),
		Top:      []string{},
	}
	SetUpLeaderboard(leaderboard)

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

	if leaderboard.Top[0] != "a" {
		t.Errorf("Expected top[0] to be a, got %s", leaderboard.Top[0])
	}

	if leaderboard.TotalCount != 5 {
		t.Errorf("Expected totalCount to be 5, got %d", leaderboard.TotalCount)
	}
}

func TestFullLeaderboard(t *testing.T) {
	leaderboard := &Leaderboard{
		counters:    make(map[string]int),
		Top:         []string{},
		NumberOfTop: 3,
	}
	SetUpLeaderboard(leaderboard)

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
	if leaderboard.Top[0] != "a" {
		t.Errorf("Expected top[0] to be a, got %s", leaderboard.Top[0])
	}

	if leaderboard.Top[1] != "b" {
		t.Errorf("Expected top[1] to be b, got %s", leaderboard.Top[1])
	}

	if leaderboard.Top[2] != "c" {
		t.Errorf("Expected top[2] to be c, got %s", leaderboard.Top[2])
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
		Top:         []string{},
		NumberOfTop: 3,
	}
	SetUpLeaderboard(leaderboard)

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
	if leaderboard.Top[0] != "a" {
		t.Errorf("Expected top[0] to be a, got %s", leaderboard.Top[0])
	}

	if leaderboard.Top[1] != "b" {
		t.Errorf("Expected top[1] to be b, got %s", leaderboard.Top[1])
	}

	if leaderboard.Top[2] != "c" {
		t.Errorf("Expected top[2] to be c, got %s", leaderboard.Top[2])
	}
}

func TestBubbleDown(t *testing.T) {
	lb := &Leaderboard{
		Top: []string{"a", "b", "c"},
	}
	SetUpLeaderboard(lb)

	bubbleDown(lb, 0, "a")
	if !Equal(lb.Top, []string{"a", "b", "c"}) {
		t.Errorf("Expected top to be [a, b, c], got %v", lb.Top)
	}

	bubbleDown(lb, 0, "a")
	if !Equal(lb.Top, []string{"a", "b", "c"}) {
		t.Errorf("Expected top to be [a, b, c], got %v", lb.Top)
	}

	bubbleDown(lb, 0, "d")
	if !Equal(lb.Top, []string{"d", "a", "b"}) {
		t.Errorf("Expected top to be [d, a, b], got %v", lb.Top)
	}

	bubbleDown(lb, 1, "i")
	if !Equal(lb.Top, []string{"d", "i", "a"}) {
		t.Errorf("Expected top to be [d, i, a], got %v", lb.Top)
	}

	bubbleDown(lb, 1, "i")
	if !Equal(lb.Top, []string{"d", "i", "a"}) {
		t.Errorf("Expected top to be [d, i, a], got %v", lb.Top)
	}

	bubbleDown(lb, 0, "i")
	if !Equal(lb.Top, []string{"i", "d", "a"}) {
		t.Errorf("Expected top to be [i, d, a], got %v", lb.Top)
	}
}

func TestDecayLeaderboard(t *testing.T) {
	leaderboard := &Leaderboard{}
	SetUpLeaderboard(leaderboard)

	AddToLeaderboard(leaderboard, "a")
	AddToLeaderboard(leaderboard, "a")
	AddToLeaderboard(leaderboard, "b")
	AddToLeaderboard(leaderboard, "c")
	AddToLeaderboard(leaderboard, "d")

	if len(leaderboard.Top) != 4 {
		t.Errorf("Expected length of top to be 4, got %v", leaderboard.Top)
	}

	DecayLeaderboard(leaderboard)

	if len(leaderboard.Top) != 1 {
		t.Errorf("Expected length of top to be 1, got %v", leaderboard.Top)
	}

	DecayLeaderboard(leaderboard)

	if len(leaderboard.Top) != 0 {
		t.Errorf("Expected top to be empty, got %v", leaderboard.Top)
	}

	AddToLeaderboard(leaderboard, "a")

	if len(leaderboard.Top) != 1 {
		t.Errorf("Expected length of top to be 1, got %v", leaderboard.Top)
	}

	if leaderboard.Top[0] != "a" {
		t.Errorf("Expected top[0] to be a, got %s", leaderboard.Top[0])
	}
}

func TestChanceDecayLeaderboard(t *testing.T) {
	leaderboard := &Leaderboard{
		DecayChance: -1,
	}
	SetUpLeaderboard(leaderboard)

	AddToLeaderboard(leaderboard, "a")
	AddToLeaderboard(leaderboard, "a")
	AddToLeaderboard(leaderboard, "b")
	AddToLeaderboard(leaderboard, "c")
	AddToLeaderboard(leaderboard, "d")

	if len(leaderboard.Top) != 4 {
		t.Errorf("Expected length of top to be 4, got %v", leaderboard.Top)
	}

	DecayLeaderboard(leaderboard)

	if len(leaderboard.Top) != 4 {
		t.Errorf("Expected length of top to be 4, got %v", leaderboard.Top)
	}
}

// Make sure that the decay function doesn't just always or never decay
func TestHalfChanceDecayLeaderboard(t *testing.T) {
	var decayed = false
	// give it 10 chances. If it can't do it then, there is likely a problem
	for i := 0; i < 10; i++ {
		leaderboard := &Leaderboard{
			DecayChance: 50,
		}
		SetUpLeaderboard(leaderboard)

		AddToLeaderboard(leaderboard, "a")
		AddToLeaderboard(leaderboard, "a")
		AddToLeaderboard(leaderboard, "b")
		AddToLeaderboard(leaderboard, "b")
		AddToLeaderboard(leaderboard, "c")
		AddToLeaderboard(leaderboard, "c")
		AddToLeaderboard(leaderboard, "d")
		AddToLeaderboard(leaderboard, "d")
		AddToLeaderboard(leaderboard, "e")
		AddToLeaderboard(leaderboard, "e")

		if len(leaderboard.Top) != 5 {
			t.Errorf("Expected length of top to be 5, got %v", leaderboard.Top)
		}

		DecayLeaderboard(leaderboard)
		DecayLeaderboard(leaderboard)

		if !(len(leaderboard.Top) == 5 || len(leaderboard.Top) == 0) {
			decayed = true
			break
		}
	}

	if !decayed {
		t.Errorf("Expected length of top to be less than 5 but more than 0")
	}
}

//a quick benchmark of adding to the leaderboard
func BenchmarkAddToLeaderboard(b *testing.B) {
	leaderboard := &Leaderboard{
		counters: map[string]int{"a": 10, "b": 20, "c": 30,
			"d": 40, "e": 50, "f": 60, "g": 70, "h": 80, "i": 90, "j": 100,
			"k": 110, "l": 120, "m": 130, "n": 140, "o": 150, "p": 160, "q": 170, "r": 180,
			"s": 190, "t": 200, "u": 210, "v": 220, "w": 230, "x": 240, "y": 250, "z": 260},
		Top:         []string{"q", "r", "s", "t", "u", "v", "w", "x", "y", "z"},
		NumberOfTop: 10,
	}
	SetUpLeaderboard(leaderboard)

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
