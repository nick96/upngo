module github.com/nick96/upngo

go 1.14

require (
	github.com/99designs/keyring v1.1.5
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/hashicorp/go-multierror v1.0.0
	github.com/kr/text v0.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 // indirect
	golang.org/x/sys v0.0.0-20200519105757-fe76b779f299 // indirect
	golang.org/x/text v0.3.2
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

// From:
// https://github.com/99designs/keyring/commit/db030e0c0975ed003e02a6049e8be23aec8daf9c.
// The version of go-keychain specified in keyring's go.mod is using some
// deprecated MacOS functions. Newer versions of go-keychain fix this so we
// force the use of one of them.
replace github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
