# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.1, 2026-10-17

- Documentation only. The agent skill explains how to look up every endpoint, parameter and model with `go doc` or on pkg.go.dev, and declares its license.
- The skill no longer lists the API credentials under OpenClaw's `requires.env`, which kept OpenClaw from loading it until they were set. It requires only `go`; the credentials are still documented as what your program reads.
- The README points to the full API reference, and the examples in the README and the skill were checked against the current code.

## v0.1.0, 2026-10-17

- Added `postgrid`, a REST client generated with oapi-codegen from PostGrid's OpenAPI document (PostGrid Address Verification API, spec version 1.0, OpenAPI 3.0.0). It includes every print and mail and address verification operation in the document.
- Unformatted number fields are `float64` instead of oapi-codegen's default `float32`.
- The spec patch also names the shared search parameter's inline string branch, avoiding undefined `N0` references while retaining text and structured search.
- `make generate` reads `openapi_spec_url` from `https://raw.githubusercontent.com/postgrid/postgrid-node/main/.stats.yml`, downloads the Stainless-hosted spec, collapses duplicate identical response unions and regenerates the client. The initial source is `https://storage.googleapis.com/stainless-sdk-openapi-specs/postgrid/postgrid-a1e82d54f9c6456622ef353b6864dd9ccc47c9b5cf6b5fbe98c123f2bdf30f34.yml`. The spec is not committed.
