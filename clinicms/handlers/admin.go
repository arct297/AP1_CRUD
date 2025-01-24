package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"gopkg.in/gomail.v2"

	"github.com/sirupsen/logrus"

	"clinicms/logger"
	"clinicms/models"
	"clinicms/tools"
)

func MakeMailing(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		tools.OperateUnsuccessfulResponse(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Ensure the temporary directory exists
	tempDir := "./temp"
	if err := ensureTempDirectoryExists(tempDir); err != nil {
		logger.Log.WithFields(logrus.Fields{
			"action": "make_mailing",
			"error":  err.Error(),
		}).Error("Failed to create temporary directory")
		tools.OperateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Bind form fields to the model
	var mailingData models.MailingRequest
	mailingData.Topic = r.FormValue("topic")
	mailingData.Message = r.FormValue("message")

	var doctors []models.Doctor

	// Process the file
	var tempFilePath string
	file, handler, err := r.FormFile("attachment")
	if err == nil { // File is provided
		defer file.Close()

		// Log the mailing action and send test email
		logger.Log.WithFields(logrus.Fields{
			"action":     "make_mailing",
			"topic":      mailingData.Topic,
			"message":    mailingData.Message,
			"recipients": len(doctors),
			"test_email": "siniov.arseniy@gmail.com",
		}).Info("Sending test email for mailing")

		// if err := sendEmail(mailingData.Topic, mailingData.Message, "debarbiest@gmail.com"); err != nil {
		// 	logger.Log.WithFields(logrus.Fields{
		// 		"action": "make_mailing",
		// 		"error":  err.Error(),
		// 		"path":   tempFilePath,
		// 	}).Error("Failed to create temporary file")
		// 	tools.OperateUnsuccessfulResponse(w, "Failed to save attachment", http.StatusInternalServerError)
		// 	return
		// }
		// defer tempFile.Close()

		// if _, err := tempFile.ReadFrom(file); err != nil {
		// 	logger.Log.WithFields(logrus.Fields{
		// 		"action": "make_mailing",
		// 		"error":  err.Error(),
		// 		"path":   tempFilePath,
		// 	}).Error("Failed to save file data")
		// 	tools.OperateUnsuccessfulResponse(w, "Failed to save attachment data", http.StatusInternalServerError)
		// 	return
		// }

		logger.Log.WithFields(logrus.Fields{
			"filename": handler.Filename,
			"size":     handler.Size,
		}).Info("Attachment saved successfully")
	}

	// Proceed with sending the email (using the extracted file if present)
	err = sendEmailWithAttachment(mailingData.Topic, mailingData.Message, "siniov.arseniy@gmail.com", tempFilePath)
	if err != nil {
		tools.OperateUnsuccessfulResponse(w, "Failed to send email", http.StatusInternalServerError)
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

func sendEmailWithAttachment(topic, text, receiver, attachmentPath string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	authEmail := "kaclinicms@gmail.com"
	authPassword := "zsas fplb seey guqo"

	logger.Log.WithFields(logrus.Fields{
		"action":     "send_email",
		"receiver":   receiver,
		"subject":    topic,
		"attachment": attachmentPath,
	}).Info("Starting email sending")

	m := gomail.NewMessage()
	m.SetHeader("From", authEmail)
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", topic)
	m.SetBody("text/plain", text)

	// Attach the file if provided
	if attachmentPath != "" {
		m.Attach(attachmentPath)
		logger.Log.WithFields(logrus.Fields{
			"action":     "send_email",
			"attachment": attachmentPath,
		}).Info("File attached successfully")
	}

	d := gomail.NewDialer(smtpHost, smtpPort, authEmail, authPassword)

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

func ensureTempDirectoryExists(tempDir string) error {
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		return os.Mkdir(tempDir, 0755)
	}
	return nil
}
