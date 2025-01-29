package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// FirebaseAuth is the Firebase Authentication client
var FirebaseAuth *auth.Client

// InitializeFirebaseAuth initializes Firebase Authentication
func InitializeFirebaseAuth() {
	// Path to your Firebase service account key JSON file
	opt := option.WithCredentialsFile("path/to/your-firebase-service-account-key.json") // Replace with your path
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v", err)
	}

	FirebaseAuth, err = app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Error initializing Firebase Auth client: %v", err)
	}
}

// LoginRequest represents the JSON request body for login
type LoginRequest struct {
	Email    string `json:"email"`    // For email/password authentication
	Password string `json:"password"` // For email/password authentication
	IDToken  string `json:"idToken"`  // For Google Sign-In authentication
}

// EmailPasswordLogin handles login using email and password
// EmailPasswordLogin handles login using email and password
func EmailPasswordLogin(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify the Firebase ID token (same as Google Login)
	token, err := FirebaseAuth.VerifyIDToken(context.Background(), req.IDToken)
	if err != nil {
		http.Error(w, "Invalid ID token", http.StatusUnauthorized)
		return
	}

	// Respond with success message and the user's UID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"uid":     token.UID,
	})
}

// GoogleLogin handles login using Google Sign-In
// GoogleLogin handles login using Google Sign-In
func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify the Google ID token
	token, err := FirebaseAuth.VerifyIDToken(context.Background(), req.IDToken)
	if err != nil {
		http.Error(w, "Invalid Google ID token", http.StatusUnauthorized)
		return
	}

	// Respond with a success message and the user's UID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Google Sign-In successful",
		"uid":     token.UID,
	})
}

	