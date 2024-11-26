package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func send(w http.ResponseWriter, r *http.Request) {
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
	// Parse the template files
	tmpl, err := template.ParseFiles("templates/send.html", "includes/footer.html", "includes/header.html", "includes/sidebar.html")
	if err != nil {
		// Handle the error properly, e.g., by returning a 500 status
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Data to pass to the template
	data := map[string]interface{}{
		"Title": "Manage Class", // Example dynamic data
	}

	// Execute the template and write to the response
	err = tmpl.Execute(w, data)
	if err != nil {
		// Handle the error properly
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
