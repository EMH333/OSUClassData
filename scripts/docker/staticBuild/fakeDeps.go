package main

import (
	"github.com/enriquebris/goconcurrentqueue"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/template/html"
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
}
