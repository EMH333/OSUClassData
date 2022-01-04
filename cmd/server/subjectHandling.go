package main

import (
	"encoding/json"
	"log"
	"net/http"

	"ethohampton.com/OSUClassData/internal/database"
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

	jsonResponse, err := json.Marshal(classList)
	if err != nil {
		http.Error(w, "Error marshaling JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "Error writing JSON response", http.StatusInternalServerError)
	}
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

	jsonResponse, err := json.Marshal(GPAPerTerm)
	if err != nil {
		http.Error(w, "Error marshaling JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "Error writing JSON response", http.StatusInternalServerError)
	}
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

	jsonResponse, err := json.Marshal(WithdrawalPerTerm)
	if err != nil {
		http.Error(w, "Error marshaling JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "Error writing JSON response", http.StatusInternalServerError)
	}
}
