package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

type API struct {
	Name  string
	Icon  string
	IName string
}

type LoginData struct {
	Name     string
	Icon     string
	Username string
	Password string
	Remember bool
}

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

func renderLoginPage(w http.ResponseWriter, api API, username string) {
	loginData := LoginData{
		Name:     api.Name,
		Icon:     api.Icon,
		Username: username,
		Password: "",
		Remember: false,
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Error loading template: %v", err)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, loginData)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}

func HandleLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == "POST" {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")
		remember := r.FormValue("remember") == "on"

		var userID int
		var foundInAdmin bool
		var adm, phone string

		// Attempt to find the user in the admin table
		queryAdmin := "SELECT ID, UserName FROM tbladmin WHERE UserName = ? AND Password = ?"
		err := db.QueryRow(queryAdmin, username, password).Scan(&userID, &username)
		if err == nil {
			log.Printf("User found in admin table: %s", username)
			foundInAdmin = true
		} else {
			// Attempt to find the user in the registration table
			queryRegistration := "SELECT id, adm, username, phone, password FROM registration WHERE username = ? AND password = ?"
			err = db.QueryRow(queryRegistration, username, password).Scan(&userID, &adm, &username, &phone, &password)
			if err != nil {
				log.Printf("Login failed for username: %s, error: %v", username, err)
				http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
				return
			}
			log.Printf("User found in registration table: %s", username)
		}

		// Log the successful login attempt
		log.Printf("Successful login attempt for user: %s", username)

		// Set cookies for the user
		http.SetCookie(w, &http.Cookie{Name: "user_login", Value: username, Path: "/", MaxAge: 86400})
		if remember {
			http.SetCookie(w, &http.Cookie{Name: "userpassword", Value: password, Path: "/", MaxAge: 86400})
		}

		// Redirect based on the user type
		if foundInAdmin {
			log.Printf("Redirecting to admin dashboard for user: %s", username)
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		} else {
			log.Printf("Redirecting to parent page for user: %s", username)
			http.Redirect(w, r, "/parent", http.StatusSeeOther)
		}
		return
	}

	// Fetch API details
	api, err := getAPIDetails(db)
	if err != nil {
		log.Printf("Error fetching API details: %v", err)
		http.Error(w, "Error fetching API details", http.StatusInternalServerError)
		return
	}

	// Get the username from the query parameters or set to empty string if not found
	username := r.URL.Query().Get("username")
	if username == "" {
		log.Println("No username provided, using default empty string")
	}

	renderLoginPage(w, api, username)
}
