package postgrid_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/psyb0t/go-postgrid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAPIKey  = "EXAMPLE-DO-NOT-USE"
	contactID   = "contact_example"
	contactBody = `{
		"id":"contact_example","object":"contact","live":false,
		"addressLine1":"123 Example St","countryCode":"US",
		"addressStatus":"verified"
	}`
)

func newTestClient(
	t *testing.T,
	handler http.HandlerFunc,
) *postgrid.ClientWithResponses {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	auth := postgrid.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("X-API-Key", testAPIKey)

			return nil
		},
	)

	client, err := postgrid.NewClientWithResponses(srv.URL, auth)
	require.NoError(t, err)

	return client
}

func TestContactsGet(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/print-mail/v1/contacts/"+contactID, r.URL.Path)
		assert.Equal(t, testAPIKey, r.Header.Get("X-API-Key"))

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(contactBody))
	})

	res, err := client.ContactsGetWithResponse(t.Context(), contactID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)

	assert.Equal(t, contactID, *res.JSON200.Id)
	assert.Equal(t, "US", res.JSON200.CountryCode)
	assert.Equal(t, "verified", string(*res.JSON200.AddressStatus))
}

func TestContactsPagination(t *testing.T) {
	t.Parallel()

	var search postgrid.ListParamsSearch

	require.NoError(t, search.FromListSearchText("Example & Co"))

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/print-mail/v1/contacts", r.URL.Path)
		assert.Equal(t, testAPIKey, r.Header.Get("X-API-Key"))
		assert.Equal(t, "1", r.URL.Query().Get("limit"))
		assert.Equal(t, "Example & Co", r.URL.Query().Get("search"))

		var skip int

		if r.URL.Query().Get("skip") != "0" {
			assert.Equal(t, "1", r.URL.Query().Get("skip"))

			skip = 1
		}

		w.Header().Set("Content-Type", "application/json")
		assert.NoError(t, json.NewEncoder(w).Encode(map[string]any{
			"object": "list", "totalCount": 2, "skip": skip, "limit": 1,
			"data": []json.RawMessage{json.RawMessage(contactBody)},
		}))
	})

	var count int

	for offset := range int32(2) {
		res, err := client.ContactsListWithResponse(
			t.Context(),
			&postgrid.ContactsListParams{
				Skip: new(offset), Limit: new(int32(1)), Search: &search,
			},
		)
		require.NoError(t, err)
		require.NotNil(t, res.JSON200)
		require.Len(t, res.JSON200.Data, 1)
		assert.Equal(t, offset, res.JSON200.Skip)
		assert.Equal(t, int32(2), res.JSON200.TotalCount)
		assert.Equal(t, contactID, *res.JSON200.Data[0].Id)
		count += len(res.JSON200.Data)
	}

	assert.Equal(t, 2, count)
}

func TestTemplatesCreate(t *testing.T) {
	t.Parallel()

	const html = "<p>Hello {{name}}</p>"

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/print-mail/v1/templates", r.URL.Path)
		assert.Equal(t, testAPIKey, r.Header.Get("X-API-Key"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var body postgrid.TemplateCreate

		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, html, *body.Html)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"id":"template_example","object":"template"
		}`))
	})

	res, err := client.TemplatesCreateWithResponse(
		t.Context(), postgrid.TemplateCreate{Html: new(html)},
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode())
	require.NotNil(t, res.JSON201)
	assert.Equal(t, "template_example", *res.JSON201.Id)
}

func TestContactsUnauthorized(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, testAPIKey, r.Header.Get("X-API-Key"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{
			"object":"error",
			"error":{"type":"invalid_api_key","message":"Invalid API key"}
		}`))
	})

	res, err := client.ContactsGetWithResponse(t.Context(), contactID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode())
	assert.Nil(t, res.JSON200)
	require.NotNil(t, res.JSON401)
}

func TestLettersCreateResponse(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/print-mail/v1/letters", r.URL.Path)
		assert.Equal(t, testAPIKey, r.Header.Get("X-API-Key"))
		assert.Equal(t, "example-idempotency", r.Header.Get("idempotency-key"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"id":"letter_example","object":"letter","live":false,
			"status":"ready","pageCount":1,"color":false
		}`))
	})

	create := client.LettersCreateJSONLettersCreateMultipartWithBodyWithResponse
	res, err := create(
		t.Context(),
		&postgrid.LettersCreateJSONLettersCreateMultipartParams{
			IdempotencyKey: new("example-idempotency"),
		},
		"application/json",
		strings.NewReader(`{
			"to":"contact_recipient","from":"contact_sender",
			"html":"<p>Example</p>"
		}`),
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode())
	require.NotNil(t, res.JSON201)
	assert.Equal(t, "letter_example", res.JSON201.Id)
}

func TestContactsMalformedResponse(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":`))
	})

	_, err := client.ContactsGetWithResponse(t.Context(), contactID)

	var syntaxError *json.SyntaxError

	require.ErrorAs(t, err, &syntaxError)
}
