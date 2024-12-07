package handlers

import (
	"database/sql"
	"log"
	"net/http"
)

// DeleteBusHandler handles the deletion of a record from the "bus" table
func DeleteBusHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Log the incoming request
	log.Printf("[INFO] Received request to delete a bus record: %s %s", r.Method, r.URL.Path)

	// Get the "bdel" query parameter (the area to delete)
	area := r.URL.Query().Get("bdel")
	if area == "" {
		http.Error(w, "Missing area parameter", http.StatusBadRequest)
		log.Println("[ERROR] Missing 'bdel' parameter in request")
		return
	}

	// Prepare the DELETE SQL query
	query := "DELETE FROM bus WHERE area = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Printf("[ERROR] Failed to prepare DELETE query: %v", err)
		http.Error(w, "Failed to prepare delete statement", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the query
	result, err := stmt.Exec(area)
	if err != nil {
		log.Printf("[ERROR] Error executing DELETE query for area '%s': %v", area, err)
		http.Error(w, "Failed to delete record", http.StatusInternalServerError)
		return
	}

	// Check the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve rows affected for area '%s': %v", area, err)
		http.Error(w, "Failed to retrieve delete status", http.StatusInternalServerError)
		return
	}

	// Handle cases where no rows were affected
	if rowsAffected == 0 {
		log.Printf("[INFO] No record found to delete for area '%s'", area)
		http.Error(w, "No record found to delete", http.StatusNotFound)
		return
	}

	// Log success and redirect the user
	log.Printf("[INFO] Successfully deleted record for area '%s'", area)
	http.Redirect(w, r, "/updatepayment", http.StatusSeeOther)
}
