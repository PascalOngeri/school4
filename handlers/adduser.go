package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// ManageUser handles adding and deleting users
func ManageUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log the HTTP method and request path
		log.Printf("Handling request: %s %s", r.Method, r.URL.Path)

		if r.Method == http.MethodPost {
			// Parse form data
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Unable to parse form data.", http.StatusBadRequest)
				log.Printf("[ERROR] Form parsing failed: %v", err)
				return
			}

			action := r.FormValue("submit") // Determine the action from the form

			// Handle "Add" action
			if action == "Add" {
				AName := r.FormValue("adminname")
				mobno := r.FormValue("mobilenumber")
				email := r.FormValue("email")
				pass := r.FormValue("password")
				username := r.FormValue("username")

				// Validate input fields
				if AName == "" || mobno == "" || email == "" || pass == "" || username == "" {
					http.Error(w, "All fields are required.", http.StatusBadRequest)
					log.Println("[ERROR] Missing required fields for adding user")
					return
				}

				// Hash the password
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
				if err != nil {
					http.Error(w, "Error while securing password.", http.StatusInternalServerError)
					log.Printf("[ERROR] Password hashing failed: %v", err)
					return
				}

				// Insert the new user into the database
				query := `INSERT INTO tblAdmin (AdminName, Email, UserName, Password, MobileNumber) VALUES (?, ?, ?, ?, ?)`
				_, err = db.Exec(query, AName, email, username, hashedPassword, mobno)
				if err != nil {
					http.Error(w, "Failed to add user to the database.", http.StatusInternalServerError)
					log.Printf("[ERROR] Database insertion failed: %v", err)
					return
				}

				log.Printf("[INFO] User '%s' successfully added", username)
				http.Redirect(w, r, "/adduser", http.StatusSeeOther)
				return
			}

			// Handle "Delete" action
			if action == "Delete" {
				username := r.FormValue("username")

				if username == "" {
					http.Error(w, "Username is required for deletion.", http.StatusBadRequest)
					log.Println("[ERROR] Username is missing for deletion")
					return
				}

				query := `DELETE FROM tblAdmin WHERE UserName = ?`
				result, err := db.Exec(query, username)
				if err != nil {
					http.Error(w, "Failed to delete user from the database.", http.StatusInternalServerError)
					log.Printf("[ERROR] Database deletion failed: %v", err)
					return
				}

				rowsAffected, _ := result.RowsAffected()
				if rowsAffected == 0 {
					http.Error(w, "No user found with the provided username.", http.StatusNotFound)
					log.Printf("[WARNING] No user found with username '%s' for deletion", username)
					return
				}

				log.Printf("[INFO] User '%s' successfully deleted", username)
				http.Redirect(w, r, "/adduser", http.StatusSeeOther)
				return
			}

			// Log if the action is not recognized
			log.Printf("[ERROR] Unknown action '%s' received", action)
			http.Error(w, "Unknown action.", http.StatusBadRequest)
			return
		}

		// Handle GET request to render the form
		tmpl, err := template.ParseFiles(
			"templates/adduser.html",
			"includes/header.html",
			"includes/sidebar.html",
			"includes/footer.html",
		)
		if err != nil {
			http.Error(w, "Failed to load the page.", http.StatusInternalServerError)
			log.Printf("[ERROR] Template parsing failed: %v", err)
			return
		}

		// Render the form
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, "Failed to render the page.", http.StatusInternalServerError)
			log.Printf("[ERROR] Template execution failed: %v", err)
			return
		}

		log.Println("[INFO] User management page rendered successfully")
	}
}
