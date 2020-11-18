package endpoint

import (
	"net/http"

	"github.com/josephspurrier/octane"
	"github.com/josephspurrier/octane/example/app"
	"github.com/josephspurrier/octane/example/app/lib/structcopy"
	"github.com/josephspurrier/octane/example/app/store"
)

// Note is a note of a user.
type Note struct {
	// example: 314445cd-e9fb-4c58-58b6-777ee06465f5
	// required: true
	UserID string `json:"id"`
	// example: This is a note.
	// required: true
	Message string `json:"message"`
}

// NoteCreate -
// swagger:route POST /api/v1/note note NoteCreate
//
// Create a note.
//
// Security:
//   token:
//
// Responses:
//   201: NoteCreateReponse
//   400: BadRequestResponse
//   401: UnauthorizedResponse
//   500: InternalServerErrorResponse
func NoteCreate(c *app.Context) (err error) {
	// swagger:parameters NoteCreate
	type Request struct {
		// in: body
		Body struct {
			Message string `json:"message"`
		}
	}

	// Request validation.
	req := new(Request)
	if err = c.Bind(req); err != nil {
		return c.BadRequestResponse(err.Error())
	}

	// Get the user ID.
	userID, ok := c.UserID()
	if !ok {
		return c.InternalServerErrorResponse("invalid user")
	}

	// Create the note.
	ID, err := store.NoteCreate(c.DB, userID, req.Body.Message)
	if err != nil {
		return c.InternalServerErrorResponse(err.Error())
	}

	// NoteCreateReponse returns a user ID.
	// swagger:response NoteCreateReponse
	type NoteCreateReponse struct {
		// in: body
		Body struct {
			octane.CreatedStatusFields
			// required: true
			Data struct {
				// RecordID contains the newly created note ID.
				// example: 314445cd-e9fb-4c58-58b6-777ee06465f5
				// required: true
				RecordID string `json:"record_id"`
			} `json:"data"`
		}
	}

	// Set the note ID.
	data := new(NoteCreateReponse).Body.Data
	data.RecordID = ID

	return c.DataResponse(http.StatusCreated, data)
}

// NoteIndex -
// swagger:route GET /api/v1/note note NoteIndex
//
// List notes.
//
// Security:
//   token:
//
// Responses:
//   200: NoteIndexResponse
//   400: BadRequestResponse
//   401: UnauthorizedResponse
//   500: InternalServerErrorResponse
func NoteIndex(c *app.Context) (err error) {
	// Get the user ID.
	userID, ok := c.UserID()
	if !ok {
		return c.InternalServerErrorResponse("invalid user")
	}

	// Get a list of notes for the user.
	group := make([]store.Note, 0)
	_, err = store.NoteFindAllByUser(c.DB, &group, userID)
	if err != nil {
		return c.InternalServerErrorResponse(err.Error())
	}

	// Copy the items to the JSON model.
	arr := make([]Note, 0)
	for _, u := range group {
		item := new(Note)
		err = structcopy.ByTag(&u, "db", item, "json")
		if err != nil {
			return c.InternalServerErrorResponse(err.Error())
		}
		arr = append(arr, *item)
	}

	// NoteIndexResponse returns an array of notes.
	// swagger:response NoteIndexResponse
	type NoteIndexResponse struct {
		// in: body
		Body struct {
			octane.OKStatusFields
			// required: true
			Data struct {
				// required: true
				Notes []Note `json:"notes"`
			} `json:"data"`
		}
	}

	// Set the notes.
	data := new(NoteIndexResponse).Body.Data
	data.Notes = arr

	return c.DataResponse(http.StatusOK, data)
}

// NoteShow -
// swagger:route GET /api/v1/note/{note_id} note NoteShow
//
// Show a note.
//
// Security:
//   token:
//
// Responses:
//   200: NoteShowResponse
//   400: BadRequestResponse
//   401: UnauthorizedResponse
//   500: InternalServerErrorResponse
func NoteShow(c *app.Context) (err error) {
	// swagger:parameters NoteShow
	type Request struct {
		// example: 314445cd-e9fb-4c58-58b6-777ee06465f5
		// in: path
		NoteID string `json:"note_id" validate:"required"`
	}

	// Request validation.
	req := new(Request)
	if err = c.Bind(req); err != nil {
		return c.BadRequestResponse(err.Error())
	}

	// Get the user ID.
	userID, ok := c.UserID()
	if !ok {
		return c.InternalServerErrorResponse("invalid user")
	}

	// Get the note for the user.
	note := new(store.Note)
	exists, err := store.FindOneByIDAndUser(c.DB, note, req.NoteID, userID)
	if err != nil {
		return c.InternalServerErrorResponse(err.Error())
	} else if !exists {
		return c.BadRequestResponse("invalid note")
	}

	// Copy the items to the JSON model.
	item := new(Note)
	err = structcopy.ByTag(&note, "db", item, "json")
	if err != nil {
		return c.InternalServerErrorResponse(err.Error())
	}

	// NoteShowResponse returns 200.
	// swagger:response NoteShowResponse
	type NoteShowResponse struct {
		// in: body
		Body struct {
			octane.OKStatusFields
			// required: true
			Data struct {
				// required: true
				Note Note `json:"note"`
			} `json:"data"`
		}
	}

	// Create the response.
	data := new(NoteShowResponse).Body.Data
	data.Note = *item

	return c.DataResponse(http.StatusOK, data)
}

// NoteUpdate -
// swagger:route PUT /api/v1/note/{note_id} note NoteUpdate
//
// Update a note.
//
// Security:
//   token:
//
// Responses:
//   200: OKResponse
//   400: BadRequestResponse
//   401: UnauthorizedResponse
//   500: InternalServerErrorResponse
func NoteUpdate(c *app.Context) (err error) {
	// swagger:parameters NoteUpdate
	type Request struct {
		// in: path
		// example: 314445cd-e9fb-4c58-58b6-777ee06465f5
		NoteID string `json:"note_id" validate:"required"`
		// in: body
		Body struct {
			// example: This is a note.
			Message string `json:"message"`
		}
	}

	// Request validation.
	req := new(Request)
	if err = c.Bind(req); err != nil {
		return c.BadRequestResponse(err.Error())
	}

	// Get the user ID.
	userID, ok := c.UserID()
	if !ok {
		return c.InternalServerErrorResponse("invalid user")
	}

	// Determine if the note exists for the user.
	note := new(store.Note)
	exists, err := store.FindOneByIDAndUser(c.DB, note, req.NoteID, userID)
	if err != nil {
		return c.InternalServerErrorResponse(err.Error())
	} else if !exists {
		return c.BadRequestResponse("invalid note")
	}

	// Update the note.
	_, err = store.NoteUpdate(c.DB, req.NoteID, userID, req.Body.Message)
	if err != nil {
		return c.InternalServerErrorResponse(err.Error())
	}

	return c.OKResponse("note updated")
}

// NoteDestroy -
// swagger:route DELETE /api/v1/note/{note_id} note NoteDestroy
//
// Delete a note.
//
// Security:
//   token:
//
// Responses:
//   200: OKResponse
//   400: BadRequestResponse
//   401: UnauthorizedResponse
//   500: InternalServerErrorResponse
func NoteDestroy(c *app.Context) (err error) {
	// swagger:parameters NoteDestroy
	type Request struct {
		// in: path
		// example: 314445cd-e9fb-4c58-58b6-777ee06465f5
		NoteID string `json:"note_id" validate:"required"`
	}

	// Request validation.
	req := new(Request)
	if err = c.Bind(req); err != nil {
		return c.BadRequestResponse(err.Error())
	}

	// Get the user ID.
	userID, ok := c.UserID()
	if !ok {
		return c.InternalServerErrorResponse("invalid user")
	}

	// Delete the note for the user.
	affected, err := store.DeleteOneByIDAndUser(c.DB, new(store.User), req.NoteID, userID)
	if err != nil {
		return c.InternalServerErrorResponse(err.Error())
	} else if affected == 0 {
		return c.BadRequestResponse("note does not exist")
	}

	return c.OKResponse("note deleted")
}
