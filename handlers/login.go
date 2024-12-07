package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"time"
)

// API structure
type API struct {
	Name  string
	Icon  string
	IName string
}

// LoginData structure for rendering the login page
type LoginData struct {
	Name     string
	Icon     string
	Username string
	Password string
}

// Get API details from the database
func getAPIDetails(db *sql.DB) (API, error) {
	var api API
	query := "SELECT name, icon, iname FROM api LIMIT 1"
	row := db.QueryRow(query)
	err := row.Scan(&api.Name, &api.Icon, &api.IName)
	if err != nil {
		log.Printf("Error fetching API details: %v", err)
		return api, err
	}
	return api, nil
}

// Render the login page
func renderLoginPage(w http.ResponseWriter, api API, username string) {
	loginData := LoginData{
		Name:     api.Name,
		Icon:     api.Icon,
		Username: username,
		Password: "",
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Error loading template: %v", err)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, loginData)
}

// HandleLogin handles login requests
func HandleLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == "POST" {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		var userID int
		var foundInAdmin bool
		var adm, phone, role string

		// Authenticate user in tbladmin
		queryAdmin := "SELECT ID, UserName FROM tbladmin WHERE UserName = ? AND Password = ?"
		err := db.QueryRow(queryAdmin, username, password).Scan(&userID, &username)
		if err == nil {
			foundInAdmin = true
			role = "admin"
		} else {
			// Authenticate user in registration table
			queryRegistration := "SELECT id, adm, username, phone, password FROM registration WHERE username = ? AND password = ?"
			err = db.QueryRow(queryRegistration, username, password).Scan(&userID, &adm, &username, &phone, &password)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			role = "user"
		}

		// Set the role in the cookie (without using JWT or sessions)
		http.SetCookie(w, &http.Cookie{
			Name:     "user_role",
			Value:    role,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(24 * time.Hour), // 1 day expiration
		})

		// Redirect based on role
		if foundInAdmin {
			// Admin should go to the dashboard
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		} else if role == "user" {
			// User (student) should go to the parent section
			http.Redirect(w, r, "/parent", http.StatusSeeOther)
		}
		return
	}

	// Render the login page for GET requests
	api, _ := getAPIDetails(db)
	renderLoginPage(w, api, "")
}

// DashboardHandler to handle the dashboard
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Read the role from the cookie
	cookie, err := r.Cookie("user_role")
	if err != nil {
		// Handle error (e.g., user not logged in or cookie expired)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Use the cookie value (role)
	role := cookie.Value
	if role == "admin" {
		// Render admin dashboard
		// Add the logic for rendering the admin dashboard
		http.ServeFile(w, r, "templates/admin_dashboard.html")
	} else {
		// Handle non-admin users (e.g., show a different dashboard)
		// Add the logic for rendering the user dashboard
		http.ServeFile(w, r, "templates/user_dashboard.html")
	}
}
