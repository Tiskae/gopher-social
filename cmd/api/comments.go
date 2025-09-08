package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tiskae/go-social/internal/store"
)

type CreateCommentPayload struct {
	UserID  int64  `json:"user_id" validate:"required"`
	Content string `json:"content" validate:"required"`
}

func (app *application) createPostCommentHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)

	// handling invalid or empty post id
	if err != nil {
		app.badRequestError(w, r, errors.New("post id is required as a valid integer"))
		return
	}

	var payload CreateCommentPayload
	// parsing payload and handling error if any
	if err = readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	err = Validate.Struct(payload)

	// handling failed payload validation
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// userID, err := strconv.ParseInt(payload.UserID, 10, 64)

	// handling failed payload validation on userID
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	comment := store.Comment{
		PostID:  postID,
		UserID:  payload.UserID,
		Content: payload.Content,
	}

	err = app.store.Comments.Create(r.Context(), &comment)

	// handling failed comment creation on DB
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonWrite(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
	}

}

func (app *application) getCommentsByPostIDHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)

	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	comments, err := app.store.Comments.GetByPostID(r.Context(), postID)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err = app.jsonWrite(w, http.StatusOK, comments); err != nil {
		app.internalServerError(w, r, err)
	}
}
