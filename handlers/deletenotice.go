package handlers

import (
	"database/sql"
	"log"
	"net/http"
)

// DeleteNotice deletes a public notice by its ID
func DeleteNotice(db *sql.DB) http.HandlerFunc {
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

		delID := r.URL.Query().Get("delID")
		if delID == "" {
			http.Error(w, "Missing ID parameter", http.StatusBadRequest)
			return
		}
		// Execute the DELETE query
		_, err = db.Exec("DELETE FROM tblpublicnotice WHERE id = ?", delID)
		if err != nil {
			log.Printf("Failed to delete notice: %v", err)
			http.Error(w, "Failed to delete notice.", http.StatusInternalServerError)
			return
		}

		// Redirect to the manage public notice page after deletion
		http.Redirect(w, r, "/manage-public-notice", http.StatusSeeOther)
	}
}
