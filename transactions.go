package upngo

import "time"

type TransactionStatus string

const (
	TransactionStatusHeld    TransactionStatus = "HELD"
	TransactionStatusSettled TransactionStatus = "SETTELED"
)

type RelatedLinksObject struct {
	Related string `json:"related"`
}

type DataObject struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type AccountObject struct {
	Links RelatedLinksObject `json:"links"`
	Data  DataObject         `json:"data"`
}

type TransactionRelationshipsObject struct {
	Account AccountObject `json:"account"`
}

type Resource struct {
	ID    string         `json:"id"`
	Type  string         `json:"type"`
	Links SelfLinkObject `json:"links"`
}

type HoldInfoObject struct {
	Amount        MoneyObject `json:"amount"`
	ForeignAmount MoneyObject `json:"foreignAmount"`
}

type RoundUpObject struct {
	Amount       MoneyObject `json:"amount"`
	BoostPortion MoneyObject `json:"boostPortion"`
}

type CashbackObject struct {
	Description string      `json:"description"`
	Amount      MoneyObject `json:"amount"`
}

type TransactionAttributes struct {
	Description   string            `json:"description"`
	Status        TransactionStatus `json:"status"`
	RawText       string            `json:"rawText"`
	Message       string            `json:"message"`
	HoldInfo      HoldInfoObject    `json:"holdInfo"`
	RoundUp       RoundUpObject     `json:"roundUp"`
	Cashback      CashbackObject    `json:"cashback"`
	Amount        MoneyObject       `json:"amount"`
	ForeignAmount MoneyObject       `json:"foreignAmount"`
	SettledAt     time.Time         `json:"settledAt"`
	CreatedAt     time.Time         `json:"createdAt"`
}

type TransactionResource struct {
	Resource
	Attributes    TransactionAttributes          `json:"attributes"`
	Relationships TransactionRelationshipsObject `json:"relationships"`
}

type TransactionsResponse struct {
	Data  []TransactionResource `json:"data"`
	Links LinksObject           `json:"links"`
}
