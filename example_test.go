package postgrid_test

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-postgrid"
)

// This is the README usage example. It needs an API key to run.
//
//nolint:testableexamples // needs real credentials, compile check only
func Example() {
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
