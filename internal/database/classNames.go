package database

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"ethohampton.com/OSUClassData/internal/util"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const classNameWaitTime = 30

var classNameTaskQueue *util.TaskQueue

func UpdateClassName(db *sql.DB, class string) {
	if classNameTaskQueue == nil {
		classNameTaskQueue = util.NewTaskQueue(db, classNameTask, classNameWaitTime*time.Second, 20)
	}

	err := classNameTaskQueue.Enqueue(class)
	// if there is an error then the queue is probably full, so we don't add another item
	// no harm in returning here, since a new request will queue it again
	if err != nil {
		return
	}
}

func GetClassNameQueueLength() int {
	return classNameTaskQueue.GetStats().CurrentQueue
}

func classNameTask(db *sql.DB, element interface{}, _ *util.TaskQueue) util.TaskQueueReturn {
	// it's possible that the same item ends up in the queue multiple times, so we need to check if it's still needed
	// just to keep our requests to the API to a minimum
	retrieveClassName, normalizeClassName := whatPartOfNameToUpdate(db, element.(string))
	if retrieveClassName {
		// we don't actually care about the API error because it will return an empty string which is fine for our purposes
		name, _ := getClassName(element.(string))
		name = normalizeName(name)
		_ = updateClassNameDatabase(db, element.(string), name)
	} else if normalizeClassName {
		// this path doesn't touch the API at all, so we don't need to wait
		name := getNameFromDB(db, element.(string))
		name = normalizeName(name)
		_ = updateClassNameDatabase(db, element.(string), name)

		// no need, since doesn't touch api
		return util.TaskQueueReturn{
			NoWait: true,
		}
	} else if !retrieveClassName && !normalizeClassName {
		//don't need to do anything so don't wait
		return util.TaskQueueReturn{
			NoWait: true,
		}
	}

	//wait by default
	return util.TaskQueueReturn{
		NoWait: false,
	}
}

// bools are should retrieve, should normalize
func whatPartOfNameToUpdate(db *sql.DB, ID string) (bool, bool) {
	query := `SELECT RetrievedClassName, NormalizedClassName FROM ClassInfo WHERE ClassIdentifier = ?`
	row := db.QueryRow(query, ID)
	var retrievedClassName, normalizedClassName bool
	err := row.Scan(&retrievedClassName, &normalizedClassName)
	if err != nil {
		return true, true
	}
	return !retrievedClassName, !normalizedClassName
}

func getClassName(class string) (string, error) {
	requestBody := `{"other":{"srcdb":"999999"}, "criteria":[{"field":"alias","value":"` + html.EscapeString(class) + `"}]}`

	resp, err := http.Post("https://classes.oregonstate.edu/api/?page=fose&route=search", "application/json", bytes.NewBufferString(requestBody))
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New("API returned non-200 status code")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	type Response struct {
		Srcdb   string                   `json:"srcdb"`
		Count   int                      `json:"count"`
		Results []map[string]interface{} `json:"results"`
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	if response.Count < 1 {
		return "", errors.New("API returned no results")
	}

	//if the API did return results, then we need to make sure that the class name is the same for at least some of them
	//if not, then we need to return an error
	var name = response.Results[0]["title"].(string)
	for i, result := range response.Results {
		if name != result["title"].(string) {
			name = ""
			break
		}
		//five should do it for now
		if i > 5 {
			break
		}
	}

	return name, nil
}

func getNameFromDB(db *sql.DB, ID string) string {
	query := `SELECT ClassName FROM ClassInfo WHERE ClassIdentifier = ?`
	row := db.QueryRow(query, ID)
	var name string
	err := row.Scan(&name)
	if err != nil {
		return ""
	}
	return name
}

func normalizeName(name string) string {
	name = strings.TrimSpace(name)
	name = properTitle(name)
	name = capitalizeRoman(name)
	name = commonSubstitutions(name)
	return name
}

func updateClassNameDatabase(db *sql.DB, ID, name string) error {
	//TODO update database
	updateQuery := `UPDATE ClassInfo SET ClassName = ?, RetrievedClassName=TRUE, NormalizedClassName=TRUE WHERE ClassIdentifier = ?`

	_, err := db.Exec(updateQuery, name, ID)
	if err != nil {
		return err
	}

	return nil
}

var caser = cases.Title(language.AmericanEnglish)

func properTitle(input string) string {
	words := strings.Split(strings.ToLower(input), " ")
	smallwords := " a an on the to in of and or for nor but yet so at by from with as if  "

	for index, word := range words {
		if strings.Contains(smallwords, " "+word+" ") && word != words[0] {
			words[index] = word
		} else {
			words[index] = caser.String(word)
		}
	}
	return strings.Join(words, " ")
}

func capitalizeRoman(input string) string {
	words := strings.Split(input, " ")
	for index, initialWord := range words {
		word := strings.ToUpper(initialWord)
		if word == "I" || word == "II" || word == "III" || word == "IV" || word == "V" || word == "VI" || word == "VII" || word == "VIII" || word == "IX" || word == "X" {
			words[index] = word
		}
	}
	return strings.Join(words, " ")
}

func commonSubstitutions(input string) string {
	//note case sensitive
	substitutions := map[string]string{
		"Computer Science": "CS",
		"*":                "", //remove the asterisk from names (why are they there? :) )
		"^":                "", //remove the caret from names (why are they there? :) )
	}

	for key, value := range substitutions {
		input = strings.Replace(input, key, value, -1)
	}

	// per word substitutions
	words := strings.Split(input, " ")
	for index, word := range words {
		word = strings.ToLower(word)
		if word == "and" {
			words[index] = "&"
		}
		if word == "introduction" {
			words[index] = "Intro"
		}
		if word == "cs" {
			words[index] = "CS"
		}
	}
	return strings.Join(words, " ")
}
