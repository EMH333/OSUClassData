package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"ethohampton.com/OSUClassData/internal/database"
	"ethohampton.com/OSUClassData/internal/util"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/template/html"
)

var db *sql.DB

var classLeaderboard = &util.Leaderboard{
	NumberOfTop: 5,
}

func main() {
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "8080")
	}

	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "OSUClassData",
	}

	// Try to connect to the database
	const maxConnectionAttempts = 10
	for i := 0; i < maxConnectionAttempts; i++ {
		err := tryToConnectToDB(&cfg)
		//if connected then start server
		if err == nil {
			break
		}

		//if not connected then wait and try again, up to 10 times, then exit with error
		if err != nil && i == maxConnectionAttempts-1 {
			log.Fatal(err)
		}

		//wait 5 seconds before trying again
		fmt.Println("Could not connect to DB. Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
	}

	fmt.Println("Connected to Database!")

	stopLeaderboard := util.SetUpLeaderboard(classLeaderboard) //make sure everything is configured

	engine := html.New("./frontend/templates", ".html")

	app := fiber.New(fiber.Config{
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"127.0.0.1"}, // localhost
		ProxyHeader:             fiber.HeaderXForwardedFor,

		Views: engine,
	})

	app.Use(etag.New())

	app.Static("/", "./frontend/dist")

	app.Get("/leaderboards", getLeaderboards)

	api := app.Group("/api/v0")
	api.Get("/status", getStatus)
	api.Get("/classes", getClasses)
	api.Get("/class", adaptor.HTTPHandlerFunc(getClass))
	api.Get("/classInfo", adaptor.HTTPHandlerFunc(getClassInfo))
	api.Get("/chart/studentsPerTerm", adaptor.HTTPHandlerFunc(getStudentsPerTerm))
	api.Get("/chart/avgGPAPerTerm", adaptor.HTTPHandlerFunc(getAvgGPAPerTerm))
	api.Get("/chart/withdrawalRatePerTerm", adaptor.HTTPHandlerFunc(getWithdrawalRatePerTerm))
	api.Get("/chart/lastTermGradeDistribution", adaptor.HTTPHandlerFunc(getLastTermGradeDistribution))

	api.Get("/subjects", adaptor.HTTPHandlerFunc(getSubjects))
	api.Get("/subject/chart/avgGPAPerTerm", adaptor.HTTPHandlerFunc(getSubjectAvgGPAPerTerm))
	api.Get("/subject/chart/withdrawalRatePerTerm", adaptor.HTTPHandlerFunc(getSubjectWithdrawalRatePerTerm))

	api.Get("/trendingClasses", getTrendingClasses)

	err := app.Listen(":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}

	close(stopLeaderboard) // stop leaderboard
}

func tryToConnectToDB(config *mysql.Config) error {
	var err error
	// Connect to the database
	db, err = sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return err
	}

	// Test the connection to the database with a ping
	err = db.Ping()
	if err != nil {
		return err
	}

	// run query to test connection
	var test int
	err = db.QueryRow("SELECT 1").Scan(&test)
	if err != nil {
		return err
	}

	return nil
}

func getClasses(c *fiber.Ctx) error {
	rows, err := db.Query("SELECT DISTINCT ClassIdentifier FROM Classes WHERE Visible=TRUE")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var classList []string

	for rows.Next() {
		var class string
		err := rows.Scan(&class)
		if err != nil {
			return util.SendError(c, http.StatusInternalServerError, "Error reading classes")
		}
		classList = append(classList, class)
	}

	return c.JSON(classList)
}

func getClass(w http.ResponseWriter, r *http.Request) {
	//TODO switch to path routing ie /api/v1/class/:classIdentifier/:term
	class := r.URL.Query().Get("class")
	if class == "" {
		http.Error(w, "Missing class parameter", http.StatusBadRequest)
		return
	}

	term := r.URL.Query().Get("term")
	if term == "" {
		http.Error(w, "Missing term parameter", http.StatusBadRequest)
		return
	}

	classData, err := database.GetTermClass(db, class, term)
	if err != nil {
		http.Error(w, "Class not found", http.StatusNotFound)
		return
	}

	util.WriteJSON(w, classData)
}

func getClassInfo(w http.ResponseWriter, r *http.Request) {
	class := r.URL.Query().Get("class")
	if class == "" {
		http.Error(w, "Missing class parameter", http.StatusBadRequest)
		return
	}

	classData, nameUpdateNeeded, err := database.GetClassInfo(db, class)
	if err != nil {
		log.Println(err)
		http.Error(w, "Class not found", http.StatusNotFound)
		return
	}

	// if we need to update the name of the class then do it
	if nameUpdateNeeded {
		database.UpdateClassName(db, class)
	}

	util.AddToLeaderboard(classLeaderboard, classData.ClassIdentifier) //start tracking this but don't do anything with it for now

	util.WriteJSON(w, classData)
}

func getStatus(c *fiber.Ctx) error {
	pingErr := db.Ping()
	if pingErr != nil {
		return util.SendError(c, http.StatusInternalServerError, "Can't ping DB")
	}
	c.Response().AppendBodyString("OK\n")

	//check queue length for getting class name
	c.Response().AppendBodyString(fmt.Sprintf("Queue length: %d\n", database.GetClassNameQueueLength()))

	return c.SendStatus(http.StatusOK)
}

// pretty simple method to get all the top trending classes
// TODO: allow trending classes per college
func getTrendingClasses(c *fiber.Ctx) error {
	return c.JSON(classLeaderboard.Top)
}

func getStudentsPerTerm(w http.ResponseWriter, r *http.Request) {
	class := r.URL.Query().Get("class")
	if class == "" {
		http.Error(w, "Missing class parameter", http.StatusBadRequest)
		return
	}

	studentsPerTerm, err := database.GetStudentsPerTerm(db, class)
	if err != nil {
		http.Error(w, "Class not found", http.StatusNotFound)
		return
	}

	var response util.ClassGraphResponse

	response.Dataset = "SpT"
	response.Terms = studentsPerTerm.Terms
	response.SpecificData = studentsPerTerm.Students
	//TODO add general subject data

	util.WriteJSON(w, response)
}

func getAvgGPAPerTerm(w http.ResponseWriter, r *http.Request) {
	class := r.URL.Query().Get("class")
	if class == "" {
		http.Error(w, "Missing class parameter", http.StatusBadRequest)
		return
	}

	GPAPerTerm, err := database.GetAvgGPAPerTerm(db, class)
	if err != nil {
		http.Error(w, "Class not found", http.StatusNotFound)
		return
	}

	var response util.ClassGraphResponse

	response.Dataset = "GpT"
	response.Terms = GPAPerTerm.Terms
	response.SpecificData = GPAPerTerm.GPA
	//TODO add general subject data

	util.WriteJSON(w, response)
}

func getWithdrawalRatePerTerm(w http.ResponseWriter, r *http.Request) {
	class := r.URL.Query().Get("class")
	if class == "" {
		http.Error(w, "Missing class parameter", http.StatusBadRequest)
		return
	}

	WithdrawalPerTerm, err := database.GetWithdrawalRatePerTerm(db, class)
	if err != nil {
		http.Error(w, "Class not found", http.StatusNotFound)
		return
	}

	var response util.ClassGraphResponse

	response.Dataset = "WpT"
	response.Terms = WithdrawalPerTerm.Terms
	response.SpecificData = WithdrawalPerTerm.WithdrawalRate
	//TODO add general subject data

	util.WriteJSON(w, response)
}

func getLastTermGradeDistribution(w http.ResponseWriter, r *http.Request) {
	class := r.URL.Query().Get("class")
	if class == "" {
		http.Error(w, "Missing class parameter", http.StatusBadRequest)
		return
	}

	latestClass, err := database.GetLastTermClass(db, class)
	if err != nil {
		http.Error(w, "Class not found", http.StatusNotFound)
		return
	}

	util.WriteJSON(w, latestClass)
}
