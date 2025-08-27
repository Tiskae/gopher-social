package main

import (
	"net/http"
)

func (a *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Health is excellent!"))

	a.store.Posts.Create(r.Context())
}
