package keyring

import (
	"fmt"

	"github.com/99designs/keyring"
)

var Config = keyring.Config{
	ServiceName: "upngo",
	// TODO: Not 100% sure what this maps to so will leave commented out for now.
	// KeychainName: "upngo",
}

// GetTokenDefaultConfig gets the upbank token using the default config
// (`Config`).
func GetTokenDefaultConfig() (string, error) {
	kr, err := keyring.Open(Config)
	if err != nil {
		return "", fmt.Errorf("failed to get open keyring with default config: %v", err)
	}

	tokenItem, err := kr.Get("upbank-token")
	if err != nil {
		return "", fmt.Errorf("failed to get token from keyring with default config: %v", err)
	}

	return string(tokenItem.Data), nil
}
