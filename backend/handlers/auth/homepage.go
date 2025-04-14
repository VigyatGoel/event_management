package auth

import (
	"encoding/json"
	"event_management/backend/utils"
	"net/http"
)

// HomepageRedirectHandler checks if the user is logged in and responds accordingly
func HomepageRedirectHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve session information
	session, _ := utils.Store.Get(r, "session")

	// Check if the session contains the necessary values
	email, emailExists := session.Values["email"].(string)
	name, nameExists := session.Values["name"].(string)

	// Respond with the login status
	w.Header().Set("Content-Type", "application/json")

	// If session is valid, respond with status logged_in
	if emailExists && nameExists && email != "" && name != "" {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "logged_in",
		})
		return
	}

	// If session is invalid or not found, respond with status not_logged_in
	json.NewEncoder(w).Encode(map[string]string{
		"status": "not_logged_in",
	})
}
