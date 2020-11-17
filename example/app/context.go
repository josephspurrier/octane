package app

import (
	"github.com/josephspurrier/octane"
	"github.com/josephspurrier/octane/example/app/lib/passhash"
	"github.com/josephspurrier/octane/example/app/lib/webtoken"
	"github.com/labstack/echo/v4"
)

// Context is a custom app context for use with handlers.
type Context struct {
	octane.ResponseJSON
	DB       IDatabase
	Passhash *passhash.Passhash
	Webtoken *webtoken.Configuration
}

// HandlerFunc allows using handlers with app.Context instead of echo Context.
func (ac *Context) HandlerFunc(next func(*Context) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := &Context{
			ResponseJSON: octane.ResponseJSON{Context: c},
			DB:           ac.DB,
			Passhash:     ac.Passhash,
			Webtoken:     ac.Webtoken,
		}

		return next(cc)
	}
}
