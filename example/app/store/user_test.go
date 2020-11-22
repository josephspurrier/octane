package store_test

import (
	"testing"

	"github.com/josephspurrier/octane/example/app/lib/testutil"
	"github.com/josephspurrier/octane/example/app/store"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	e := echo.New()
	db := testutil.LoadDatabase(e.Logger)
	defer testutil.TeardownDatabase(db)

	// Create a user.
	_, err := store.CreateUser(db, "John", "Smith", "jsmith@example.com", "password")
	assert.NoError(t, err)

	// Verify user.
	u := new(store.User)
	exists, err := store.FindOneByField(db, u, "email", "jsmith@example.com")
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.Equal(t, "John", u.FirstName)

	// Test fail user.
	u = new(store.User)
	exists, err = store.FindOneByField(db, u, "email", "bad email")
	assert.NoError(t, err)
	assert.False(t, exists)
}
