package models

// import (

// )

type MailingRequest struct {
	Topic          string `json:"topic"`
	Message        string `json:"message"`
	ReceivingGroup string `json:"receiving_group"`
}
