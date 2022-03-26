package main

import (
	"database/sql"
	"fmt"
	"log"

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
	return c.Render("leaderboard", fiber.Map{
		"Leaderboards": []LeaderboardDisplay{
			{
				Name:    "Largest Classes Last Term",
				Entries: getLargestClassesLastTerm(db, 10),
			},
			{
				Name:    "Highest GPA Classes Last Term",
				Entries: getHighestGPAClassesLastTerm(db, 10),
			},
			{
				Name:    "Highest Withdrawal Classes Last Term",
				Entries: getHighestWithdrawalClassesLastTerm(db, 10),
			},
		},
	})
}

func getLargestClassesLastTerm(db *sql.DB, num int) []LeaderboardEntry {
	rows, err := db.Query("SELECT ClassIdentifier, Students FROM Classes WHERE TermID=(SELECT TermID FROM Classes WHERE Visible=TRUE ORDER BY TermID DESC LIMIT 1) AND Visible=TRUE ORDER BY Students DESC LIMIT ?", num)
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

func getHighestGPAClassesLastTerm(db *sql.DB, num int) []LeaderboardEntry {
	rows, err := db.Query("SELECT ClassIdentifier, ClassGPA FROM Classes WHERE TermID=(SELECT TermID FROM Classes WHERE Visible=TRUE ORDER BY TermID DESC LIMIT 1) AND Visible=TRUE ORDER BY ClassGPA DESC LIMIT ?", num)
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

func getHighestWithdrawalClassesLastTerm(db *sql.DB, num int) []LeaderboardEntry {
	rows, err := db.Query("SELECT ClassIdentifier, (W / Students) AS WithdrawalRate FROM Classes WHERE TermID=(SELECT TermID FROM Classes WHERE Visible=TRUE ORDER BY TermID DESC LIMIT 1) AND Visible=TRUE ORDER BY WithdrawalRate DESC LIMIT ?", num)
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
