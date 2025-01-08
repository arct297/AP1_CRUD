package handlers

import (
	"net/http"
	"os"

	"task3/logger"
	"task3/models"
	"task3/tools"

	"gopkg.in/gomail.v2"

	"github.com/sirupsen/logrus"
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

	// Process the file
	var tempFilePath string
	file, handler, err := r.FormFile("attachment")
	if err == nil { // File is provided
		defer file.Close()

		tempFilePath = tempDir + "/" + handler.Filename
		tempFile, err := os.Create(tempFilePath)
		if err != nil {
			logger.Log.WithFields(logrus.Fields{
				"action": "make_mailing",
				"error":  err.Error(),
				"path":   tempFilePath,
			}).Error("Failed to create temporary file")
			tools.OperateUnsuccessfulResponse(w, "Failed to save attachment", http.StatusInternalServerError)
			return
		}
		defer tempFile.Close()

		if _, err := tempFile.ReadFrom(file); err != nil {
			logger.Log.WithFields(logrus.Fields{
				"action": "make_mailing",
				"error":  err.Error(),
				"path":   tempFilePath,
			}).Error("Failed to save file data")
			tools.OperateUnsuccessfulResponse(w, "Failed to save attachment data", http.StatusInternalServerError)
			return
		}

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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Mailing completed successfully"))
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
