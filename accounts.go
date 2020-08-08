package upngo

import (
	"fmt"
	"strconv"
	"time"

	"golang.org/x/text/currency"
)

// AccountType represents an enum of the possible account type. Currently that
// is only saver and transactional but I guess there could potentially be more
// supported in the future.
type AccountType string

const (
	AccountTypeSaver         AccountType = "SAVER"
	AccountTypeTransactional AccountType = "TRANSACTIONAL"
)

type MoneyObject struct {
	CurrencyCode     string `json:"currencyCode"`
	Value            string `json:"value"`
	ValueInBaseUnits int64  `json:"valueInBaseUnits"`
}

func (m MoneyObject) Format() string {
	unit := currency.MustParseISO(m.CurrencyCode)
	balance, _ := strconv.ParseFloat(m.Value, 64)
	amount := unit.Amount(balance)
	fmtdAmount := currency.NarrowSymbol(amount)
	return fmt.Sprint(fmtdAmount)
}

type AttributesObject struct {
	DisplayName string      `json:"displayName"`
	AccountType AccountType `json:"accountType"`
	Balance     MoneyObject `json:"balance"`
	CreatedAt   time.Time   `json:"createdAt"`
}

type SelfLinkObject struct {
	Self string `json:"self"`
}

type TransactionsLinksObject struct {
	Related string `json:"related"`
}

type TransactionsObject struct {
	Links TransactionsLinksObject `json:"links"`
}

type RelationshipsObject struct {
	Transactions TransactionsObject `json:"transactions"`
}

type AccountResource struct {
	ID            string              `json:"id"`
	Type          string              `json:"type"`
	Attributes    AttributesObject    `json:"attributes"`
	Links         SelfLinkObject      `json:"links"`
	Relationships RelationshipsObject `json:"relationships"`
}

type LinksObject struct {
	Prev string `json:"prev"`
	Next string `json:"next"`
}

type AccountsResponse struct {
	Data  []AccountResource `json:"data"`
	Links LinksObject       `json:"links"`
}

// AccountResponse represents a response from the account endpoint.
type AccountResponse struct {
	Data AccountResource `json:"data"`
}
