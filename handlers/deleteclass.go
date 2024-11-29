package handlers

import (
	"database/sql"
	"log"
	"net/http"
)

// DeleteClass deletes a class by ID
func DeleteClass(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("DeleteClass handler invoked.")

		// Retrieve the `delid` query parameter
		delID := r.URL.Query().Get("delid")
		if delID == "" {
			log.Println("Missing 'delid' query parameter.")
			http.Error(w, "Missing class ID to delete.", http.StatusBadRequest)
			return
		}
		log.Printf("Received request to delete class with ID: %s", delID)

		// Attempt to delete the class from the database
		result, err := db.Exec("DELETE FROM classes WHERE id = ?", delID)
		if err != nil {
			log.Printf("Failed to execute DELETE query for class ID %s: %v", delID, err)
			http.Error(w, "Failed to delete class.", http.StatusInternalServerError)
			return
		}

		// Check if any rows were affected
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Printf("Error retrieving rows affected for class ID %s: %v", delID, err)
			http.Error(w, "Failed to delete class.", http.StatusInternalServerError)
			return
		}
		if rowsAffected == 0 {
			log.Printf("No class found with ID: %s", delID)
			http.Error(w, "No class found with the provided ID.", http.StatusNotFound)
			return
		}

		log.Printf("Successfully deleted class with ID: %s. Rows affected: %d", delID, rowsAffected)

		// Redirect back to the manage page
		log.Println("Redirecting to /manage.")
		http.Redirect(w, r, "/manage", http.StatusSeeOther)
	}
}
