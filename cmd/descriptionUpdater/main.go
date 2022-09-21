package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: description_updater <csv file>")
		os.Exit(1)
	}

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
	_, err = csvReader.Read()      //skip header
	if err != nil {
		log.Fatal(err)
	}
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	hadToSkip := false
	for _, v := range data {
		if strings.Contains(v[1], "'") { //|| strings.Contains(v[1], "\"") || strings.Contains(v[1], "\\") || strings.Contains(v[1], ";") {
			hadToSkip = true
			continue
		}
		fmt.Printf("UPDATE ClassInfo SET ClassDescription='%s' WHERE ClassIdentifier='%s';\n", v[1], v[0])
	}

	if hadToSkip {
		fmt.Println("Had to skip some classes because they contained special characters")
	}
}
