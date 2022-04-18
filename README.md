# Octane

This project contains a library and a template for a Go REST API using [Echo](https://echo.labstack.com/) and [Swagger](https://github.com/go-swagger/go-swagger) so you can:
- generate a Swagger spec from annotations in your code - [go-swagger](https://github.com/go-swagger/go-swagger)
- unmarshal HTML forms and JSON into structs - [go-playground/form](https://github.com/go-playground/form)
- validate struct fields using annotations - [go-playground/validator](https://github.com/go-playground/validator)

This project designed to be more of an example of how to create a custom binder for Echo.

## Usage

You can add Octane to your Echo application like this:

```go
package main

import (
	"github.com/josephspurrier/octane"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Use Octane for binding and validation.
	e.Binder = octane.NewBinder()
}
```

You can then create endpoints with [Swagger annotations that will generate a Swagger spec](https://goswagger.io/generate/spec.html) as well as [validate the incoming data using annotations](https://pkg.go.dev/github.com/go-playground/validator/v10). 

```go
package endpoint

import (
	"net/http"

	"github.com/josephspurrier/octane"
	"github.com/josephspurrier/octane/example/app"
	"github.com/josephspurrier/octane/example/app/store"
)

// Login - the annotations below can be used to generate a Swagger spec.
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
			Email string `json:"email" validate:"required,email"` // Binding and validation annotations.
			// Password.
			// example: password
			// required: true
			Password string `json:"password" validate:"required"` // Binding and validation annotations.
		}
	}

	// Request binding and validation.
	req := new(Request)
	if err = c.Bind(req); err != nil {
		return c.BadRequestResponse(err.Error())
	}

    // ... Logic to check password.

	return c.DataResponse(http.StatusOK, data)
}
```