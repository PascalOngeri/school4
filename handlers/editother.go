package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// EditOtherPaymentHandler handles updating "Other" payment details
func EditOtherPaymentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log the incoming request
		log.Printf("[INFO] Received request to edit 'Other' payment details: %s %s", r.Method, r.URL.Path)

		// Retrieve the ID based on the request method
		var id string
		if r.Method == "POST" {
			id = r.FormValue("id") // ID sent via POST form
		} else {
			id = r.URL.Query().Get("id") // ID sent via query parameters for GET requests
		}

		// Validate the ID
		if id == "" {
			log.Println("[ERROR] Missing payment ID")
			http.Error(w, "Missing payment ID", http.StatusBadRequest)
			return
		}

		// Handle POST request to update payment details
		if r.Method == "POST" {
			// Retrieve form values
			paymentType := r.FormValue("fname")
			term1 := r.FormValue("lname")
			term2 := r.FormValue("stuemail")
			term3 := r.FormValue("dob")

			// Convert terms to integers
			term1Int, err := strconv.Atoi(term1)
			if err != nil {
				log.Printf("[ERROR] Invalid value for Term 1: %v", err)
				http.Error(w, "Invalid value for Term 1", http.StatusBadRequest)
				return
			}
			term2Int, err := strconv.Atoi(term2)
			if err != nil {
				log.Printf("[ERROR] Invalid value for Term 2: %v", err)
				http.Error(w, "Invalid value for Term 2", http.StatusBadRequest)
				return
			}
			term3Int, err := strconv.Atoi(term3)
			if err != nil {
				log.Printf("[ERROR] Invalid value for Term 3: %v", err)
				http.Error(w, "Invalid value for Term 3", http.StatusBadRequest)
				return
			}

			// Update the record in the "other" table
			query := `
				UPDATE other 
				SET type = ?, t1 = ?, t2 = ?, t3 = ? 
				WHERE id = ?`
			_, err = db.Exec(query, paymentType, term1Int, term2Int, term3Int, id)
			if err != nil {
				log.Printf("[ERROR] Error updating payment details: %v", err)
				http.Error(w, "Error updating payment details", http.StatusInternalServerError)
				return
			}

			// Redirect to the success page or list of payments
			http.Redirect(w, r, "/updatepayment", http.StatusSeeOther)
			return
		}

		// Handle GET request to fetch current payment details
		var payment struct {
			ID                  string
			Type                string
			Term1, Term2, Term3 int
		}

		query := `SELECT id, type, t1, t2, t3 FROM other WHERE id = ?`
		err := db.QueryRow(query, id).Scan(&payment.ID, &payment.Type, &payment.Term1, &payment.Term2, &payment.Term3)
		if err != nil {
			log.Printf("[ERROR] Failed to fetch payment details for ID %s: %v", id, err)
			http.Error(w, "Failed to fetch payment details", http.StatusInternalServerError)
			return
		}

		// Parse the HTML template
		tmpl, err := template.ParseFiles(
			"templates/editO.html",
			"includes/header.html",
			"includes/sidebar.html",
			"includes/footer.html",
		)
		if err != nil {
			log.Printf("[ERROR] Error loading template: %v", err)
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}

		// Render the template with payment data
		err = tmpl.Execute(w, payment)
		if err != nil {
			log.Printf("[ERROR] Error rendering template: %v", err)
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}

		// Log successful page load
		log.Printf("[INFO] Successfully loaded 'Other' payment edit page for ID %s", id)
	}
}
