package config

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/josephspurrier/octane/example/app/lib/database"
	"github.com/josephspurrier/rove"
	"github.com/josephspurrier/rove/pkg/adapter/mysql"
	"github.com/labstack/echo/v4"
)

// Database migrates the database and then returns the database connection.
func Database(l echo.Logger) *database.DBW {
	// If the host env var is set, use it.
	host := os.Getenv("MYSQL_HOST")
	if len(host) == 0 {
		host = "127.0.0.1"
	}

	// If the password env var is set, use it.
	password := os.Getenv("MYSQL_ROOT_PASSWORD")

	// Set the database connection information.
	con := &mysql.Connection{
		Hostname:  host,
		Username:  "admin",
		Password:  password,
		Name:      "main",
		Port:      3306,
		Parameter: "collation=utf8mb4_unicode_ci&parseTime=true&multiStatements=true",
	}

	// Migrate the database.
	dbx, err := migrate(l, con, Changesets)
	if err != nil {
		l.Fatalf(err.Error())
	}

	return database.New(dbx, con.Name)
}

// migrate will run the database migrations and will create the database if it
// does not exist.
func migrate(l echo.Logger, con *mysql.Connection, changesets string) (*sqlx.DB, error) {
	// Connect to the database.
	db, err := mysql.New(con)
	if err != nil {
		// Attempt to connect without the database name.
		d, err := con.Connect(false)
		if err != nil {
			return nil, err
		}

		// Create the database.
		_, err = d.Query(fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS %v
		DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;`, con.Name))
		if err != nil {
			return nil, err
		}
		if l != nil {
			l.Printf("Database created.")
		}

		// Attempt to reconnect with the database name.
		db, err = mysql.New(con)
		if err != nil {
			return nil, err
		}
	}

	// Perform all migrations against the database.
	r := rove.NewChangesetMigration(db, changesets)
	r.Verbose = false
	return db.DB, r.Migrate(0)
}
