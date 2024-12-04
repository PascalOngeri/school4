package handlers

import (
	"net/http"
	"time"
)

// LogoutHandler handles the logout process for JWT-based authentication
func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Clear the JWT by setting an expired cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",    // The cookie where JWT is stored
			Value:    "",              // Clear the cookie
			Expires:  time.Unix(0, 0), // Expire immediately
			HttpOnly: true,
			Secure:   true, // Set to true if using HTTPS
			Path:     "/",
		})

		// Optionally, you can add any other steps needed before logout (like logging out from sessions, etc.)

		// Redirect to the login page
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
