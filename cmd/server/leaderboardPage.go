package main

import (
	"database/sql"
	"fmt"
	"log"

	"ethohampton.com/OSUClassData/internal/database"
	"ethohampton.com/OSUClassData/internal/util"
	"github.com/gofiber/fiber/v2"
)

type LeaderboardDisplay struct {
	Name    string
	Entries []LeaderboardEntry
}

type LeaderboardEntry struct {
	Name  string
	Link  string
	Score string
}

func getLeaderboards(c *fiber.Ctx) error {
	termID, err := database.GetLatestTerm(db)
	if err != nil {
		return err
	}
	latestTerm := util.TermIDToName(termID)

	return c.Render("leaderboard", fiber.Map{
		"Leaderboards": []LeaderboardDisplay{
			{
				Name:    "Largest Classes " + latestTerm,
				Entries: getLargestClassesLastTerm(db, 10, termID),
			},
			{
				Name:    "Highest GPA Classes " + latestTerm,
				Entries: getHighestGPAClassesLastTerm(db, 10, termID),
			},
			{
				Name:    "Highest Withdrawal Rate Classes " + latestTerm,
				Entries: getHighestWithdrawalClassesLastTerm(db, 10, termID),
			},
			{
				Name:    "Highest Pass Rate Classes " + latestTerm,
				Entries: getHighestPassRateClassesLastTerm(db, 10, termID),
			},
			{
				Name:    "Highest GPA Subjects " + latestTerm,
				Entries: getHighestGPASubjectsLastTerm(db, 10, termID),
			},
			{
				Name:    "Highest Withdrawal Rate Subjects " + latestTerm,
				Entries: getHighestWithdrawalSubjectsLastTerm(db, 10, termID),
			},
		},
	})
}

func getLargestClassesLastTerm(db *sql.DB, num int, term int) []LeaderboardEntry {
	rows, err := db.Query("SELECT ClassIdentifier, Students FROM Classes WHERE TermID=? AND Visible=TRUE ORDER BY Students DESC LIMIT ?", term, num)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var list []LeaderboardEntry

	for rows.Next() {
		var name string
		var score int
		err = rows.Scan(&name, &score)
		if err != nil {
			log.Fatal(err)
		}

		list = append(list, LeaderboardEntry{
			Name:  name,
			Score: fmt.Sprintf("%d", score),
			Link:  util.GetClassLink(name),
		})
	}

	return list
}

func getHighestGPAClassesLastTerm(db *sql.DB, num int, term int) []LeaderboardEntry {
	rows, err := db.Query("SELECT ClassIdentifier, ClassGPA FROM Classes WHERE TermID=? AND Visible=TRUE ORDER BY ClassGPA DESC LIMIT ?", term, num)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var list []LeaderboardEntry

	for rows.Next() {
		var name string
		var score float32
		err = rows.Scan(&name, &score)
		if err != nil {
			log.Fatal(err)
		}

		list = append(list, LeaderboardEntry{
			Name:  name,
			Score: fmt.Sprintf("%0.2f", score),
			Link:  util.GetClassLink(name),
		})
	}

	return list
}

func getHighestWithdrawalClassesLastTerm(db *sql.DB, num int, term int) []LeaderboardEntry {
	rows, err := db.Query("SELECT ClassIdentifier, (W / Students) AS WithdrawalRate FROM Classes WHERE TermID=? AND Visible=TRUE ORDER BY WithdrawalRate DESC LIMIT ?", term, num)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var list []LeaderboardEntry

	for rows.Next() {
		var name string
		var score float32
		err = rows.Scan(&name, &score)
		if err != nil {
			log.Fatal(err)
		}

		list = append(list, LeaderboardEntry{
			Name:  name,
			Score: fmt.Sprintf("%0.2f%%", score*100),
			Link:  util.GetClassLink(name),
		})
	}

	return list
}

func getHighestPassRateClassesLastTerm(db *sql.DB, num int, term int) []LeaderboardEntry {
	rows, err := db.Query("SELECT ClassIdentifier, ((A+AMinus+B+BPlus+BMinus+C+CPlus) / Students) AS PassRate FROM Classes WHERE TermID=? AND Visible=TRUE ORDER BY PassRate DESC LIMIT ?", term, num)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var list []LeaderboardEntry

	for rows.Next() {
		var name string
		var score float32
		err = rows.Scan(&name, &score)
		if err != nil {
			log.Fatal(err)
		}

		list = append(list, LeaderboardEntry{
			Name:  name,
			Score: fmt.Sprintf("%0.2f%%", score*100),
			Link:  util.GetClassLink(name),
		})
	}

	return list
}

func getHighestGPASubjectsLastTerm(db *sql.DB, num int, term int) []LeaderboardEntry {
	rows, err := db.Query(`SELECT DISTINCT REPLACE
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
	'9', '') as Subject, AVG(ClassGPA) AS GPA FROM Classes WHERE TermID=? AND Visible=TRUE GROUP BY Subject ORDER BY GPA DESC LIMIT ?`, term, num)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var list []LeaderboardEntry

	for rows.Next() {
		var name string
		var score float32
		err = rows.Scan(&name, &score)
		if err != nil {
			log.Fatal(err)
		}

		list = append(list, LeaderboardEntry{
			Name:  name,
			Score: fmt.Sprintf("%0.2f", score),
			Link:  util.GetSubjectLink(name),
		})
	}

	return list
}

func getHighestWithdrawalSubjectsLastTerm(db *sql.DB, num int, term int) []LeaderboardEntry {
	rows, err := db.Query(`SELECT DISTINCT REPLACE
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
	'9', '') as Subject, AVG(W / Students) AS WithdrawalRate FROM Classes WHERE TermID=? AND Visible=TRUE GROUP BY Subject ORDER BY WithdrawalRate DESC LIMIT ?`, term, num)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var list []LeaderboardEntry

	for rows.Next() {
		var name string
		var score float32
		err = rows.Scan(&name, &score)
		if err != nil {
			log.Fatal(err)
		}

		list = append(list, LeaderboardEntry{
			Name:  name,
			Score: fmt.Sprintf("%0.2f%%", score*100),
			Link:  util.GetSubjectLink(name),
		})
	}

	return list
}
