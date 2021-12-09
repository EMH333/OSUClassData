package util

import (
	"regexp"
	"strconv"
)

//function to convert a string to a int
func StringToIntPanic(s string) int {
	//convert the string to a number
	num, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	//return the number
	return num
}

//function to convert a string to a float
func StringToFloatPanic(s string) float64 {
	//convert the string to a number
	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	//return the number
	return num
}

var onlyNum *regexp.Regexp = regexp.MustCompile("[0-9]+")

//converts a messy string (w/ other characters) to a number
func IntFromMessyString(s string) int {
	//convert the string to a number
	num, err := strconv.Atoi(onlyNum.FindString(s))
	if err != nil {
		return 0
	}
	//return the number
	return num
}
