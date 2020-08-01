# upngo

CLI and library for working with [UpBank](https://up.com.au) via its
[API](https://developer.up.com.au/).

## TODO

- [x] Ping
- [ ] List accounts
- [ ] List transactions
- [ ] Get specific account
- [ ] Get specific transaction
- [ ] Get webhooks
- [ ] Create webhooks
- [ ] Get speific webhook
- [ ] Delete webhook
- [ ] Ping webhook
- [ ] List webhook logs
- [ ] Generate completion

## CLI

The CLI lets you do anything you can do via the library or API it just provides
a nice CLI wrapper.

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
