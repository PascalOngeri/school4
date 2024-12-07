package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

// AddClass handles adding a new class to the database
func AddClass(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			log.Printf("ERROR: Unable to parse form data: %v", err)
			http.Error(w, "Bad Request: Unable to parse form data", http.StatusBadRequest)
			return
		}

		// Get the class name from the form
		className := r.FormValue("cname")
		if className == "" {
			log.Printf("ERROR: Class name not provided in form submission")
			http.Error(w, "Class name is required", http.StatusBadRequest)
			return
		}

		// Insert the class into the database
		_, err = db.Exec("INSERT INTO classes (class, fee, t1, t2, t3) VALUES (?, ?, ?, ?, ?)", className, 0, 0, 0, 0)
		if err != nil {
			log.Printf("ERROR: Failed to insert class '%s' into database: %v", className, err)
			http.Error(w, "Internal Server Error: Failed to add class", http.StatusInternalServerError)
			return
		}

		log.Printf("INFO: Successfully added class '%s' to the database", className)

		// Redirect to the form or a confirmation page
		http.Redirect(w, r, "/addclass", http.StatusSeeOther)
		return
	}

	// Render the add class form for GET requests
	tmpl, err := template.ParseFiles(
		"templates/addclass.html",
		"includes/header.html",
		"includes/sidebar.html",
		"includes/footer.html",
	)
	if err != nil {
		log.Printf("ERROR: Failed to parse template files for AddClass page: %v", err)
		http.Error(w, "Internal Server Error: Unable to load page", http.StatusInternalServerError)
		return
	}

	// Execute the template
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("ERROR: Failed to render AddClass page: %v", err)
		http.Error(w, "Internal Server Error: Unable to render page", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: Successfully rendered AddClass page")
}
