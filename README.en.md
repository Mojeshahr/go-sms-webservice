<div align="center">

<a href="https://payam-resan.com">
  <img src=".github/assets/logo.svg" width="64" height="64" alt="Payam Resan">
</a>

<h1>Go examples for the Payam Resan SMS web service</h1>

Talk to the <a href="https://payam-resan.com"><b>Payam Resan SMS panel</b></a> from Go<br>
One runnable file per API method, with no dependencies

[![API](https://img.shields.io/badge/API-V3-0a7cbd)](https://payam-resan.com)
[![Go](https://img.shields.io/badge/Go-1.21%2B-00add8)](https://go.dev)
[![Dependencies](https://img.shields.io/badge/dependencies-none-2ea44f)](#quick-start)
[![License](https://img.shields.io/badge/license-MIT-6e7781)](LICENSE)

<a href="README.md">فارسی</a> · <b>English</b>

</div>

<sub>Looking for another language? The same examples exist for the others at
[github.com/Mojeshahr](https://github.com/Mojeshahr).</sub>

---

## Quick start

```bash
git clone https://github.com/Mojeshahr/go-sms-webservice.git
cd go-sms-webservice

export PAYAM_RESAN_API_KEY='123456-XXXXXXXXXXXXXXX'
export PAYAM_RESAN_SENDER='30004040'

go run examples/v3/account-info.go
```

No `go mod` and no packages. `net/http` and `encoding/json` are both in the Go
standard library.

Start with `account-info.go`. It sends nothing, spends no credit, and if it
answers then both the key and the connection are fine.

## Before sending anything real

There is a sandbox server that answers exactly like production but sends no
message and spends no credit. Swap `V3` for `V3SandBox` in the URL. The one
exception is `TokenList`, which the sandbox does not implement.

## The methods

| Example | Method | What it does |
|---|---|---|
| [account-info.go](examples/v3/account-info.go) | `AccountInfo` | Credit and active lines |
| [send.go](examples/v3/send.go) | `Send` | Simple send over `GET` |
| [send-bulk.go](examples/v3/send-bulk.go) | `SendBulk` | One text to many recipients, with tracking ids |
| [send-multiple.go](examples/v3/send-multiple.go) | `SendMultiple` | A separate text per recipient |
| [token-list.go](examples/v3/token-list.go) | `TokenList` | The account's templates |
| [send-token-single.go](examples/v3/send-token-single.go) | `SendTokenSingle` | Send a template to one number |
| [send-token-single-get.go](examples/v3/send-token-single-get.go) | `SendTokenSingle` | The same, over `GET` |
| [send-token-multi.go](examples/v3/send-token-multi.go) | `SendTokenMulti` | One template, many recipients |
| [status-by-id.go](examples/v3/status-by-id.go) | `StatusById` | Status by the service's id |
| [status-by-user-trace-id.go](examples/v3/status-by-user-trace-id.go) | `StatusByUserTraceId` | Status by your own id |
| [get-inbox.go](examples/v3/get-inbox.go) | `GetInbox` | Messages people sent to your lines |

## Why there is no go.mod

Each file is a standalone `package main`, run with `go run <file>`. Several
files in one directory each declaring `func main()` is fine, because `go run`
compiles only the file you name.

If your own project is a module, copy the body of `main` into your own function;
the code works unchanged.

## Things that will save you time

**Do not drop `UseNumber`.** Go decodes every JSON number into a `float64`, and
the service's ids are large enough that they then print in exponential form:
`6.8733095017201e+14` instead of `687330950172013`. With `UseNumber` every
number stays exactly as the service sent it.

**Do not read the HTTP status code.** The service answers `200` to everything,
including a wrong key, so `answer.StatusCode` proves nothing. Decide on the
`Success` field.

**Recipient numbers carry no leading zero.** Use `9121112222`, or
`989121112222` with the country code. A number that does not start with `9` or
`989` returns error code `13`.

**Do not encode the text twice.** In `send.go`, `Encode` on `url.Values` already
does it once. Call `QueryEscape` beforehand and the message arrives full of
`%D8`.

**Send a unique `UserTraceId` per recipient.** After a timeout it is the only
way to learn whether the message was registered.

## Key safety

The key is a secret. It does not belong in a code repository, in browser
JavaScript, or in a mobile app bundle. It belongs in an environment variable,
which is where every example here reads it from.

If a key leaks, issue a new one from the panel. A deleted key never comes back.

## Layout

| Path | What it holds |
|---|---|
| `examples/v3/` | One self-contained example per service operation |
| `.env.example` | The environment variables the examples read |

The `v3` in the path is deliberate. A new service version means a new
`examples/v<n>/`, with the existing folder left alone.

## Documentation and support

The full guide is at [docs.payam-resan.com](https://docs.payam-resan.com). The
machine-readable OpenAPI description is in
[sms-webservice-spec](https://github.com/Mojeshahr/sms-webservice-spec).

## License

MIT. Full text in [`LICENSE`](LICENSE).
