package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	router.Use(SessionLoad)

	router.Get("/", app.HomeHandler)
	router.Get("/signup", app.Authorization)
	router.Get("/login", app.Authentication)
	router.Get("/receipt", app.Receipt)
	
	router.Get("/donate", app.Donate)
	
	router.Post("/succeeded-registration", app.ProcessRegisterData)

	return router
} 