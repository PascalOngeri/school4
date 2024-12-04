package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func edelete(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		// If the cookie is not found, redirect to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Validate JWT token from the cookie
	claims, err := ValidateJWT(cookie.Value)
	if err != nil {
		// If the token is invalid or expired, redirect to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if claims.Role != "admin" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// Log authenticated user info for debugging
	log.Printf("Authenticated user: %s, Role: %s", claims.Username, claims.Role)

	// Parse the template files
	tmpl, err := template.ParseFiles("templates/edelete.html", "includes/footer.html", "includes/header.html", "includes/sidebar.html")
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
