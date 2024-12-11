package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Message struct {
	Text string `json:"text"`
}

func getMessage(w http.ResponseWriter, r *http.Request) {
	msg := Message{Text: "Hello, world!"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

func postMessage(w http.ResponseWriter, r *http.Request) {
	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

func main() {
	http.HandleFunc("/message", getMessage)
	// http.HandleFunc("/message", postMessage)
	log.Println("Server starting on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
