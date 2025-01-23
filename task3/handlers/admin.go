package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"gorm.io/gorm"

	"task3/logger"
	"task3/models"
	"task3/tools"

	"gopkg.in/gomail.v2"

	"github.com/sirupsen/logrus"
)

func MakeMailing(w http.ResponseWriter, r *http.Request) {
	var mailingData models.MailingRequest

	logger.Log.WithFields(logrus.Fields{
		"action": "make_mailing",
		"method": r.Method,
	}).Info("Received mailing request")

	// Decode the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&mailingData); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"action":  "make_mailing",
			"error":   err.Error(),
			"details": "Invalid JSON received",
		}).Error("Failed to decode request body")
		tools.OperateUnsuccessfulResponse(w, "Bad request: Invalid JSON received", http.StatusBadRequest)
		return
	}

	// Validate the receiving group
	if mailingData.ReceivingGroup == "doctors" {
		logger.Log.WithFields(logrus.Fields{
			"action":          "make_mailing",
			"receiving_group": "doctors",
			"topic":           mailingData.Topic,
		}).Info("Fetching doctors for mailing")

		var doctors []models.Doctor

		// Fetch doctors from the database
		result := tools.DB.Find(&doctors)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				logger.Log.WithFields(logrus.Fields{
					"action": "make_mailing",
					"error":  "No doctors found",
				}).Warn("No doctors found for mailing")
				tools.OperateUnsuccessfulResponse(w, "No doctors found", http.StatusNotFound)
			} else {
				logger.Log.WithFields(logrus.Fields{
					"action": "make_mailing",
					"error":  result.Error.Error(),
				}).Error("Failed to fetch doctors from database")
				tools.OperateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		// Log the mailing action and send test email
		logger.Log.WithFields(logrus.Fields{
			"action":     "make_mailing",
			"topic":      mailingData.Topic,
			"message":    mailingData.Message,
			"recipients": len(doctors),
			"test_email": "siniov.arseniy@gmail.com",
		}).Info("Sending test email for mailing")

		if err := sendEmail(mailingData.Topic, mailingData.Message, "debarbiest@gmail.com"); err != nil {
			logger.Log.WithFields(logrus.Fields{
				"action": "make_mailing",
				"error":  err.Error(),
			}).Error("Failed to send test email")
			tools.OperateUnsuccessfulResponse(w, "Failed to send test email", http.StatusInternalServerError)
			return
		}

	} else {
		logger.Log.WithFields(logrus.Fields{
			"action":          "make_mailing",
			"receiving_group": mailingData.ReceivingGroup,
		}).Warn("Unsupported receiving group for mailing")
		tools.OperateUnsuccessfulResponse(w, "Bad request: Unsupported receiving group", http.StatusBadRequest)
		return
	}

	logger.Log.WithFields(logrus.Fields{
		"action": "make_mailing",
		"status": "success",
	}).Info("Mailing completed successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Mailing completed successfully",
	})
}

func sendEmail(topic string, text string, receiver string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	authEmail := "kaclinicms@gmail.com"
	authPassword := "zsas fplb seey guqo"

	logger.Log.WithFields(logrus.Fields{
		"action":   "send_email",
		"receiver": receiver,
		"subject":  topic,
	}).Info("Starting email sending")

	m := gomail.NewMessage()
	m.SetHeader("From", authEmail)
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", topic)
	m.SetBody("text/plain", text)

	d := gomail.NewDialer(smtpHost, smtpPort, authEmail, authPassword)

	// Handle email sending errors
	if err := d.DialAndSend(m); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"action":   "send_email",
			"receiver": receiver,
			"error":    err.Error(),
		}).Error("Failed to send email")
		return err
	}

	logger.Log.WithFields(logrus.Fields{
		"action":   "send_email",
		"receiver": receiver,
	}).Info("Email sent successfully")
	return nil
}
