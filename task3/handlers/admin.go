package handlers

import (
	"encoding/json"
	"errors"
	// "fmt"
	"log"
	"net/http"

	"gorm.io/gorm"

	"gopkg.in/gomail.v2"
	"task3/models"
	"task3/tools"
)

func MakeMailing(w http.ResponseWriter, r *http.Request) {
	var mailingData models.MailingRequest
	if err := json.NewDecoder(r.Body).Decode(&mailingData); err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"message": "Bad request: Invalid JSON received"})
        return
    } 
	if mailingData.ReceivingGroup == "doctors" {
        var doctors []models.Doctor

        result := tools.DB.Find(&doctors)
        if result.Error != nil {
            w.Header().Set("Content-Type", "application/json")
            if errors.Is(result.Error, gorm.ErrRecordNotFound) {
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(map[string]string{"message": "No doctors found"})
            } else {
                w.WriteHeader(http.StatusInternalServerError)
                json.NewEncoder(w).Encode(map[string]string{"message": "Internal server error"})
            }
            return
        }

        // Send email (just for testing, sending to one email)
        sendEmail(mailingData.Topic, mailingData.Message, "debarbiest@gmail.com")

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"success": "true", "message": "Email sent successfully!"})
    } else {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"message": "Bad request: Unsupported receiving group"})
    }
}

	

func sendEmail(topic string, text string, receiver string) {
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	authEmail := "kaclinicms@gmail.com"
	authPassword := "zsas fplb seey guqo"

	m := gomail.NewMessage()
	m.SetHeader("From", authEmail)
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", topic)
	m.SetBody("text/plain", text)

	d := gomail.NewDialer(smtpHost, smtpPort, authEmail, authPassword)

	if err := d.DialAndSend(m); err != nil {
		log.Fatalf("Failed to send email: %v", err)
	} else {
		log.Println("Email sent successfully!")
	}
}
