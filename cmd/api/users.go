package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tiskae/go-social/internal/store"
)

type UserContextKey string

const userKey UserContextKey = "user"

// ActivateUser godoc
//
//	@Summary		Activate/Register a user
//	@Description	Activate/Register a user by invitation toke
//	@Tags			users
//	@Produce		json
//	@Param			token	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"User activated"
//	@Failure		404		{string}	error
//	@Failure		500		{string}	error
//	@Router			/users/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	err := app.store.Users.Activate(r.Context(), token)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.badRequestErrorResponse(w, r, err)
		default:
			app.internalServerErrorResponse(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}

// GetUserByID godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	store.User
//	@Failure		400	{string}	error	"Invalid user ID"
//	@Failure		404	{string}	error	"User not found"
//	@Failure		500	{string}	error	"Internal server error"
//	@Router			/users/{id} [get]
func (app *application) getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)

	// handling invalid or empty userID URL param
	if err != nil {
		app.badRequestErrorResponse(w, r, errors.New("user id must be a valid integer"))
		return
	}

	user, err := app.store.Users.GetByID(r.Context(), userID)

	// handling DB query errors
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

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}

// FollowUser godoc
//
//	@Summary		Follow a user
//	@Description	Follow a user with the ID provided
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		204	{nil}		nil		"User followed successfully"
//	@Failure		400	{string}	error	"Invalid user ID"
//	@Failure		404	{string}	error	"User not found"
//	@Failure		500	{string}	error	"Internal server error"
//	@Router			/users/{user_id}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerID := getUserFromContext(r).ID

	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)

	// handling invalid or empty userID URL param
	if err != nil {
		app.badRequestErrorResponse(w, r, errors.New("user id must be a valid integer"))
		return
	}

	err = app.store.Followers.Follow(r.Context(), followerID, userID)

	if err != nil {
		switch err {
		case store.ErrConflict:
			app.conflictErrorResponse(w, r, err)
			return
		case store.ErrNotFound:
			app.notFoundErrorResponse(w, r, err)
			return
		default:
			app.internalServerErrorResponse(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}

// UnfollowUser godoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user with the ID provided
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user_id	body		int		true	"ID of the user to unfollow"
//	@Success		204		{nil}		nil		"User unfollowed successfully"
//	@Failure		400		{string}	error	"Invalid user ID"
//	@Failure		404		{string}	error	"User not found"
//	@Failure		500		{string}	error	"Internal server error"
//	@Router			/users/{user_id}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followerID := getUserFromContext(r).ID

	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestErrorResponse(w, r, errors.New("user id must be a valid integer"))
		return
	}

	err = app.store.Followers.Unfollow(r.Context(), followerID, userID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundErrorResponse(w, r, err)
			return
		default:
			app.internalServerErrorResponse(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerErrorResponse(w, r, err)
	}
}

func getUserFromContext(r *http.Request) *store.User {
	user := r.Context().Value(userKey).(*store.User)
	return user
}
