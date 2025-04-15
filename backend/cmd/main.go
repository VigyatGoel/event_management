package main

import (
	"event_management/backend/database"
	"event_management/backend/handlers"
	"event_management/backend/handlers/auth"
	"fmt"
	"log"
	"net/http"
)

func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h(w, r)
	}
}

func main() {
	fmt.Println("Starting the server...")
	database.InitDB()

	router := http.NewServeMux()

	router.HandleFunc("POST /signup", withCORS(auth.SignupHandler))
	router.HandleFunc("POST /login", withCORS(auth.LoginHandler))
	router.HandleFunc("GET /session", withCORS(auth.SessionHandler))
	router.HandleFunc("POST /logout", withCORS(auth.LogoutHandler))
	router.HandleFunc("GET /users", withCORS(auth.SessionMiddleware(handlers.GetAllUsersHandler)))
	router.HandleFunc("/users/deactivate", withCORS(auth.SessionMiddleware(handlers.DeactivateUserHandler)))

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("Server running on http://localhost:8080")

	// Build the application using: go build -o app
	// Run the application using: ./app
	log.Fatal(server.ListenAndServe())
}
