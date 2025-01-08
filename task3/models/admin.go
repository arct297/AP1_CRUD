package models

// import (

// )

type MailingRequest struct {
	Topic          string `json:"topic"`           // Subject of the email
	Message        string `json:"message"`         // Content of the email
	ReceivingGroup string `json:"receiving_group"` // Target group for the email
	// AttachmentName string `json:"attachment_name,omitempty"` // Original name of the uploaded file
	// AttachmentPath string `json:"attachment_path,omitempty"` // Path where the file is stored
	// AttachmentSize int64  `json:"attachment_size,omitempty"` // Size of the attachment (bytes)
	// AttachmentType string `json:"attachment_type,omitempty"` // MIME type of the attachment
}
