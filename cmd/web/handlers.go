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
func (app *application) ProcessRegisterData(w http.ResponseWriter, r *http.Request) (models.User, []string, error) {
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
		if unicode.IsLetter(v) == false && unicode.IsNumber(v) == false {
			ok = false
		}
	}
	if ok == false {
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

	fmt.Println("Password is: \t", password)
	fmt.Println("Username is: \t", username)
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
	if errorData != nil {
		// app.RegisterRenderWithError(w, r, errorData)
		app.RegisterRenderWithError(w, r, errorData)
		// return u, errorData, nil
	}
	fmt.Println("PROCESSING THIS------------------")
	// Creating password hash
	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		app.errorLog.Println(err)
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
	fmt.Println("Error number A lot: \t", err)
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

// This function jsut returns all data, which was inputed in 
// forms(if it is inputed correctly, by given instructions)
// And redirect to succeeded registration page 
func (app *application) SucceededRegistration(w http.ResponseWriter, r *http.Request) {
	var u models.User
	var errorData = []string{}
	// var emptyArr = []string{}
	var err error
	u, errorData, err = app.ProcessRegisterData(w, r)
	fmt.Println(errorData)
	fmt.Println(err)
	// if errorData != nil || err != nil {
	// 	app.errorLog.Println(err)
	// 	return 
	// }
	

	// app.Session.Put(r.Context(), "user-data", u)
	// http.Redirect(w, r, "/succeeded-registration", http.StatusOK)
	data := make(map[string]interface{})
	data["user"] = u
	// http.Redirect(w, r, "/succeeded-regstration", http.StatusOK)
	if err := app.renderTemplate(w, r, "succeededRegistration", &templateData{
		Data: data,
	}); err != nil {
		app.errorLog.Println(err)
	}

}

// ENDED HERE -> NEXT STEPS:
// FINISH REGISTRATION FUNCTION
// CREATE LOGIN FUNCTION
// CREATE AUTH TOKEN
// CREATE FUNCTION, WHICH CHECKS IF USER TOKEN IS NOT EXPIRED
// HI, {{USERNAME}}|




// Login, add token to session
func (app *application) RenderSuccess(w http.ResponseWriter, r *http.Request) {
	var user models.User = app.Session.Get(r.Context(), "user-data").(models.User)
	data := make(map[string]interface{})
	data["userData"] = user
	app.Session.Remove(r.Context(), "user-data")
	if err := app.renderTemplate(w, r, "succeededRegistration", &templateData{
		Data: data,
	}); err != nil {
		app.errorLog.Println(err)
	}
}