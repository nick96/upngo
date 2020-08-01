package upngo

import "errors"

var ErrNotImplemented = errors.New("not implemented")

// Account represents an UpBank account.
type Account struct{}

// Transaction represents an UpBank transaction.
type Transaction struct{}

type Client struct {
	token string
}

func NewClient(token string) *Client {
	return &Client{token: token}
}

// Ping pings the UpBank API and returns an error if there is an problem.
func Ping() error {
	return ErrNotImplemented
}

// Accounts lists all the accounts associated with the authenticated account.
func (c *Client) Accounts() ([]Account, error) {
	return []Account{}, ErrNotImplemented
}

// Transactions lists all the transactions associated with the authenticated
// account.
func (c *Client) Transactions() ([]Transaction, error) {
	return []Transaction{}, ErrNotImplemented
}
