package upngo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

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

type PingErrorResponse struct {
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

func NewClient(token string) *Client {
	return &Client{
		token:   token,
		baseURL: "https://api.up.com.au",
		client:  http.DefaultClient,
	}
}

// Ping pings the UpBank API and returns an error if there is an problem.
func (c *Client) Ping() error {
	url := fmt.Sprintf("%s/api/v1/util/ping", c.baseURL)
	log.Printf("Sending request to %s", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

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
		var pingErrorResponse PingErrorResponse
		if err := json.Unmarshal(respBody, &pingErrorResponse); err != nil {
			return fmt.Errorf("failed to unmarshal ping error response: %w", err)
		}
		// It doesn't say it but it looks like form the docs (and makes sense)
		// that the ping `errors` field will always be of length 1. If this
		// turns out to be wrong (by blowing up in my face :D) then we can just
		// add a bit more detail here.
		return fmt.Errorf("ping failed: %s", pingErrorResponse.Errors[0].Detail)
	}

	var pingResponse PingResponse
	if err := json.Unmarshal(respBody, &pingResponse); err != nil {
		return fmt.Errorf("failed to unmarshal ping response: %w", err)
	}

	// Everything is okay so we just return nil.
	return nil
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
