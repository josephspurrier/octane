package config

import (
	"github.com/josephspurrier/octane/bind"
	"github.com/josephspurrier/octane/example/app/endpoint"
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

	// Use app context.
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &endpoint.Context{Context: c}
			return next(cc)
		}
	})

	// // Set the error page.
	// func customHTTPErrorHandler(err error, c echo.Context) {
	// 	code := http.StatusInternalServerError
	// 	if he, ok := err.(*echo.HTTPError); ok {
	// 		code = he.Code
	// 	}
	// 	errorPage := fmt.Sprintf("%d.html", code)
	// 	if err := c.File(errorPage); err != nil {
	// 		c.Logger().Error(err)
	// 	}
	// 	c.Logger().Error(err)
	// }

	// e.HTTPErrorHandler = customHTTPErrorHandler

	// Database.
	_ = Database(e.Logger)

	// Endpoints.
	e.GET("/", endpoint.Hello)
	//e.GET("/2", endpoint.Hello2)

	e.POST("/api/v1/login", endpoint.Login)

	// Static routes.
	e.Static("/swagger/*", "swaggerui")

	e.Binder = bind.New()

	return e
}
