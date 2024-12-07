package handlers

import (
	"database/sql"
	"log"
	"net/http"
)

// DeleteOtherHandler handles the deletion of a record from the "other" table
func DeleteOtherHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Log the incoming request
	log.Printf("[INFO] Received request to delete 'other' record: %s %s", r.Method, r.URL.Path)

	// Retrieve the `otherdel` query parameter
	paymentName := r.URL.Query().Get("otherdel")
	if paymentName == "" {
		http.Error(w, "Missing payment name", http.StatusBadRequest)
		log.Println("[ERROR] Missing 'otherdel' parameter in request")
		return
	}

	// Prepare the DELETE statement to remove the record
	query := "DELETE FROM other WHERE type = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println("[ERROR] Error preparing DELETE query:", err)
		http.Error(w, "Failed to prepare delete statement", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the DELETE statement
	result, err := stmt.Exec(paymentName)
	if err != nil {
		log.Println("[ERROR] Error executing DELETE query:", err)
		http.Error(w, "Failed to delete record", http.StatusInternalServerError)
		return
	}

	// Check how many rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("[ERROR] Error getting rows affected:", err)
		http.Error(w, "Failed to retrieve delete status", http.StatusInternalServerError)
		return
	}

	// Respond to the client
	if rowsAffected == 0 {
		http.Error(w, "No record found to delete", http.StatusNotFound)
		log.Printf("[INFO] No record found for deletion with type '%s'", paymentName)
		return
	}

	// Log successful deletion
	log.Printf("[INFO] Successfully deleted record with type '%s'", paymentName)

	// Redirect to the update payment page after successful deletion
	http.Redirect(w, r, "/updatepayment", http.StatusSeeOther)
}
