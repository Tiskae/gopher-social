package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tiskae/go-social/internal/store"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	err := Validate.Struct(payload)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	post := store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		// TODO: Change after auth
		UserID: 1,
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, &post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPostByIDHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)

	if err != nil {
		app.badRequestError(w, r, errors.New("post id is required as a valid integer"))
		return
	}

	post, err := app.store.Posts.GetByID(r.Context(), postID)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	comments, err := app.store.Comments.GetByPostID(r.Context(), postID)

	if err != nil {
		app.internalServerError(w, r, err)
	}

	post.Comments = comments

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)

	// handling invalid or empty postID
	if err != nil {
		app.badRequestError(w, r, errors.New("post id is required as a valid integer"))
		return
	}

	err = app.store.Posts.Delete(r.Context(), postID)

	// handling failed deletion
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err) // post not found, so nothing was deleted
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err = writeJSON(w, http.StatusOK, map[string]string{"message": "post deleted successfully!"}); err != nil {
		// handling failed JSON write
		app.internalServerError(w, r, err)
	}
}

type UpdatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags" validate:"required,dive,required"`
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)

	// handling invalid or empty postID
	if err != nil {
		app.badRequestError(w, r, errors.New("post id is required as a valid integer"))
		return
	}

	var payload UpdatePostPayload

	err = readJSON(w, r, &payload)

	// handling bad payload
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	err = Validate.Struct(payload)

	// handling payload validation error
	if err != nil {
		app.badRequestError(w, r, err)
	}

	updatedPost := store.Post{
		ID:      postID,
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
	}

	err = app.store.Posts.UpdateOne(r.Context(), postID, &updatedPost)

	// handling failed update
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

	if err = writeJSON(w, http.StatusOK, updatedPost); err != nil {
		// handling failed JSON write
		app.internalServerError(w, r, err)
	}
}
