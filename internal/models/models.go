package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type DBModel struct {
	Db *sql.DB
}

type Models struct {
	Database DBModel
}

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
	CreatedAt 	 	   time.Time		`json:"-"`
	UpdatedAt 	 	   time.Time 		`json:"-"`
}

// Order is the type for all orders
type Order struct {
	ID            int       `json:"id"`
	ProductID      int       `json:"widget_id"`
	TransactionID int       `json:"transaction_id"`
	CustomerID    int       `json:"customer_id"`
	StatusID      int       `json:"status_id"`
	Quantity      int       `json:"quantity"`
	Amount        int       `json:"amount"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

// Status is the type for order statuses
type Status struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// TransactionStatus is the type for transaction statuses
type TransactionStatus struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// Transaction is the type for transactions
type Transaction struct {
	ID                  int       `json:"id"`
	Amount              int       `json:"amount"`
	Currency            string    `json:"currency"`
	LastFour            string    `json:"last_four"`
	ExpiryMonth         int       `json:"expiry_month"`
	ExpiryYear          int       `json:"expiry_year"`
	PaymentIntent       string    `json:"payment_intent"`
	PaymentMethod       string    `json:"payment_method"`
	BankReturnCode      string    `json:"bank_return_code"`
	TransactionStatusID int       `json:"transaction_status_id"`
	CreatedAt           time.Time `json:"-"`
	UpdatedAt           time.Time `json:"-"`
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
	Currency  		Currency		`json:"currency"`
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


// GetUserByEmail gets a user by email address
func (m *DBModel) GetUserByEmail(email string) (User, error) {
	_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	// email = strings.ToLower(email)
	var u User

	email=strings.ToLower(email)

	row := m.Db.QueryRow(`
		SELECT
			(id, firstname, lastname, username, password, hashedpassword, email, createdat, updatedat, address1, address2)
		FROM
			users
		WHERE email = '$1'`, email)

	err := row.Scan(
		&u.ID,
		&u.Firstname,
		&u.Lastname,
		&u.Username,
		&u.Password,
		&u.UserHashedPassword,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.Address1,
		&u.Address2,
	)
	fmt.Println(err)

	if err != nil {
		return u, err
	}

	return u, nil
}


func (m *DBModel) InsertTX(tx Transaction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	stmt := `
		INSERT INTO
			transactions (amount, currency, last_four, bank_return_code, expiry_month, expiry_year,
			payment_intent, payment_method, 
			transaction_status_id, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	row, err := m.Db.ExecContext(ctx, stmt, 
		tx.Amount,
		tx.Currency,
		tx.LastFour,
		tx.BankReturnCode,
		tx.ExpiryMonth,
		tx.ExpiryYear,
		tx.PaymentIntent,
		tx.PaymentMethod,
		tx.TransactionStatusID,
		tx.CreatedAt,
		tx.UpdatedAt,
	)


	if err != nil {
		fmt.Println("Error inserting tx")
		return 0, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		fmt.Println("Error getting last insert id")
		return 0, err
	}
	return int(id), nil
}