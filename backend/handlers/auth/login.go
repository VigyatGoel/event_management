package auth

import (
	"context"
	"encoding/json"
	"event_management/backend/database"
	"event_management/backend/models"
	"event_management/backend/utils"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// This is an in-memory session map for demonstration purposes (not suitable for production)
var sessions = make(map[string]string) // Maps sessionID to userID

// Utility function to send errors in JSON format
func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
}

// Middleware to validate the session for protected routes
func SessionMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the session cookie
		cookie, err := r.Cookie("session_id")
		if err != nil {
			writeJSONError(w, "Session not found. Please log in.", http.StatusUnauthorized)
			return
		}

		sessionID := cookie.Value
		userID, exists := sessions[sessionID]
		if !exists {
			writeJSONError(w, "Invalid session. Please log in again.", http.StatusUnauthorized)
			return
		}

		// Store the user ID in the request context for downstream handlers
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", userID)
		r = r.WithContext(ctx)

		h(w, r) // Call the next handler
	}
}

// Login handler to authenticate the user and create a session
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

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password))
	if err != nil {
		writeJSONError(w, "Invalid login credentials", http.StatusUnauthorized)
		return
	}

	// Create a new session and store the user data
	session, _ := utils.Store.Get(r, "session")
	session.Values["email"] = dbUser.Email
	session.Values["name"] = dbUser.Name
	session.Values["user_id"] = dbUser.ID
	session.Save(r, w)

	// Set a cookie with the session ID
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Second), // 1 day expiration
		HttpOnly: true,
		Secure:   false, // Set to true in production for HTTPS
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful!",
		"name":    dbUser.Name,
		"email":   dbUser.Email,
	})
}

// Session handler to check if the user is logged in
func SessionHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the session from the store
	session, _ := utils.Store.Get(r, "session")

	// Check if the session has user data
	email, ok := session.Values["email"].(string)
	name, ok2 := session.Values["name"].(string)

	if !ok || !ok2 {
		writeJSONError(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	// Return session details as JSON
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Session active",
		"email":   email,
		"name":    name,
	})
}

// Logout handler to clear the session
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := utils.Store.Get(r, "session")
	session.Options.MaxAge = -1 // Delete the session by setting the MaxAge to a negative value
	session.Save(r, w)

	// Delete the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}
