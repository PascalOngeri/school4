package handlers

import (
	"database/sql"
	"log"
	"net/http"
)

// Function ya kufuta mwanafunzi
func DeleteStudent(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

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
		id := r.URL.Query().Get("id")
		_, err = db.Exec("DELETE FROM registration WHERE id = ?", id)

		if err != nil {
			http.Error(w, "Error deleting user", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/managestudent", http.StatusSeeOther)
	}
}
