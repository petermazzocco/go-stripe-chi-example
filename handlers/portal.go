package handlers

import (
	"net/http"
)

func Portal(r *http.Request, w http.ResponseWriter) {
	// Check the query param ?status=success
	// If it's "success", then we update the db to active sub true and render the portal page
	// If it's "failed", then we redirect back to the /subscribe page
	if status := r.URL.Query().Get("status"); status == "success" {
		// Update the db to active sub true
		// Render the portal page
	} else if status == "failed" {
		// Redirect back to the /subscribe page
		http.Redirect(w, r, "/subscribe", http.StatusFound)
	}
}
