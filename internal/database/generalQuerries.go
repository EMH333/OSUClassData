package database

import (
	"database/sql"
	"strconv"
)

func GetLatestTerm(db *sql.DB) (int, error) {
	var query = "SELECT TermID FROM Classes WHERE Visible=TRUE ORDER BY TermID DESC LIMIT 1"

	rows, err := db.Query(query)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	var term string
	for rows.Next() {
		err := rows.Scan(&term)
		if err != nil {
			return -1, err
		}
	}

	termID, err := strconv.Atoi(term)
	if err != nil {
		return -1, err
	}

	return termID, nil
}
