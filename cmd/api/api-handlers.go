package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"web-application/internal/models"
)

type stripePayload struct {
	Currency      string `json:"currency"`
	Amount        string `json:"amount"`
	PaymentMethod string `json:"payment_method"`
	Email         string `json:"email"`
	CardBrand     string `json:"card_brand"`
	ExpiryMonth   int 	 `json:"exp_month"`
	ExpiryYear    int 	 `json:"exp_year"`
	LastFour      string `json:"last_four"`
	Plan          string `json:"plan"`
	ProductID     string `json:"product_id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
}

type TransactionData struct {
	FirstName       string
	LastName        string
	Email           string
	PaymentIntentID string
	PaymentMethodID string
	PaymentAmount   int
	PaymentCurrency string
	LastFour        string
	ExpiryMonth     int
	ExpiryYear      int
	BankReturnCode  string
}

func (app *application) CreateAuthToken(w http.ResponseWriter, r *http.Request) {
	var request struct {
		// Username string `json:"username"`
		Email 	 string `json:"email"`
 		Password string `json:"password"`
	}	

	err := app.ReadJSON(w, r, &request)
	if err != nil {
		app.BadRequest(w, r, err)
		return
	}


	// get the user from database by email or username 
	user, err := app.database.GetUserByEmailOrUsername(request.Email)
	if err != nil {
		app.InvalidCredentials(w)
		return
	}

	// validate the password, with one from database
	validPassword, err := app.ValidatePassword(user.Password, request.Password)

	if !validPassword {
		app.InvalidCredentials(w)
		return
	}

	if err != nil {
		app.InvalidCredentials(w)
		return
	}

	// Generate token for user
	token, err := models.GenerateToken(user.ID, 24 * time.Hour, models.ScopeAuthentication)
	if err != nil {
		app.BadRequest(w, r, err)
		return 
	}
	
	// Save to database
	err = app.database.InsertToken(token, user)
	if err != nil {
		app.BadRequest(w, r, err)
		return
	}


	var payload struct {
		Error 	bool 		  `json:"error"`
		Message string 		  `json:"message"`
		Token	*models.Token `json:"token"`
	}


	payload.Error = false
	payload.Message = "Success"
	payload.Token = token

	out, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)


}	

func (app *application) AuthenticateToken(r *http.Request) (*models.User, error) {
	authorizationToken := r.Header.Get("Authorization")
	if authorizationToken == "" {
		return nil, errors.New("no authorization token recieved")
	}

	headerParts := strings.Split(authorizationToken, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("no authorization header recieved")
	}

	token := headerParts[1]
	if len(token) != 26 {
		return nil, errors.New("authorization token wrong size")
	}

	user, err := app.database.GetUserForToken(token)
	if err != nil {
		return nil, errors.New("no matching users found")
	}


	return user, nil
}

func (app *application) CheckAuthentication(w http.ResponseWriter, r *http.Request) {
	user, err := app.AuthenticateToken(r)
	if err != nil {
		app.InvalidCredentials(w)
		return
	}

	var payload struct {
		Error 	bool `json:"error"`
		Message string `json:"message"`
	}

	payload.Error = false
	payload.Message = fmt.Sprintf("authenticated user %s", user.Email)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	out, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		return 
	}
	w.Write(out)
}