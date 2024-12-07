package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

// Struct kwa data ya mwanafunzi
type SelectStudent struct {
	ID    int
	Adm   string
	Class string
	Fname string
	Mname string
	Lname string
	Fee   float64
	Email string
	Phone string
}

// Function ya kusimamia wanafunzi
func ManageStudent(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log incoming request for managing students
		log.Println("Received request to manage students")

		// Query to fetch students
		rows, err := db.Query("SELECT id, adm, class, fname, mname, lname, fee, email, phone FROM registration")
		if err != nil {
			// Log database query error
			log.Printf("Database query failed: %v", err)
			http.Error(w, "Database query failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var sele []SelectStudent
		// Process rows from the database
		for rows.Next() {
			var student SelectStudent
			if err := rows.Scan(&student.ID, &student.Adm, &student.Class, &student.Fname, &student.Mname, &student.Lname, &student.Fee, &student.Email, &student.Phone); err != nil {
				// Log error during row scanning
				log.Printf("Error scanning row: %v", err)
				http.Error(w, "Error scanning data from the database.", http.StatusInternalServerError)
				return
			}
			sele = append(sele, student)
		}

		// Check for errors during row iteration
		if err := rows.Err(); err != nil {
			// Log error during row iteration
			log.Printf("Error iterating rows: %v", err)
			http.Error(w, "Error reading students from the database.", http.StatusInternalServerError)
			return
		}

		// Log number of students fetched
		log.Printf("Successfully fetched %d students from the database", len(sele))

		// Parse template files
		tmpl, err := template.ParseFiles(
			"templates/managestudent.html",
			"includes/footer.html",
			"includes/header.html",
			"includes/sidebar.html",
		)
		if err != nil {
			// Log error during template parsing
			log.Printf("Template parsing failed: %v", err)
			http.Error(w, "Template parsing failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Pass data to the template
		data := map[string]interface{}{
			"Title":    "Manage Students",
			"Students": sele,
		}

		// Render the template
		if err := tmpl.Execute(w, data); err != nil {
			// Log error during template execution
			log.Printf("Template execution failed: %v", err)
			http.Error(w, "Template execution failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Log successful rendering
		log.Println("Successfully rendered the manage students page")
	}
}
