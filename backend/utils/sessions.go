package utils

import "github.com/gorilla/sessions"

var Store = sessions.NewCookieStore([]byte("09bc8bd828263e4985faf8d8b6b575071a8aede7efbd8a877ed56997c0ac672a"))

func init() {
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   false,
		SameSite: 1,
	}
}
