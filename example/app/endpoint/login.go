package endpoint

import (
	"net/http"

	"github.com/josephspurrier/octane"
	"github.com/josephspurrier/octane/example/app"
	"github.com/josephspurrier/octane/example/app/store"
)

// Login -
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
			// Email address.
			// example: jsmith@example.com
			// required: true
			Email string `json:"email" validate:"required,email"`
			// Password.
			// example: password
			// required: true
			Password string `json:"password" validate:"required"`
		}
	}

	// Request validation.
	req := new(Request)
	if err = c.Bind(req); err != nil {
		return c.BadRequestResponse(err.Error())
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

	// LoginResponse returns a token.
	// swagger:response LoginResponse
	type LoginResponse struct {
		// in: body
		Body struct {
			octane.OKStatusFields
			Data struct {
				// Token contains the API token for authentication
				// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIwNWE3ZjlmYS1mN2ViLTIzNmItYjJiYi1iYTE0NWUwYTRhMmQiLCJleHAiOjE2MDU2MTQ1NzEsImp0aSI6IjA0MjQ0Yzc4LTU5MzItYTBjZS1lMjAzLTc3MmNiMDVhYmFhZiIsImlhdCI6MTYwNTU4NTc3MSwibmJmIjoxNjA1NTg1NzcxfQ.kAeCynxCh35moPf5OEsn7LW0oHNEBVWxVOiZ6RdyUwk
				// required: true
				Token string `json:"token"`
			} `json:"data"`
		}
	}

	data := new(LoginResponse).Body.Data

	// Generate a token for the user.
	data.Token, err = c.Webtoken.Generate(user.ID)
	if err != nil {
		return c.InternalServerErrorResponse(err.Error())
	}

	return c.DataResponse(http.StatusOK, data)
}
