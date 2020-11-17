package store

import (
	"fmt"

	"github.com/josephspurrier/octane/example/app"
)

// *****************************************************************************
// Find
// *****************************************************************************

// FindOneByID will find a record by string ID.
func FindOneByID(db app.IDatabase, dest app.IRecord, ID string) (exists bool, err error) {

	err = db.Get(dest, fmt.Sprintf(`
		SELECT * FROM %s
		WHERE %s = ?
		LIMIT 1`, dest.Table(), dest.PrimaryKey()),
		ID)
	return db.RecordExists(err)
}

// FindOneByField will find a record by a specified field.
func FindOneByField(db app.IDatabase, dest app.IRecord, field string, value string) (exists bool, err error) {
	err = db.Get(dest, fmt.Sprintf(`
		SELECT * FROM %s
		WHERE %s = ?
		LIMIT 1`, dest.Table(), field),
		value)
	return db.RecordExists(err)
}

// FindAll returns all users.
func FindAll(db app.IDatabase, dest app.IRecord) (total int, err error) {
	err = db.QueryRowScan(&total, fmt.Sprintf(`
		SELECT COUNT(DISTINCT %s)
		FROM %s
		`, dest.PrimaryKey(), dest.Table()))

	if err != nil {
		return total, db.SuppressNoRowsError(err)
	}

	err = db.Select(dest, fmt.Sprintf(`SELECT * FROM %s`, dest.Table()))
	return total, err
}

// *****************************************************************************
// Delete
// *****************************************************************************

// DeleteOneByID removes one record by ID.
func DeleteOneByID(db app.IDatabase, dest app.IRecord, ID string) (affected int, err error) {
	result, err := db.Exec(fmt.Sprintf("DELETE FROM %s WHERE %s = ? LIMIT 1",
		dest.Table(), dest.PrimaryKey()), ID)
	if err != nil {
		return 0, err
	}

	return db.AffectedRows(result), err
}

// DeleteAll removes all records.
func DeleteAll(db app.IDatabase, dest app.IRecord) (affected int, err error) {
	result, err := db.Exec(fmt.Sprintf(`DELETE FROM %s`, dest.Table()))
	if err != nil {
		return 0, err
	}

	return db.AffectedRows(result), err
}

// *****************************************************************************
// Exists
// *****************************************************************************

// ExistsByID determines if a records exists by ID.
func ExistsByID(db app.IDatabase, dest app.IRecord, value string) (found bool, err error) {
	err = db.Get(dest, fmt.Sprintf(`
		SELECT %s FROM %s
		WHERE %s = ?
		LIMIT 1`, dest.PrimaryKey(), dest.Table(), dest.PrimaryKey()),
		value)
	return db.RecordExists(err)
}

// ExistsByField determines if a records exists by a specified field and
// returns the ID.
func ExistsByField(db app.IDatabase, dest app.IRecord, field string, value string) (found bool, ID string, err error) {
	err = db.QueryRowScan(&ID, fmt.Sprintf(`
		SELECT %s FROM %s
		WHERE %s = ?
		LIMIT 1`, dest.PrimaryKey(), dest.Table(), field),
		value)

	return db.RecordExistsString(err, ID)
}
