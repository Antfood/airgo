package airtable

import (
	"encoding/json"
	"os"
	"testing"

	. "github.com/pbotsaris/airgo/testutils/testutils"
)

func TestTypes(t *testing.T) {
	t.Run("AttachmentUnmarshal", testAttachmentUnmarshal)
	t.Run("AttachmentNoThumbnails", testAttachmentNoThumbnails)
	t.Run("AttachmentMarshal", testAttachmentMarshal)
	t.Run("CollaboratorUnmarshal", testCollaboratorUnmarshal)
}

func testAttachmentUnmarshal(t *testing.T) {
	data, err := os.ReadFile("./testdata/attachment.json")
	Ok(t, err)

	var att Attachment
	err = json.Unmarshal(data, &att)
	Ok(t, err)

	Assert(t, att.ID == "attABC123XYZ", "Expected ID 'attABC123XYZ', got '%s'", att.ID)
	Assert(t, att.URL == "https://dl.airtable.com/.attachments/abc123/document.pdf", "Expected URL mismatch")
	Assert(t, att.Filename == "document.pdf", "Expected Filename 'document.pdf', got '%s'", att.Filename)
	Assert(t, att.Size == 12345, "Expected Size 12345, got %d", att.Size)
	Assert(t, att.Type == "application/pdf", "Expected Type 'application/pdf', got '%s'", att.Type)

	Assert(t, att.Thumbnails != nil, "Expected Thumbnails to be present")
	Assert(t, att.Thumbnails.Small != nil, "Expected Small thumbnail to be present")
	Assert(t, att.Thumbnails.Small.URL == "https://dl.airtable.com/.attachmentThumbnails/abc123/small", "Small URL mismatch")
	Assert(t, att.Thumbnails.Small.Width == 36, "Expected Small width 36, got %d", att.Thumbnails.Small.Width)
	Assert(t, att.Thumbnails.Small.Height == 36, "Expected Small height 36, got %d", att.Thumbnails.Small.Height)

	Assert(t, att.Thumbnails.Large != nil, "Expected Large thumbnail to be present")
	Assert(t, att.Thumbnails.Large.Width == 512, "Expected Large width 512, got %d", att.Thumbnails.Large.Width)

	Assert(t, att.Thumbnails.Full != nil, "Expected Full thumbnail to be present")
	Assert(t, att.Thumbnails.Full.Width == 1024, "Expected Full width 1024, got %d", att.Thumbnails.Full.Width)
	Assert(t, att.Thumbnails.Full.Height == 768, "Expected Full height 768, got %d", att.Thumbnails.Full.Height)
}

func testAttachmentNoThumbnails(t *testing.T) {
	data, err := os.ReadFile("./testdata/attachment_no_thumbnails.json")
	Ok(t, err)

	var att Attachment
	err = json.Unmarshal(data, &att)
	Ok(t, err)

	Assert(t, att.ID == "attDEF456ABC", "Expected ID 'attDEF456ABC', got '%s'", att.ID)
	Assert(t, att.Filename == "data.csv", "Expected Filename 'data.csv', got '%s'", att.Filename)
	Assert(t, att.Type == "text/csv", "Expected Type 'text/csv', got '%s'", att.Type)
	Assert(t, att.Thumbnails == nil, "Expected Thumbnails to be nil for non-image attachment")
}

func testAttachmentMarshal(t *testing.T) {
	// When creating an attachment, only URL is required
	att := Attachment{
		URL: "https://example.com/file.pdf",
	}

	data, err := json.Marshal(att)
	Ok(t, err)

	// Should only contain url field (others are omitempty)
	expected := `{"url":"https://example.com/file.pdf"}`
	Assert(t, string(data) == expected, "Expected '%s', got '%s'", expected, string(data))
}

func testCollaboratorUnmarshal(t *testing.T) {
	data, err := os.ReadFile("./testdata/collaborator.json")
	Ok(t, err)

	var collab Collaborator
	err = json.Unmarshal(data, &collab)
	Ok(t, err)

	Assert(t, collab.ID == "usrUluXd4j2EEgZnt", "Expected ID 'usrUluXd4j2EEgZnt', got '%s'", collab.ID)
	Assert(t, collab.Email == "pedro@antfood.com", "Expected Email 'pedro@antfood.com', got '%s'", collab.Email)
	Assert(t, collab.Name == "Pedro Botsaris", "Expected Name 'Pedro Botsaris', got '%s'", collab.Name)
}
