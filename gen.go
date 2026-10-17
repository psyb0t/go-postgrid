// Package postgrid is a Go client for the PostGrid Public API,
// generated with oapi-codegen from the OpenAPI document PostGrid publishes.
package postgrid

//go:generate bash scripts/fetch-spec.sh
//go:generate bash scripts/patch-spec.sh
//go:generate go tool oapi-codegen --config=config.yaml api/openapi.patched.yaml
