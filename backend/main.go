package main

import (
	"event_management/backend/database"
	"event_management/backend/handlers"
	"log"
	"net/http"
)

// CORS middleware
func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call actual handler
		h(w, r)
	}
}

func main() {
	database.Connect()

	http.HandleFunc("/signup", withCORS(handlers.SignupPage))
	http.HandleFunc("/signup/submit", withCORS(handlers.SignupHandler))

	http.HandleFunc("/login", withCORS(handlers.LoginPage))
	http.HandleFunc("/login/submit", withCORS(handlers.LoginHandler))

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
