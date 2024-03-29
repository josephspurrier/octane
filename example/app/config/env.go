package config

import (
	"github.com/josephspurrier/octane/example/app/lib/env"
	"github.com/labstack/echo/v4"
)

// Settings holds the variables for the application and the defaults.
type Settings struct {
	Port           int    `env:"API_PORT" default:"8080"`
	Secret         string `env:"API_SECRET" default:"TA8tALZAvLVLo4ToI44xF/nF6IyrRNOR6HSfpno/81M="`
	SessionTimeout int    `env:"API_SESSION_TIMEOUT" default:"480"` // 480 min = 8 hours.
}

// LoadEnv will load the settings from the environment variables or defaults.
func LoadEnv(l echo.Logger, prefix string) *Settings {
	s := new(Settings)
	err := env.Unmarshal(s, prefix)
	if err != nil {
		l.Fatalf("error getting environment variables: %v", err.Error())
	}

	return s
}
