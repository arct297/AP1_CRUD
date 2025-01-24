package tools

import (
	"net/http"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from cookie
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/login", http.StatusFound)
				// OperateUnsuccessfulResponse(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			OperateUnsuccessfulResponse(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		tokenString := cookie.Value

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Redirect(w, r, "/login", http.StatusFound)
			// OperateUnsuccessfulResponse(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusFound)
			// OperateUnsuccessfulResponse(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add claims to the request header for downstream handlers
		r.Header.Set("UserID", strconv.Itoa(int(claims.UserID)))
		r.Header.Set("Role", claims.Role)

		next.ServeHTTP(w, r)
	})
}

func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := r.Header.Get("Role")
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					next.ServeHTTP(w, r)
					return
				}
			}
			OperateUnsuccessfulResponse(w, "Forbidden", http.StatusForbidden)
		})
	}
}
