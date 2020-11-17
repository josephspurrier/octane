package store

import "time"

// User is a person who can login to the application.
type User struct {
	ID        string     `db:"id"`
	FirstName string     `db:"first_name"`
	LastName  string     `db:"last_name"`
	Email     string     `db:"email"`
	Password  string     `db:"password"`
	StatusID  uint8      `db:"status_id"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// Table returns the table name.
func (x *User) Table() string {
	return "user"
}

// PrimaryKey returns the primary key field.
func (x *User) PrimaryKey() string {
	return "id"
}
