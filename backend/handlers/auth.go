package handlers

import (
    "event_management/backend/database"
    "event_management/backend/models"
    "golang.org/x/crypto/bcrypt"
    "html/template"
    "net/http"
)

var templates = template.Must(template.ParseGlob("backend/templates/*.html"))

func SignupPage(w http.ResponseWriter, r *http.Request) {
    templates.ExecuteTemplate(w, "signup.html", nil)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
    templates.ExecuteTemplate(w, "login.html", nil)
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/signup", http.StatusSeeOther)
        return
    }

    name := r.FormValue("name")
    email := r.FormValue("email")
    password := r.FormValue("password")

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

    _, err := database.DB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)",
        name, email, hashedPassword)

    if err != nil {
        http.Error(w, "Email already registered", http.StatusBadRequest)
        return
    }

    http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    email := r.FormValue("email")
    password := r.FormValue("password")

    var dbUser models.User
    err := database.DB.QueryRow("SELECT id, name, email, password FROM users WHERE email=?", email).
        Scan(&dbUser.ID, &dbUser.Name, &dbUser.Email, &dbUser.Password)

    if err != nil || bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password)) != nil {
        http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
        return
    }

    w.Write([]byte("Welcome, " + dbUser.Name + "!"))
}