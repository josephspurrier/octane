package endpoint

import (
	"net/http"

	"github.com/josephspurrier/octane"
	"github.com/josephspurrier/octane/example/app"
	"github.com/josephspurrier/octane/example/app/store"
)

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

/*
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
	userID, ok := p.Context.UserID(r)
	if !ok {
		return http.StatusInternalServerError, errors.New("invalid user")
	}

	// Get a list of notes for the user.
	group := p.Store.Note.NewGroup()
	_, err := p.Store.Note.FindAllByUser(&group, userID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Copy the items to the JSON model.
	arr := make([]model.Note, 0)
	for _, u := range group {
		item := new(model.Note)
		err = structcopy.ByTag(&u, "db", item, "json")
		if err != nil {
			return http.StatusInternalServerError, err
		}
		arr = append(arr, *item)
	}

	// Create the response.
	m := new(model.NoteIndexResponse).Body
	m.Notes = arr

	return p.Response.JSON(w, m)
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
	type request struct {
		// in: path
		NoteID string `json:"note_id" validate:"required"`
	}

	// Request validation.
	req := new(request)
	if err := p.Bind.UnmarshalAndValidate(req, r); err != nil {
		return http.StatusBadRequest, err
	}

	// Get the user ID.
	userID, ok := p.Context.UserID(r)
	if !ok {
		return http.StatusInternalServerError, errors.New("invalid user")
	}

	// Get the note for the user.
	note := p.Store.Note.New()
	exists, err := p.Store.Note.FindOneByIDAndUser(&note, req.NoteID, userID)
	if err != nil {
		return http.StatusInternalServerError, err
	} else if !exists {
		return http.StatusBadRequest, errors.New("invalid note")
	}

	// Copy the items to the JSON model.
	item := new(model.Note)
	err = structcopy.ByTag(&note, "db", item, "json")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Create the response.
	m := new(model.NoteShowResponse).Body
	m.Note = *item

	return p.Response.JSON(w, m)
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
	type request struct {
		// in: path
		NoteID string `json:"note_id" validate:"required"`
		// in: body
		Body struct {
			Message string `json:"message"`
		}
	}

	// Request validation.
	req := new(request)
	if err := p.Bind.UnmarshalAndValidate(req, r); err != nil {
		return http.StatusBadRequest, err
	}

	// Get the user ID.
	userID, ok := p.Context.UserID(r)
	if !ok {
		return http.StatusInternalServerError, errors.New("invalid user")
	}

	// Determine if the note exists for the user.
	note := p.Store.Note.New()
	exists, err := p.Store.Note.FindOneByIDAndUser(&note, req.NoteID, userID)
	if err != nil {
		return http.StatusInternalServerError, err
	} else if !exists {
		return http.StatusBadRequest, errors.New("note does not exist")
	}

	// Update the note.
	_, err = p.Store.Note.Update(req.NoteID, userID, req.Body.Message)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return p.Response.OK(w, "note updated")
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
	type request struct {
		// in: path
		NoteID string `json:"note_id" validate:"required"`
	}

	// Request validation.
	req := new(request)
	if err := p.Bind.UnmarshalAndValidate(req, r); err != nil {
		return http.StatusBadRequest, err
	}

	// Get the user ID.
	userID, ok := p.Context.UserID(r)
	if !ok {
		return http.StatusInternalServerError, errors.New("invalid user")
	}

	// Get a the note for the user.
	affected, err := p.Store.Note.DeleteOneByIDAndUser(req.NoteID, userID)
	if err != nil {
		return http.StatusInternalServerError, err
	} else if affected == 0 {
		return http.StatusBadRequest, errors.New("note does not exist")
	}

	return p.Response.OK(w, "note deleted")
}
*/
