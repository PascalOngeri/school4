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
		// Fetch the ID of the compulsory payment to be edited from the URL

		var id string
		if r.Method == "POST" {
			id = r.FormValue("id") // Fetch ID from the form
		} else {
			id = r.URL.Query().Get("id") // Fetch ID from URL query for GET requests
		}

		// If it's a POST request, update the payment details in the database
		if r.Method == "POST" {

			id := r.FormValue("id")
			paymentName := r.FormValue("fname")
			className := r.FormValue("mname")
			term1 := r.FormValue("lname")
			term2 := r.FormValue("stuemail")
			term3 := r.FormValue("dob")

			// Convert term1, term2, and term3 to integers
			amount := term1 + term2 + term3

			// Get the current payment details
			var currentAmount, currentTerm1, currentTerm2, currentTerm3 int
			var currentClassName string
			query := `SELECT amount, term1, term2, term3, form FROM feepay WHERE id = ?`
			err := db.QueryRow(query, id).Scan(&currentAmount, &currentTerm1, &currentTerm2, &currentTerm3, &currentClassName)
			if err != nil {
				http.Error(w, "Error fetching payment details", http.StatusInternalServerError)
				log.Printf("Error fetching payment details: %v", err)
				return
			}

			// First, subtract the old payment data from the classes table
			_, err = db.Exec(`
				UPDATE classes 
				SET fee = fee - ?, t1 = t1 - ?, t2 = t2 - ?, t3 = t3 - ? 
				WHERE class = ?`, currentAmount, currentTerm1, currentTerm2, currentTerm3, currentClassName)
			if err != nil {
				http.Error(w, "Error updating classes table", http.StatusInternalServerError)
				log.Printf("Error updating classes table: %v", err)
				return
			}

			// Now, update the payment details in the feepay table
			_, err = db.Exec(`
				UPDATE feepay 
				SET paymentname = ?, form = ?, term1 = ?, term2 = ?, term3 = ?, amount = ? 
				WHERE id = ?`, paymentName, className, term1, term2, term3, amount, id)
			if err != nil {
				http.Error(w, "Error updating compulsory payment", http.StatusInternalServerError)
				log.Printf("Error updating compulsory payment: %v", err)
				return
			}

			// Finally, add the new payment data to the classes table
			_, err = db.Exec(`
				UPDATE classes 
				SET fee = fee + ?, t1 = t1 + ?, t2 = t2 + ?, t3 = t3 + ? 
				WHERE class = ?`, amount, term1, term2, term3, className)
			if err != nil {
				http.Error(w, "Error updating classes table", http.StatusInternalServerError)
				log.Printf("Error updating classes table with new data: %v", err)
				return
			}

			// Redirect to a success page or another handler
			http.Redirect(w, r, "/updatepayment", http.StatusSeeOther)
			return
		}

		// If it's a GET request, fetch the payment details to edit
		var payment Payment
		query := `SELECT id, paymentname, form, term1, term2, term3, amount FROM feepay WHERE id = ?`
		err = db.QueryRow(query, id).Scan(&payment.ID, &payment.PaymentName, &payment.ClassName, &payment.Term1, &payment.Term2, &payment.Term3, &payment.Amount)

		// Parse and render the template with the payment details
		// Parse template files
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
	}
}
