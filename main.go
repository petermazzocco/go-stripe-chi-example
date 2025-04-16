package main

import (
	"context"
	"fmt"
	handlers "go-stripe-chi-example/handlers/stripe"
	"go-stripe-chi-example/initializers"
	"go-stripe-chi-example/internal/auth"
	"go-stripe-chi-example/views"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/markbates/goth/gothic"
	"github.com/stripe/stripe-go/v82"
)

func init() {
	initializers.InitENV()
	initializers.InitDatabase()
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	auth.NewAuth()
}

func main() {

	r := chi.NewRouter()

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*", "*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Views
	r.Get("/", templ.Handler(views.Index()).ServeHTTP)
	r.Group(func(r chi.Router) {
		r.Use(auth.RedirectMiddleware("/"))
		r.Get("/portal", templ.Handler(views.Portal()).ServeHTTP)
	})

	//Auth
	r.Get("/auth/{provider}/callback", func(w http.ResponseWriter, r *http.Request) {
		provider := chi.URLParam(r, "provider")

		r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		// Store user info in session
		session, _ := gothic.Store.Get(r, gothic.SessionName)
		session.Values["userid"] = user.UserID
		session.Values["name"] = user.Name
		session.Values["email"] = user.Email
		session.Save(r, w)
		http.Redirect(w, r, "http://localhost:3000/portal", http.StatusFound)
	})

	r.Get("/logout/{provider}", func(w http.ResponseWriter, r *http.Request) {
		gothic.Logout(w, r)
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)

	})

	r.Get("/auth/{provider}", func(w http.ResponseWriter, r *http.Request) {
		// try to get the user without re-authenticating
		if _, err := gothic.CompleteUserAuth(w, r); err == nil {

			http.Redirect(w, r, "http://localhost:3000/portal", http.StatusFound)
		} else {
			gothic.BeginAuthHandler(w, r)
		}
	})

	// Routes
	r.Post("/subscribe", handlers.CreateSubscription)

	http.ListenAndServe(":3000", r)
}
