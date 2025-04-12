package utils

import (
	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte("your-secret-key")) // Replace with env var in production
