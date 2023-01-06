package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	router.Get("/login", app.Authentication)
	router.Get("/", app.HomeHandler)
	// router.Get("/process-register-data", app.ProcessRegisterData)
	// router.Get("/process-register-data", app.)
	router.Get("/signup", app.Authorization)
	router.Post("/succeeded-registration", app.SucceededRegistration)
	

	// mux.Get("/virtual-terminal", app.VirtualTerminal)
	// mux.Post("/virtual-terminal-payment-succeeded", app.VirtualTerminalPaymentSucceeded)
	// mux.Get("/virtual-terminal-receipt", app.VirtualTerminalReceipt)

	return router
} 