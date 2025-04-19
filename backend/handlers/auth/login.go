package auth

import (
	"encoding/json"
	"event_management/backend/database"
	"event_management/backend/utils"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
}

func JWTMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSONError(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			writeJSONError(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		tokenString := tokenParts[1]

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			writeJSONError(w, "Unauthorized. Invalid or expired token.", http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-User-ID", string(rune(claims.UserID)))
		r.Header.Set("X-User-Email", claims.Email)
		r.Header.Set("X-User-Name", claims.Name)
		r.Header.Set("X-User-Role", claims.Role)

		h(w, r)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		writeJSONError(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	role := r.FormValue("role")

	if email == "" || password == "" || role == "" {
		writeJSONError(w, "Email, password and role are required", http.StatusBadRequest)
		return
	}

	var searchQuery string
	var userId int
	var name, storedPassword, userEmail, phone string

	switch role {
	case "admin":
		searchQuery = `
			SELECT admin_id, name, email, phone, password FROM admin WHERE email=? AND isalive = 1
		`
	case "organiser":
		searchQuery = `
			SELECT organiser_id, name, email, phone, password FROM organiser WHERE email=? AND isalive = 1
		`
	case "attendee":
		searchQuery = `
			SELECT attendee_id, name, email, phone, password FROM attendee WHERE email=? AND isalive = 1
		`
	default:
		writeJSONError(w, "Invalid role specified", http.StatusBadRequest)
		return
	}

	err = database.DB.QueryRow(searchQuery, email).Scan(&userId, &name, &userEmail, &phone, &storedPassword)
	if err != nil {
		writeJSONError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		writeJSONError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(userId, email, name, role)
	if err != nil {
		writeJSONError(w, "Failed to generate authentication token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"token":   token,
		"name":    name,
		"email":   email,
		"role":    role,
	})
}

func ValidateTokenHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeJSONError(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		writeJSONError(w, "Invalid authorization format", http.StatusUnauthorized)
		return
	}

	tokenString := tokenParts[1]

	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		writeJSONError(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Session active",
		"email":   claims.Email,
		"name":    claims.Name,
		"role":    claims.Role,
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}
