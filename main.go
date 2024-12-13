package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Request struct {
	Message string `json:"message"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func returnJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Status:  "fail",
		Message: message,
	})
}

func operateMessage(w http.ResponseWriter, r *http.Request) {
	var request Request

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		returnJSONError(w, "Incorrect JSON message", http.StatusBadRequest)
		return
	}

	if request.Message == "" {
		returnJSONError(w, "Incorrect JSON message", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(Response{
		Status:  "success",
		Message: "data has been successfully received",
	})
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		operateMessage(w, r)
	case "POST":
		operateMessage(w, r)
	default:
		http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/message", handleMessage)
	log.Println("Server starting on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
