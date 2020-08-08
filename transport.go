package upngo

import (
	"fmt"
	"net/http"
)


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
