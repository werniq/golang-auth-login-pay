package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
	"unicode"
	"web-application/internal/models"
	driver "web-application/internal/models/drivers"

	"golang.org/x/crypto/bcrypt"
)

func (app *application) Authentication(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "login", nil); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) Authorization(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "register", nil); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "home", nil); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) Register(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "regisre", nil); err != nil {
		app.errorLog.Println(err)
	}
}


func (app *application) RegisterRenderWithError(w http.ResponseWriter, r *http.Request, errMessage []string) error {
	if err := app.renderTemplate(w, r, "register", &templateData{
			ErrorData: errMessage,
		}); err != nil {
			app.errorLog.Println(err)
			return err
		}
	return nil
}



// ProcessRegisterData gets all value from request forms, and check if username is in database
// And if it is inputed correctly(by given instructions)
// Checks password also, and if all "Tests" are passed, save user in database
// Returns user(if all inputs are correct), errorData, which represents all errors, or mistakes
// while registrating. Error, if any exist. Also creates hashed password, and also stored it 
// In database

func (app *application) CheckUserData(r *http.Request) (models.User, []string, error) {
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
	// http.Redirect(w, r, r.Header.Get("Referer"), 302)
	return u, errorData, nil
}


func (app *application) ProcessRegisterData(w http.ResponseWriter, r *http.Request)  {
	// var u models.User
	var err error
	var errorData = []string{}

	_, errorData, err = app.CheckUserData(r)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	if errorData != nil {
		app.RegisterRenderWithError(w, r, errorData)
	}

	// data := make(map[string]interface{})

	// app.Session.Put(r.Context(), "receipt", u)
	http.Redirect(w, r, "/receipt", http.StatusSeeOther)
	fmt.Println("Redirected")
}

// Receipt returns rendered page, with request data, or data which was in session.
func (app *application) Receipt(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "receipt", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}

// Donate renders template, for choosing whether you want donate in crypto or with credit card
func (app *application) Donate(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "donate", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}


// ChargeCreditCard renders template, for donating with credit card
func (app *application) ChargeCreditCard(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "exec-donate", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) CryptoAuthentication(w http.ResponseWriter, r *http.Request) {	
	if err := app.renderTemplate(w, r, "crypto-login", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}