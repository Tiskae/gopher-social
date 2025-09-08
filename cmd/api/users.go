package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tiskae/go-social/internal/store"
)

func (app *application) getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)

	// handling invalid or empty userID URL param
	if err != nil {
		app.badRequestError(w, r, errors.New("user id must be a valid integer"))
		return
	}

	user, err := app.store.Users.GetByID(r.Context(), userID)

	// handling DB query errors
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

	if err := app.jsonWrite(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}
