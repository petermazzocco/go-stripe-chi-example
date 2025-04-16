package auth

import (
	"fmt"
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const (
	key    = "9vI5odu3w+OGZg=="
	MaxAge = 86400 * 30
	IsProd = false
)

var store *sessions.CookieStore

func NewAuth() {
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(MaxAge)
	gothic.Store = store

	goth.UseProviders(
		google.New(googleClientID, googleClientSecret, "http://localhost:3000/auth/google/callback"),
	)

	fmt.Println("Auth initialized")
}
