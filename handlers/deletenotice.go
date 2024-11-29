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
		delID := r.URL.Query().Get("delID")
		if delID == "" {
			log.Printf("Delete failed: Missing ID parameter")
			http.Error(w, "Missing ID parameter", http.StatusBadRequest)
			return
		}
		log.Printf("Attempting to delete notice with ID: %s", delID)

		// Execute the DELETE query
		_, err := db.Exec("DELETE FROM tblpublicnotice WHERE id = ?", delID)
		if err != nil {
			log.Printf("Failed to delete notice with ID %s: %v", delID, err)
			http.Error(w, "Failed to delete notice.", http.StatusInternalServerError)
			return
		}

		// Log successful deletion
		log.Printf("Successfully deleted notice with ID: %s", delID)

		// Redirect to the manage public notice page after deletion
		log.Println("Redirecting to /manage-public-notice page.")
		http.Redirect(w, r, "/manage-public-notice", http.StatusSeeOther)
	}
}
