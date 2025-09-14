package main

import (
	"net/http"
)

// APIHealth godoc
//
//	@Summary		API health
//	@Description	Checks the health of the API
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	interface{}
//	@Failure		500	{string}	error	"Internal server error"
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "Health is excellent",
		"env":     app.config.env,
		"version": app.config.version,
	}

	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
	}
}
