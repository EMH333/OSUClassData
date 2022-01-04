package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"ethohampton.com/OSUClassData/internal/database"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "8080")
	}

	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "OSUClassData",
	}

	// Try to connect to the database
	const maxConnectionAttempts = 10
	for i := 0; i < maxConnectionAttempts; i++ {
		err := tryToConnectToDB(&cfg)
		//if connected then start server
		if err == nil {
			break
		}

		//if not connected then wait and try again, up to 10 times, then exit with error
		if err != nil && i == maxConnectionAttempts-1 {
			log.Fatal(err)
		}

		//wait 5 seconds before trying again
		fmt.Println("Could not connect to DB. Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
	}

	fmt.Println("Connected to Database!")

	//hello world web server
	http.Handle("/", http.FileServer(http.Dir("frontend/dist")))
	http.HandleFunc("/api/v0/status", getStatus)
	http.HandleFunc("/api/v0/classes", getClasses)
	http.HandleFunc("/api/v0/class", getClass)
	http.HandleFunc("/api/v0/classInfo", getClassInfo)
	http.HandleFunc("/api/v0/chart/studentsPerTerm", getStudentsPerTerm)
	http.HandleFunc("/api/v0/chart/avgGPAPerTerm", getAvgGPAPerTerm)
	http.HandleFunc("/api/v0/chart/withdrawalRatePerTerm", getWithdrawalRatePerTerm)
	http.HandleFunc("/api/v0/chart/lastTermGradeDistribution", getLastTermGradeDistribution)

	http.HandleFunc("/api/v0/subjects", getSubjects)
	http.HandleFunc("/api/v0/subject/chart/avgGPAPerTerm", getSubjectAvgGPAPerTerm)
	http.HandleFunc("/api/v0/subject/chart/withdrawalRatePerTerm", getSubjectWithdrawalRatePerTerm)

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)

}

func tryToConnectToDB(config *mysql.Config) error {
	var err error
	// Connect to the database
	db, err = sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return err
	}

	// Test the connection to the database with a ping
	err = db.Ping()
	if err != nil {
		return err
	}

	// run query to test connection
	var test int
	err = db.QueryRow("SELECT 1").Scan(&test)
	if err != nil {
		return err
	}

	return nil
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
			http.Error(w, "Error reading classes", http.StatusInternalServerError)
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

	classData, err := database.GetTermClass(db, class, term)
	if err != nil {
		http.Error(w, "Class not found", http.StatusNotFound)
		return
	}

	jsonResponse, err := json.Marshal(classData)
	if err != nil {
		http.Error(w, "Error marshaling JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func getClassInfo(w http.ResponseWriter, r *http.Request) {
	class := r.URL.Query().Get("class")
	if class == "" {
		http.Error(w, "Missing class parameter", http.StatusBadRequest)
		return
	}

	classData, err := database.GetClassInfo(db, class)
	if err != nil {
		http.Error(w, "Class not found", http.StatusNotFound)
		return
	}

	jsonResponse, err := json.Marshal(classData)
	if err != nil {
		http.Error(w, "Error marshaling JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	pingErr := db.Ping()
	if pingErr != nil {
		http.Error(w, "Can't ping DB", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("OK"))
}

func getStudentsPerTerm(w http.ResponseWriter, r *http.Request) {
	class := r.URL.Query().Get("class")
	if class == "" {
		http.Error(w, "Missing class parameter", http.StatusBadRequest)
		return
	}

	studentsPerTerm, err := database.GetStudentsPerTerm(db, class)
	if err != nil {
		http.Error(w, "Class not found", http.StatusNotFound)
		return
	}

	jsonResponse, err := json.Marshal(studentsPerTerm)
	if err != nil {
		http.Error(w, "Error marshaling JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func getAvgGPAPerTerm(w http.ResponseWriter, r *http.Request) {
	class := r.URL.Query().Get("class")
	if class == "" {
		http.Error(w, "Missing class parameter", http.StatusBadRequest)
		return
	}

	GPAPerTerm, err := database.GetAvgGPAPerTerm(db, class)
	if err != nil {
		http.Error(w, "Class not found", http.StatusNotFound)
		return
	}

	jsonResponse, err := json.Marshal(GPAPerTerm)
	if err != nil {
		http.Error(w, "Error marshaling JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func getWithdrawalRatePerTerm(w http.ResponseWriter, r *http.Request) {
	class := r.URL.Query().Get("class")
	if class == "" {
		http.Error(w, "Missing class parameter", http.StatusBadRequest)
		return
	}

	WithdrawalPerTerm, err := database.GetWithdrawalRatePerTerm(db, class)
	if err != nil {
		http.Error(w, "Class not found", http.StatusNotFound)
		return
	}

	jsonResponse, err := json.Marshal(WithdrawalPerTerm)
	if err != nil {
		http.Error(w, "Error marshaling JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func getLastTermGradeDistribution(w http.ResponseWriter, r *http.Request) {
	class := r.URL.Query().Get("class")
	if class == "" {
		http.Error(w, "Missing class parameter", http.StatusBadRequest)
		return
	}

	latestClass, err := database.GetLastTermClass(db, class)
	if err != nil {
		http.Error(w, "Class not found", http.StatusNotFound)
		return
	}

	jsonResponse, err := json.Marshal(latestClass)
	if err != nil {
		http.Error(w, "Error marshaling JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
