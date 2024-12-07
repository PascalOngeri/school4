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
		// Log the incoming request
		log.Printf("[INFO] Received request to edit class details: %s %s", r.Method, r.URL.Path)

		// Fetch the edit ID from the URL query
		editID := r.URL.Query().Get("editid")
		if editID == "" {
			log.Println("[ERROR] Missing class ID (editid)")
			http.Error(w, "Missing class ID", http.StatusBadRequest)
			return
		}

		// Fetch class details from the database
		var class Class
		err := db.QueryRow("SELECT id, class, t1, t2, t3, fee FROM classes WHERE id = ?", editID).
			Scan(&class.ID, &class.Class, &class.T1, &class.T2, &class.T3, &class.Fee)
		if err != nil {
			log.Printf("[ERROR] Failed to fetch class details for ID %s: %v", editID, err)
			http.Error(w, "Failed to fetch class details", http.StatusInternalServerError)
			return
		}

		// Parse the template files
		tmpl, err := template.ParseFiles(
			"templates/edit-class.html",
			"includes/header.html",
			"includes/sidebar.html",
			"includes/footer.html",
		)
		if err != nil {
			log.Printf("[ERROR] Failed to parse template files: %v", err)
			http.Error(w, "Failed to load page templates", http.StatusInternalServerError)
			return
		}

		// Prepare data to pass to the template
		data := map[string]interface{}{
			"Title": "Edit Class",
			"Class": class,
		}

		// Execute the template and send the response
		err = tmpl.Execute(w, data)
		if err != nil {
			log.Printf("[ERROR] Failed to execute template: %v", err)
			http.Error(w, "Failed to render the page", http.StatusInternalServerError)
			return
		}

		// Log successful page load
		log.Printf("[INFO] Successfully loaded class edit page for ID %s", editID)
	}
}
