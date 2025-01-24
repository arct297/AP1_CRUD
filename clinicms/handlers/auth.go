package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	// "gorm.io/gorm"

	"clinicms/models"
	"clinicms/tools"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// JWT claims
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginData models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		tools.OperateUnsuccessfulResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := tools.DB.Where("login = ?", loginData.Login).First(&user).Error; err != nil {
		tools.OperateUnsuccessfulResponse(w, "Invalid login or password", http.StatusUnauthorized)
		return
	}

	if err := tools.CheckPassword(loginData.Password, user.Password); err != nil {
		tools.OperateUnsuccessfulResponse(w, "Invalid login or password", http.StatusUnauthorized)
		return
	}

	if !user.IsVerified {
		tools.OperateUnsuccessfulResponse(w, "Email is not verified", http.StatusForbidden)
		return
	}

	token, err := generateJWT(user.ID, user.Role)
	if err != nil {
		tools.OperateUnsuccessfulResponse(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Set the token in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false, // Prevents JavaScript access to the cookie
		Secure:   false, // Use true in production with HTTPS
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Response{
		Code:    http.StatusOK,
		Status:  "success",
		Message: "Login successful",
	})
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	// Clear the token in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, "/login", http.StatusFound)
}

func generateJWT(userID uint, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
