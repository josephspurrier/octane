package jwt

import (
	"fmt"
	"strings"

	"github.com/josephspurrier/octane"
	"github.com/josephspurrier/octane/example/app"
	"github.com/labstack/echo/v4"
)

// Config contains the dependencies for the handler.
type Config struct {
	whitelist []string
	webtoken  app.IToken
	ctx       app.Context
}

// New returns a new loq request middleware.
func New(whitelist []string, webtoken app.IToken, ctx app.Context) *Config {
	return &Config{
		whitelist: whitelist,
		ctx:       ctx,
		webtoken:  webtoken,
	}
}

// Handler will require a JWT.
func (c *Config) Handler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			r := ctx.Request()
			c.ctx.ResponseJSON = octane.ResponseJSON{Context: ctx}

			// Determine if the page is in the JWT whitelist.
			if !IsWhitelisted(r.Method, r.URL.Path, c.whitelist) {
				// Require JWT on all routes.
				bearer := r.Header.Get("Authorization")

				// If the token is missing, show an error.
				if len(bearer) < 8 || !strings.HasPrefix(bearer, "Bearer ") {
					return c.ctx.UnauthorizedResponse("authorization token is missing")
				}

				userID, err := c.webtoken.Verify(bearer[7:])
				if err != nil {
					return c.ctx.UnauthorizedResponse("authorization token is invalid")
				}

				c.ctx.SetUserID(userID)
			}

			return next(ctx)
		}
	}
}

// IsWhitelisted returns true if the request is in the whitelist. If only an
// asterisk is found in the whitelist, allow all routes. If an asterisk is
// found in the page string, then whitelist only the matching paths.
func IsWhitelisted(method string, path string, arr []string) (found bool) {
	s := fmt.Sprintf("%v %v", method, path)
	for _, i := range arr {
		if i == "*" || s == i {
			return true
		} else if strings.Contains(i, "*") {
			if strings.HasPrefix(s, i[:strings.Index(i, "*")]) {
				return true
			}
		}
	}
	return
}
