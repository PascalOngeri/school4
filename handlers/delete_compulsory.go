package handlers

import (
	"database/sql"
	"log"
	"net/http"
)

// DeleteCompulsoryHandler handles the deletion of a record from the "feepay" table
func DeleteCompulsoryHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Log the incoming request
	log.Printf("[INFO] Received request to delete a payment record: %s %s", r.Method, r.URL.Path)

	// Get the payment ID from query parameters
	paymentID := r.URL.Query().Get("delid")
	if paymentID == "" {
		http.Error(w, "Missing payment ID", http.StatusBadRequest)
		log.Println("[ERROR] Missing 'delid' parameter in request")
		return
	}

	// Prepare the DELETE SQL query
	query := "DELETE FROM feepay WHERE id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Printf("[ERROR] Failed to prepare DELETE query: %v", err)
		http.Error(w, "Failed to prepare delete statement", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the query
	result, err := stmt.Exec(paymentID)
	if err != nil {
		log.Printf("[ERROR] Error executing DELETE query for payment ID '%s': %v", paymentID, err)
		http.Error(w, "Failed to delete record", http.StatusInternalServerError)
		return
	}

	// Check the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve rows affected for payment ID '%s': %v", paymentID, err)
		http.Error(w, "Failed to retrieve delete status", http.StatusInternalServerError)
		return
	}

	// Handle cases where no rows were affected
	if rowsAffected == 0 {
		log.Printf("[INFO] No record found to delete for payment ID '%s'", paymentID)
		http.Error(w, "No record found to delete", http.StatusNotFound)
		return
	}

	// Log success and redirect the user
	log.Printf("[INFO] Successfully deleted payment record with ID '%s'", paymentID)
	http.Redirect(w, r, "/updatepayment", http.StatusSeeOther)
}
