package middleware

import (
    "event_management/backend/utils"
    "net/http"
)

// AuthMiddleware checks if user is authenticated
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        session, err := utils.Store.Get(r, "user-session")
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Check if user is authenticated
        auth, ok := session.Values["authenticated"].(bool)
        if !ok || !auth {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    }
}

// AdminMiddleware checks if user is admin
func AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        session, err := utils.Store.Get(r, "user-session")
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Check if user is authenticated
        auth, ok := session.Values["authenticated"].(bool)
        if !ok || !auth {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Check if user is admin
        role, ok := session.Values["user_role"].(string)
        if !ok || role != "admin" {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }

        next.ServeHTTP(w, r)
    }
}