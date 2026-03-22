package airtable

import (
	"testing"

	. "github.com/Antfood/airgo/testutils/testutils"
)

func TestConfig(t *testing.T) {

	c := newMockClient(200, []byte{}, nil)
	Configure(c, "mock_token")

	Equals(t, c, client)
	Assert(t, token == "mock_token", "token should be set")

}
