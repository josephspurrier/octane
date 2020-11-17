package app_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/josephspurrier/octane"
	"github.com/josephspurrier/octane/example/app"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUserID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	ctx := &app.Context{ResponseJSON: octane.ResponseJSON{Context: c}}

	ctx.SetUserID("foo")
	s, b := ctx.UserID()

	assert.Equal(t, "foo", s)
	assert.Equal(t, true, b)
}
