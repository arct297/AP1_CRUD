package tools

import (
	"gopkg.in/gomail.v2"

	"github.com/sirupsen/logrus"

	"clinicms/logger"
)

func SendEmail(topic string, text string, receiver string) error {
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
