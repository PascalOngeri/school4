package handlers

import (
	"log"
	"net/http"
)

func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log the request to log out
		log.Println("Logout request received")

		// Assuming the logout doesn't use a session, handle it directly
		// Log the user's attempt to log out
		username := r.URL.Query().Get("username")
		if username == "" {
			log.Println("Logout attempt without username parameter")
		} else {
			log.Printf("Logging out user: %s", username)
		}

		// Redirect to the login page
		log.Println("Redirecting to the login page")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
