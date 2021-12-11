package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"ethohampton.com/OSUClassData/internal/database"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "OSUClassData",
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected to Database!")

	//hello world web server
	http.Handle("/", http.FileServer(http.Dir("frontend/dist")))
	http.HandleFunc("/api/v1/classes", getClasses)
	http.HandleFunc("/api/v1/class", getClass)

	http.ListenAndServe(":8080", nil)

}

func getClasses(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT DISTINCT ClassIdentifier FROM Classes WHERE Visible=TRUE")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var classList []string

	for rows.Next() {
		var class string
		err := rows.Scan(&class)
		if err != nil {
			//just ignore any weird errors
			log.Println(err)
			continue
		}
		classList = append(classList, class)
	}

	jsonResponse, err := json.Marshal(classList)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func getClass(w http.ResponseWriter, r *http.Request) {
	//TODO switch to path routing ie /api/v1/class/:classIdentifier/:term
	class := r.URL.Query().Get("class")
	if class == "" {
		http.Error(w, "Missing class parameter", http.StatusBadRequest)
		return
	}

	term := r.URL.Query().Get("term")
	if term == "" {
		http.Error(w, "Missing term parameter", http.StatusBadRequest)
		return
	}

	row := db.QueryRow(`SELECT 
	ClassIdentifier,TermID,Students,Credits,ClassGPA,
	A,AMinus,B,BPlus,BMinus,C,CPlus,CMinus,D,DPlus,DMinus,F,S,U,W,I
	FROM Classes WHERE ClassIdentifier=? AND TermID=? AND Visible=TRUE`, class, term)
	if row == nil {
		http.Error(w, "Class not found", http.StatusNotFound)
		return
	}

	var classData database.Class
	row.Scan(&classData.ClassIdentifier, &classData.TermID, &classData.Students, &classData.Credits, &classData.ClassGPA,
		&classData.A, &classData.AMinus,
		&classData.B, &classData.BPlus, &classData.BMinus,
		&classData.C, &classData.CPlus, &classData.CMinus,
		&classData.D, &classData.DPlus, &classData.DMinus,
		&classData.F,
		&classData.S, &classData.U, &classData.W, &classData.I)
	classData.Visible = true

	if classData.ClassIdentifier == "" {
		http.Error(w, "Class not found", http.StatusNotFound)
		return
	}

	jsonResponse, err := json.Marshal(classData)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
