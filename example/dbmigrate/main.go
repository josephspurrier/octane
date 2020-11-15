package main

import (
	"github.com/josephspurrier/octane/example/app/config"
	"github.com/labstack/echo/v4"
)

func main() {
	// Migrate the database.
	e := echo.New()
	config.Database(e.Logger)
	e.Logger.Printf("Database migration complete.")
}
