package util

import "testing"

func TestTermIDToName(t *testing.T) {
	//test cases
	testCases := []struct {
		termID int
		name   string
	}{
		{202100, "Summer 2020"},
		{202101, "Fall 2020"},
		{202102, "Winter 2021"},
		{202103, "Spring 2021"},
	}

	//run test cases
	for _, testCase := range testCases {
		name := TermIDToName(testCase.termID)
		if name != testCase.name {
			t.Errorf("TermIDToName(%d) = %s, want %s", testCase.termID, name, testCase.name)
		}
	}
}
