package auth

import (
	"encoding/json"
	"event_management/backend/database"
	"event_management/backend/models"
	"event_management/backend/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		writeJSONError(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		writeJSONError(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	var dbUser models.User
	err = database.DB.QueryRow("SELECT id, name, email, password FROM users WHERE email=?", email).
		Scan(&dbUser.ID, &dbUser.Name, &dbUser.Email, &dbUser.Password)

	if err != nil {
		writeJSONError(w, "Invalid login credentials", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password))
	if err != nil {
		writeJSONError(w, "Invalid login credentials", http.StatusUnauthorized)
		return
	}

	tokenString, err := utils.GenerateToken(dbUser.Email)
	if err != nil {
		writeJSONError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.AuthResponse{
		Message: "Login successful!",
		Name:    dbUser.Name,
		Email:   dbUser.Email,
		Token:   tokenString,
	})
}