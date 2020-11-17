package app

import "database/sql"

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
