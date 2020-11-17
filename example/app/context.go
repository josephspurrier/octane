package app

import (
	"database/sql"

	"github.com/josephspurrier/octane"
	"github.com/josephspurrier/octane/example/app/lib/passhash"
	"github.com/labstack/echo/v4"
)

// Context -
type Context struct {
	octane.ResponseJSON
	DB       IDatabase
	Passhash *passhash.Passhash
}

// HandlerFunc -
func (ac *Context) HandlerFunc(next func(*Context) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := &Context{
			ResponseJSON: octane.ResponseJSON{Context: c},
			DB:           ac.DB,
		}
		return next(cc)
	}
}

// IDatabase provides query capabilities for database handling.
type IDatabase interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	QueryRowScan(dest interface{}, query string, args ...interface{}) error
	RecordExists(err error) (bool, error)
	AffectedRows(result sql.Result) int
	RecordExistsString(err error, s string) (bool, string, error)
	SuppressNoRowsError(err error) error
}

// IRecord provides table information for use with query helpers.
type IRecord interface {
	Table() string
	PrimaryKey() string
}
