package airtable

// Thumbnail represents a thumbnail image for an Airtable attachment.
// Thumbnails are generated automatically by Airtable for image attachments.
type Thumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// Thumbnails contains the different thumbnail sizes that Airtable generates.
// All fields are pointers as not all size variants may be present.
type Thumbnails struct {
	Small *Thumbnail `json:"small,omitempty"`
	Large *Thumbnail `json:"large,omitempty"`
	Full  *Thumbnail `json:"full,omitempty"`
}

// Attachment represents a file attached to an Airtable record.
// When reading attachments, all fields will be populated.
// When creating/updating attachments, only URL is required - Airtable will
// fetch the file and populate the other fields automatically.
//
// Example usage in a schema:
//
//	type MySchema struct {
//	    Documents []airtable.Attachment `json:"Documents"`
//	    Photo     []airtable.Attachment `json:"Photo"`
//	}
type Attachment struct {
	ID         string      `json:"id,omitempty"`
	URL        string      `json:"url"`
	Filename   string      `json:"filename,omitempty"`
	Size       int         `json:"size,omitempty"`
	Type       string      `json:"type,omitempty"`
	Thumbnails *Thumbnails `json:"thumbnails,omitempty"`
}

// Collaborator represents an Airtable user or collaborator.
// This is used for fields like "Created by", "Last modified by",
// and user/collaborator field types.
//
// Example usage in a schema:
//
//	type MySchema struct {
//	    CreatedBy  *airtable.Collaborator   `json:"Created By"`
//	    AssignedTo []airtable.Collaborator  `json:"Assigned To"`
//	}
//
// Note: Collaborator fields are read-only. You cannot set or modify
// collaborator values through the API - they are populated automatically
// by Airtable based on user actions.
type Collaborator struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
