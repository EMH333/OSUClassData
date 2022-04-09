package main

import (
	"ethohampton.com/OSUClassData/internal/database"
	"ethohampton.com/OSUClassData/internal/util"
	"github.com/gofiber/fiber/v2"
)

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

	subjectGPA, err := database.GetSubjectAvgGPA(db, resp.ClassIdentifier, resp.LastTerm)
	if err != nil {
		return err
	}

	_ = subjectGPA - resp.AverageGPALastTerm //TODO
	return c.JSON(nil)
}
