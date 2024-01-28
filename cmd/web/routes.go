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
	mux.Use(app.SessionLoad)

	// define applicaiton routes
	mux.Get("/", app.HomePage)
	mux.Get("/login", app.LoginPage)
	mux.Post("/login", app.PostLoginPage)
	mux.Get("/logout", app.Logout)
	mux.Get("/register", app.RegisterPage)
	mux.Post("/register", app.PostRegisterPage)
	mux.Get("/activate", app.ActivateAccount)

	mux.Mount("/members", app.authRouter())

	return mux
}

func (app *Config) authRouter() http.Handler {

	// Create router
	mux := chi.NewRouter()

	// set up Middleware
	mux.Use(app.Auth)

	// define applicaiton routes
	mux.Get("/plans", app.ChooseSubscription)
	mux.Get("/subscribe", app.SubscribeToPlan)

	return mux
}
