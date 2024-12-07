package handlers

import (
	"html/template"
	"log"
	"net/http"
)

// edelete handles the display and deletion of a class
func edelete(w http.ResponseWriter, r *http.Request) {
	// Log the incoming request
	log.Printf("[INFO] Received request to delete class: %s %s", r.Method, r.URL.Path)

	// Retrieve the class ID from the query parameter
	classID := r.URL.Query().Get("id")
	if classID == "" {
		log.Println("[ERROR] Missing 'id' query parameter")
		http.Error(w, "Missing class ID", http.StatusBadRequest)
		return
	}

	// Log the class ID to be deleted
	log.Printf("[INFO] Attempting to delete class with ID: %s", classID)

	// Parse the template files
	tmpl, err := template.ParseFiles("templates/edelete.html", "includes/footer.html", "includes/header.html", "includes/sidebar.html")
	if err != nil {
		log.Printf("[ERROR] Failed to parse template files: %v", err)
		http.Error(w, "Error loading page", http.StatusInternalServerError)
		return
	}

	// Data to pass to the template
	data := map[string]interface{}{
		"Title":   "Manage Class", // Example dynamic data
		"ClassID": classID,        // Pass the class ID to the template
	}

	// Execute the template and write to the response
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("[ERROR] Failed to execute template: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}

	// Log successful page load
	log.Printf("[INFO] Successfully loaded delete page for class ID: %s", classID)
}
