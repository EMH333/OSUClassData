package database

import "database/sql"

type AvgGPAPerTermResponse struct {
	Terms []string
	GPA   []float64
}

func GetSubjectAvgGPAPerTerm(db *sql.DB, id string) (AvgGPAPerTermResponse, error) {
	var query = "SELECT TermID, AVG(ClassGPA) FROM Classes WHERE ClassIdentifier LIKE ? AND Visible=TRUE GROUP BY TermID ORDER BY TermID"
	var response AvgGPAPerTermResponse
	response.Terms = make([]string, 0)
	response.GPA = make([]float64, 0)

	rows, err := db.Query(query, id+"%") //note the added % which allows us to ignore the numbers at end of each class id
	if err != nil {
		return AvgGPAPerTermResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var term string
		var GPA float64
		err := rows.Scan(&term, &GPA)
		if err != nil {
			return AvgGPAPerTermResponse{}, err
		}

		response.Terms = append(response.Terms, term)
		response.GPA = append(response.GPA, GPA)
	}
	return response, nil
}

type WithdrawalRatePerTermResponse struct {
	Terms          []string
	WithdrawalRate []float64
}

func GetSubjectWithdrawalRatePerTerm(db *sql.DB, id string) (WithdrawalRatePerTermResponse, error) {
	var query = "SELECT TermID, AVG(W / Students) AS WithdrawalRate FROM Classes WHERE ClassIdentifier LIKE ? AND Visible=TRUE GROUP BY TermID ORDER BY TermID"
	var response WithdrawalRatePerTermResponse
	response.Terms = make([]string, 0)
	response.WithdrawalRate = make([]float64, 0)

	rows, err := db.Query(query, id+"%") //note the added % which allows us to ignore the numbers at end of each class id
	if err != nil {
		return WithdrawalRatePerTermResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var term string
		var WithdrawalRate float64
		err := rows.Scan(&term, &WithdrawalRate)
		if err != nil {
			return WithdrawalRatePerTermResponse{}, err
		}

		response.Terms = append(response.Terms, term)
		response.WithdrawalRate = append(response.WithdrawalRate, WithdrawalRate)
	}
	return response, nil
}

func GetSubjectAvgGPA(db *sql.DB, id string, term string) (float64, error) {
	var query = "SELECT AVG(ClassGPA) FROM Classes WHERE ClassIdentifier LIKE ? AND TermID=? AND Visible=TRUE GROUP BY TermID ORDER BY TermID"

	var GPA float64

	rows, err := db.Query(query, id+"%", term) //note the added % which allows us to ignore the numbers at end of each class id
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&GPA)
		if err != nil {
			return -1, err
		}
	}

	return GPA, nil
}
