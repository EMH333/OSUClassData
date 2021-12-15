package database

import "database/sql"

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
		rows.Scan(&term, &GPA)
		response.Terms = append(response.Terms, term)
		response.GPA = append(response.GPA, GPA)
	}
	return response, nil
}

func GetSubjectWithdrawlRatePerTerm(db *sql.DB, id string) (WithdrawlRatePerTermResponse, error) {
	var query = "SELECT TermID, AVG(W / Students) AS WithdrawlRate FROM Classes WHERE ClassIdentifier LIKE ? AND Visible=TRUE GROUP BY TermID ORDER BY TermID"
	var response WithdrawlRatePerTermResponse
	response.Terms = make([]string, 0)
	response.WithdrawlRate = make([]float64, 0)

	rows, err := db.Query(query, id+"%") //note the added % which allows us to ignore the numbers at end of each class id
	if err != nil {
		return WithdrawlRatePerTermResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var term string
		var WithdrawlRate float64
		rows.Scan(&term, &WithdrawlRate)
		response.Terms = append(response.Terms, term)
		response.WithdrawlRate = append(response.WithdrawlRate, WithdrawlRate)
	}
	return response, nil
}
