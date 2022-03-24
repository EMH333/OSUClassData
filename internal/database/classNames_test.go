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
	name := normalizeName("INTRODUCTION TO COMPUTER SCIENCE II")
	if name != "Intro to Computer Science II" {
		t.Error("Expected 'Intro to Computer Science II', got ", name)
	}
}
