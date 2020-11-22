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
// Return a token after verifying the login information.
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
				// Token contains the API token for authentication.
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

// Register -
// swagger:route POST /api/v1/register authentication UserRegister
//
// Create a user in the system.
//
// Responses:
//   201: RegisterResponse
//   400: BadRequestResponse
//   500: InternalServerErrorResponse
func Register(c *app.Context) (err error) {
	// swagger:parameters UserRegister
	type Request struct {
		// in: body
		Body struct {
			// First name.
			// example: John
			// required: true
			FirstName string `json:"first_name" validate:"required"`
			// Last name.
			// example: Smith
			// required: true
			LastName string `json:"last_name" validate:"required"`
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
	found, _, err := store.ExistsByField(c.DB, user, "email", req.Body.Email)
	if err != nil {
		return c.InternalServerErrorResponse(err.Error())
	} else if found {
		return c.BadRequestResponse("user already exists")
	}

	// Encrypt the password.
	password, err := c.Passhash.Hash(req.Body.Password)
	if err != nil {
		return c.InternalServerErrorResponse(err.Error())
	}

	// Create the user.
	ID, err := store.CreateUser(c.DB, req.Body.FirstName,
		req.Body.LastName, req.Body.Email, password)
	if err != nil {
		return c.InternalServerErrorResponse(err.Error())
	}

	// RegisterResponse returns a user ID.
	// swagger:response RegisterResponse
	type RegisterResponse struct {
		// in: body
		Body struct {
			octane.CreatedStatusFields
			Data struct {
				// RecordID contains the newly created user ID.
				// example: 314445cd-e9fb-4c58-58b6-777ee06465f5
				// required: true
				RecordID string `json:"record_id"`
			} `json:"data"`
		}
	}

	// Set the user ID.
	data := new(RegisterResponse).Body.Data
	data.RecordID = ID

	return c.DataResponse(http.StatusCreated, data)
}
