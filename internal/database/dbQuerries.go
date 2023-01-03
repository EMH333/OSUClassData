package database

import (
	"database/sql"
	"errors"
	"strings"
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
	err := row.Scan(&classData.ClassIdentifier, &classData.TermID, &classData.Students, &classData.Credits, &classData.ClassGPA,
		&classData.A, &classData.AMinus,
		&classData.B, &classData.BPlus, &classData.BMinus,
		&classData.C, &classData.CPlus, &classData.CMinus,
		&classData.D, &classData.DPlus, &classData.DMinus,
		&classData.F,
		&classData.S, &classData.U, &classData.W, &classData.I)
	if err != nil {
		return Class{}, err
	}

	classData.Visible = true
	return classData, nil
}

func GetLastTermClass(db *sql.DB, id string) (Class, error) {
	row := db.QueryRow(`SELECT 
	ClassIdentifier,TermID,Students,Credits,ClassGPA,
	A,AMinus,B,BPlus,BMinus,C,CPlus,CMinus,D,DPlus,DMinus,F,S,U,W,I
	FROM Classes WHERE ClassIdentifier=? AND Visible=TRUE ORDER BY TermID DESC LIMIT 1`, id)
	if row == nil {
		return Class{}, ErrNotFound
	}

	var classData Class
	err := row.Scan(&classData.ClassIdentifier, &classData.TermID, &classData.Students, &classData.Credits, &classData.ClassGPA,
		&classData.A, &classData.AMinus,
		&classData.B, &classData.BPlus, &classData.BMinus,
		&classData.C, &classData.CPlus, &classData.CMinus,
		&classData.D, &classData.DPlus, &classData.DMinus,
		&classData.F,
		&classData.S, &classData.U, &classData.W, &classData.I)
	if err != nil {
		return Class{}, err
	}

	classData.Visible = true
	return classData, nil
}

type ClassInfoResponse struct {
	ClassIdentifier    string
	ClassName          string
	ClassDescription   string
	LastTerm           string // Last term the class was taught in data we have
	Credits            int    //TODO deal with variable credit classes
	AverageGPA         float64
	AverageGPALastTerm float64
	AverageStudents    float64
	StudentsLastTerm   int
	WithdrawalRate     float64
}

//TODO remove withdrawn students from total count
//TODO combine into one query
//returns classinfo response, if the class name should be updated and error
func GetClassInfo(db *sql.DB, id string) (ClassInfoResponse, bool, error) {
	var classInfoQuery = "SELECT Credits, RetrievedClassName, NormalizedClassName, ClassName, ClassDescription FROM ClassInfo WHERE ClassIdentifier=?"
	var lastTermQuery = "SELECT TermID FROM Classes WHERE ClassIdentifier=? AND Visible=TRUE ORDER BY TermID DESC LIMIT 1"
	var lastTermInfo = "SELECT ClassGPA, Students FROM Classes WHERE ClassIdentifier=? AND TermID=? AND Visible=TRUE"
	var averageInfo = "SELECT AVG(ClassGPA), AVG(Students), SUM(W)/SUM(Students) AS WithdrawalRate FROM Classes WHERE ClassIdentifier=? AND Visible=TRUE"

	var classNamedRetrieved bool
	var classNameNormalized bool

	var classData ClassInfoResponse
	classData.ClassIdentifier = strings.Clone(id) // needed to prevent fasthttp from reusing the underlying buffer

	// Get name and credits
	row := db.QueryRow(classInfoQuery, id)
	if row.Err() != nil {
		return ClassInfoResponse{}, false, row.Err()
	}
	_ = row.Scan(&classData.Credits, &classNamedRetrieved, &classNameNormalized, &classData.ClassName, &classData.ClassDescription) // we expect errors here

	// Get the last term the class was taught in
	row = db.QueryRow(lastTermQuery, id)
	if row.Err() != nil {
		return ClassInfoResponse{}, false, row.Err()
	}
	err := row.Scan(&classData.LastTerm)
	if err != nil {
		return ClassInfoResponse{}, false, err
	}

	// Get the last term info
	row = db.QueryRow(lastTermInfo, id, classData.LastTerm)
	if row.Err() != nil {
		return ClassInfoResponse{}, false, row.Err()
	}
	err = row.Scan(&classData.AverageGPALastTerm, &classData.StudentsLastTerm)
	if err != nil {
		return ClassInfoResponse{}, false, err
	}

	// Get all the other averages and info
	row = db.QueryRow(averageInfo, id)
	if row.Err() != nil {
		return ClassInfoResponse{}, false, row.Err()
	}
	err = row.Scan(&classData.AverageGPA, &classData.AverageStudents, &classData.WithdrawalRate)
	if err != nil {
		return ClassInfoResponse{}, false, err
	}

	return classData, !(classNameNormalized && classNamedRetrieved), nil
}

type CombinedClassResponse struct {
	Terms          []string
	WithdrawalRate []float64
	GPA            []float64
	Students       []float64
}

func GetCombinedClassStats(db *sql.DB, id string) (CombinedClassResponse, error) {
	var query = "SELECT TermID, (W / Students) AS WithdrawalRate, ClassGPA, Students FROM Classes WHERE ClassIdentifier=? AND Visible=TRUE"
	var response CombinedClassResponse
	response.Terms = make([]string, 0)
	response.WithdrawalRate = make([]float64, 0)
	response.GPA = make([]float64, 0)
	response.Students = make([]float64, 0)

	rows, err := db.Query(query, id)
	if err != nil {
		return CombinedClassResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var term string
		var withdrawalRate float64
		var gpa float64
		var students float64
		err := rows.Scan(&term, &withdrawalRate, &gpa, &students)
		if err != nil {
			return CombinedClassResponse{}, err
		}

		response.Terms = append(response.Terms, term)
		response.WithdrawalRate = append(response.WithdrawalRate, withdrawalRate)
		response.GPA = append(response.GPA, gpa)
		response.Students = append(response.Students, students)
	}
	return response, nil
}
