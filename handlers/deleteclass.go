package handlers

import (
	"database/sql"
	"log"
	"net/http"
)

// DeleteClass deletes a class by ID
func DeleteClass(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log the incoming request
		log.Printf("[INFO] Received request to delete class: %s %s", r.Method, r.URL.Path)

		// Retrieve the `delid` query parameter
		delID := r.URL.Query().Get("delid")
		if delID == "" {
			http.Error(w, "Missing class ID to delete.", http.StatusBadRequest)
			log.Println("[ERROR] Missing 'delid' parameter in request")
			return
		}

		// Delete the class from the database
		_, err := db.Exec("DELETE FROM classes WHERE id = ?", delID)
		if err != nil {
			log.Printf("[ERROR] Failed to delete class with ID %s: %v", delID, err)
			http.Error(w, "Failed to delete class.", http.StatusInternalServerError)
			return
		}

		// Log successful deletion
		log.Printf("[INFO] Successfully deleted class with ID %s", delID)

		// Redirect back to the manage page
		http.Redirect(w, r, "/manage", http.StatusSeeOther)
	}
}
