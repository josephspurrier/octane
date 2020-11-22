package config

import (
	"time"

	"github.com/josephspurrier/octane"
	"github.com/josephspurrier/octane/example/app"
	"github.com/josephspurrier/octane/example/app/endpoint"
	"github.com/josephspurrier/octane/example/app/lib/passhash"
	"github.com/josephspurrier/octane/example/app/lib/webtoken"
	"github.com/josephspurrier/octane/example/app/middleware/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Config .
func Config() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Logger.SetPrefix("api")

	// Load the environment variables.
	settings := LoadEnv(e.Logger, "")

	// Middleware.
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Use Go Playground Validator.
	e.Binder = octane.NewBinder()

	// Connect the services.
	// Any changes here need to be also be made in the app/context.go file.
	ac := new(app.Context)
	ac.DB = Database(e.Logger)
	ac.Passhash = passhash.New()
	ac.Webtoken = webtoken.New([]byte(settings.Secret),
		time.Duration(settings.SessionTimeout)*time.Minute)

	// Set up the webtoken.
	token := jwt.New([]string{ // JWT whitelist.
		"GET /favicon.ico",
		"GET /api/v1/healthcheck",
		"POST /api/v1/login",
		"POST /api/v1/register",
		"GET /swagger/*",
	}, ac.Webtoken, *ac)
	e.Use(token.Handler())

	// Set the default error handler so all errors are handled outside of this function.
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		c.Logger().Error(err)
	}

	// Endpoints.
	e.GET("/api/v1/healthcheck", ac.HandlerFunc(endpoint.Healthcheck))
	e.POST("/api/v1/login", ac.HandlerFunc(endpoint.Login))
	e.POST("/api/v1/register", ac.HandlerFunc(endpoint.Register))
	e.POST("/api/v1/note", ac.HandlerFunc(endpoint.NoteCreate))
	e.GET("/api/v1/note", ac.HandlerFunc(endpoint.NoteIndex))
	e.GET("/api/v1/note/:note_id", ac.HandlerFunc(endpoint.NoteShow))
	e.PUT("/api/v1/note/:note_id", ac.HandlerFunc(endpoint.NoteUpdate))
	e.DELETE("/api/v1/note/:note_id", ac.HandlerFunc(endpoint.NoteDestroy))

	// Static routes.
	e.Static("/swagger/*", "swaggerui")

	return e
}
