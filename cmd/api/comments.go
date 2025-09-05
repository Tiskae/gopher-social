package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tiskae/go-social/internal/store"
)

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
