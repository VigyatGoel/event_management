package auth

import (
	"encoding/json"
	"event_management/backend/utils"
	"net/http"
)

func HomepageRedirectHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := utils.Store.Get(r, "session")

	email, emailExists := session.Values["email"].(string)
	name, nameExists := session.Values["name"].(string)

	w.Header().Set("Content-Type", "application/json")

	if emailExists && nameExists && email != "" && name != "" {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "logged_in",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": "not_logged_in",
	})
}
