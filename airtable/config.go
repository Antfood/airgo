package airtable

const (
	baseUrl        = "https://api.airtable.com/v0"
	bearer         = "Bearer "
	maxPageSize    = 100
	DateTimeLayout = "2006-01-02T15:04:05.000Z"
	sortField      = "field"
	sortDirection  = "direction"
)

var token string
var client Client = NewAirtableClient()

/*
SetToken sets the Airtable API token used for authentication.

Example:

	airtable.SetToken(os.Getenv("AIRTABLE_TOKEN"))
*/
func SetToken(airtableToken string) {
	token = airtableToken
}

/*
Configure sets both the HTTP client and the Airtable API token.
Use this when you need a custom HTTP client (e.g., for testing or custom timeouts).

Example:

	client := airtable.NewAirtableClient()
	airtable.Configure(client, os.Getenv("AIRTABLE_TOKEN"))
*/
func Configure(c Client, airtableToken string) {
	token = airtableToken
	client = c
}
