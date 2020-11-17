package endpoint

import (
	"net/http"

	"github.com/josephspurrier/octane"
	"github.com/josephspurrier/octane/example/app"
	"github.com/josephspurrier/octane/example/app/store"
)

// Register -
// swagger:route POST /api/v1/register authentication UserRegister
//
// Register a user.
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
