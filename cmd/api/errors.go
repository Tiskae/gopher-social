package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusInternalServerError, "internal server error")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusNotFound, "not found")
}

func (app *application) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("conflict error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) unauthorizedBasicError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	writeJSONError(w, http.StatusUnauthorized, err.Error())
}
