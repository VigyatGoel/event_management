package auth

import (
	"encoding/json"
	"event_management/backend/database"
	"event_management/backend/models"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		writeJSONError(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	phone := r.FormValue("phone")
	role := strings.ToLower(r.FormValue("role"))

	fmt.Println(role)

	if name == "" || email == "" || password == "" || phone == "" || role == "" {
		writeJSONError(w, "All fields are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		writeJSONError(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	createdAt := time.Now()
	isAlive := true

	var insertQuery string

	switch role {
	case "admin":
		insertQuery = `
			INSERT INTO admin (name, email, phone, password, isalive, created_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`
	case "organiser":
		insertQuery = `
			INSERT INTO organiser (name, email, phone, password, isalive, created_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`
	case "attendee":
		insertQuery = `
			INSERT INTO attendee (name, email, phone, password, isalive, created_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`
	default:
		writeJSONError(w, "Invalid role specified", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec(insertQuery, name, email, phone, hashedPassword, isAlive, createdAt)
	if err != nil {
		writeJSONError(w, "Email already registered or DB error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.AuthResponse{
		Message: "Signup successful!",
		Name:    name,
		Email:   email,
		Role:    role,
	})
}
