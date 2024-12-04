package handlers

import (
	"database/sql"
	"log"
	"net/http"
)

// Function ya kufuta mwanafunzi
func DeleteStudent(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("auth_token")
		if err != nil {
			// If the cookie is not found, redirect to login
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Validate JWT token from the cookie
		claims, err := ValidateJWT(cookie.Value)
		if err != nil {
			// If the token is invalid or expired, redirect to login
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if claims.Role != "admin" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		// Log authenticated user info for debugging
		log.Printf("Authenticated user: %s, Role: %s", claims.Username, claims.Role)

		id := r.URL.Query().Get("id")
		_, err = db.Exec("DELETE FROM registration WHERE id = ?", id)

		if err != nil {
			http.Error(w, "Error deleting user", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/managestudent", http.StatusSeeOther)
	}
}
