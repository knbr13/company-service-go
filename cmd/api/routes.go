package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *app) SetupRoutes() http.Handler {
	r := chi.NewRouter()

	r.Post("/register", app.hndlrs.Users.Register)
	r.Post("/login", app.hndlrs.Users.Login)

	return r
}
