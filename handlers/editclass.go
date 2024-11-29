package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

// EditClass handler for editing class details
func EditClass(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log the incoming request to edit class details
		log.Println("Received request to edit class details")

		// Fetch the 'editid' from URL query parameters
		editID := r.URL.Query().Get("editid")
		if editID == "" {
			log.Println("No edit ID provided in URL")
			http.Error(w, "Missing class ID", http.StatusBadRequest)
			return
		}

		// Initialize a Class object to store fetched class details
		var class Class

		// Fetch class details from the database
		log.Printf("Fetching class details for class ID: %s", editID)
		err := db.QueryRow("SELECT id, class, t1, t2, t3, fee FROM classes WHERE id = ?", editID).
			Scan(&class.ID, &class.Class, &class.T1, &class.T2, &class.T3, &class.Fee)

		if err != nil {
			log.Printf("Failed to fetch class details for ID %s: %v", editID, err)
			http.Error(w, "Failed to fetch class details.", http.StatusInternalServerError)
			return
		}
		log.Printf("Fetched class details: %+v", class)

		// Parse the template files
		log.Println("Parsing templates for editing class")
		tmpl, err := template.ParseFiles(
			"templates/edit-class.html",
			"includes/header.html",
			"includes/sidebar.html",
			"includes/footer.html",
		)
		if err != nil {
			log.Printf("Template parsing failed: %v", err)
			http.Error(w, "Failed to load page templates.", http.StatusInternalServerError)
			return
		}

		// Prepare data for the template execution
		data := map[string]interface{}{
			"Title": "Edit Class",
			"Class": class,
		}

		// Execute the template with the class data
		log.Println("Rendering edit class template")
		err = tmpl.Execute(w, data)
		if err != nil {
			log.Printf("Template execution failed: %v", err)
			http.Error(w, "Failed to render the page.", http.StatusInternalServerError)
			return
		}

		// Log successful rendering
		log.Println("Successfully rendered the edit class page")
	}
}
