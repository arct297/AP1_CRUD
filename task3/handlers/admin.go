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
		tools.OperateUnsuccessfulResponse(w, "Bad request: Invalid JSON received", http.StatusBadRequest)
		return
	}

	if mailingData.ReceivingGroup == "doctors" {
		var doctors []models.Doctor

		result := tools.DB.Find(&doctors)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				tools.OperateUnsuccessfulResponse(w, "No doctors found", http.StatusNotFound)
			} else {
				tools.OperateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		// Commented to avoid random emails to generated doctors
		// for _, doctor := range doctors {
		// 	sendEmail(mailingData.Topic, mailingData.Message, doctor.Email)
		// }

		sendEmail(mailingData.Topic, mailingData.Message, "siniov.arseniy@gmail.com")

	} else {
		tools.OperateUnsuccessfulResponse(w, "Bad request: Unsupported receiving group", http.StatusBadRequest)
		return
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
