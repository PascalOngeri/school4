package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

// Send handles sending SMS without session management
func Send(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Handle POST request
	if r.Method == http.MethodPost {
		phone := r.FormValue("phone")
		message := r.FormValue("message")

		// Validate form input
		if phone == "" || message == "" {
			http.Error(w, "Phone number and message are required", http.StatusBadRequest)
			return
		}

		// Send the SMS
		err := SendSmsHandler(phone, message) // Use the renamed function
		if err != nil {
			log.Printf("Failed to send SMS: %v", err)
			http.Error(w, "Failed to send SMS", http.StatusInternalServerError)
			return
		}

		// Redirect after sending the SMS
		http.Redirect(w, r, "/send", http.StatusSeeOther)
		return
	}

	// Parse the template files
	tmpl, err := template.ParseFiles(
		"templates/send.html",
		"includes/footer.html",
		"includes/header.html",
		"includes/sidebar.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Data to pass to the template
	data := map[string]interface{}{
		"Title": "Send SMS",
	}

	// Execute the template and write to the response
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// SendSmsHandler sends an SMS to the provided phone number with the given message
func SendSms(phone string, message string) error {
	// Example SMS sending implementation - replace with your SMS service logic
	log.Printf("Sending SMS to %s: %s", phone, message)
	// Integrate your SMS API here and return any potential errors
	return nil // Replace with actual error handling from your SMS API
}
