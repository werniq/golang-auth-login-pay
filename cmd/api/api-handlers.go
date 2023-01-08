package main

import (
	"encoding/json"
	"net/http"
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

	// At this point, all "tests" should be passed, and if username is in database, and
	// password equal to one from db, user have to be right
	// So, we can generate authentication token
	token, err := 

	var payload struct {
		Error bool `json:"error"`
		Message string `json:"message"`
	}

	payload.Error = false
	payload.Message = "Success"

	out, _ := json.MarshalIndent(payload, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}	

// func (app *application) CheckAuthentication(w http.ResponseWriter, r *http.Request) {
	
// }