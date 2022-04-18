package octane

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/form/v4"
	"github.com/labstack/echo/v4"
	validator "gopkg.in/go-playground/validator.v9"
)

// Set a Decoder instance as a package global, because it caches
// meta-data about structs, and an instance can be shared safely.
var decoder = form.NewDecoder()

// IRouter extracts a URL parameter value.
type IRouter interface {
	Param(param string) string
}

// Binder contains the request bind an validator objects.
type Binder struct {
	validator *validator.Validate
}

// NewBinder returns a new binder for request bind and validation.
func NewBinder() *Binder {
	decoder.SetTagName("json")

	return &Binder{
		validator: validator.New(),
	}
}

// Bind will unmarshal and validate a struct from a request.
func (b *Binder) Bind(i interface{}, c echo.Context) (err error) {
	return b.unmarshalAndValidate(i, c.Request(), c)
}

// UnmarshalAndValidate will unmarshal and validate a struct using the validator.
func (b *Binder) unmarshalAndValidate(s interface{}, r *http.Request, router IRouter) (err error) {
	if err = b.Unmarshal(s, r, router); err != nil {
		return
	} else if err = b.validate(s); err != nil {
		return
	}

	return
}

// Validate will validate a struct using the validator.
func (b *Binder) validate(s interface{}) error {
	return b.validator.Struct(s)
}

// Unmarshal will perform an unmarshal on an interface using: form or JSON.
func (b *Binder) Unmarshal(iface interface{}, r *http.Request, router IRouter) (err error) {
	// Check for errors.
	v := reflect.ValueOf(iface)
	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value")
	}

	// Load the map.
	m := make(map[string]interface{})

	// Try to auto detect data type based on on the header.
	// Header can having multiple values separated by a semicolon.
	ct := r.Header.Get("Content-Type")
	switch true {
	case ct == "", strings.Contains(ct, "application/x-www-form-urlencoded"):
		err := r.ParseForm()
		if err != nil {
			return fmt.Errorf("body could not be read: %v", err.Error())
		}

		// Loop through each field to extract the URL parameter.
		elem := reflect.Indirect(v.Elem())
		keys := elem.Type()
		for j := 0; j < elem.NumField(); j++ {
			tag := keys.Field(j).Tag
			tagvalue := tag.Get("json")
			pathParam := router.Param(tagvalue)
			if len(pathParam) > 0 {
				r.PostForm.Add(tagvalue, pathParam)
			}
		}

		// r.PostForm is a map of our POST form values.
		err = decoder.Decode(iface, r.PostForm)
		if err != nil {
			return fmt.Errorf("form could not be decoded: %v", err.Error())
		}

		return nil
	case strings.Contains(ct, "application/json"):
		// Decode to the interface.
		_ = json.NewDecoder(r.Body).Decode(&m)
		r.Body.Close()
		// if err != nil {
		// No longer fail on an unmarshal error. This is so users can submit
		// empty data for GET requests, yet we can still map the URL
		// parameter by using the same logic.
		//}

		// Copy the map items to a new map.
		mt := make(map[string]interface{})
		for key, value := range m {
			mt[key] = value
		}

		// Save the map to the body to handle cases where there is a body
		// defined.
		m["body"] = mt
	}

	// Loop through each field to extract the URL parameter.
	elem := reflect.Indirect(v.Elem())
	keys := elem.Type()
	for j := 0; j < elem.NumField(); j++ {
		tag := keys.Field(j).Tag
		tagvalue := tag.Get("json")
		pathParam := router.Param(tagvalue)
		if len(pathParam) > 0 {
			m[tagvalue] = pathParam
		}
	}

	// Convert to JSON.
	var data []byte
	data, err = json.Marshal(m)
	if err != nil {
		return
	}

	// Unmarshal to the interface from JSON.
	return json.Unmarshal(data, &iface)
}
