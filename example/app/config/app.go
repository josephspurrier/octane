package config

import (
	"fmt"
	"net/http"

	"github.com/josephspurrier/octane"
	"github.com/josephspurrier/octane/example/app"
	"github.com/josephspurrier/octane/example/app/endpoint"
	"github.com/josephspurrier/octane/example/app/lib/passhash"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Config .
func Config() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	// Middleware.
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Use Go Playground Validator.
	e.Binder = octane.NewBinder()

	// Connect the services.
	ac := new(app.Context)
	ac.DB = Database(e.Logger)
	ac.Passhash = passhash.New()

	// Use app context.
	// e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
	// 	return func(c echo.Context) error {
	// 		cc := &app.Context{
	// 			ResponseJSON: octane.ResponseJSON{Context: c},
	// 			DB:           db,
	// 		}
	// 		return next(cc)
	// 	}
	// })

	// Set the default error handler so all errors use the standard format.
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		cc := &app.Context{
			ResponseJSON: octane.ResponseJSON{Context: c},
			DB:           ac.DB,
		}

		code := http.StatusInternalServerError
		message := ""
		if he, ok := err.(*echo.HTTPError); ok {
			// Send a response when we know the message.
			code = he.Code
			message = fmt.Sprint(he.Message)
			c.Logger().Error(err)
			cc.MessageResponse(message, code)
			return
		}

		// If we don't know the message, send the whole error as the message
		// and use the Internal Server Error.
		c.Logger().Error(err)
		cc.MessageResponse(err.Error(), code)
	}

	// Endpoints.
	e.GET("/", ac.HandlerFunc(endpoint.Healthcheck))
	e.POST("/api/v1/login", ac.HandlerFunc(endpoint.Login))

	// Static routes.
	e.Static("/swagger/*", "swaggerui")

	return e
}
