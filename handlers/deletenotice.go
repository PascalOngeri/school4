package handlers

import (
	"database/sql"
	"log"
	"net/http"
)

// DeleteNotice deletes a public notice by its ID
func DeleteNotice(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log the incoming request
		log.Printf("[INFO] Received request to delete notice: %s %s", r.Method, r.URL.Path)

		// Retrieve the 'delID' query parameter
		delID := r.URL.Query().Get("delID")
		if delID == "" {
			http.Error(w, "Missing ID parameter", http.StatusBadRequest)
			log.Println("[ERROR] Missing 'delID' parameter in request")
			return
		}

		// Execute the DELETE query
		_, err := db.Exec("DELETE FROM tblpublicnotice WHERE id = ?", delID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete notice with ID %s: %v", delID, err)
			http.Error(w, "Failed to delete notice.", http.StatusInternalServerError)
			return
		}

		// Log successful deletion
		log.Printf("[INFO] Successfully deleted notice with ID %s", delID)

		// Redirect to the manage public notice page after deletion
		http.Redirect(w, r, "/manage-public-notice", http.StatusSeeOther)
	}
}
