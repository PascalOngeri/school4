package handlers

import (
	"database/sql"
	"log"
	"net/http"
)

// DeleteNotice deletes a public notice by its ID
func DeleteNotice(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the notice ID from the URL query string

		session, err := store.Get(r, "store")
		if err != nil {
			log.Printf("Failed to retrieve session: %v", err)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}

		// Check if user is logged in
		if session.Values["sturecmsaid"] == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		
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
