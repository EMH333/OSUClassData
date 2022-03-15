package util

import (
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// This includes the data specific to the API response as well as the data generalized over the whole college (CS, PH, etc.)
type ClassGraphResponse struct {
	Dataset      string
	Terms        []string
	SpecificData []float64
	OverallData  []float64
}

// Write data to the response, handles all errors
func WriteJSON(w http.ResponseWriter, v interface{}) {
	jsonResponse, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "Error marshaling JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		http.Error(w, "Error writing JSON response", http.StatusInternalServerError)
	}
}

func SendError(c *fiber.Ctx, status int, message string) error {
	err := c.SendString(message)
	if err != nil {
		return err
	}
	return c.SendStatus(status)
}
