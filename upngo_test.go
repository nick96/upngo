package upngo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPingOk(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		require.Equal(t, "/api/v1/utils/ping", req.URL.Path)
		response, _ := json.Marshal(PingResponse{
			Meta: PingResponseMeta{
				ID:          "c0ee698b-6707-4d87-a1b3-80393f1f8571",
				StatusEmoji: "⚡️",
			},
		})
		_, err := rw.Write(response)
		require.NoError(t, err)
	}))
	defer server.Close()

	client := Client{token: "", baseURL: server.URL, client: server.Client()}
	require.NoError(t, client.Ping())
}

func TestPingErr(t *testing.T) {
	detail := "The request was not authenticated because no valid credential was found in the Authorization header, or the Authorization header was not present."
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		require.Equal(t, "/api/v1/utils/ping", req.URL.Path)
		rw.WriteHeader(http.StatusUnauthorized)
		response, _ := json.Marshal(PingErrorResponse{
			Errors: []ErrorObject{
				{
					Status: "401",
					Title:  "Not Authorized",
					Detail: detail,
				},
			},
		})
		_, err := rw.Write(response)
		require.NoError(t, err)
	}))
	defer server.Close()

	client := Client{token: "", baseURL: server.URL, client: server.Client()}
	err := client.Ping()
	require.Error(t, err)
	expected := fmt.Sprintf("ping failed: %s", detail)
	require.Equal(t, expected, err.Error())
}
