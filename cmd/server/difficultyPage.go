package main

import (
	"ethohampton.com/OSUClassData/internal/database"
	"ethohampton.com/OSUClassData/internal/util"
	"github.com/gofiber/fiber/v2"
)

type DifficultyResponse struct {
	Difficulty int8
	ClassName  string
}

func getRelativeClassDifficulty(c *fiber.Ctx) error {
	//get GPA of class
	//get average GPA of subject
	//determine how many standard deviations away from the mean
	//normalize into difficulty
	class := c.Query("class")
	if class == "" {
		return util.SendError(c, fiber.StatusBadRequest, "No class specified")
	}

	resp, _, err := database.GetClassInfo(db, class)
	if err != nil {
		return err
	}

	subjectGPA, subjectSD, err := database.GetSubjectAvgStdGPA(db, util.ClassIDToSubject(resp.ClassIdentifier), resp.LastTerm)
	if err != nil {
		return err
	}

	standardDeviation := (subjectGPA - resp.AverageGPALastTerm) / subjectSD

	const same = 0.1
	const sightly = 0.5    //below this is considered slightly easier/harder
	const moderately = 1.5 //below this is considered easier/harder. Above this is significinally easier/harder

	var difficulty int8

	if standardDeviation > -same && standardDeviation < same {
		difficulty = 3 //about average difficult
	} else {
		if standardDeviation < 0 {
			standardDeviation = -standardDeviation //make positive
			if standardDeviation < sightly {
				difficulty = 2
			} else if standardDeviation < moderately {
				difficulty = 1
			} else {
				difficulty = 0
			}
		} else {
			if standardDeviation < sightly {
				difficulty = 4
			} else if standardDeviation < moderately {
				difficulty = 5
			} else {
				difficulty = 6
			}
		}
	}

	return c.JSON(&DifficultyResponse{
		Difficulty: difficulty,
		ClassName:  resp.ClassIdentifier,
	})
}
