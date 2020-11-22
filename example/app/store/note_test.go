package store_test

import (
	"testing"

	"github.com/josephspurrier/octane/example/app/lib/testutil"
	"github.com/josephspurrier/octane/example/app/store"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNote(t *testing.T) {
	e := echo.New()
	db := testutil.LoadDatabase(e.Logger)
	defer testutil.TeardownDatabase(db)

	// Create a user.
	userID, err := store.CreateUser(db, "first", "last", "email", "password")
	assert.NoError(t, err)

	// Create a note.
	noteID, err := store.NoteCreate(db, userID, "note1")
	assert.NoError(t, err)

	// Read the note.
	n := new(store.Note)
	exists, err := store.FindOneByIDAndUser(db, n, noteID, userID)
	assert.NoError(t, err)
	assert.Equal(t, true, exists)
	assert.Equal(t, "note1", n.Message)

	// Update the note.
	affected, err := store.NoteUpdate(db, noteID, userID, "note2")
	assert.NoError(t, err)
	assert.Equal(t, 1, affected)

	// Verify the note.
	exists, err = store.FindOneByIDAndUser(db, n, noteID, userID)
	assert.NoError(t, err)
	assert.Equal(t, true, exists)
	assert.Equal(t, "note2", n.Message)

	// Delete the node.
	affected, err = store.DeleteOneByIDAndUser(db, new(store.Note), noteID, userID)
	assert.NoError(t, err)
	assert.Equal(t, 1, affected)
}
