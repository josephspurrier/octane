package app

import (
	"context"

	"github.com/josephspurrier/octane"
	"github.com/josephspurrier/octane/example/app/lib/passhash"
	"github.com/josephspurrier/octane/example/app/lib/webtoken"
	"github.com/labstack/echo/v4"
)

var (
	// KeyUserID -
	KeyUserID = contextKey("user_id")
)

type contextKey string

// Context is a custom app context for use with handlers.
type Context struct {
	octane.ResponseJSON
	DB       IDatabase
	Passhash *passhash.Passhash
	Webtoken *webtoken.Configuration
}

// HandlerFunc allows using handlers with app.Context instead of echo Context.
func (ctx *Context) HandlerFunc(next func(*Context) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := &Context{
			ResponseJSON: octane.ResponseJSON{Context: c},
			DB:           ctx.DB,
			Passhash:     ctx.Passhash,
			Webtoken:     ctx.Webtoken,
		}

		return next(cc)
	}
}

// SetUserID will set the user ID in the context.
func (ctx *Context) SetUserID(val string) {
	r := ctx.Request()
	*r = *r.WithContext(context.WithValue(r.Context(), KeyUserID, val))
}

// UserID gets the user ID from the context.
func (ctx *Context) UserID() (string, bool) {
	val, ok := ctx.Request().Context().Value(KeyUserID).(string)
	return val, ok
}
