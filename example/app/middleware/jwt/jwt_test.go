package jwt_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/josephspurrier/octane"
	"github.com/josephspurrier/octane/example/app"
	"github.com/josephspurrier/octane/example/app/lib/webtoken"
	"github.com/josephspurrier/octane/example/app/middleware/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestWhitelistAllowed(t *testing.T) {
	for _, v := range []string{
		"GET /v1",
		"GET /v1/auth",
	} {
		arr := strings.Split(v, " ")

		whitelist := []string{
			"GET /v1",
			"GET /v1/auth",
		}

		e := echo.New()
		r := httptest.NewRequest(arr[0], arr[1], nil)
		r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		w := httptest.NewRecorder()
		c := e.NewContext(r, w)
		ctx := &app.Context{ResponseJSON: octane.ResponseJSON{Context: c}}

		webtoken := webtoken.New([]byte("secret"), 1*time.Minute)
		token := jwt.New(whitelist, webtoken, *ctx)

		e.Use(token.Handler())
		e.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	}
}
func TestWhitelistNotAllowed(t *testing.T) {
	for _, v := range []string{
		"POST /v1",
		"POST /v1/auth",
		"POST /v1/user",
		"DELETE /v1/user/1",
	} {
		arr := strings.Split(v, " ")

		whitelist := []string{
			"GET /v1",
			"GET /v1/auth",
		}

		e := echo.New()
		r := httptest.NewRequest(arr[0], arr[1], nil)
		r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		w := httptest.NewRecorder()
		c := e.NewContext(r, w)
		ctx := &app.Context{ResponseJSON: octane.ResponseJSON{Context: c}}

		webtoken := webtoken.New([]byte("secret"), 1*time.Minute)
		token := jwt.New(whitelist, webtoken, *ctx)

		e.Use(token.Handler())
		e.ServeHTTP(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), `authorization token is missing`)
	}
}

func TestWhitelistBadBearer(t *testing.T) {
	whitelist := []string{
		"GET /v1",
		"GET /v1/auth",
	}

	e := echo.New()
	r := httptest.NewRequest("POST", "/v1/user", nil)
	r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	w := httptest.NewRecorder()
	c := e.NewContext(r, w)
	ctx := &app.Context{ResponseJSON: octane.ResponseJSON{Context: c}}

	webtoken := webtoken.New([]byte("secret"), 1*time.Minute)
	token := jwt.New(whitelist, webtoken, *ctx)

	r.Header.Set("Authorization", "Bearer bad")
	e.Use(token.Handler())
	e.ServeHTTP(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), `authorization token is invalid`)
}

func TestIsWhitelisted(t *testing.T) {
	assert.Equal(t, true, jwt.IsWhitelisted("GET", "/v1", []string{
		"GET /v1",
	}))

	assert.Equal(t, true, jwt.IsWhitelisted("GET", "/v1", []string{
		"*",
	}))

	assert.Equal(t, true, jwt.IsWhitelisted("GET", "/v1", []string{
		"POST /v1",
		"*",
	}))

	// Allow weird spacing.
	assert.Equal(t, true, jwt.IsWhitelisted("GET", "/v1", []string{
		"POST /v1",
		"* ",
	}))

	// Not in the list.
	assert.Equal(t, false, jwt.IsWhitelisted("GET", "/v2", []string{
		"POST /v1",
		"GET /v1",
	}))
}
