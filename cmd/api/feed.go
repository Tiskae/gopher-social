package main

import (
	"net/http"

	"github.com/tiskae/go-social/internal/store"
)

// GetUserFeed godoc
//
//	@Summary		Get user feed
//	@Description	Get the feed for the user with the auth token
//	@Tags			users
//	@Produce		json
//	@Param			limit	query		int			false	"How many posts to return"
//	@Param			offset	query		int			false	"Offset to start from"
//	@Param			limit	query		string		false	"Whether to sort in ascending (asc) or descending(desc, default)"
//	@Param			tags	query		[]string	false	"Tags to filter by"
//	@Success		200		{object}	[]store.PostWithMetadata
//	@Failure		400		{string}	error	"Invalid query params"
//	@Failure		500		{string}	error	"Internal server error"
//	@Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {

	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
		Tags:   []string{},
	}

	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	err = Validate.Struct(fq)

	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	feed, err := app.store.Posts.GetUserFeed(ctx, int64(99), fq) // TODO: replace with authctd userID

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
