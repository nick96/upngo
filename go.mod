module github.com/nick96/upngo

go 1.14

require (
	github.com/99designs/keyring v1.1.5
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.3.0
	golang.org/x/text v0.3.2
)

// From:
// https://github.com/99designs/keyring/commit/db030e0c0975ed003e02a6049e8be23aec8daf9c.
// The version of go-keychain specified in keyring's go.mod is using some
// deprecated MacOS functions. Newer versions of go-keychain fix this so we
// force the use of one of them.
replace github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
