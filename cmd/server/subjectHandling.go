package main

import (
	"log"
	"net/http"

	"ethohampton.com/OSUClassData/internal/database"
	"ethohampton.com/OSUClassData/internal/util"
)

func getSubjects(w http.ResponseWriter, r *http.Request) {
	//The substring is to remove the H for honors classes
	const query = `SELECT DISTINCT REPLACE
	(REPLACE
	(REPLACE
	(REPLACE
	(REPLACE
	(REPLACE
	(REPLACE
	(REPLACE
	(REPLACE
	(REPLACE (
		SUBSTR(ClassIdentifier FROM 1 FOR 4),
	'0', ''),
	'1', ''),
	'2', ''),
	'3', ''),
	'4', ''),
	'5', ''),
	'6', ''),
	'7', ''),
	'8', ''),
	'9', '') FROM Classes WHERE Visible=TRUE`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var classList []string

	for rows.Next() {
		var class string
		err := rows.Scan(&class)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Error reading subjects", http.StatusInternalServerError)
			return
		}
		classList = append(classList, class)
	}

	util.WriteJSON(w, classList)
}

func getSubjectAvgGPAPerTerm(w http.ResponseWriter, r *http.Request) {
	subject := r.URL.Query().Get("subject")
	if subject == "" {
		http.Error(w, "Missing subject parameter", http.StatusBadRequest)
		return
	}

	if subject == "University-wide" {
		subject = ""
	}

	GPAPerTerm, err := database.GetSubjectAvgGPAPerTerm(db, subject)
	if err != nil {
		http.Error(w, "Subject not found", http.StatusNotFound)
		return
	}

	util.WriteJSON(w, GPAPerTerm)
}

func getSubjectWithdrawalRatePerTerm(w http.ResponseWriter, r *http.Request) {
	subject := r.URL.Query().Get("subject")
	if subject == "" {
		http.Error(w, "Missing subject parameter", http.StatusBadRequest)
		return
	}

	if subject == "University-wide" {
		subject = ""
	}

	WithdrawalPerTerm, err := database.GetSubjectWithdrawalRatePerTerm(db, subject)
	if err != nil {
		http.Error(w, "Subject not found", http.StatusNotFound)
		return
	}

	util.WriteJSON(w, WithdrawalPerTerm)
}
