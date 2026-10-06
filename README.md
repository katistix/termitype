# termitype

Monkeytype in your terminal.

## Prerequisites

- Go installed ([go.dev/dl](https://go.dev/dl))

## Install

```sh
go install github.com/katistix/termitype/cmd/termitype@latest
```

This installs the `termitype` binary (make sure `$(go env GOPATH)/bin` is on your `PATH`).

## Use

```sh
termitype           # 15 second test
termitype -t 30     # 30 second test
termitype -w 50     # 50 word test
termitype --help    # all options
```

Keys: `esc` restart, `ctrl+c` quit, `ctrl+w` delete word.
