package util

import "strconv"

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
