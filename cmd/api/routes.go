package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/knbr13/company-service-go/cmd/api/middlewares"
)

func (app *app) SetupRoutes() http.Handler {
	r := chi.NewRouter()
	mdlws := middlewares.NewMiddlewares(app.cfg)

	r.Post("/register", app.hndlrs.Users.Register)
	r.Post("/login", app.hndlrs.Users.Login)

	r.Get("/companies/{id}", app.hndlrs.Companies.GetCompany)
	r.Group(func(r chi.Router) {
		r.Use(mdlws.JWTMiddleware)
		r.Post("/companies", app.hndlrs.Companies.Create)
		r.Patch("/companies/{id}", app.hndlrs.Companies.Update)
		r.Delete("/companies/{id}", app.hndlrs.Companies.Delete)
	})
	return r
}
