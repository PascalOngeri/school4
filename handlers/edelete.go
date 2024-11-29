package handlers

import (
	"html/template"
	"log"
	"net/http"
)

// edelete renders the delete page without session management
func edelete(w http.ResponseWriter, r *http.Request) {
	// Log the incoming request to render the delete page
	log.Println("Received request to render delete student page")

	// Parse the template files
	tmpl, err := template.ParseFiles("templates/edelete.html", "includes/footer.html", "includes/header.html", "includes/sidebar.html")
	if err != nil {
		// Log the error and return an internal server error response
		log.Printf("Error parsing templates: %v", err)
		http.Error(w, "Error parsing templates", http.StatusInternalServerError)
		return
	}

	// Data to pass to the template
	data := map[string]interface{}{
		"Title": "Delete Student", // Dynamic data for the title
	}

	// Execute the template and write to the response
	err = tmpl.Execute(w, data)
	if err != nil {
		// Log the error and return an internal server error response
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}

	// Log successful rendering of the page
	log.Println("Successfully rendered delete student page")
}
