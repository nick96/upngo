package upngo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
)

// unmarshal is a custom JSON unmarshaller that differs from the standard
// library's JSON.Unmarshal in that it does not allow unknown fields.
func unmarshal(data []byte, v interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(v); err != nil {
		return err
	}
	return nil
}

var ErrNotImplemented = errors.New("not implemented")

type PingResponseMeta struct {
	ID          string `json:"id"`
	StatusEmoji string `json:"statusEmoji"`
}

type PingResponse struct {
	Meta PingResponseMeta `json:"meta"`
}

type ErrorObject struct {
	Status string `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

type SourceObject struct {
	Parameter string `json:"parameter"`
	Pointer   string `json:"pointer"`
}

type ErrorResponse struct {
	Errors []ErrorObject `json:"errors"`
	Source SourceObject  `json:"source"`
}

// Account represents an UpBank account.
type Account struct{}

// Transaction represents an UpBank transaction.
type Transaction struct{}

type Client struct {
	token   string
	baseURL string
	client  *http.Client
}

type addAuthorizationHeaderTransport struct {
	roundTripper http.RoundTripper
	token        string
}

func (t *addAuthorizationHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.token))
	return t.roundTripper.RoundTrip(req)
}

func newAddAuthorizationHeaderTransport(roundTripper http.RoundTripper, token string) *addAuthorizationHeaderTransport {
	if roundTripper == nil {
		roundTripper = http.DefaultTransport
	}
	return &addAuthorizationHeaderTransport{roundTripper, token}
}

func NewClient(token string) *Client {
	httpClient := http.DefaultClient
	httpClient.Transport = newAddAuthorizationHeaderTransport(httpClient.Transport, token)

	return &Client{
		token:   token,
		baseURL: "https://api.up.com.au",
		client:  http.DefaultClient,
	}
}

// Ping pings the UpBank API and returns an error if there is an problem.
func (c *Client) Ping() error {
	url := c.buildURL("util/ping")
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	// Read the body out here because in either case (error or happy), we're
	// going to want the body. The only different will be in the structure we
	// parse the JSON into.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read ping body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var pingErrorResponse ErrorResponse
		if err := unmarshal(respBody, &pingErrorResponse); err != nil {
			return fmt.Errorf("failed to unmarshal ping error response: %w", err)
		}
		// It doesn't say it but it looks like form the docs (and makes sense)
		// that the ping `errors` field will always be of length 1. If this
		// turns out to be wrong (by blowing up in my face :D) then we can just
		// add a bit more detail here.
		return fmt.Errorf("ping failed: %s", pingErrorResponse.Errors[0].Detail)
	}

	var pingResponse PingResponse
	if err := unmarshal(respBody, &pingResponse); err != nil {
		return fmt.Errorf("failed to unmarshal ping response: %w", err)
	}

	// Everything is okay so we just return nil.
	return nil
}

// accountsOption represents an option (URL param) for the Accounts endpoint.
//
// This type isn't exposed so we instead expose constructors to build it with
// the given values. Once it gets to the point that this struct is constructed,
// the type of value doesn't really matter because it's just going into the URL
// param so it's a string for convenience. It's the job of the constructor to
// convert the given value to a string.
//
// For constructors, the `name` should be the key in the URL param and `value`
// should be the value, i.e. options will be formatted as `<name>=<value>` in
// the URL.
type accountsOption struct {
	name  string
	value string
}

// WithPageSize specifies that the API should return `size` number of accounts.
func WithPageSize(size int) accountsOption {
	return accountsOption{
		// This is the key in the URL param that dictates the paging size. It's
		// kind of weird, I've never seen URL params putting stuff in square
		// brackets before but there you go.
		name: "page[size]",
		// Using FormatInt is better than Sprintf because Go knows that we're
		// converting an _int_ to a string so the parsing of format strings and
		// stuff isn't required.
		value: strconv.FormatInt(int64(size), 10),
	}
}

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

type AccountsResource struct {
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
	Data  []AccountsResource `json:"data"`
	Links LinksObject        `json:"links"`
}

// Accounts lists all the accounts associated with the authenticated account.
func (c *Client) Accounts(options ...accountsOption) (AccountsResponse, error) {
	url := c.buildURL("accounts")
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return AccountsResponse{}, fmt.Errorf("failed to get accounts request: %w", err)
	}

	query := req.URL.Query()
	for _, option := range options {
		// Using `Add`, not `Set` is important because it meanst that if and
		// option is supplied twice then both values are included in the query.
		query.Add(option.name, option.value)
	}
	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return AccountsResponse{}, fmt.Errorf("failed to get accounts: %w", err)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return AccountsResponse{}, fmt.Errorf("failed to read account response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		if err := unmarshal(responseBody, &errorResponse); err != nil {
			return AccountsResponse{}, fmt.Errorf("failed to unmarshall get accounts error: %w", err)
		}

		var err error
		for _, responseError := range errorResponse.Errors {
			err = multierror.Append(
				err,
				errors.New(responseError.Detail),
			)
		}
		return AccountsResponse{}, err
	}

	var accountsResponse AccountsResponse
	if err := unmarshal(responseBody, &accountsResponse); err != nil {
		return AccountsResponse{}, fmt.Errorf("failed to unmarshal get accounts response: %w", err)
	}

	return accountsResponse, nil
}

// Transactions lists all the transactions associated with the authenticated
// account.
func (c *Client) Transactions() ([]Transaction, error) {
	return []Transaction{}, ErrNotImplemented
}

func (c *Client) buildURL(endpoint string) string {
	return fmt.Sprintf("%s/api/v1/%s", c.baseURL, endpoint)
}
