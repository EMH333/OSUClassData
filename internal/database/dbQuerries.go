package database

import (
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("item not found")

func GetTermClass(db *sql.DB, id string, term string) (Class, error) {
	row := db.QueryRow(`SELECT 
	ClassIdentifier,TermID,Students,Credits,ClassGPA,
	A,AMinus,B,BPlus,BMinus,C,CPlus,CMinus,D,DPlus,DMinus,F,S,U,W,I
	FROM Classes WHERE ClassIdentifier=? AND TermID=? AND Visible=TRUE`, id, term)
	if row == nil {
		return Class{}, ErrNotFound
	}

	var classData Class
	row.Scan(&classData.ClassIdentifier, &classData.TermID, &classData.Students, &classData.Credits, &classData.ClassGPA,
		&classData.A, &classData.AMinus,
		&classData.B, &classData.BPlus, &classData.BMinus,
		&classData.C, &classData.CPlus, &classData.CMinus,
		&classData.D, &classData.DPlus, &classData.DMinus,
		&classData.F,
		&classData.S, &classData.U, &classData.W, &classData.I)
	classData.Visible = true
	return classData, nil
}

type ClassInfoResponse struct {
	ClassIdentifier    string
	ClassName          string
	LastTerm           string // Last term the class was taught in data we have
	Credits            int    //TODO deal with variable credit classes
	AverageGPA         float64
	AverageGPALastTerm float64
	AverageStudents    float64
	StudentsLastTerm   int
	WithdrawlRate      float64
}

//TODO remove withdrawled students from total count
//TODO combine into one query
func GetClassInfo(db *sql.DB, id string) (ClassInfoResponse, error) {
	var classInfoQuery = "SELECT Credits, ClassName FROM ClassInfo WHERE ClassIdentifier=?"
	var lastTermQuery = "SELECT TermID FROM Classes WHERE ClassIdentifier=? AND Visible=TRUE ORDER BY TermID DESC LIMIT 1"
	var lastTermInfo = "SELECT ClassGPA, Students FROM Classes WHERE ClassIdentifier=? AND TermID=? AND Visible=TRUE"
	var averageInfo = "SELECT AVG(ClassGPA), AVG(Students), SUM(W)/SUM(Students) AS WithdrawlRate FROM Classes WHERE ClassIdentifier=? AND Visible=TRUE"

	var classData ClassInfoResponse
	classData.ClassIdentifier = id

	// Get name and credits
	row := db.QueryRow(classInfoQuery, id)
	if row.Err() != nil {
		return ClassInfoResponse{}, row.Err()
	}
	row.Scan(&classData.Credits, &classData.ClassName)

	// Get the last term the class was taught in
	row = db.QueryRow(lastTermQuery, id)
	if row.Err() != nil {
		return ClassInfoResponse{}, row.Err()
	}
	row.Scan(&classData.LastTerm)

	// Get the last term info
	row = db.QueryRow(lastTermInfo, id, classData.LastTerm)
	if row.Err() != nil {
		return ClassInfoResponse{}, row.Err()
	}
	row.Scan(&classData.AverageGPALastTerm, &classData.StudentsLastTerm)

	// Get all the other averages and info
	row = db.QueryRow(averageInfo, id)
	if row.Err() != nil {
		return ClassInfoResponse{}, row.Err()
	}
	row.Scan(&classData.AverageGPA, &classData.AverageStudents, &classData.WithdrawlRate)

	return classData, nil
}

type StudentsPerTermResponse struct {
	Terms    []string
	Students []int
}

func GetStudentsPerTerm(db *sql.DB, id string) (StudentsPerTermResponse, error) {
	var query = "SELECT TermID, Students FROM Classes WHERE ClassIdentifier=? AND Visible=TRUE"
	var response StudentsPerTermResponse
	response.Terms = make([]string, 0)
	response.Students = make([]int, 0)

	rows, err := db.Query(query, id)
	if err != nil {
		return StudentsPerTermResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var term string
		var students int
		rows.Scan(&term, &students)
		response.Terms = append(response.Terms, term)
		response.Students = append(response.Students, students)
	}
	return response, nil
}

type AvgGPAPerTermResponse struct {
	Terms []string
	GPA   []float64
}

func GetAvgGPAPerTerm(db *sql.DB, id string) (AvgGPAPerTermResponse, error) {
	var query = "SELECT TermID, ClassGPA FROM Classes WHERE ClassIdentifier=? AND Visible=TRUE"
	var response AvgGPAPerTermResponse
	response.Terms = make([]string, 0)
	response.GPA = make([]float64, 0)

	rows, err := db.Query(query, id)
	if err != nil {
		return AvgGPAPerTermResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var term string
		var GPA float64
		rows.Scan(&term, &GPA)
		response.Terms = append(response.Terms, term)
		response.GPA = append(response.GPA, GPA)
	}
	return response, nil
}
