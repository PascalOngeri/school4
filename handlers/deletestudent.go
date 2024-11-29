package handlers

import (
	"database/sql"
	"log"
	"net/http"
)

// Function to delete a student
func DeleteStudent(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the student ID from the URL query string
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Missing student ID", http.StatusBadRequest)
			return
		}

		// Execute the DELETE query
		_, err := db.Exec("DELETE FROM registration WHERE id = ?", id)
		if err != nil {
			log.Printf("Error deleting user: %v", err)
			http.Error(w, "Error deleting user", http.StatusInternalServerError)
			return
		}

		// Redirect to the manage students page after successful deletion
		http.Redirect(w, r, "/managestudent", http.StatusSeeOther)
	}
}
