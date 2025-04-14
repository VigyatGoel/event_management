package auth

import (
	"context"
	"encoding/json"
	"event_management/backend/database"
	"event_management/backend/models"
	"event_management/backend/utils"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
}

func SessionMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := utils.Store.Get(r, "session")
		userID, ok := session.Values["user_id"].(int)
		if !ok {
			writeJSONError(w, "Unauthorized. Please log in.", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		r = r.WithContext(ctx)
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

	switch role {
	case "admin":
		searchQuery = `
			SELECT admin_id, name, email, phone, password FROM admin WHERE email=?
		`
	case "organiser":
		searchQuery = `
			SELECT organiser_id, name, email, phone, password FROM organiser WHERE email=?
		`
	case "attendee":
		searchQuery = `
			SELECT attendee_id, name, email, phone, password FROM attendee WHERE email=?
		`
	default:
		writeJSONError(w, "Invalid role specified", http.StatusBadRequest)
		return
	}

	var dbUser models.User
	err = database.DB.QueryRow(searchQuery, email).
		Scan(&dbUser.ID, &dbUser.Name, &dbUser.Email, &dbUser.Phone, &dbUser.Password)

	if err != nil {
		writeJSONError(w, "Invalid login credentials", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password))
	if err != nil {
		writeJSONError(w, "Invalid login credentials", http.StatusUnauthorized)
		return
	}

	session, _ := utils.Store.Get(r, "session")
	session.Values["user_id"] = dbUser.ID
	session.Values["email"] = dbUser.Email
	session.Values["name"] = dbUser.Name
	session.Values["role"] = role
	session.Save(r, w)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful!",
		"name":    dbUser.Name,
		"email":   dbUser.Email,
		"role":    role,
	})
}

func SessionHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := utils.Store.Get(r, "session")

	email, ok := session.Values["email"].(string)
	name, ok2 := session.Values["name"].(string)
	role, ok3 := session.Values["role"].(string)

	if !ok || !ok2 || !ok3 {
		writeJSONError(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Session active",
		"email":   email,
		"name":    name,
		"role":    role,
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := utils.Store.Get(r, "session")
	session.Options.MaxAge = -1
	session.Save(r, w)

	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}
