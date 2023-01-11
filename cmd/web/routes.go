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
	router.Get("/charge-credit-card", app.ChargeCreditCard)
	// router.Post("/process-card-data", app.ProcessCardData)
	// router.Get("/exec-")

	router.Get("/donate", app.Donate)
	router.Get("/crypto-authentication", app.CryptoAuthentication)
	router.Post("/succeeded-registration", app.ProcessRegisterData)

	fileServer := http.FileServer(http.Dir("./static"))
	router.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return router
} 