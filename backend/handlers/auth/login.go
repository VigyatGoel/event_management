package auth

import (
	"context"
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

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSONError(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			writeJSONError(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ValidateJWT(parts[1])
		if err != nil {
			writeJSONError(w, "Unauthorized. Invalid or expired token.", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), utils.UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, utils.UserEmailKey, claims.Email)
		ctx = context.WithValue(ctx, utils.UserNameKey, claims.Name)
		ctx = context.WithValue(ctx, utils.UserRoleKey, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// parse form
	if err := r.ParseForm(); err != nil {
		writeJSONError(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	role    := strings.ToLower(r.FormValue("role"))

	// validate inputs
	if email == "" || password == "" || role == "" {
		writeJSONError(w, "Email, password and role are required", http.StatusBadRequest)
		return
	}
	if role != "admin" && role != "organiser" && role != "attendee" {
		writeJSONError(w, "Invalid role specified", http.StatusBadRequest)
		return
	}

	// unified lookup in user table
	var (
		userId         int
		name, userEmail, phone, storedHash, dbRole string
	)
	query := `
		SELECT user_id, name, email, phone, password, role
		FROM user
		WHERE email = ? AND role = ? AND isalive = 1
	`
	err := database.DB.QueryRow(query, email, role).
		Scan(&userId, &name, &userEmail, &phone, &storedHash, &dbRole)
	if err != nil {
		writeJSONError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// verify password
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)); err != nil {
		writeJSONError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// issue JWT using the DB-confirmed role
	token, err := utils.GenerateJWT(userId, userEmail, name, dbRole)
	if err != nil {
		writeJSONError(w, "Failed to generate authentication token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"token":   token,
		"name":    name,
		"email":   userEmail,
		"role":    dbRole,
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
