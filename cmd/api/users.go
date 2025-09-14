package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tiskae/go-social/internal/store"
)

type UserContextKey string

const userKey UserContextKey = "user"

func (app *application) getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

type FollowPayload struct {
	UserID int64 `json:"user_id" validate:"required"`
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	userToFollow := getUserFromContext(r)

	// TODO: revert to auth userID later
	var payload FollowPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	err := app.store.Followers.Follow(r.Context(), userToFollow.ID, payload.UserID)

	if err != nil {
		switch err {
		case store.ErrConflict:
			app.conflictError(w, r, err)
			return
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	userToUnfollow := getUserFromContext(r)

	// TODO: revert to auth userID later
	var payload FollowPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	err := app.store.Followers.Unfollow(r.Context(), userToUnfollow.ID, payload.UserID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		// injecting fetched user into the request context
		newCtx := context.WithValue(r.Context(), userKey, &user)

		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}

func getUserFromContext(r *http.Request) *store.User {
	user := r.Context().Value(userKey).(*store.User)
	return user
}
