package auth

import (
	"encoding/json"
	"event_management/backend/database"
	"event_management/backend/models"

	"net/http"
	"strings"
	"time"
	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		writeJSONError(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	phone := r.FormValue("phone")
	role := strings.ToLower(r.FormValue("role"))

	// Validate inputs
	if name == "" || email == "" || password == "" || phone == "" || role == "" {
		writeJSONError(w, "All fields are required", http.StatusBadRequest)
		return
	}

	if role != "admin" && role != "organiser" && role != "attendee" {
		writeJSONError(w, "Invalid role specified", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		writeJSONError(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	createdAt := time.Now()
	isAlive := true

	// Single INSERT into unified `user` table
	insertQuery := `
		INSERT INTO user (name, email, phone, password, role, isalive, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err = database.DB.Exec(insertQuery,
		name, email, phone, hashedPassword, role, isAlive, createdAt,
	)
	if err != nil {
		// Could check for duplicate-email error via MySQL error code 1062 if desired
		writeJSONError(w, "Email already registered or DB error", http.StatusBadRequest)
		return
	}

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.AuthResponse{
		Message: "Signup successful!",
		Name:    name,
		Email:   email,
		Role:    role,
	})
}

// writeJSONError is assumed to be defined elsewhere in this package:
// func writeJSONError(w http.ResponseWriter, message string, code int) { … }
