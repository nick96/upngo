package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/nick96/upngo/keyring"
)

func abort(msg string, args ...interface{}) {
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(1)
}

func getToken() string {
	// The token can either be in an env var or a keyring. If someone has
	// set the env var they probably intend to do that over the keychain
	// because it's useful for testing and stuff.
	var token string
	if value, ok := os.LookupEnv("UPBANK_TOKEN"); ok {
		token = value
	} else {
		// Not sure if it's just me being dumb but it looks like we can't
		// assign `token` (an existing variable) and define `err` in a
		// single statement. So, instead, we just assign a temporary
		// `keyringToken` and set token to that value once we know there
		// were no issues getting it out of the keyring.
		keyringToken, err := keyring.GetTokenDefaultConfig()
		if err != nil {
			abort("Failed to get UpBank token from keyring: %v", err)
		}

		token = keyringToken
	}

	// We won't be able to use the API without a token so just blow up here.
	if token == "" {
		abort("Failed to find UpBank token in UPBANK_TOKEN environment variable and keyring")
	}

	return token
}
