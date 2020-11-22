package octane

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// OKStatusFields should be included in all requests.
type OKStatusFields struct {
	// Code contains the HTTP status code.
	// example: 200
	// required: true
	StatusCode int `json:"status_code"`
	// Status contains the string of the HTTP status.
	// example: OK
	// required: true
	StatusMessage string `json:"status_message"`
}

// OKResponse is a success.
// swagger:response OKResponse
type OKResponse struct {
	// in: body
	Body struct {
		// Message contains a user friendly message.
		// example: The operation was successful.
		// required: true
		Message string `json:"message"`
		OKStatusFields
	}
}

// CreatedStatusFields should be included in all requests.
type CreatedStatusFields struct {
	// Code contains the HTTP status code.
	// example: 201
	// required: true
	StatusCode int `json:"status_code"`
	// Status contains the string of the HTTP status.
	// example: Created
	// required: true
	StatusMessage string `json:"status_message"`
}

// BadRequestResponse is a failure.
// swagger:response BadRequestResponse
type BadRequestResponse struct {
	// in: body
	Body struct {
		// Message contains a user friendly message.
		// example: The data submitted was invalid.
		// required: true
		Message string `json:"message"`
		// Code contains the HTTP status code.
		// example: 400
		// required: true
		StatusCode int `json:"status_code"`
		// Status contains the string of the HTTP status.
		// example: Bad Request
		// required: true
		StatusMessage string `json:"status_message"`
	}
}

// UnauthorizedResponse is a failure.
// swagger:response UnauthorizedResponse
type UnauthorizedResponse struct {
	// in: body
	Body struct {
		// Message contains a user friendly message.
		// example: You are not authorized to view this page.
		// required: true
		Message string `json:"message"`
		// Code contains the HTTP status code.
		// example: 401
		// required: true
		StatusCode int `json:"status_code"`
		// Status contains the string of the HTTP status.
		// example: Unauthorized
		// required: true
		StatusMessage string `json:"status_message"`
	}
}

// NotFoundResponse is a failure.
// swagger:response NotFoundResponse
type NotFoundResponse struct {
	// in: body
	Body struct {
		// Message contains a user friendly message.
		// example: The page was not found.
		// required: true
		Message string `json:"message"`
		// Code contains the HTTP status code.
		// example: 404
		// required: true
		StatusCode int `json:"status_code"`
		// Status contains the string of the HTTP status.
		// example: Not Found
		// required: true
		StatusMessage string `json:"status_message"`
	}
}

// InternalServerErrorResponse is a failure.
// swagger:response InternalServerErrorResponse
type InternalServerErrorResponse struct {
	// in: body
	Body struct {
		// Message contains a user friendly message.
		// example: An unexpected error occurred in the application.
		// required: true
		Message string `json:"message"`
		// Code contains the HTTP status code.
		// example: 500
		// required: true
		StatusCode int `json:"status_code"`
		// Status contains the string of the HTTP status.
		// example: Internal Server Error
		// required: true
		StatusMessage string `json:"status_message"`
	}
}

// ResponseJSON -
type ResponseJSON struct {
	echo.Context
}

// MessageResponse sends a JSON message with a status code.
func (c *ResponseJSON) MessageResponse(message string, statusCode int) error {
	resp := new(OKResponse)
	resp.Body.Message = message
	resp.Body.StatusCode = statusCode
	resp.Body.StatusMessage = http.StatusText(statusCode)
	err := c.JSON(resp.Body.StatusCode, resp.Body)
	if statusCode >= 400 {
		err = errors.New(message)
	}
	return err
}

// OKResponse sends 200.
func (c *ResponseJSON) OKResponse(message string) error {
	return c.MessageResponse(message, http.StatusOK)
}

// BadRequestResponse sends 400.
func (c *ResponseJSON) BadRequestResponse(message string) error {
	return c.MessageResponse(message, http.StatusBadRequest)
}

// UnauthorizedResponse sends 401.
func (c *ResponseJSON) UnauthorizedResponse(message string) error {
	return c.MessageResponse(message, http.StatusUnauthorized)
}

// NotFoundResponse sends 404.
func (c *ResponseJSON) NotFoundResponse(message string) error {
	return c.MessageResponse(message, http.StatusNotFound)
}

// InternalServerErrorResponse sends 500.
func (c *ResponseJSON) InternalServerErrorResponse(message string) error {
	return c.MessageResponse(message, http.StatusInternalServerError)
}

// DataResponse sends content with a status_code and a status_message to the response writer.
func (c *ResponseJSON) DataResponse(code int, i interface{}) error {
	c.Response().WriteHeader(code)
	c.Response().Header().Set(echo.HeaderContentType, "application/json")

	f := map[string]interface{}{
		"data":           i,
		"status_code":    code,
		"status_message": http.StatusText(code),
	}

	b, err := json.Marshal(f)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	fmt.Fprint(c.Response().Writer, string(b))

	return nil
}
