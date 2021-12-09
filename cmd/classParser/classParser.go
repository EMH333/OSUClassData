package main

import (
	"encoding/csv"
	"log"
	"os"
	"strings"

	"ethohampton.com/OSUClassData/internal/database"
	"ethohampton.com/OSUClassData/internal/util"
)

func main() {
	//TODO get term we are reading from
	var term = "202001"

	// open file
	f, err := os.Open("classes.csv")
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	csvReader.FieldsPerRecord = -1 //there will be varying number of fields in each record
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	//array of Class
	var classes []database.Class

	// print the data
	for _, v := range data {
		//if v[0] starts with "Course" then print the row
		if strings.HasPrefix(v[0], "Course") {
			//log.Println(v[0])
			var header []string = strings.Split(v[0], " ")

			if header[10] == "Variable" {
				//TODO we don't know how to handle this
				continue
			}

			var c database.Class
			c.ClassIdentifier = header[1] + header[2]
			c.Students = util.StringToIntPanic(header[6])
			c.Credits = util.StringToIntPanic(header[10])
			c.ClassGPA = util.StringToFloatPanic(header[21])
			c.TermID = term

			classes = append(classes, c)
		}
	}

	//print the classes
	for _, v := range classes {
		log.Println(v)
	}

	println("Hello, World!")
}
