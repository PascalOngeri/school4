package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

// EditCompulsoryPaymentHandler handles the update of compulsory payments without session handling
func EditCompulsoryPaymentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log the incoming request
		log.Println("Received request to edit compulsory payment")

		// Fetch the ID of the compulsory payment to be edited from the URL
		var id string
		if r.Method == "POST" {
			id = r.FormValue("id") // Fetch ID from the form for POST requests
		} else {
			id = r.URL.Query().Get("id") // Fetch ID from URL query for GET requests
		}

		// If it's a POST request, update the payment details in the database
		if r.Method == "POST" {
			// Get form values for updating payment details
			paymentName := r.FormValue("fname")
			className := r.FormValue("mname")
			term1 := r.FormValue("lname")
			term2 := r.FormValue("stuemail")
			term3 := r.FormValue("dob")

			// Combine term values to calculate the total amount (assuming terms are integers)
			amount := term1 + term2 + term3

			// Get the current payment details to be updated
			var currentAmount, currentTerm1, currentTerm2, currentTerm3 int
			var currentClassName string
			query := `SELECT amount, term1, term2, term3, form FROM feepay WHERE id = ?`
			err := db.QueryRow(query, id).Scan(&currentAmount, &currentTerm1, &currentTerm2, &currentTerm3, &currentClassName)
			if err != nil {
				log.Printf("Error fetching current payment details: %v", err)
				http.Error(w, "Error fetching payment details", http.StatusInternalServerError)
				return
			}

			// Subtract old payment data from the classes table
			log.Printf("Updating classes table by subtracting old payment details for class: %s", currentClassName)
			_, err = db.Exec(`
				UPDATE classes 
				SET fee = fee - ?, t1 = t1 - ?, t2 = t2 - ?, t3 = t3 - ? 
				WHERE class = ?`, currentAmount, currentTerm1, currentTerm2, currentTerm3, currentClassName)
			if err != nil {
				log.Printf("Error updating classes table with old data: %v", err)
				http.Error(w, "Error updating classes table", http.StatusInternalServerError)
				return
			}

			// Update the payment details in the feepay table
			log.Printf("Updating payment details for payment ID: %s", id)
			_, err = db.Exec(`
				UPDATE feepay 
				SET paymentname = ?, form = ?, term1 = ?, term2 = ?, term3 = ?, amount = ? 
				WHERE id = ?`, paymentName, className, term1, term2, term3, amount, id)
			if err != nil {
				log.Printf("Error updating compulsory payment: %v", err)
				http.Error(w, "Error updating compulsory payment", http.StatusInternalServerError)
				return
			}

			// Add the new payment data to the classes table
			log.Printf("Adding new payment data to classes table for class: %s", className)
			_, err = db.Exec(`
				UPDATE classes 
				SET fee = fee + ?, t1 = t1 + ?, t2 = t2 + ?, t3 = t3 + ? 
				WHERE class = ?`, amount, term1, term2, term3, className)
			if err != nil {
				log.Printf("Error updating classes table with new data: %v", err)
				http.Error(w, "Error updating classes table with new data", http.StatusInternalServerError)
				return
			}

			// Redirect to the update payment page after successful update
			log.Println("Payment successfully updated, redirecting to /updatepayment")
			http.Redirect(w, r, "/updatepayment", http.StatusSeeOther)
			return
		}

		// If it's a GET request, fetch the payment details to display for editing
		log.Printf("Fetching payment details for editing, payment ID: %s", id)
		var payment Payment
		query := `SELECT id, paymentname, form, term1, term2, term3, amount FROM feepay WHERE id = ?`
		err := db.QueryRow(query, id).Scan(&payment.ID, &payment.PaymentName, &payment.ClassName, &payment.Term1, &payment.Term2, &payment.Term3, &payment.Amount)
		if err != nil {
			log.Printf("Error fetching payment details from database: %v", err)
			http.Error(w, "Error fetching payment details", http.StatusInternalServerError)
			return
		}

		// Parse and render the template with the payment details
		log.Println("Rendering edit compulsory payment template")
		tmpl, err := template.ParseFiles(
			"templates/editC.html",
			"includes/header.html",
			"includes/sidebar.html",
			"includes/footer.html",
		)
		if err != nil {
			log.Printf("Error loading templates: %v", err) // Log detailed error
			http.Error(w, "Error loading templates", http.StatusInternalServerError)
			return
		}

		// Pass the payment struct to the template for rendering
		err = tmpl.Execute(w, payment)
		if err != nil {
			log.Printf("Error rendering template: %v", err)
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}

		// Log the successful rendering of the template
		log.Println("Successfully rendered edit payment template")
	}
}
