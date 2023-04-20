package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"ethohampton.com/OSUClassData/internal/database"
	"ethohampton.com/OSUClassData/internal/util"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/template/html"
)

var db *sql.DB

var classLeaderboard = &util.Leaderboard{
	NumberOfTop: 5,
}

var templateEngine = html.New("./frontend/distTemplates", ".html")

var dev = os.Getenv("DEV") == "true"

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

	app := fiber.New(fiber.Config{
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"127.0.0.1"}, // localhost
		ProxyHeader:             fiber.HeaderXForwardedFor,

		Views: templateEngine,
	})

	app.Use(etag.New())

	app.Static("/", "./frontend/dist")

	// Can use a cache for all requests right now since we don't have any dynamic content per user
	// for 10 minutes
	app.Use(cache.New(cache.Config{
		Expiration:   30 * time.Minute,
		CacheControl: true,
		MaxBytes:     1024 * 1024 * 20, // 20 MB
		KeyGenerator: func(c *fiber.Ctx) string {
			return utils.CopyString(c.Path()) + utils.CopyString(c.Query("class")) + utils.CopyString(c.Query("term")) + utils.CopyString(c.Query("subject"))
		},
		Next: func(c *fiber.Ctx) bool {
			// Don't cache if getting trending classes
			shouldSkip := c.Path() == "/api/v0/trendingClasses" || c.Path() == "/api/v0/status"
			if dev {
				log.Println(c.Path(), " can be cached: ", !shouldSkip)
			}
			return shouldSkip
		},
	}))

	app.Get("/leaderboards", getLeaderboards)
	app.Get("/class/:class", serveClass)
	app.Get("/class.html", redirectClass) //redirect that can go away eventually
	app.Get("/sitemap.xml", getSitemap)

	api := app.Group("/api/v0")
	api.Get("/status", getStatus)
	api.Get("/classes", getClasses)
	api.Get("/class", adaptor.HTTPHandlerFunc(getClass))
	api.Get("/classInfo/:class", getClassInfo)
	api.Get("/classInfo", func(c *fiber.Ctx) error {
		// if the class query parameter exists, then redirect to the one using path routing
		if c.Query("class") != "" {
			return c.Redirect("/api/v0/classInfo/"+c.Query("class"), fiber.StatusMovedPermanently)
		}
		return c.SendStatus(http.StatusBadRequest)
	})
	api.Get("/chart/lastTermGradeDistribution", adaptor.HTTPHandlerFunc(getLastTermGradeDistribution))
	api.Get("/chart/combinedData/:class", getCombinedClassStats)

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

func serveClass(c *fiber.Ctx) error {
	// confirm class exists
	var class string
	err := db.QueryRow("SELECT ClassIdentifier FROM Classes WHERE ClassIdentifier=?", c.Params("class")).Scan(&class)
	if err != nil {
		return util.SendError(c, http.StatusNotFound, "Class not found")
	}

	return c.Render("class", fiber.Map{
		"Class": c.Params("class"),
	})
}

//TODO remove once SEO is done
func redirectClass(c *fiber.Ctx) error {
	// don't redirect if class query parameter is missing
	if c.Query("class") == "" {
		return util.SendError(c, http.StatusBadRequest, "Missing class query parameter")
	}

	// confirm class exists
	var class string
	err := db.QueryRow("SELECT ClassIdentifier FROM Classes WHERE ClassIdentifier=?", c.Query("class")).Scan(&class)
	if err != nil {
		return util.SendError(c, http.StatusNotFound, "Class not found")
	}

	return c.Redirect("/class/" + c.Query("class"))
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

func getClassInfo(c *fiber.Ctx) error {
	class := c.Params("class")
	if class == "" {
		return util.SendError(c, http.StatusBadRequest, "Missing class parameter")
	}

	classData, nameUpdateNeeded, err := database.GetClassInfo(db, class)
	if err != nil {
		log.Println(err)
		return util.SendError(c, http.StatusNotFound, "Class not found")
	}

	// if we need to update the name of the class then do it
	if nameUpdateNeeded {
		database.UpdateClassName(db, class)
	}

	util.AddToLeaderboard(classLeaderboard, classData.ClassIdentifier) //start tracking this but don't do anything with it for now

	return c.JSON(classData)
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

func getCombinedClassStats(c *fiber.Ctx) error {
	class := c.Params("class")
	if class == "" {
		return util.SendError(c, fiber.StatusBadRequest, "Missing class")

	}

	CombinedClassStats, err := database.GetCombinedClassStats(db, class)
	if err != nil {
		return util.SendError(c, fiber.StatusNotFound, "Class not found")
	}

	var response util.CombinedClassGraphResponse

	response.Terms = CombinedClassStats.Terms
	response.SpecificData = make(map[string][]float64, 3)
	response.SpecificData["WR"] = CombinedClassStats.WithdrawalRate
	response.SpecificData["GPA"] = CombinedClassStats.GPA
	response.SpecificData["S"] = CombinedClassStats.Students
	//TODO add general subject data

	return c.JSON(response)
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

func getSitemap(c *fiber.Ctx) error {
	rows, err := db.Query("SELECT DISTINCT ClassIdentifier FROM Classes WHERE Visible=TRUE ORDER BY ClassIdentifier ASC")
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

	// xml sitemap can't have more than 50,000 entries, else need to split into multiple
	if len(classList) > 50000 {
		return util.SendError(c, http.StatusInternalServerError, "Too many classes to generate sitemap")
	}

	err = c.Render("sitemap", fiber.Map{
		"Classes": classList,
		"Others": []string{
			"/",
			"/about.html",
		},
		"Header": template.HTML(xml.Header),
	})
	if err != nil {
		return util.SendError(c, http.StatusInternalServerError, "Error rendering sitemap")
	}

	return c.Type("xml").SendStatus(http.StatusOK)
}
