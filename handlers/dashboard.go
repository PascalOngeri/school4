package handlers

import (
	"html/template"
	"net/http"
)

// Dashboard handles the /dashboard route
func Dashboard(w http.ResponseWriter, r *http.Request) {
	// Read the role from the cookie
	cookie, err := r.Cookie("user_role")
	if err != nil {
		// Handle error (e.g., user not logged in or cookie expired)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Use the cookie value (role)
	role := cookie.Value

	// If role is "admin", show the dashboard
	if role == "admin" {
		// Parse templates
		tmpl, err := template.ParseFiles(
			"templates/dashboard.html",
			"includes/footer.html",
			"includes/header.html",
			"includes/sidebar.html",
		)
		if err != nil {
			// Handle the error properly
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Data to pass to the template
		data := map[string]interface{}{
			"Title": "Admin Dashboard", // Admin-specific title
		}

		// Execute the template and write to the response
		err = tmpl.Execute(w, data)
		if err != nil {
			// Handle the error properly
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if role == "user" {
		// If the role is "user", redirect to the parent section
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		// If role is not recognized, redirect to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
