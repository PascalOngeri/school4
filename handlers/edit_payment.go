package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

// EditCompulsoryPaymentHandler handles the update of compulsory payments
func EditCompulsoryPaymentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log the incoming request
		log.Printf("[INFO] Received request to edit compulsory payment: %s %s", r.Method, r.URL.Path)

		// Fetch the ID of the compulsory payment to be edited from the URL or form
		var id string
		if r.Method == "POST" {
			id = r.FormValue("id") // Fetch ID from the form
		} else {
			id = r.URL.Query().Get("id") // Fetch ID from URL query for GET requests
		}

		if id == "" {
			log.Println("[ERROR] Missing compulsory payment ID")
			http.Error(w, "Missing compulsory payment ID", http.StatusBadRequest)
			return
		}

		// If it's a POST request, update the payment details in the database
		if r.Method == "POST" {
			paymentName := r.FormValue("fname")
			className := r.FormValue("mname")
			term1 := r.FormValue("lname")
			term2 := r.FormValue("stuemail")
			term3 := r.FormValue("dob")

			// Calculate the total amount (term1 + term2 + term3) as integers
			amount := term1 + term2 + term3 // You may need to parse these terms into integers

			// Get the current payment details
			var currentAmount, currentTerm1, currentTerm2, currentTerm3 int
			var currentClassName string
			query := `SELECT amount, term1, term2, term3, form FROM feepay WHERE id = ?`
			err := db.QueryRow(query, id).Scan(&currentAmount, &currentTerm1, &currentTerm2, &currentTerm3, &currentClassName)
			if err != nil {
				log.Printf("[ERROR] Failed to fetch payment details for ID %s: %v", id, err)
				http.Error(w, "Error fetching payment details", http.StatusInternalServerError)
				return
			}

			// Subtract the old payment data from the classes table
			_, err = db.Exec(`
				UPDATE classes 
				SET fee = fee - ?, t1 = t1 - ?, t2 = t2 - ?, t3 = t3 - ? 
				WHERE class = ?`, currentAmount, currentTerm1, currentTerm2, currentTerm3, currentClassName)
			if err != nil {
				log.Printf("[ERROR] Failed to update classes table after subtraction for class %s: %v", currentClassName, err)
				http.Error(w, "Error updating classes table", http.StatusInternalServerError)
				return
			}

			// Update the payment details in the feepay table
			_, err = db.Exec(`
				UPDATE feepay 
				SET paymentname = ?, form = ?, term1 = ?, term2 = ?, term3 = ?, amount = ? 
				WHERE id = ?`, paymentName, className, term1, term2, term3, amount, id)
			if err != nil {
				log.Printf("[ERROR] Failed to update compulsory payment with ID %s: %v", id, err)
				http.Error(w, "Error updating compulsory payment", http.StatusInternalServerError)
				return
			}

			// Add the new payment data to the classes table
			_, err = db.Exec(`
				UPDATE classes 
				SET fee = fee + ?, t1 = t1 + ?, t2 = t2 + ?, t3 = t3 + ? 
				WHERE class = ?`, amount, term1, term2, term3, className)
			if err != nil {
				log.Printf("[ERROR] Failed to update classes table with new data for class %s: %v", className, err)
				http.Error(w, "Error updating classes table", http.StatusInternalServerError)
				return
			}

			// Redirect to the success page or another handler
			log.Printf("[INFO] Successfully updated compulsory payment with ID %s", id)
			http.Redirect(w, r, "/updatepayment", http.StatusSeeOther)
			return
		}

		// If it's a GET request, fetch the payment details to edit
		var payment Payment
		query := `SELECT id, paymentname, form, term1, term2, term3, amount FROM feepay WHERE id = ?`
		err := db.QueryRow(query, id).Scan(&payment.ID, &payment.PaymentName, &payment.ClassName, &payment.Term1, &payment.Term2, &payment.Term3, &payment.Amount)
		if err != nil {
			log.Printf("[ERROR] Failed to fetch payment details for ID %s: %v", id, err)
			http.Error(w, "Error fetching payment details", http.StatusInternalServerError)
			return
		}

		// Parse and render the template with the payment details
		tmpl, err := template.ParseFiles(
			"templates/editC.html",
			"includes/header.html",
			"includes/sidebar.html",
			"includes/footer.html",
		)
		if err != nil {
			log.Printf("[ERROR] Failed to load template files: %v", err)
			http.Error(w, "Error loading templates", http.StatusInternalServerError)
			return
		}

		// Pass the payment struct to the template for rendering
		err = tmpl.Execute(w, payment)
		if err != nil {
			log.Printf("[ERROR] Failed to render template: %v", err)
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}

		// Log successful page load
		log.Printf("[INFO] Successfully loaded payment edit page for ID %s", id)
	}
}
