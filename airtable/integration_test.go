package airtable

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/Antfood/airgo/testutils/testutils"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

// =============================================================================
// Configuration - Update these values for your test base
// =============================================================================

const (
	integrationBaseID  = "appojF8QxBOqb1Jvc"
	integrationTableID = "tblbO7Sj13weO8rdB"
)

// Button represents an Airtable button field (read-only)
type Button struct {
	Label string `json:"label"`
	URL   string `json:"url,omitempty"`
}

// Barcode represents an Airtable barcode field
type Barcode struct {
	Text string `json:"text,omitempty"`
	Type string `json:"type,omitempty"`
}

// IntegrationSchema represents the test table with all Airtable field types.
// This schema matches the integration test table exactly.
type IntegrationSchema struct {
	// Editable fields
	Name         string         `json:"name"`
	Notes        string         `json:"notes,omitempty"`
	Assignee     *Collaborator  `json:"assignee,omitempty"`
	SingleSelect string         `json:"single select,omitempty"`
	Attachments  []Attachment   `json:"attachments,omitempty"`
	Checkbox     bool           `json:"checkbox,omitempty"`
	Date         string         `json:"date,omitempty"`
	MultiSelect  []string       `json:"multi select,omitempty"`
	PhoneNumber  string         `json:"phone number,omitempty"`
	Email        string         `json:"email,omitempty"`
	URL          string         `json:"url,omitempty"`
	Number       float64        `json:"number,omitempty"`
	Currency     float64        `json:"currency,omitempty"`
	Percent      float64        `json:"percent,omitempty"`
	Duration     int            `json:"duration,omitempty"`
	Rating       int            `json:"rating,omitempty"`
	BarcodeField *Barcode       `json:"barcode,omitempty"`
	Link         []string       `json:"link,omitempty"`

	// Read-only / computed fields
	Formula        string        `json:"formula,omitempty" update:"ignore"`
	CreatedBy      *Collaborator `json:"created by,omitempty" update:"ignore"`
	CreatedTime    string        `json:"created time,omitempty" update:"ignore"`
	LastModifiedBy *Collaborator `json:"Last Modified By,omitempty" update:"ignore"`
	LastModified   string        `json:"Last Modified,omitempty" update:"ignore"`
	ButtonField    *Button       `json:"button,omitempty" update:"ignore"`
	Lookup         []any         `json:"lookup,omitempty" update:"ignore"`
	Rollup         any           `json:"rollup,omitempty" update:"ignore"`
	Count          int           `json:"count,omitempty" update:"ignore"`
	Autonumber     int           `json:"autonumber,omitempty" update:"ignore"`
}

// =============================================================================
// Test Setup
// =============================================================================

// setupRecorder creates a VCR recorder for integration tests.
// By default, it replays from existing cassettes.
// Set AIRTABLE_RECORD=1 to record new interactions.
func setupRecorder(t *testing.T, cassetteName string) (*recorder.Recorder, func()) {
	t.Helper()

	cassettePath := filepath.Join("testdata", "fixtures", cassetteName)

	// Determine mode based on environment
	mode := recorder.ModeReplayOnly
	if os.Getenv("AIRTABLE_RECORD") != "" {
		mode = recorder.ModeRecordOnly
	}

	// Skip if in replay mode and cassette doesn't exist
	if mode == recorder.ModeReplayOnly {
		if _, err := os.Stat(cassettePath + ".yaml"); os.IsNotExist(err) {
			t.Skipf("No fixture found: %s.yaml. Set AIRTABLE_RECORD=1 to record.", cassettePath)
		}
	}

	opts := &recorder.Options{
		CassetteName:       cassettePath,
		Mode:               mode,
		SkipRequestLatency: true,
	}

	r, err := recorder.NewWithOptions(opts)
	if err != nil {
		t.Fatalf("Failed to create recorder: %v", err)
	}

	// Add hook to scrub sensitive data before saving
	r.AddHook(scrubAuthHeader, recorder.BeforeSaveHook)

	// Configure airtable client with recorder's HTTP client
	token := os.Getenv("AIRTABLE_KEY")
	if token == "" && mode == recorder.ModeRecordOnly {
		t.Fatal("AIRTABLE_KEY required when recording")
	}
	if token == "" {
		token = "test_token" // Placeholder for replay mode
	}

	Configure(r.GetDefaultClient(), token)

	return r, func() {
		if err := r.Stop(); err != nil {
			t.Errorf("Failed to stop recorder: %v", err)
		}
	}
}

// scrubAuthHeader removes sensitive data from recorded cassettes
func scrubAuthHeader(i *cassette.Interaction) error {
	delete(i.Request.Headers, "Authorization")
	return nil
}

// newIntegrationTable creates a table for integration tests
func newIntegrationTable() *Table[IntegrationSchema] {
	return NewTable[IntegrationSchema](integrationBaseID, integrationTableID)
}

// =============================================================================
// Integration Tests
// =============================================================================

func TestIntegration(t *testing.T) {
	t.Run("List", testIntegrationList)
	t.Run("ListWithFilter", testIntegrationListWithFilter)
	t.Run("ListWithSort", testIntegrationListWithSort)
	t.Run("Get", testIntegrationGet)
	t.Run("CreateUpdateDelete", testIntegrationCreateUpdateDelete)
}

func testIntegrationList(t *testing.T) {
	rec, teardown := setupRecorder(t, "list")
	defer teardown()

	table := newIntegrationTable()
	records, err := table.List()

	Ok(t, err)
	Assert(t, len(records) > 0, "Expected at least one record, got %d", len(records))

	// Log for debugging during recording
	if rec.IsRecording() {
		t.Logf("Found %d records", len(records))
		for _, r := range records {
			t.Logf("  - %s: %s", r.Id, r.Fields.Name)
		}
	}
}

func testIntegrationListWithFilter(t *testing.T) {
	rec, teardown := setupRecorder(t, "list_filter")
	defer teardown()

	table := newIntegrationTable()
	records, err := table.
		WithFilter("{single select} = 'selection 1'").
		WithLimit(10).
		List()

	Ok(t, err)

	// All returned records should have single select = "selection 1"
	for _, r := range records {
		Assert(t, r.Fields.SingleSelect == "selection 1",
			"Expected SingleSelect='selection 1', got '%s'", r.Fields.SingleSelect)
	}

	if rec.IsRecording() {
		t.Logf("Found %d records with selection 1", len(records))
	}
}

func testIntegrationListWithSort(t *testing.T) {
	rec, teardown := setupRecorder(t, "list_sort")
	defer teardown()

	table := newIntegrationTable()
	records, err := table.
		WithSort(Sorts{{Field: "name", Direction: "asc"}}).
		WithLimit(10).
		List()

	Ok(t, err)
	Assert(t, len(records) > 0, "Expected at least one record")

	// Verify ascending order
	for i := 1; i < len(records); i++ {
		Assert(t, records[i-1].Fields.Name <= records[i].Fields.Name,
			"Expected ascending order: '%s' should come before '%s'",
			records[i-1].Fields.Name, records[i].Fields.Name)
	}

	if rec.IsRecording() {
		t.Logf("Retrieved %d sorted records", len(records))
	}
}

func testIntegrationGet(t *testing.T) {
	rec, teardown := setupRecorder(t, "get")
	defer teardown()

	table := newIntegrationTable()

	// First, list to get a valid record ID
	records, err := table.WithLimit(1).List()
	Ok(t, err)
	Assert(t, len(records) > 0, "Need at least one record to test Get")

	recordID := records[0].Id

	// Now fetch that specific record
	record, err := table.Get(recordID)
	Ok(t, err)
	Assert(t, record.Id == recordID, "Expected ID '%s', got '%s'", recordID, record.Id)

	if rec.IsRecording() {
		t.Logf("Retrieved record: %s - %s", record.Id, record.Fields.Name)
	}
}

func testIntegrationCreateUpdateDelete(t *testing.T) {
	_, teardown := setupRecorder(t, "crud")
	defer teardown()

	table := newIntegrationTable()

	// CREATE
	record := table.NewRecord()
	record.Fields = IntegrationSchema{
		Name:         "Integration Test Record",
		Notes:        "Created by integration test",
		SingleSelect: "selection 1",
		MultiSelect:  []string{"selection 1", "selection 2"},
		Number:       42.5,
		Currency:     99.99,
		Checkbox:     true,
		Date:         "2026-03-22",
		Rating:       4,
		Email:        "test@example.com",
		PhoneNumber:  "+1 555 123 4567",
		URL:          "https://example.com",
	}

	err := record.Save()
	Ok(t, err)
	Assert(t, record.Id != "", "Expected record ID after create")
	createdID := record.Id
	t.Logf("Created record: %s", createdID)

	// UPDATE
	record.Fields.Name = "Integration Test Record (Updated)"
	record.Fields.SingleSelect = "selection 2"
	record.Fields.Number = 100.0

	err = record.Save()
	Ok(t, err)
	Assert(t, record.Id == createdID, "ID should remain the same after update")
	t.Logf("Updated record: %s", record.Id)

	// VERIFY UPDATE
	fetched, err := table.Get(createdID)
	Ok(t, err)
	Assert(t, fetched.Fields.Name == "Integration Test Record (Updated)",
		"Expected updated name, got '%s'", fetched.Fields.Name)
	Assert(t, fetched.Fields.Number == 100.0,
		"Expected updated number=100, got %f", fetched.Fields.Number)

	// DELETE
	_, err = record.Destroy()
	Ok(t, err)
	t.Logf("Deleted record: %s", createdID)

	// VERIFY DELETE - should fail to fetch
	_, err = table.Get(createdID)
	Assert(t, err != nil, "Expected error when fetching deleted record")
}

// =============================================================================
// Find Tests (search by field value)
// =============================================================================

func TestIntegrationFind(t *testing.T) {
	rec, teardown := setupRecorder(t, "find")
	defer teardown()

	table := newIntegrationTable()

	// Find records by a specific field value
	// Make sure you have at least one record with single select = "selection 1"
	found, err := table.Find("single select", "selection 1")
	Ok(t, err)

	// All found records should have single select = "selection 1"
	for _, r := range found {
		Assert(t, r.Fields.SingleSelect == "selection 1",
			"Expected SingleSelect='selection 1', got '%s'", r.Fields.SingleSelect)
	}

	if rec.IsRecording() {
		t.Logf("Found %d records with SingleSelect='selection 1'", len(found))
	}
}
