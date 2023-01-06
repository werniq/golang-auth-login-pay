package models

import (
	"context"
	"database/sql"
	"time"
)

type DBModel struct {
	Db *sql.DB
}

type Models struct {
	Database DBModel
}
/*
	username := r.Form.Get("reg-username")
	email := r.Form.Get("reg-email")
	password := r.Form.Get("reg-password")
	address1 := r.Form.Get("reg-address1")
	address2 := r.Form.Get("reg-address1")
	city := r.Form.Get("reg-city")
*/
// User type is used for storing user's data
type User struct {
	ID 				   int 				`json:"id"`
	Firstname 		   string 			`json:"firstname"`
	Lastname 		   string 			`json:"lastname"`
	Username 		   string 			`json:"username"`
	Password		   string 			`json:"password"`
	UserHashedPassword []byte   		`json:"hashed_password"`
	Address1 	 	   string 			`json:"address1"`
	Address2 	 	   string 			`json:"address2"`
	Email 		 	   string 			`json:"email"`
	DateOfBirth  	   time.Time 		`json:"date_of_birth"`
	CreatedAt 	 	   time.Time		`json:"-"`
	UpdatedAt 	 	   time.Time 		`json:"-"`
}

// Tx represents normal transaction, using debit card
type Tx struct {
	ID 		  		int 			`json:"tx_id"`
	Amount	  		int				`json:"amount"`
	LastFour		int 			`json:"last_four"`
	BankReturnCode  string 			`json:"bankReturnCode"`
	Message 		string 			`json:"message"`
	Recipient 		User 			`json:"recipient"`
	Sender    		User			`json:"tx_organizer"`
	Currency  		Currency		`json:"cryptocurrency"`
}

// CryptoTx represents crypto transaction. However, it is stored in blockchain, I decided to also store in database.
type CryptoTx struct {
	ID 		  		int 			`json:"tx_id"`
	Amount	  		int				`json:"amount"`
	Message 		string 			`json:"message"`
	Recipient 		User 			`json:"recipient"`
	Sender    		User			`json:"tx_organizer"`
	Cryptocurrency  Cryptocurrency	`json:"cryptocurrency"`
}

// Cryptocurrency such as: ETH, BTC, DOGE
type Cryptocurrency struct {
	Id 	 	 int 	`json:"cryptocurrency_id"`
	UsdValue int 	`json:"usd_value"`
	Name 	 string `json:"cryptocurrency"`
}

// Currency for instance:  USD, UAH, CAD
type Currency struct {
	Id 	 	 int 	`json:"currency_id"`
	Name 	 string `json:"currency"`
	UsdValue int 	`json:"usd_value"`
}

// Product is something for sale
type Product struct {
	ID 	 		int     `json:"id"`
	Name 		string  `json:"name"`
	Description string  `json:"description"`
	Value 		int     `json:"price"`
	Image 		string 	`json:"image"`
	Seller 		User 	`json:"owner"`
}

// SaveUser func saves user to the database, and returns userId, and error, if any
func (m *DBModel) SaveUser(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	statement := `
		INSERT INTO users
			(username, firstname, lastname, password, hashedpassword, address1, address2, email, createdat, updatedat)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`

	_, err := m.Db.ExecContext(ctx, statement, 
		user.Firstname,
		user.Lastname,
		user.Username,
		user.Password,
		user.UserHashedPassword,
		user.Address1,
		user.Address2,
		user.Email,
		// user.DateOfBirth,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return err
	}

	// id, err := res.LastInsertId()
	// if err != nil {
	// 	return 0, err
	// }

	return nil
}