package endpoint

import (
	"github.com/josephspurrier/octane"
	"github.com/labstack/echo/v4"
)

// Hello -
// swagger:route GET / hello HelloGET
//
// Return a hello message.
//
// Security:
//   token:
//
// Responses:
//   200: OKResponse
func Hello(c echo.Context) error {
	cc := c.(*octane.Context)
	return cc.OKResponse("Hello World!")
}

// func Hello2(c echo.Context) error {
// 	//return errors.New("Ok dude")
// 	return echo.NewHTTPError(http.StatusNotFound, "Please provide valid credentials")
// 	//cc := c.(*Context)
// 	//return cc.OKResponse("Hello World!")
// }
