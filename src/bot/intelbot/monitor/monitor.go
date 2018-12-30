package monitor

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	pusher "github.com/pusher/pusher-http-go"
)

var Client = pusher.Client{
	AppId:   "681531",
	Key:     "ba844c624003f02c6c0f",
	Secret:  "78d4a04ab77b874f4116",
	Cluster: "ap1",
	Secure:  true,
}

// visitsData is a struct
type VisitsData struct {
	Count int
	Time  string
}

func Launch() {

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Define the HTTP routes
	e.File("/", "public/index.html")
	e.File("/style.css", "public/style.css")
	e.File("/app.js", "public/app.js")

	// Start server
	e.Logger.Fatal(e.Start(":9012"))

}
