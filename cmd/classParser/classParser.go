package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"ethohampton.com/OSUClassData/internal/database"
	"ethohampton.com/OSUClassData/internal/util"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: classParser <csv file> <term>")
		os.Exit(1)
	}

	var term = os.Args[2]

	// open file
	f, err := os.Open(os.Args[1])
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
			var header = strings.Split(v[0], " ")

			var c database.Class
			c.ClassIdentifier = strings.TrimSpace(header[1] + header[2])
			c.Students = util.StringToIntPanic(header[6])
			c.ClassGPA = util.StringToFloatPanic(header[21])
			c.TermID = term

			if header[10] == "Variable" {
				//TODO we don't know how to handle this
				//continue
				log.Printf("Variable Class: %s\n", strings.TrimSpace(header[1]+header[2]))
				c.Credits = -1
			} else {
				c.Credits = util.StringToIntPanic(header[10])
			}

			classes = append(classes, c)
		}

		if strings.HasPrefix(v[0], "Grade: ") {
			var c = &classes[len(classes)-1]
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
				continue
			}
		}
	}

	//TODO also create SQL for ClassInfo table (could also do that in a follow up query)

	var upsertString = "REPLACE INTO Classes (classIdentifier, termID, students, credits, a, aMinus, b, bPlus, bMinus, c, cPlus, cMinus, d, dPlus, dMinus, f, s, u, w, i, classGPA) VALUES\n"
	//print the classes
	for index, v := range classes {
		//log.Println(v)
		upsertString = upsertString + " ROW('" + v.ClassIdentifier + "', '" + v.TermID + "', " + fmt.Sprint(v.Students) + ", " + fmt.Sprint(v.Credits) + ", " +
			fmt.Sprint(v.A) + ", " + fmt.Sprint(v.AMinus) + ", " +
			fmt.Sprint(v.B) + ", " + fmt.Sprint(v.BPlus) + ", " + fmt.Sprint(v.BMinus) + ", " +
			fmt.Sprint(v.C) + ", " + fmt.Sprint(v.CPlus) + ", " + fmt.Sprint(v.CMinus) + ", " +
			fmt.Sprint(v.D) + ", " + fmt.Sprint(v.DPlus) + ", " + fmt.Sprint(v.DMinus) + ", " +
			fmt.Sprint(v.F) + ", " +
			fmt.Sprint(v.S) + ", " + fmt.Sprint(v.U) + ", " + fmt.Sprint(v.W) + ", " + fmt.Sprint(v.I) + ", " + strconv.FormatFloat(v.ClassGPA, 'f', -1, 64) + ")"
		if index != len(classes)-1 {
			upsertString = upsertString + ",\n"
		} else {
			upsertString = upsertString + ";"
		}
	}

	output, err := os.Create(string(term) + "Classes.sql")
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	// write the upsert string to the file
	if _, err := output.WriteString(upsertString); err != nil {
		log.Fatal(err)
	}

}
