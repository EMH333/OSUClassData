package database

import "testing"

func TestGetClassName(t *testing.T) {
	name, err := getClassName("CS162")
	if err != nil {
		t.Error(err)
	}

	if name != "INTRODUCTION TO COMPUTER SCIENCE II" {
		t.Error("Expected 'INTRODUCTION TO COMPUTER SCIENCE II', got ", name)
	}
}

func TestNormalizeName(t *testing.T) {
	testCases := map[string]string{
		"Computer Science":                    "CS",
		"Computer Science II":                 "CS II",
		"INTRODUCTION TO COMPUTER SCIENCE II": "Intro to CS II",
		"*SOMETHING":                          "Something",
		"ETHICS IN CS":                        "Ethics in CS",
		"THE INTRODUCTION CLASS TO A TIME ABOUT THE THINGS": "The Intro Class to a Time About the Things",
		"AN EPIC CLASS IN AN EPIC PLACE":                    "An Epic Class in an Epic Place",
	}

	for input, expected := range testCases {
		if actual := normalizeName(input); actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	}
}
