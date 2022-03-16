package main

import (
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
)

/*
Create a file that will be used to generate the dependencies for the static build.
*/
func main() {
	_ = mysql.Config{}
	_ = adaptor.FiberApp(fiber.New())
	_ = fiber.New()
}
