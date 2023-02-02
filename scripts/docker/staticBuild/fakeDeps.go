package main

import (
	"container/heap"

	"github.com/enriquebris/goconcurrentqueue"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/template/html"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

/*
Create a file that will be used to generate the dependencies for the static build.
*/
func main() {
	_ = mysql.Config{}
	_ = adaptor.FiberApp(fiber.New())
	_ = fiber.New()
	_ = goconcurrentqueue.NewFixedFIFO(20)
	_ = html.New("", "")
	_ = etag.New()
	_ = cache.New()
	_ = heap.Interface(nil)
	_ = cases.Title(language.AmericanEnglish)
}
