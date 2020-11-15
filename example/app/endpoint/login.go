package endpoint

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Login .
// swagger:route POST /api/v1/login authentication UserLogin
//
// Authenticate a user.
//
// Responses:
//   200: LoginResponse
//   400: BadRequestResponse
//   500: InternalServerErrorResponse
func Login(c echo.Context) (err error) {
	cc := c.(*Context)

	// swagger:parameters UserLogin
	type Request struct {
		// in: body
		Body struct {
			// required: true
			Email string `json:"email" validate:"required,email"`
			// required: true
			Password string `json:"password" validate:"required"`
		}
	}

	// Request validation.
	req := new(Request)
	if err = c.Bind(req); err != nil {
		return cc.BadRequestResponse(err.Error())
	}

	// LoginResponse returns a token.
	// swagger:response LoginResponse
	type LoginResponse struct {
		// in: body
		Body struct {
			CommonStatusFields
			Data struct {
				// Token contains the API token for authentication
				// example: api-123456
				// required: true
				Token string `json:"token"`
			} `json:"data"`
		}
	}

	m := new(LoginResponse).Body.Data
	m.Token = "random"

	return cc.DataResponse(http.StatusOK, m)

	// // Determine if the user exists.
	// user := p.Store.User.New()
	// found, err := p.Store.User.FindOneByField(&user, "email", req.Body.Email)
	// if err != nil {
	// 	return http.StatusInternalServerError, err
	// } else if !found {
	// 	return http.StatusBadRequest, errors.New("login information does not match")
	// }

	// // Ensure the user's password matches. Use the same error message to prevent
	// // brute-force from finding usernames.
	// if !p.Password.Match(user.Password, req.Body.Password) {
	// 	return http.StatusBadRequest, errors.New("login information does not match")
	// }

	// // Create the response.
	// m := new(model.LoginResponse).Body
	// m.Status = http.StatusText(http.StatusOK)

	// // Generate the access token.
	// m.Token, err = p.Token.Generate(user.ID)
	// if err != nil {
	// 	return http.StatusInternalServerError, err
	// }

	// return p.Response.JSON(w, m)
}
