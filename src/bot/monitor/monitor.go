package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Launch() {

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Define the HTTP routes
	e.File("/", "/usr/local/cmx_bot/current/public/index.html")
	e.File("/style.css", "/usr/local/cmx_bot/current/public/style.css")
	e.File("/app.js", "/usr/local/cmx_bot/current/public/app.js")

	// Start server
	e.Logger.Fatal(e.Start(":9012"))

}
