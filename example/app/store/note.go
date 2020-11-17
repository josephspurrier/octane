package store

import (
	"time"

	"github.com/josephspurrier/octane/example/app"
	"github.com/josephspurrier/octane/example/app/lib/securegen"
)

// Note is a note belonging a user.
type Note struct {
	ID        string     `db:"id"`
	UserID    string     `db:"user_id"`
	Message   string     `db:"message"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// Table returns the table name.
func (x *Note) Table() string {
	return "note"
}

// PrimaryKey returns the primary key field.
func (x *Note) PrimaryKey() string {
	return "id"
}

// NoteCreate creates a new note.
func NoteCreate(db app.IDatabase, userID, message string) (string, error) {
	uuid, err := securegen.UUID()
	if err != nil {
		return "", err
	}

	_, err = db.Exec(`
		INSERT INTO note
		(id, user_id, message)
		VALUES
		(?,?,?)
		`,
		uuid, userID, message)

	return uuid, err
}

// NoteFindAllByUser returns all notes for a user.
func NoteFindAllByUser(db app.IDatabase, dest *[]Note, userID string) (
	total int, err error) {
	err = db.Select(dest, `
		SELECT *
		FROM note
		WHERE user_id = ?
		ORDER BY message ASC
		`,
		userID)
	return len(*dest), db.SuppressNoRowsError(err)
}
