package upngo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/currency"
)

func newServerClientForURL(
	t *testing.T,
	token string,
	url string,
	status int,
	expectedResponse interface{},
	checks ...func(*http.Request),
) (*httptest.Server, *Client) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		require.Equal(t, url, req.URL.Path)

		authzHeader := req.Header.Get("Authorization")
		expectedHeader := fmt.Sprintf("Bearer %s", token)
		require.Equal(t, expectedHeader, authzHeader)

		for _, check := range checks {
			check(req)
		}

		rw.WriteHeader(status)

		response, _ := json.Marshal(expectedResponse)
		_, err := rw.Write(response)
		require.NoError(t, err)
	}))

	// NewClient inserts the transport middleware that sets the authorization
	// header for us, so we want that. Then we just smash in the values that we
	// need to the private fields because the test is in the same page.
	client := NewClient(token)
	client.baseURL = server.URL

	// We need to use the test servers client but we want the transport that
	// NewClient set for the client it created, so we hackily put it back in.
	// Looking at the source of httptest, it looks like it's just using the
	// default HTTP transport so we should be okay.
	oldTransport := client.client.Transport
	client.client = server.Client()
	client.client.Transport = oldTransport

	return server, client
}

func TestPingOk(t *testing.T) {
	token := "token"
	expectedResponse := PingResponse{
		Meta: PingResponseMeta{
			ID:          "c0ee698b-6707-4d87-a1b3-80393f1f8571",
			StatusEmoji: "⚡️",
		},
	}
	server, client := newServerClientForURL(t, token, "/api/v1/util/ping", http.StatusOK, expectedResponse)
	defer server.Close()
	require.NoError(t, client.Ping())
}

func TestPingErr(t *testing.T) {
	token := "token"
	detail := "The request was not authenticated because no valid credential was found in the Authorization header, or the Authorization header was not present."
	expectedResponse := ErrorResponse{
		Errors: []ErrorObject{
			{
				Status: "401",
				Title:  "Not Authorized",
				Detail: detail,
			},
		},
	}
	server, client := newServerClientForURL(t, token, "/api/v1/util/ping", http.StatusUnauthorized, expectedResponse)
	defer server.Close()

	err := client.Ping()
	require.Error(t, err)
	expected := fmt.Sprintf("ping failed: %s", detail)
	require.Equal(t, expected, err.Error())
}

func TestAccountsNoPageSize(t *testing.T) {
	expectedResponse := AccountsResponse{
		Data: []AccountResource{
			{
				ID:   "id",
				Type: "accounts",
				Attributes: AttributesObject{
					DisplayName: "test",
					AccountType: AccountTypeSaver,
					Balance: MoneyObject{
						CurrencyCode:     "AUD",
						Value:            "1.00",
						ValueInBaseUnits: 100,
					},
					CreatedAt: time.Date(2020, 8, 2, 15, 20, 22, 100, time.UTC),
				},
				Links: SelfLinkObject{
					Self: "https://blahblahblah",
				},
				Relationships: RelationshipsObject{
					Transactions: TransactionsObject{
						Links: TransactionsLinksObject{
							Related: "https:/blahblahblah",
						},
					},
				},
			},
		},
		Links: LinksObject{
			Prev: "https:/blahblahblah",
			Next: "https:/blahblahblah",
		},
	}

	token := "token"
	server, client := newServerClientForURL(
		t,
		token,
		"/api/v1/accounts",
		http.StatusOK,
		expectedResponse,
		func(req *http.Request) {
			require.Empty(t, req.URL.Query().Get("page[size]"))
		},
	)
	defer server.Close()
	accounts, err := client.Accounts()
	require.NoError(t, err)
	require.Equal(t, expectedResponse, accounts)
}

func TestAccountsPageSize(t *testing.T) {
	expectedResponse := AccountsResponse{
		Data: []AccountResource{
			{
				ID:   "id",
				Type: "accounts",
				Attributes: AttributesObject{
					DisplayName: "test",
					AccountType: AccountTypeSaver,
					Balance: MoneyObject{
						CurrencyCode:     "AUD",
						Value:            "1.00",
						ValueInBaseUnits: 100,
					},
					CreatedAt: time.Date(2020, 8, 2, 15, 20, 22, 100, time.UTC),
				},
				Links: SelfLinkObject{
					Self: "https://blahblahblah",
				},
				Relationships: RelationshipsObject{
					Transactions: TransactionsObject{
						Links: TransactionsLinksObject{
							Related: "https:/blahblahblah",
						},
					},
				},
			},
		},
		Links: LinksObject{
			Prev: "https:/blahblahblah",
			Next: "https:/blahblahblah",
		},
	}

	token := "token"
	server, client := newServerClientForURL(
		t,
		token,
		"/api/v1/accounts",
		http.StatusOK,
		expectedResponse,
		func(req *http.Request) {
			require.Equal(t, "10", req.URL.Query().Get("page[size]"))
		},
	)
	defer server.Close()
	accounts, err := client.Accounts(WithPageSize(10))
	require.NoError(t, err)
	require.Equal(t, expectedResponse, accounts)
}

func TestAccountsError(t *testing.T) {
	detail := "spilling the tea"
	expectedResponse := ErrorResponse{
		Errors: []ErrorObject{
			{
				Status: "status",
				Title:  "title",
				Detail: detail,
			},
		},
	}

	token := "token"
	server, client := newServerClientForURL(
		t,
		token,
		"/api/v1/accounts",
		http.StatusInternalServerError,
		expectedResponse,
	)
	defer server.Close()
	_, err := client.Accounts()
	var expectedErr error
	expectedErr = multierror.Append(expectedErr, errors.New(detail))
	require.Equal(t, expectedErr, err)
}

func TestAccountsMultipleError(t *testing.T) {
	detail1 := "spilling the tea"
	detail2 := "stirring the pot"
	expectedResponse := ErrorResponse{
		Errors: []ErrorObject{
			{
				Status: "status",
				Title:  "title",
				Detail: detail1,
			},
			{
				Status: "status",
				Title:  "title",
				Detail: detail2,
			},
		},
	}

	token := "token"
	server, client := newServerClientForURL(
		t,
		token,
		"/api/v1/accounts",
		http.StatusInternalServerError,
		expectedResponse,
	)
	defer server.Close()
	_, err := client.Accounts()
	var expectedErr error
	expectedErr = multierror.Append(expectedErr, errors.New(detail1))
	expectedErr = multierror.Append(expectedErr, errors.New(detail2))
	require.Equal(t, expectedErr, err)
}

func TestTransactionsNoPageSize(t *testing.T) {
	expectedResponse := TransactionsResponse{
		Data: []TransactionResource{
			{
				Resource: Resource{
					ID:   "id",
					Type: "transactions",
					Links: SelfLinkObject{
						Self: "https:/blahblablah",
					},
				},
				Attributes: TransactionAttributes{
					Description: "description",
					Status:      TransactionStatusHeld,
					RawText:     "raw text",
					Message:     "message",
					HoldInfo: HoldInfoObject{
						Amount: MoneyObject{
							CurrencyCode:     currency.AUD.String(),
							Value:            "1.00",
							ValueInBaseUnits: 100,
						},
						ForeignAmount: MoneyObject{
							CurrencyCode:     currency.CAD.String(),
							Value:            "1.00",
							ValueInBaseUnits: 100,
						},
					},
					RoundUp: RoundUpObject{
						Amount: MoneyObject{
							CurrencyCode:     currency.AUD.String(),
							Value:            "1.00",
							ValueInBaseUnits: 100,
						},
						BoostPortion: MoneyObject{
							CurrencyCode:     currency.AUD.String(),
							Value:            "1.00",
							ValueInBaseUnits: 100,
						},
					},
					Cashback: CashbackObject{
						Description: "description",
						Amount: MoneyObject{
							CurrencyCode:     currency.AUD.String(),
							Value:            "1.00",
							ValueInBaseUnits: 100,
						},
					},
					Amount: MoneyObject{
						CurrencyCode:     currency.AUD.String(),
						Value:            "1.00",
						ValueInBaseUnits: 100,
					},
					ForeignAmount: MoneyObject{
						CurrencyCode:     currency.CAD.String(),
						Value:            "1.00",
						ValueInBaseUnits: 100,
					},
					SettledAt: time.Date(2020, 8, 2, 15, 20, 22, 100, time.UTC),
					CreatedAt: time.Date(2020, 8, 2, 15, 20, 22, 100, time.UTC),
				},
				Relationships: TransactionRelationshipsObject{
					Account: AccountObject{
						Links: RelatedLinksObject{
							Related: "https://blahblahblah",
						},
						Data: DataObject{
							Type: "type",
							ID:   "id",
						},
					},
					Tags: TagObject{
						Links: SelfLinkObject{
							Self: "https://blahblahblah",
						},
					},
				},
			},
		},
		Links: LinksObject{},
	}
	token := "token"
	server, client := newServerClientForURL(
		t,
		token,
		"/api/v1/transactions",
		http.StatusOK,
		expectedResponse,
		func(req *http.Request) {
			require.Empty(t, req.URL.Query().Get("page[size]"))
		},
	)
	defer server.Close()
	transactions, err := client.Transactions()
	require.NoError(t, err)
	require.Equal(t, expectedResponse, transactions)
}

func TestTransactionsPageSize(t *testing.T) {
	expectedResponse := TransactionsResponse{
		Data: []TransactionResource{
			{
				Resource: Resource{
					ID:   "id",
					Type: "transactions",
					Links: SelfLinkObject{
						Self: "https:/blahblablah",
					},
				},
				Attributes: TransactionAttributes{
					Description: "description",
					Status:      TransactionStatusHeld,
					RawText:     "raw text",
					Message:     "message",
					HoldInfo: HoldInfoObject{
						Amount: MoneyObject{
							CurrencyCode:     currency.AUD.String(),
							Value:            "1.00",
							ValueInBaseUnits: 100,
						},
						ForeignAmount: MoneyObject{
							CurrencyCode:     currency.CAD.String(),
							Value:            "1.00",
							ValueInBaseUnits: 100,
						},
					},
					RoundUp: RoundUpObject{
						Amount: MoneyObject{
							CurrencyCode:     currency.AUD.String(),
							Value:            "1.00",
							ValueInBaseUnits: 100,
						},
						BoostPortion: MoneyObject{
							CurrencyCode:     currency.AUD.String(),
							Value:            "1.00",
							ValueInBaseUnits: 100,
						},
					},
					Cashback: CashbackObject{
						Description: "description",
						Amount: MoneyObject{
							CurrencyCode:     currency.AUD.String(),
							Value:            "1.00",
							ValueInBaseUnits: 100,
						},
					},
					Amount: MoneyObject{
						CurrencyCode:     currency.AUD.String(),
						Value:            "1.00",
						ValueInBaseUnits: 100,
					},
					ForeignAmount: MoneyObject{
						CurrencyCode:     currency.CAD.String(),
						Value:            "1.00",
						ValueInBaseUnits: 100,
					},
					SettledAt: time.Date(2020, 8, 2, 15, 20, 22, 100, time.UTC),
					CreatedAt: time.Date(2020, 8, 2, 15, 20, 22, 100, time.UTC),
				},
				Relationships: TransactionRelationshipsObject{
					Account: AccountObject{
						Links: RelatedLinksObject{
							Related: "https://blahblahblah",
						},
						Data: DataObject{
							Type: "type",
							ID:   "id",
						},
					},
					Tags: TagObject{
						Links: SelfLinkObject{
							Self: "https://blahblahblah",
						},
					},
				},
			},
		},
		Links: LinksObject{},
	}
	token := "token"
	server, client := newServerClientForURL(
		t,
		token,
		"/api/v1/transactions",
		http.StatusOK,
		expectedResponse,
		func(req *http.Request) {
			require.Equal(t, req.URL.Query().Get("page[size]"), "10")
		},
	)
	defer server.Close()
	transactions, err := client.Transactions(WithTransactionPageSize(10))
	require.NoError(t, err)
	require.Equal(t, expectedResponse, transactions)
}

func TestTransactionsError(t *testing.T) {
	detail := "spilling the tea"
	expectedResponse := ErrorResponse{
		Errors: []ErrorObject{
			{
				Status: "status",
				Title:  "title",
				Detail: detail,
			},
		},
	}

	token := "token"
	server, client := newServerClientForURL(
		t,
		token,
		"/api/v1/transactions",
		http.StatusInternalServerError,
		expectedResponse,
	)
	defer server.Close()
	_, err := client.Transactions()
	var expectedErr error
	expectedErr = multierror.Append(expectedErr, errors.New(detail))
	require.Equal(t, expectedErr, err)
}

func TestTransactionsMultipleError(t *testing.T) {
	detail1 := "spilling the tea"
	detail2 := "stirring the pot"
	expectedResponse := ErrorResponse{
		Errors: []ErrorObject{
			{
				Status: "status",
				Title:  "title",
				Detail: detail1,
			},
			{
				Status: "status",
				Title:  "title",
				Detail: detail2,
			},
		},
	}

	token := "token"
	server, client := newServerClientForURL(
		t,
		token,
		"/api/v1/transactions",
		http.StatusInternalServerError,
		expectedResponse,
	)
	defer server.Close()
	_, err := client.Transactions()
	var expectedErr error
	expectedErr = multierror.Append(expectedErr, errors.New(detail1))
	expectedErr = multierror.Append(expectedErr, errors.New(detail2))
	require.Equal(t, expectedErr, err)
}
