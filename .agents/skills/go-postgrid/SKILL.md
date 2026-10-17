---
name: go-postgrid
description: "Use the unofficial PostGrid Go client to manage print and mail resources or verify addresses. Applies when writing Go code that imports github.com/psyb0t/go-postgrid or calls PostGrid through this module."
homepage: https://github.com/psyb0t/go-postgrid
user-invocable: true
metadata:
  openclaw:
    emoji: "📬"
    requires:
      bins:
        - go
      env:
        - POSTGRID_API_KEY
permissions:
  filesystem: "Guides edits to Go source in the current project. The library reads and writes no files."
  shell: "The Go toolchain (`go get`, `go build`, `go test`) to add and compile the dependency."
  network: "`go get` fetches the module from the Go module proxy. At runtime the written code calls https://api.postgrid.com."
---

# go-postgrid

Unofficial, typed Go client for PostGrid print and mail and address verification, in package `postgrid` at the module root.

Read [references/setup.md](references/setup.md) for installation, operation groups, pagination and response handling. Inspect the generated method signatures before writing a call.

## Security & safety

- Set `X-API-Key` through `WithRequestEditorFn`. Read the key from `POSTGRID_API_KEY`; never hardcode, commit or log it.
- Use `https://api.postgrid.com`. The spec has separate print and mail and address verification security schemes. Select the key for the API being called.
- Print and mail test or live behavior depends on the key. Use a test key while developing. Creating mail with a live key can send physical mail and incur charges.
- Create, send, cancel or delete resources only when the user requested that action. Do not retry mail creation after an uncertain response without checking whether the first request succeeded.

## When to use

- Go code needs PostGrid contacts, templates, letters, postcards, cheques, self mailers, campaigns, mailing lists, reports or webhooks.
- Go code needs standard, international, batch or bulk address verification.
- The project already imports `github.com/psyb0t/go-postgrid`.

## When NOT to use

- The code is not Go. This module ships only a Go package.
- A PostGrid-maintained client is required. This client is unofficial.
- The task changes the generated client itself. Read `AGENTS.md` in the repository.

## Usage

```go
import (
    "context"
    "net/http"
    "os"

    "github.com/psyb0t/go-postgrid"
)

auth := postgrid.WithRequestEditorFn(
    func(_ context.Context, req *http.Request) error {
        req.Header.Set("X-API-Key", os.Getenv("POSTGRID_API_KEY"))

        return nil
    },
)

client, err := postgrid.NewClientWithResponses(
    "https://api.postgrid.com",
    auth,
)
if err != nil {
    return err
}

res, err := client.ContactsListWithResponse(
    context.Background(),
    &postgrid.ContactsListParams{},
)
if err != nil {
    return err
}

_ = res.JSON200
```

Check `StatusCode()` and the status-specific response fields before consuming a response. The Go error covers transport and decoding failures; an HTTP error can return a nil Go error. Optional model fields may be pointers, so check them before dereferencing.

Finish by compiling the calling code and testing it against a local HTTP server. Tests must not create live mail or call PostGrid.
