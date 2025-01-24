package handlers

import (
	"encoding/json"
	// "fmt"
	"net/http"

	"clinicms/models"
	"clinicms/tools"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var patient models.Patient

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		tools.OperateUnsuccessfulResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	hashedPassword, err := tools.HashPassword(user.Password)
	if err != nil {
		tools.OperateUnsuccessfulResponse(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user.Password = hashedPassword
	user.IsVerified = false
	user.Role = "patient" // Default role for new users

	if err := tools.DB.Create(&user).Error; err != nil {
		tools.OperateUnsuccessfulResponse(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	patient.UserID = user.ID
	patient.Name = "Not specified"
	patient.Gender = "O"
	patient.Contact = "Not specified"
	patient.Address = "Not specified"
	if err := tools.DB.Create(&patient).Error; err != nil {
		tools.OperateUnsuccessfulResponse(w, "Failed to create patient", http.StatusInternalServerError)
		return
	}

	// Generate confirmation link
	// confirmationLink := fmt.Sprintf("http://localhost:8080/api/confirm?email=%s", user.Login)

	// Send confirmation email
	// go tools.SendEmail("Confirm your registration", fmt.Sprintf("Click the link to confirm: %s", confirmationLink), user.Login)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Link to verification sent",
	})
}

func ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	var user models.User
	if err := tools.DB.Where("login = ?", email).First(&user).Error; err != nil {
		tools.OperateUnsuccessfulResponse(w, "User not found", http.StatusNotFound)
		return
	}

	if user.IsVerified {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	user.IsVerified = true
	if err := tools.DB.Save(&user).Error; err != nil {
		tools.OperateUnsuccessfulResponse(w, "Failed to update user status", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}
