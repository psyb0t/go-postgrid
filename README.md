# go-postgrid

[![Go Reference](https://pkg.go.dev/badge/github.com/psyb0t/go-postgrid.svg)](https://pkg.go.dev/github.com/psyb0t/go-postgrid)
[![CI](https://github.com/psyb0t/go-postgrid/actions/workflows/pipeline.yml/badge.svg?branch=main)](https://github.com/psyb0t/go-postgrid/actions/workflows/pipeline.yml)
[![version](https://raw.githubusercontent.com/psyb0t/go-postgrid/badges/version.svg)](https://github.com/psyb0t/go-postgrid/tags)
[![license](https://raw.githubusercontent.com/psyb0t/go-postgrid/badges/license.svg)](LICENSE)
[![imported by](https://raw.githubusercontent.com/psyb0t/go-postgrid/badges/importers.svg)](https://github.com/psyb0t/go-postgrid/blob/badges/importers.md)

Unofficial Go client for the [PostGrid](https://www.postgrid.com) print and mail API, generated with [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) from the OpenAPI document PostGrid publishes. I am not affiliated with PostGrid.

The package is `postgrid`, in the repository root. Everything in `client.gen.go` is generated. Don't edit it by hand, regenerate it.

## Install

```sh
go get github.com/psyb0t/go-postgrid
```

## Usage

Every endpoint in the spec has a typed method on `ClientWithResponses`. Responses decode into `JSON200` and the other status fields.

```go
package main

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-postgrid"
)

func main() {
	auth := postgrid.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("X-API-Key", os.Getenv("POSTGRID_API_KEY"))

			return nil
		},
	)

	client, err := postgrid.NewClientWithResponses(
		"https://api.postgrid.com", auth,
	)
	if err != nil {
		panic(err)
	}

	res, err := client.ContactsListWithResponse(
		context.Background(),
		&postgrid.ContactsListParams{},
	)
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
```

The same code is in `example_test.go`, so `make test` checks that it compiles.

The full API is in the [Go reference](https://pkg.go.dev/github.com/psyb0t/go-postgrid); from inside your module, `go doc github.com/psyb0t/go-postgrid` prints the same.

The client covers all operations in the published document, including contacts, templates, letters, postcards, cheques, self mailers, campaigns, mailing lists, reports, webhooks and address verification. Method names come from the spec's operation IDs.

## Auth

PostGrid uses an API key in the `X-API-Key` header. Set it with `WithRequestEditorFn`, as in the example above, and read the key from `POSTGRID_API_KEY`.

The spec defines separate `PrintMailApiKeyAuth` and `AddressVerificationApiKeyAuth` security schemes. Both use the same header, but select the key for the API you call. The base URL is `https://api.postgrid.com`. Print and mail test or live behavior depends on the key. Use a test key while developing; live mail creation can send physical mail and incur charges.

## Regenerating

The spec is not committed. `make generate` downloads the current OpenAPI document into the gitignored `api/` directory and regenerates the client:

```sh
make generate
```

`gen.go` declares the steps, so `make generate` is just `go generate ./...`:

- `scripts/fetch-spec.sh` downloads [PostGrid's Node SDK metadata](https://raw.githubusercontent.com/postgrid/postgrid-node/main/.stats.yml), reads `openapi_spec_url` and downloads that document to `api/openapi.yaml`. Stainless hosts it at a URL containing a content hash, so the URL changes as the spec changes.
- `scripts/patch-spec.sh` runs the YAML tool in `internal/cmd/patch-spec` through `go tool`, which writes `api/openapi.patched.yaml` from the OpenAPI 3.0.0 document and collapses repeated identical references in four mail creation response unions. Without this patch, oapi-codegen emits duplicate Go methods. The wire schema and all operations are preserved.
- The patch also moves the shared search parameter's inline string branch into a named schema, avoiding an undefined `N0` type while preserving text and structured search.
- oapi-codegen runs with `config.yaml` against `api/openapi.patched.yaml`. The config maps `number` to `float64` instead of oapi-codegen's default `float32`.

`.github/workflows/regen.yml` does the same every day in CI. When the new spec changes the client, it runs `make lint`, `make test` and `make sec`. If a check fails, Claude edits the code and the checks run again, up to three times, and a review step reads the result. A passing change is committed with a CHANGELOG entry and tagged as the next release: a patch version for fixes and additions, a minor version for breaking changes.

## Agent integrations

The [skill](.agents/skills/go-postgrid) teaches a coding agent how to use go-postgrid when it writes Go. It works in any agent that reads `.agents/skills/`, and installs natively in the clients below.

### Claude Code

```bash
claude plugin marketplace add psyb0t/agents
claude plugin install go-postgrid@psyb0t
```

### Codex

```bash
codex plugin marketplace add psyb0t/agents
codex plugin add go-postgrid@psyb0t
```

Codex also picks the skill up with no install in any checkout of this repo.

### OpenClaw

The skill is published to ClawHub on every release:

```bash
openclaw skills install @psyb0t/go-postgrid
```

## Development

Every command runs in the `Dockerfile.dev` image. The host needs only Docker and Make.

```sh
make help      # list targets
make dep       # go mod tidy + vendor
make lint      # shellcheck, shfmt, go fix, golangci-lint
make test      # go test -race
make sec       # gitleaks, govulncheck, semgrep
make generate  # download the spec and regenerate
```

## License

[MIT](LICENSE). The generated code comes from PostGrid's published API specification. The specification itself is not redistributed here. The downloaded document has no license or terms-of-service declaration; this repository's MIT license does not grant rights to PostGrid's service or specification.
