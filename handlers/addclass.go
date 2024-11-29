package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)



func AddClass(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Println("AddClass handler invoked")

	// Uncomment this section if sessions are implemented
	/*
		session, err := store.Get(r, "store")
		if err != nil {
			log.Printf("Failed to retrieve session: %v", err)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}

		// Check if user is logged in
		if session.Values["sturecmsaid"] == nil {
			log.Println("Unauthorized access attempt: user not logged in")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	*/

	if r.Method == http.MethodPost {
		log.Println("Processing POST request to add a class")

		// Parse the form data
		if err := r.ParseForm(); err != nil {
			log.Printf("Form parsing failed: %v", err)
			http.Error(w, "Unable to parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Get the class name from the form
		className := r.FormValue("cname")
		if className == "" {
			log.Println("Class name is missing in the form submission")
			http.Error(w, "Class name is required", http.StatusBadRequest)
			return
		}

		log.Printf("Class name received: %s", className)

		// Insert the class into the database
		query := "INSERT INTO classes (class,fee,t1,t2,t3) VALUES (?,?,?,?,?)"
		log.Printf("Executing query: %s with values (%s, 0, 0, 0, 0)", query, className)
		_, err := db.Exec(query, className, 0, 0, 0, 0)
		if err != nil {
			log.Printf("Database insert failed for class %s: %v", className, err)
			http.Error(w, "Failed to add class: "+err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Class %s successfully added to the database", className)

		// Redirect to a confirmation page or reload the form with a success message
		http.Redirect(w, r, "/addclass", http.StatusSeeOther)
		return
	}

	log.Println("Rendering addclass.html template")

	// Render the template
	tmpl, err := template.ParseFiles(
		"templates/addclass.html",
		"includes/header.html",
		"includes/sidebar.html",
		"includes/footer.html",
	)
	if err != nil {
		log.Printf("Template parsing failed: %v", err)
		http.Error(w, "Template parsing failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template
	log.Println("Executing template with no data")
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Template execution failed: %v", err)
		http.Error(w, "Template execution failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Template successfully rendered")
}
