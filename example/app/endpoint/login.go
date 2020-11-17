package endpoint

import (
	"net/http"

	"github.com/josephspurrier/octane"
	"github.com/josephspurrier/octane/example/app"
	"github.com/josephspurrier/octane/example/app/lib/securegen"
	"github.com/josephspurrier/octane/example/app/store"
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
func Login(c *app.Context) (err error) {
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
		return c.BadRequestResponse(err.Error())
	}

	// LoginResponse returns a token.
	// swagger:response LoginResponse
	type LoginResponse struct {
		// in: body
		Body struct {
			octane.CommonStatusFields
			Data struct {
				// Token contains the API token for authentication
				// example: api-123456
				// required: true
				Token string `json:"token"`
			} `json:"data"`
		}
	}

	// Check if user exists.
	user := new(store.User)
	found, err := store.FindOneByField(c.DB, user, "email", req.Body.Email)
	if err != nil {
		return c.InternalServerErrorResponse(err.Error())
	} else if !found {
		return c.BadRequestResponse("login information does not match")
	}

	// Check user password.
	if !c.Passhash.Match(user.Password, req.Body.Password) {
		return c.BadRequestResponse("login information does not match")
	}

	data := new(LoginResponse).Body.Data

	// Generate a token for the user.
	data.Token, err = securegen.UUID()
	if err != nil {
		return c.InternalServerErrorResponse(err.Error())
	}

	return c.DataResponse(http.StatusOK, data)
}
