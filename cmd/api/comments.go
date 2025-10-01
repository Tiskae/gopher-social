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

// CreateComment godoc
//
//	@Summary		Create a comment
//	@Description	Create a new comment
//	@Tags			comments
//	@Produce		json
//	@Param			post_id	path		int		true	"Post ID"
//	@Param			user_id	body		int		true	"User ID"
//	@Param			content	body		string	true	"Comment content"
//	@Success		201		{object}	store.Comment
//	@Failure		400		{string}	error	"Invalid body"
//	@Failure		500		{string}	error	"Internal server error"
//	@Router			/posts/{post_id}/comments [post]
func (app *application) createPostCommentHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)

	// handling invalid or empty post id
	if err != nil {
		app.badRequestErrorResponse(w, r, errors.New("post id is required as a valid integer"))
		return
	}

	var payload CreateCommentPayload
	// parsing payload and handling error if any
	if err = readJSON(w, r, &payload); err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	err = Validate.Struct(payload)

	// handling failed payload validation
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
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
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}

}

// GetPostComments godoc
//
//	@Summary		Get post comments
//	@Description	Create a new comment
//	@Tags			comments
//	@Produce		json
//	@Param			post_id	path		int	true	"Post ID"
//	@Success		200		{object}	[]store.Comment
//	@Failure		400		{string}	error	"Invalid post ID"
//	@Failure		404		{string}	error	"Post not found"
//	@Failure		500		{string}	error	"Internal server error"
//	@Router			/posts/id/comments [get]
func (app *application) getCommentsByPostIDHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)

	if err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	comments, err := app.store.Comments.GetByPostID(r.Context(), postID)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundErrorResponse(w, r, err)
			return
		default:
			app.internalServerErrorResponse(w, r, err)
			return
		}
	}

	if err = app.jsonResponse(w, http.StatusOK, comments); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}
