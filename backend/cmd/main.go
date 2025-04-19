package main

import (
	"event_management/backend/database"
	"event_management/backend/handlers"
	"event_management/backend/handlers/auth"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting the server...")
	database.InitDB()

	router := http.NewServeMux()

	router.HandleFunc("POST /signup", auth.SignupHandler)
	router.HandleFunc("POST /login", auth.LoginHandler)
	router.HandleFunc("GET /validate_token", auth.ValidateTokenHandler)
	router.HandleFunc("POST /logout", auth.LogoutHandler)
	router.HandleFunc("GET /users", auth.JWTMiddleware(handlers.GetAllUsersHandler))
	router.HandleFunc("POST /users/deactivate", auth.JWTMiddleware(handlers.DeactivateUserHandler))

	corsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		router.ServeHTTP(w, r)
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: corsHandler,
	}

	log.Println("Server running on http://localhost:8080")

	log.Fatal(server.ListenAndServe())
}
