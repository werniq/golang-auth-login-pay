package cards

import (
	"github.com/stripe/stripe-go/v72"
)

type Card struct {
	Secret   string
	Key      string
	Currency string
}

// Transaction stores data for a given tx
type Transaction struct {
	TxStatusId     int
	Amount         int
	Currency       string
	LastFour       string
	BankReturnCode string
}

// Charge is an another name to CreatePaymentIntent
func (c *Card) Charge(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	return c.CreatePaymentIntent(currency, amount)
}

// CreatePaymentIntent attemts to get a pyment intent object from Stripe 
func (c *Card) CreatePaymentIntent(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	stripe.Key = c.Secret

	params := &stripe.PaymentIntentParams{
		Amount: stripe.Int64(int64(amount)),
		Currency: stripe.String(currency),
	}
}