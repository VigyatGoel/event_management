package utils

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
    "fmt"
    "net/http"
    "context"
)

var jwtKey = []byte("b51c7729ad6e4c6ea4470b3ec55d7aaa4e8b7985cbd20c7617d3bd63168659bc")

type Claims struct {
    Email string `json:"email"`
    jwt.RegisteredClaims
}

func GenerateToken(email string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        Email: email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    return tokenString, err
}

func ValidateToken(tokenStr string) (*Claims, error) {
    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    if err != nil || !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }

    return claims, nil
}

func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        tokenStr := r.Header.Get("Authorization")
        if tokenStr == "" {
            http.Error(w, "Missing token", http.StatusUnauthorized)
            return
        }

        claims, err := ValidateToken(tokenStr)
        if err != nil {
            http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
            return
        }

		r = r.WithContext(context.WithValue(r.Context(), UserEmailKey, claims.Email))
        next(w, r)
    }
}