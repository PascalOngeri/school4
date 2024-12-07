package handlers

import (
	"database/sql"
	"log"
	"net/http"
)

// DeleteStudent deletes a student from the "registration" table by ID
func DeleteStudent(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log the incoming request
		log.Printf("[INFO] Received request to delete student: %s %s", r.Method, r.URL.Path)

		// Retrieve the student ID from the query parameter
		id := r.URL.Query().Get("id")
		if id == "" {
			log.Println("[ERROR] Missing 'id' query parameter")
			http.Error(w, "Missing student ID", http.StatusBadRequest)
			return
		}

		// Log the student ID to be deleted
		log.Printf("[INFO] Attempting to delete student with ID: %s", id)

		// Prepare the DELETE statement
		_, err := db.Exec("DELETE FROM registration WHERE id = ?", id)
		if err != nil {
			log.Printf("[ERROR] Failed to delete student with ID %s: %v", id, err)
			http.Error(w, "Error deleting student", http.StatusInternalServerError)
			return
		}

		// Log successful deletion
		log.Printf("[INFO] Successfully deleted student with ID: %s", id)

		// Redirect to the manage student page after successful deletion
		http.Redirect(w, r, "/managestudent", http.StatusSeeOther)
	}
}
