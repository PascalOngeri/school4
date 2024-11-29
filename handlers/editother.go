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
		// Log the incoming request to edit "Other" payment details
		log.Println("Received request to edit 'Other' payment details")

		// Retrieve the ID based on the request method (POST or GET)
		var id string
		if r.Method == "POST" {
			id = r.FormValue("id") // ID sent via POST form
		} else {
			id = r.URL.Query().Get("id") // ID sent via query parameters for GET requests
		}

		// Validate the ID parameter
		if id == "" {
			log.Println("No ID provided in the request")
			http.Error(w, "Missing payment ID", http.StatusBadRequest)
			return
		}

		// Handle POST request: Update payment details in the "other" table
		if r.Method == "POST" {
			log.Println("Handling form submission for 'Other' payment update")

			// Get values from the form
			paymentType := r.FormValue("fname")
			term1 := r.FormValue("lname")
			term2 := r.FormValue("stuemail")
			term3 := r.FormValue("dob")

			// Convert terms to integers
			term1Int, err := strconv.Atoi(term1)
			if err != nil {
				log.Printf("Invalid value for Term 1: %v", err)
				http.Error(w, "Invalid value for Term 1", http.StatusBadRequest)
				return
			}
			term2Int, err := strconv.Atoi(term2)
			if err != nil {
				log.Printf("Invalid value for Term 2: %v", err)
				http.Error(w, "Invalid value for Term 2", http.StatusBadRequest)
				return
			}
			term3Int, err := strconv.Atoi(term3)
			if err != nil {
				log.Printf("Invalid value for Term 3: %v", err)
				http.Error(w, "Invalid value for Term 3", http.StatusBadRequest)
				return
			}

			// Update the record in the "other" table
			log.Printf("Updating 'Other' payment details for ID: %s", id)
			query := `
				UPDATE other 
				SET type = ?, t1 = ?, t2 = ?, t3 = ? 
				WHERE id = ?`
			_, err = db.Exec(query, paymentType, term1Int, term2Int, term3Int, id)
			if err != nil {
				log.Printf("Error updating payment for ID %s: %v", id, err)
				http.Error(w, "Error updating payment details", http.StatusInternalServerError)
				return
			}

			// Redirect to the success page or list of payments
			log.Println("Redirecting to /updatepayment after successful update")
			http.Redirect(w, r, "/updatepayment", http.StatusSeeOther)
			return
		}

		// Handle GET request: Fetch the current payment details to display in the form
		log.Printf("Fetching payment details for 'Other' payment ID: %s", id)
		var payment struct {
			ID                  string
			Type                string
			Term1, Term2, Term3 int
		}

		// Query the database for the existing payment details
		query := `SELECT id, type, t1, t2, t3 FROM other WHERE id = ?`
		err := db.QueryRow(query, id).Scan(&payment.ID, &payment.Type, &payment.Term1, &payment.Term2, &payment.Term3)
		if err != nil {
			log.Printf("Error fetching payment details for ID %s: %v", id, err)
			http.Error(w, "Error fetching payment details", http.StatusInternalServerError)
			return
		}

		// Log the fetched payment details
		log.Printf("Fetched payment details: %+v", payment)

		// Parse the HTML template files
		log.Println("Parsing the edit payment template")
		tmpl, err := template.ParseFiles(
			"templates/editO.html",
			"includes/header.html",
			"includes/sidebar.html",
			"includes/footer.html",
		)
		if err != nil {
			log.Printf("Error loading template files: %v", err)
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}

		// Render the template with the payment data
		log.Println("Rendering the edit payment template with payment data")
		err = tmpl.Execute(w, payment)
		if err != nil {
			log.Printf("Error rendering template: %v", err)
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}

		// Log successful rendering
		log.Println("Successfully rendered the edit 'Other' payment page")
	}
}
