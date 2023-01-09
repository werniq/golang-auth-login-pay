package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"
	"web-application/internal/models"
	driver "web-application/internal/models/drivers"

	"golang.org/x/crypto/bcrypt"
)

// type stripePayload struct {
// 	Currency      string `json:"currency"`
// 	Amount        string `json:"amount"`
// 	PaymentMethod string `json:"payment_method"`
// 	Email         string `json:"email"`
// 	CardBrand     string `json:"card_brand"`
// 	ExpiryMonth   int 	 `json:"exp_month"`
// 	ExpiryYear    int 	 `json:"exp_year"`
// 	LastFour      string `json:"last_four"`
// 	Plan          string `json:"plan"`
// 	ProductID     string `json:"product_id"`
// 	FirstName     string `json:"first_name"`
// 	LastName      string `json:"last_name"`
// }

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


// ProcessRegisterData gets all value from request forms, and check if username is in database
// And if it is inputed correctly(by given instructions)
// Checks password also, and if all "Tests" are passed, save user in database
// Returns user(if all inputs are correct), errorData, which represents all errors, or mistakes
// while registrating. Error, if any exist. Also creates hashed password, and also stored it 
// In database
func (app *application) ProcessRegisterData(r *http.Request) (models.User, []string, error) {
	var u models.User
	var errorData = []string{}


	fmt.Println("Processing register data...")
	
	// Data := make(map[string]interface{})

	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
		return u, errorData, err
	}
	
	passwordContainLowercase, passwordContainNumber, passwordContainSpecial, passwordContainUppercase := false, false, false, false


	username := r.Form.Get("reg-username")
	password := r.Form.Get("reg_password")
	repeatPassword := r.Form.Get("repeat_password")
	
	var ok bool = true
	// Checking, if Username is valid
	for _, v := range username {
		if !unicode.IsLetter(v) && !unicode.IsNumber(v) {
			ok = false
		}
	}
	if !ok {
		errorData = append(errorData, "Username should contain only symbols or numbers")
	}

	if len(username) > 50 || len(username) < 3 {
		errorData = append(errorData, "Username should contain more than 3 symbols, and less than 55")
	}

	// Check whether user is already in database..
	conn, err := driver.OpenDB(app.cfg.db.dsn)
	if err != nil {
		app.errorLog.Println(err)
		return u, errorData, err
	}
	// "SELECT UserID FROM bcrypt WHERE username = ?"
	stmt := "SELECT id FROM users WHERE username = $1"
	row := conn.QueryRow(stmt, username)
	var uID int
	err = row.Scan(&uID)
	if err != sql.ErrNoRows {
		fmt.Println(err)
		errorData = append(errorData, "Username already exists")
	}
	
	// Audit password for patterns
	if password != repeatPassword {
		errorData = append(errorData, "Password1 and password2 should be the same")
	}

	fmt.Println("Password is: 		 \t", password)
	fmt.Println("Username is: 		 \t", username)
	fmt.Println("Repeat password is: \t", repeatPassword)
	for _, v := range password {
		switch {
			case unicode.IsLower(v):
				passwordContainLowercase = true
			case unicode.IsUpper(v):
				passwordContainUppercase = true
			case unicode.IsNumber(v):
				passwordContainNumber = true
			case unicode.IsPunct(v) || unicode.IsSymbol(v):
				passwordContainSpecial = true
			}
	}
	if len(password) < 4 || len(password) > 15 {
		errorData = append(errorData, "Password length have to be more than 4 and less than 15")
	}

	// passwordContainLowercase, passwordContainNumber, 
	// passwordContainSpecial, passwordContainUppercase 
	for i := 0; i < 4; i++ {
		switch {
			case !passwordContainLowercase: 
				errorData = append(errorData, "Password must contain at least 1 lowercase letter")
				passwordContainLowercase = true
			case !passwordContainNumber:
				errorData = append(errorData, "Password must contain at least 1 integer")
				passwordContainNumber = true
			case !passwordContainSpecial:
				errorData = append(errorData, "Password must contain at least 1 special symbol")
				passwordContainSpecial = true
			case !passwordContainUppercase:
				errorData = append(errorData, "Password must contain at least 1 uppercase letter")
				passwordContainUppercase = true
		}
	}
	// if errorData != nil {
		// app.RegisterRenderWithError(w, r, errorData)
		// app.RegisterRenderWithError(w, r, errorData)
		// return u, errorData, nil
	// }
	fmt.Println("PROCESSING THIS------------------")
	// Creating password hash
	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		app.errorLog.Println(err)
		return u, errorData, err
	}
	if errorData != nil {
		// app.errorLog.Println(err)
		return u, errorData, nil
	}

	//
	// CREATE FOLLOWING COLUMNS!
	// 
	stmt = `
		INSERT INTO
			users
			(username, firstname, lastname, email, 
			password, hashed_password, address1, address2, date_of_birth)
		VALUES
			(?, ?, ?, ?, ?, ?, ?, ?, ?)
	`


	firstname := r.Form.Get("reg-firstname")
	lastname := r.Form.Get("reg-surname")
	email := r.Form.Get("reg-email")
	address1 := r.Form.Get("reg-address1")
	address2 := r.Form.Get("reg-address1")
	// dateOfBirth := r.Form.Get("reg-date-of-birth")
	// date, err := time.Parse("2006-01-02", dateOfBirth)

	// if err != nil {
	// 	// app.errorLog.Printf("Error converting time: %s", err)
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), 7 * time.Second)
	// defer cancel()

	// row = conn.QueryRowContext(
	// 	ctx, 
	// 	stmt, 
	// 	username, firstname, lastname, email, 
	// 	password, hash, address1, address2, dateOfBirth)
	
	// Data["username"] = username
	// Data["firstname"] = firstname 
	// Data["lastname"] = lastname
	// Data["email"] = email
	// Data["address1"] = address1 
	// Data["address2"] = address2
	// Data["dateOfBirth"] = dateOfBirth 
	
	user := models.User{
		Username: 	 username,
		Firstname:   firstname,
		Lastname: 	 lastname,
		Password: 	 password,
		UserHashedPassword: hash,
		Address1: 	 address1,
		Address2: 	 address2,
		Email: 		 email,
		// DateOfBirth: date,
		CreatedAt: 	 time.Now(),
		UpdatedAt: 	 time.Now(),
	}

	err = app.database.SaveUser(user)	
	if err != nil {
		app.errorLog.Fatal(err)
		return u, errorData, err
	}

	// user.ID = uID
	// Data["id"] = uID
	
	// http.Redirect(w, r, "/succeeded-registration", http.StatusCreated)
	fmt.Println("Congradulations, user registered")
	return u, errorData, nil
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