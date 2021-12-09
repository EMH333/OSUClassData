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
			c.ClassIdentifier = strings.TrimSpace(header[1] + header[2])
			c.Students = util.StringToIntPanic(header[6])
			c.Credits = util.StringToIntPanic(header[10])
			c.ClassGPA = util.StringToFloatPanic(header[21])
			c.TermID = term

			classes = append(classes, c)
		}

		if strings.HasPrefix(v[0], "Grade: ") {
			var c *database.Class = &classes[len(classes)-1]
			if util.StringToIntPanic(string(v[1][len(v[1])-1])) == c.Credits {
				var letter = strings.Split(v[0], " ")[1]
				var students = util.IntFromMessyString(v[2])
				//log.Printf("%s %d\n", letter, students)
				switch letter {
				case "A":
					c.A += students
				case "A-":
					c.AMinus += students
				case "B":
					c.B += students
				case "B+":
					c.BPlus += students
				case "B-":
					c.BMinus += students
				case "C":
					c.C += students
				case "C+":
					c.CPlus += students
				case "C-":
					c.CMinus += students
				case "D":
					c.D += students
				case "D+":
					c.DPlus += students
				case "D-":
					c.DMinus += students
				case "F":
					c.F += students
				case "S":
					c.S += students
				case "U":
					c.U += students
				case "W":
					c.W += students
				}

				//deal with all incomplete grades as one
				if strings.HasPrefix(letter, "I") {
					c.I += students
				}
			} else {
				//TODO handle credits for variable credit classes
			}
		}
	}

	//print the classes
	for _, v := range classes {
		log.Println(v)
	}
}
