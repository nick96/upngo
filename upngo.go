package upngo

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
)

type Client struct {
	token   string
	baseURL string
	client  *http.Client
}

func (c *Client) buildURL(endpoint string) string {
	return fmt.Sprintf("%s/api/v1/%s", c.baseURL, endpoint)
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
		name:  "page[size]",
		value: strconv.Itoa(size),
	}
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

type transactionsOption struct {
	name  string
	value string
}

// TODO: This is a horrible name but otherwise it clashes with the accounts page
// size option. Anyway, I'm thinking either have to test out whether we can have
// a base option then alias it for the specific ones but I'm not sure if that
// will still be typesafe. Alternatively, it might be good, in general, to split
// out each resoure into a subclient, then they could have different namespaces.
func WithTransactionPageSize(size int) transactionsOption {
	return transactionsOption{
		name:  "page[size]",
		value: strconv.Itoa(size),
	}
}

func WithFilterSince(since time.Time) transactionsOption {
	return transactionsOption{
		name:  "filter[since]",
		value: since.Format(time.RFC3339),
	}
}

func WithFilterUntil(until time.Time) transactionsOption {
	return transactionsOption{
		name:  "filter[until]",
		value: until.Format(time.RFC3339),
	}
}

// Transactions lists all the transactions associated with the authenticated
// account.
func (c *Client) Transactions(options ...transactionsOption) (TransactionsResponse, error) {
	url := c.buildURL("transactions")
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return TransactionsResponse{}, fmt.Errorf("failed to create transactions request: %w", err)
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
		return TransactionsResponse{}, fmt.Errorf("failed to get transactions: %w", err)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return TransactionsResponse{}, fmt.Errorf("failed to read transactions response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		if err := unmarshal(responseBody, &errorResponse); err != nil {
			return TransactionsResponse{}, fmt.Errorf("failed to unmarshall get transactions error: %w", err)
		}

		var err error
		for _, responseError := range errorResponse.Errors {
			err = multierror.Append(
				err,
				errors.New(responseError.Detail),
			)
		}
		return TransactionsResponse{}, err
	}

	var transactionsResponse TransactionsResponse
	if err := unmarshal(responseBody, &transactionsResponse); err != nil {
		// fmt.Println(string(responseBody))
		return TransactionsResponse{}, fmt.Errorf("failed to unmarshal get transactions response: %w", err)
	}

	return transactionsResponse, nil
}
