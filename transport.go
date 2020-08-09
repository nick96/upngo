package upngo

import (
	"fmt"
	"log"
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

type logTransport struct {
	rt http.RoundTripper
}

func (t *logTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	log.Printf("--> %s %s", req.Method, req.URL.String())
	resp, err := t.rt.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	log.Printf("<-- %s %s", resp.Status, resp.Request.URL)
	return resp, err
}

func newLogTransport(rt http.RoundTripper) *logTransport {
	if rt == nil {
		rt = http.DefaultTransport
	}
	return &logTransport{rt}
}
