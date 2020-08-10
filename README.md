# upngo

CLI and library for working with [UpBank](https://up.com.au) via its
[API](https://developer.up.com.au/).

## TODO

- [x] Ping
- [x] List accounts
- [x] List transactions
- [x] Get specific account
- [x] Get specific transaction
- [ ] Get webhooks
- [ ] Create webhooks
- [ ] Get speific webhook
- [ ] Delete webhook
- [ ] Ping webhook
- [ ] List webhook logs
- [ ] Generate completion
- [ ] Move raw API wrapper into `api` subdirectory and create nicer high-level
      wrapper

## CLI

The CLI lets you do anything you can do via the library or API it just provides
a nice CLI wrapper.

``` sh
$ upngo help
Talk to your bank from the CLI!

Usage:
  upngo [command]

Available Commands:
  add         A brief description of your command
  completion  Generate completion script
  get         Get accounts or transactions.
  help        Help about any command
  init        Initialise the UpBank CLI for ease of use by adding the token to your keyring.
  list        List accounts or transactions
  ping        Ping UpBank. Useful to test your token is correct.

Flags:
  -h, --help      help for upngo
  -v, --verbose   Verbose logging

Use "upngo [command] --help" for more information about a command.
```

### Completion

You can generate completions for `upgngoctl` with:

``` sh
upngoctl completion zsh
```

If you use a different shell, replace `zsh` with the shell you use.

This command will output the completion script the console. You can save this to
the directory you source your completions from. See the output of:

``` sh
upngoctl help completion
```

for more details on how to get completion working.

### Webhook management

The other parts of the CLI are cool but I think the really useful bit will be
being able to manage webhooks.

## Library

Currently the library is just a raw wrapper around the API. There is nothing
done to distill the API output into more idiomatic Go. For instance, the
structures directly map to the JSON and are highly nested, which is
unnecessarily difficult to work with because we can wrap them up into some nice
Go types. In particular anything to do with money can probably be wrapped up
into a `currency.Amount`and things like that.
