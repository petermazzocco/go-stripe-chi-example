package handlers

import (
	"fmt"

	"net/http"
	"os"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
)

func CreateSubscription(w http.ResponseWriter, r *http.Request) {

	// Create a Checkout Session
	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String("http://localhost:3000/portal"),
		CancelURL:  stripe.String("http://localhost:3000"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(os.Getenv("PRICE_ID")),
				Quantity: stripe.Int64(1),
			},
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
	}
	result, err := session.New(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Redirect", result.URL)
	fmt.Fprintf(w, "Redirecting to Stripe...")
}
