package keyring

import (
	"fmt"

	"github.com/99designs/keyring"
)

const (
	upbankTokenKey = "upbank-token"
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
		return "", fmt.Errorf("failed to open keyring with default config: %w", err)
	}

	tokenItem, err := kr.Get(upbankTokenKey)
	if err != nil {
		return "", fmt.Errorf("failed to get token from keyring with default config: %w", err)
	}

	return string(tokenItem.Data), nil
}

func SetTokenDefaultconfig(token string) error {
	kr, err := keyring.Open(Config)
	if err != nil {
		return fmt.Errorf("failed to open keyring with default config: %w", err)
	}

	err = kr.Set(keyring.Item{
		Key:  upbankTokenKey,
		Data: []byte(token),
	})
	if err != nil {
		return fmt.Errorf("failed to set key '%s' in keyring to UpBank token: %w", upbankTokenKey, err)
	}

	return nil
}
