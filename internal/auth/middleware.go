package auth

import (
	"fmt"
	"net/http"

	"github.com/markbates/goth/gothic"
)

func RedirectMiddleware(redirectURL string) func(http.Handler) http.Handler {
	fmt.Println("RedirectMiddleware initialized")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Try to get the session using gothic's store
			session, err := gothic.Store.Get(r, gothic.SessionName)
			if err != nil {
				fmt.Println("Error getting session:", err)
				http.Redirect(w, r, redirectURL, http.StatusSeeOther)
				return
			}

			// Check if user is authenticated
			if auth, ok := session.Values["userid"]; !ok || auth == "" {
				fmt.Println("User not authenticated")
				http.Redirect(w, r, redirectURL, http.StatusSeeOther)
				return
			}

			// User is authenticated, proceed
			next.ServeHTTP(w, r)
		})
	}
}
