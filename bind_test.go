package octane_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/josephspurrier/octane"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// StringInt create a type alias for type int
type StringInt int

// UnmarshalJSON create a custom unmarshal for the StringInt
/// this helps us check the type of our value before unmarshalling it

func (st *StringInt) UnmarshalJSON(b []byte) error {
	//convert the bytes into an interface
	//this will help us check the type of our value
	//if it is a string that can be converted into an int we convert it
	///otherwise we return an error
	var item interface{}
	if err := json.Unmarshal(b, &item); err != nil {
		return err
	}
	switch v := item.(type) {
	case int:
		*st = StringInt(v)
	case float64:
		*st = StringInt(int(v))
	case string:
		///here convert the string into
		///an integer
		i, err := strconv.Atoi(v)
		if err != nil {
			///the string might not be of integer type
			///so return an error
			return err

		}
		*st = StringInt(i)

	}
	return nil
}

func TestFormSuccess(t *testing.T) {
	called := false

	e := echo.New()
	cb := octane.NewBinder()
	e.Binder = cb

	e.POST("/user/:user_id", func(c echo.Context) error {
		called = true

		// swagger:parameters UserCreate
		type request struct {
			// in: path
			UserID string `form:"user_id" validate:"required"`
			// in: formData
			// Required: true
			FirstName string `form:"first_name" validate:"required"`
			// in: formData
			// Required: true
			LastName string `form:"last_name"  validate:"required"`
			// in: formData
			Age uint8 `form:"age,omitempty" bson:"age,omitempty" validate:"required"`
			// in: formData
			//Count StringInt `form:"count" validate:"required" bson:"count,omitempty" json:"count,omitempty"`
		}

		req := new(request)
		assert.NoError(t, c.Bind(req))

		assert.Equal(t, "10", req.UserID)
		assert.Equal(t, "john", req.FirstName)
		assert.Equal(t, "smith", req.LastName)
		assert.Equal(t, uint8(3), req.Age)
		//assert.Equal(t, 3, req.Count)
		return nil
	})

	form := url.Values{}
	form.Add("first_name", "john")
	form.Add("last_name", "smith")
	form.Add("age", "3")
	//form.Add("count", "3")

	r := httptest.NewRequest("POST", "/user/10", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, r)

	assert.Equal(t, true, called)
}

func TestFormFailValidate(t *testing.T) {
	called := false

	e := echo.New()
	cb := octane.NewBinder()
	e.Binder = cb

	e.POST("/user/:user_id", func(c echo.Context) error {
		called = true

		// swagger:parameters UserCreate
		type request struct {
			// in: path
			UserID string `form:"user_id" validate:"required"`
			// in: formData
			// Required: true
			FirstName string `form:"first_name" validate:"required"`
			// in: formData
			// Required: true
			LastName string `form:"last_name" validate:"required"`
		}

		req := new(request)
		assert.NotNil(t, c.Bind(req))

		assert.Equal(t, "10", req.UserID)
		assert.Equal(t, "john", req.FirstName)
		assert.Equal(t, "", req.LastName)
		return nil
	})

	form := url.Values{}
	form.Add("first_name", "john")
	//form.Add("last_name", "smith")

	r := httptest.NewRequest("POST", "/user/10", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, r)

	assert.Equal(t, true, called)
}

func TestFormNil(t *testing.T) {
	called := false

	e := echo.New()
	cb := octane.NewBinder()
	e.Binder = cb

	e.POST("/user/:user_id", func(c echo.Context) error {
		called = true

		// swagger:parameters UserCreate
		type request struct {
			// in: path
			UserID string `form:"user_id" validate:"required"`
			// in: formData
			// Required: true
			FirstName string `form:"first_name" validate:"required"`
			// in: formData
			// Required: true
			LastName string `form:"last_name" validate:"required"`
		}

		req := new(request)
		assert.NotNil(t, c.Bind(req))

		assert.Equal(t, "10", req.UserID)
		assert.Equal(t, "", req.FirstName)
		assert.Equal(t, "", req.LastName)
		return nil
	})

	r := httptest.NewRequest("POST", "/user/10", nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, r)

	assert.Equal(t, true, called)
}

func TestFormMissingPointer(t *testing.T) {
	called := false

	e := echo.New()
	cb := octane.NewBinder()
	e.Binder = cb

	e.POST("/user/:user_id", func(c echo.Context) error {
		called = true

		// swagger:parameters UserCreate
		type request struct {
			// in: path
			UserID string `form:"user_id" validate:"required"`
			// in: formData
			// Required: true
			FirstName string `form:"first_name" validate:"required"`
			// in: formData
			// Required: true
			LastName string `form:"last_name" validate:"required"`
		}

		req := request{}

		assert.NotNil(t, c.Bind(req))

		assert.Equal(t, "", req.UserID)
		assert.Equal(t, "", req.FirstName)
		assert.Equal(t, "", req.LastName)
		return nil
	})

	form := url.Values{}
	form.Add("first_name", "john")
	form.Add("last_name", "smith")

	r := httptest.NewRequest("POST", "/user/10", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, r)

	assert.Equal(t, true, called)
}

func TestJSONSuccess(t *testing.T) {
	called := false

	e := echo.New()
	cb := octane.NewBinder()
	e.Binder = cb

	e.POST("/user/:user_id", func(c echo.Context) error {
		called = true

		// swagger:parameters UserCreate
		type request struct {
			// in: body
			Body struct {
				// in: path
				UserID string `json:"user_id" validate:"required"`
				// Required: true
				FirstName string `json:"first_name" validate:"required"`
				// Required: true
				LastName string `json:"last_name" validate:"required"`
				// Required: true
				Age uint8 `json:"age,omitempty" bson:"age,omitempty" validate:"required"`
			}
		}

		req := new(request).Body
		assert.Nil(t, c.Bind(&req))

		assert.Equal(t, "10", req.UserID)
		assert.Equal(t, "john", req.FirstName)
		assert.Equal(t, "smith", req.LastName)
		assert.Equal(t, uint8(3), req.Age)
		return nil
	})

	form := make(map[string]interface{})
	form["first_name"] = "john"
	form["last_name"] = "smith"
	form["age"] = uint8(3)
	jf, _ := json.Marshal(form)

	r := httptest.NewRequest("POST", "/user/10", bytes.NewReader(jf))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, r)

	assert.Equal(t, true, called)
}

func TestJSONFailure(t *testing.T) {
	called := false

	e := echo.New()
	cb := octane.NewBinder()
	e.Binder = cb

	e.POST("/user/:user_id", func(c echo.Context) error {
		called = true

		// swagger:parameters UserCreate
		type request struct {
			// in: body
			Body struct {
				// in: path
				UserID string `json:"user_id" validate:"required"`
				// Required: true
				FirstName string `json:"first_name" validate:"required"`
				// Required: true
				LastName string `json:"last_name" validate:"required"`
			}
		}

		req := new(request).Body

		assert.NotNil(t, c.Bind(&req))

		assert.Equal(t, "10", req.UserID)
		assert.Equal(t, "", req.FirstName)
		assert.Equal(t, "", req.LastName)
		return nil
	})

	form := url.Values{}
	form.Add("first_name", "john")
	form.Add("last_name", "smith")

	r := httptest.NewRequest("POST", "/user/10", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, r)

	assert.Equal(t, true, called)
}

func TestJSONFailureNil(t *testing.T) {
	called := false

	e := echo.New()
	cb := octane.NewBinder()
	e.Binder = cb

	e.POST("/user/:user_id", func(c echo.Context) error {
		called = true

		// swagger:parameters UserCreate
		type request struct {
			// in: body
			Body struct {
				// in: path
				UserID string `json:"user_id" validate:"required"`
				// Required: true
				FirstName string `json:"first_name" validate:"required"`
				// Required: true
				LastName string `json:"last_name" validate:"required"`
			}
		}

		req := new(request).Body

		assert.NotNil(t, c.Bind(&req))

		assert.Equal(t, "10", req.UserID)
		assert.Equal(t, "", req.FirstName)
		assert.Equal(t, "", req.LastName)
		return nil
	})

	r := httptest.NewRequest("POST", "/user/10", nil)
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, r)

	assert.Equal(t, true, called)
}

func TestJSONFailureDataType(t *testing.T) {
	called := false

	e := echo.New()
	cb := octane.NewBinder()
	e.Binder = cb

	e.POST("/user/:user_id", func(c echo.Context) error {
		called = true

		// swagger:parameters UserCreate
		type request struct {
			// in: body
			Body struct {
				// in: path
				UserID string `json:"user_id" validate:"required"`
				// Required: true
				FirstName string `json:"first_name" validate:"required"`
				// Required: true
				LastName string `json:"last_name" validate:"required"`
			}
		}

		req := new(request).Body

		assert.NotNil(t, c.Bind(&req))

		assert.Equal(t, "10", req.UserID)
		assert.Equal(t, "", req.FirstName)
		assert.Equal(t, "", req.LastName)
		return nil
	})

	// Try to parse to JSON and should not fail.
	r := httptest.NewRequest("POST", "/user/10", strings.NewReader("food text"))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	e.ServeHTTP(w, r)

	assert.Equal(t, true, called)
}
