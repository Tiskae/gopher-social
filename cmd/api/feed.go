package main

import "net/http"

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: paginatio, filters, searching?

	ctx := r.Context()

	feed, err := app.store.Posts.GetUserFeed(ctx, int64(99)) // TODO: replace with authctd userID

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
