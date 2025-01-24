package models

// import (

// )

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type SignUpRequest struct {
	Login       string  `json:"login"`
	Password    string  `json:"password"`
	PatientData Patient `json:"patient"`
}
