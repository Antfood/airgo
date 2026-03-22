package airtable

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	. "github.com/pbotsaris/airgo/testutils/testutils"
)

func TestMockClient(t *testing.T) {
	msg := "Not Found"
	inBody := testErrorBody{}
	inBody.Error.Message = msg

	jsonBody, err := json.Marshal(inBody)
	Ok(t, err)

	client := newMockClient(404, jsonBody, errors.New(msg))

	req, err := http.NewRequest(http.MethodPost, "http://example.com", nil)
	Ok(t, err)

	resp, err := client.Do(req)
	Assert(t, resp.StatusCode == 404, "Expected '%d', got '%d'", 404, resp.StatusCode)
	Assert(t, err.Error() == msg, "Expected '%s', got '%s'", msg, err.Error())

	defer resp.Body.Close()

	outBody := testErrorBody{}

	err = json.NewDecoder(resp.Body).Decode(&outBody)
	Ok(t, err)

	Assert(t, outBody.Error.Message == msg, "Expected '%s', got '%s'", msg, outBody.Error.Message)
}
