package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/tiskae/go-social/internal/store"
)

type PostKey string

const postKey PostKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

// CreatePost godoc
//
//	@Summary		Create a post
//	@Description	Create a post for the user with the auth token
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			title	body		string		true	"Post title"	maxlength(100)
//	@Param			content	body		string		true	"Post content"	maxlength(1000)
//	@Param			tags	body		[]string	false	"Post tags"
//	@Success		201		{object}	store.Post
//	@Failure		400		{string}	error	"Invalid body"
//	@Failure		500		{string}	error	"Internal server error"
//	@Router			/posts [post]
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	err := Validate.Struct(payload)
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	user := getUserFromContext(r)

	post := store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  user.ID,
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, &post); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}
}

// GetPostByID godoc
//
//	@Summary		Fetches a post
//	@Description	Fetches a post by ID
//	@Tags			posts
//	@Produce		json
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	store.Post
//	@Failure		400	{string}	error	"Invalid post ID"
//	@Failure		404	{string}	error	"post not found"
//	@Failure		500	{string}	error	"Internal server error"
//	@Router			/posts/{id} [get]
func (app *application) getPostByIDHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	// handling absent post from ctx
	if post == nil {
		app.internalServerErrorResponse(w, r, errors.New("post not fetched"))
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}

// DeletePost godoc
//
//	@Summary		Delete a post
//	@Description	Delete the post with the ID provided
//	@Tags			posts
//	@Produce		json
//	@Param			id	path		int	true	"ID of the user to follow"
//	@Success		200	{object}	interface{}
//	@Failure		400	{string}	error	"Invalid post ID"
//	@Failure		404	{string}	error	"Post not found"
//	@Failure		500	{string}	error	"Internal server error"
//	@Router			/posts/{id} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)

	// handling invalid or empty postID
	if err != nil {
		app.badRequestErrorResponse(w, r, errors.New("post id is required as a valid integer"))
		return
	}

	err = app.store.Posts.Delete(r.Context(), postID)

	// handling failed deletion
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundErrorResponse(w, r, err) // post not found, so nothing was deleted
			return
		default:
			app.internalServerErrorResponse(w, r, err)
			return
		}
	}

	if err = app.jsonResponse(w, http.StatusOK, map[string]string{"message": "post deleted successfully!"}); err != nil {
		// handling failed JSON write
		app.internalServerErrorResponse(w, r, err)
	}
}

type UpdatePostPayload struct {
	Title   *string   `json:"title" validate:"omitempty,max=100"`
	Content *string   `json:"content" validate:"omitempty,max=1000"`
	Tags    *[]string `json:"tags" validate:"omitempty,dive,required"`
}

// godoc UpdatePost
//
//	@Summary		Update a post
//	@Description	Update a post with the post body
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int			true	"Post ID"
//	@Param			title	body		string		false	"Post title"	maxlength(100)
//	@Param			content	body		string		false	"Post body"		maxlength(1000)
//	@Param			tags	body		[]string	false	"Post tags"
//	@Success		200		{object}	store.Post
//	@Failure		400		{string}	error	"Invalid body"
//	@Failure		404		{string}	error	"Post not found"
//	@Failure		500		{string}	error	"Internal server error"
//	@Router			/posts/{id} [patch]
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)

	// handling invalid or empty postID
	if err != nil {
		app.badRequestErrorResponse(w, r, errors.New("post id is required as a valid integer"))
		return
	}

	post := getPostFromCtx(r)

	// handling failed post fetch from ctx
	if post == nil {
		app.internalServerErrorResponse(w, r, errors.New("post not fetched"))
		return
	}

	var payload UpdatePostPayload

	err = readJSON(w, r, &payload)

	// handling bad payload
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
		return
	}

	err = Validate.Struct(payload)

	// handling payload validation error
	if err != nil {
		app.badRequestErrorResponse(w, r, err)
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}
	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Tags != nil {
		post.Tags = *payload.Tags
	}

	err = app.store.Posts.UpdateOne(r.Context(), postID, post)

	// handling failed update
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

	if err = app.jsonResponse(w, http.StatusOK, post); err != nil {
		// handling failed JSON write
		app.internalServerErrorResponse(w, r, err)
	}
}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip this middleware for /comments route
		if strings.HasPrefix(r.URL.Path, "/comments") {
			next.ServeHTTP(w, r)
			return
		}

		postID, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)

		// handling invalid of missing postID
		if err != nil {
			app.badRequestErrorResponse(w, r, errors.New("post id is required as a valid integer"))
			return
		}

		post, err := app.store.Posts.GetByID(r.Context(), postID)

		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound): // post not found
				app.notFoundErrorResponse(w, r, err)
			default:
				app.internalServerErrorResponse(w, r, err) // other error
			}
			return
		}

		comments, err := app.store.Comments.GetByPostID(r.Context(), postID)

		if err != nil {
			app.internalServerErrorResponse(w, r, err)
			return
		}

		post.Comments = comments

		// injecting fetched post (with comments) into the request context
		newCtx := context.WithValue(r.Context(), postKey, &post)

		// calling the next handlerFunc
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post {
	post := r.Context().Value(postKey).(*store.Post)
	return post
}
