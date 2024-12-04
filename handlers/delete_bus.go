package handlers

import (
	"database/sql"
	"log"
	"net/http"
)

// DeleteBusHandler handles the deletion of a record from the "bus" table
func DeleteBusHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		// If the cookie is not found, redirect to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Validate JWT token from the cookie
	claims, err := ValidateJWT(cookie.Value)
	if err != nil {
		// If the token is invalid or expired, redirect to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if claims.Role != "admin" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// Log authenticated user info for debugging
	log.Printf("Authenticated user: %s, Role: %s", claims.Username, claims.Role)

	area := r.URL.Query().Get("bdel")
	if area == "" {
		http.Error(w, "Missing area", http.StatusBadRequest)
		return
	}

	// Prepare the delete statement
	query := "DELETE FROM bus WHERE area = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println("Error preparing query:", err)
		http.Error(w, "Failed to prepare delete statement", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Execute the delete statement
	result, err := stmt.Exec(area)
	if err != nil {
		log.Println("Error executing delete:", err)
		http.Error(w, "Failed to delete record", http.StatusInternalServerError)
		return
	}

	// Check how many rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		http.Error(w, "Failed to retrieve delete status", http.StatusInternalServerError)
		return
	}

	// Respond to the client
	if rowsAffected == 0 {
		http.Error(w, "No record found to delete", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/updatepayment", http.StatusSeeOther)
}
