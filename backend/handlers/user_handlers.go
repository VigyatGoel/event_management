package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"event_management/backend/database"
)

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := database.GetAllUserRoles()
	if err != nil {
		http.Error(w, "Failed to retrieve user data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func DeactivateUserHandler(w http.ResponseWriter, r *http.Request) {

	var requestData struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestData.Email == "" || requestData.Role == "" {
		http.Error(w, "Email and role are required", http.StatusBadRequest)
		return
	}

	err := database.DeactivateUser(requestData.Email, requestData.Role)
	if err != nil {
		if strings.Contains(err.Error(), "no user found") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Printf("Error deactivating user: %v", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User deactivated successfully",
		"email":   requestData.Email,
		"role":    requestData.Role,
	})

}
