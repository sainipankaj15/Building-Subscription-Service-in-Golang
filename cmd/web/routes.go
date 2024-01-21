package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Config) routes() http.Handler {

	// Create router
	mux := chi.NewRouter()

	// set up Middleware
	mux.Use(middleware.Recoverer)

	// define applicaiton routes
	mux.Get("/", app.HomePage)

	return mux
}
