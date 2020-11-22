package testutil

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/josephspurrier/octane/example/app/config"
	"github.com/josephspurrier/octane/example/app/lib/database"
	"github.com/josephspurrier/rove/pkg/adapter/mysql"
	"github.com/labstack/echo/v4"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// LoadDatabase will set up the DB and apply migrations for the tests.
func LoadDatabase(ml echo.Logger) *database.DBW {
	unique := "T" + fmt.Sprint(rand.Intn(999999999))

	// If the host env var is set, use it.
	host := os.Getenv("MYSQL_HOST")
	if len(host) == 0 {
		host = "127.0.0.1"
	}

	// If the user env var is set, use it.
	username := os.Getenv("MYSQL_USER")
	if len(username) == 0 {
		username = "root"
	}

	// If the password env var is set, use it.
	password := os.Getenv("MYSQL_ROOT_PASSWORD")

	// Set the database connection information.
	con := &mysql.Connection{
		Hostname:  host,
		Username:  username,
		Password:  password,
		Name:      "maintest" + unique,
		Port:      3306,
		Parameter: "parseTime=true&allowNativePasswords=true&collation=utf8mb4_unicode_ci&multiStatements=true",
	}

	db, err := config.Migrate(ml, con, config.Changesets)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return database.New(db, con.Name)
}

// TeardownDatabase will destroy the test database.
func TeardownDatabase(db *database.DBW) {
	_, err := db.Exec(`DROP DATABASE IF EXISTS ` + db.Name())
	if err != nil {
		fmt.Println("DB DROP TEARDOWN Error:", err)
	}
}
