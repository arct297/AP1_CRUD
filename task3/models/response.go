package models

// import (

// )

type Response struct {
	Code    int      `json:"code"`
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Content *Patient `json:"content,omitempty"`
}

type ListResponseMeta struct {
	Page   int `json:"page"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type ListResponse[T any] struct {
	Code    int              `json:"code"`
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Content *T               `json:"content,omitempty"`
	Meta    ListResponseMeta `json:"meta"`
}
