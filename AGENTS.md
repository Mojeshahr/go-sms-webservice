# Agent guide

Runnable Go examples for the Payam Resan SMS web service. One file per API
method, and every file has to work on its own.

## Rule one: no dependencies, and no go.mod either

`net/http` and `encoding/json` are both in the standard library, so these
examples need nothing at all. There is deliberately no `go.mod`: each file is a
standalone `package main` that runs with

```bash
go run examples/v3/send-bulk.go
```

Several files in one directory each declaring `func main()` is fine, because
`go run <file>` compiles only the file named. Do not add a module or a
per-method subdirectory to "fix" it; both would break the file naming that ties
each example to its documentation page.

The same rule rules out helpers shared between the examples. A function defined
next door means nothing to somebody reading one file on the documentation site.

## Rule two: decode with UseNumber

Go turns every JSON number into a `float64` by default, and the service's ids
are large enough that the default `%v` then prints them in exponential form:
`6.8733095017201e+14` instead of `687330950172013`. A reader who copies that
into their own code gets corrupted ids and a long afternoon.

```go
decoder := json.NewDecoder(answer.Body)
decoder.UseNumber()
```

With `UseNumber` each number arrives as a `json.Number`, which prints exactly
what the service sent.

## Rule three: the examples are the documentation

Each file carries `// docs:start` and `// docs:end`. The region between them is
lifted verbatim into the method's page on docs.payam-resan.com, so it is read by
people who have never seen this repository.

Two consequences:

- **Full-line comments are stripped** when the region is lifted. Anything the
  reader must see has to be code. The `Success` check is an `if`, not a note.
- The file name matches the reference page slug exactly: `send-bulk.go`,
  `status-by-user-trace-id.go`. A path with two variants gets two files, the
  plain name for `POST` and a `-get` suffix for `GET`.

The full contract lives in the `handbook` repository, section `docs-site`, file
`code-samples.md`.

## Rule four: check Success, never the status code

The service answers `200` to everything, including a wrong key and an empty
account, so `answer.StatusCode` proves nothing:

```go
if response["Success"] != true {
    fmt.Fprintf(os.Stderr, "ناموفق. کد %v: %v\n", response["ErrorCode"], response["Error"])
    os.Exit(1)
}
```

Type assertions on the result use the two-value form (`result, _ := ...`) on
purpose. A body that is not the usual envelope then prints a readable message
instead of panicking.

## Rule five: a version is a folder

A new service version means a new `examples/v<n>/`. No file inside an existing
version folder is moved or renamed; older versions still have users.

## Secrets

The key comes from `PAYAM_RESAN_API_KEY` in the environment. No key, no real
phone number and no customer name goes into a file here, not even a dead one.
Example numbers are `9121112222` upward and the example key is
`123456-XXXXXXXXXXXXXXX`.

## Layout

| Path | What it holds |
|---|---|
| `examples/v3/` | one self-contained file per service operation |
| `.env.example` | the environment variables the examples read |

## Before every commit

```bash
gofmt -l examples/v3
for f in examples/v3/*.go; do go run "$f" || echo "FAILED $f"; done
```

Point them at `api/V3SandBox/` first so no real message goes out.

## Git

Semantic messages, `type(scope): subject`, with no explanatory body and no
attribution trailer. Commits here are authored as Payam Resan.
