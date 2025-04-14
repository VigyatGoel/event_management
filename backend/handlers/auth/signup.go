package auth

import (
	"encoding/json"
	"event_management/backend/database"
	"event_management/backend/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
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

	if name == "" || email == "" || password == "" {
		writeJSONError(w, "All fields are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		writeJSONError(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	_, err = database.DB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)",
		name, email, hashedPassword)
	if err != nil {
		writeJSONError(w, "Email already registered", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.AuthResponse{
		Message: "Signup successful!",
		Name:    name,
		Email:   email,
	})
}