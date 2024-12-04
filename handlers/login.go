package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Secret key for signing JWT tokens
var jwtSecretKey = []byte("your-secret-key")

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

// Claims struct for JWT claims

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

func HandleLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == "POST" {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")

		var userID int
		var foundInAdmin bool
		var adm, phone, role string // Define the fee variable

		// Authenticate user in tbladmin
		queryAdmin := "SELECT ID, UserName FROM tbladmin WHERE UserName = ? AND Password = ?"
		err := db.QueryRow(queryAdmin, username, password).Scan(&userID, &username)
		if err == nil {
			foundInAdmin = true
			role = "admin"
		} else {
			// Authenticate user in registration table
			queryRegistration := "SELECT id, adm, username, phone, password FROM registration WHERE username = ? AND password = ?"
			err = db.QueryRow(queryRegistration, username, password).Scan(&userID, &adm, &username, &phone, &password) // Pass address of fee
			if err != nil {
				http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
				return
			}
			role = "user"
		}

		// Create JWT token
		expirationTime := time.Now().Add(8 * time.Hour) // 8 hours expiration

		claims := &Claims{
			UserID:   userID,
			Username: username,
			Adm:      adm,
			Phone:    phone,
			Role:     role,
			password: password,
			// Password included (consider security implications)
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtSecretKey)
		if err != nil {
			http.Error(w, "Error creating JWT token", http.StatusInternalServerError)
			return
		}

		// Set JWT in cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    tokenString,
			Path:     "/",
			Expires:  expirationTime,
			HttpOnly: true,
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

// HandleLogin handles login requests
// func HandleLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
// 	if r.Method == "POST" {
// 		r.ParseForm()
// 		username := r.FormValue("username")
// 		password := r.FormValue("password")

// 		var userID int
// 		var foundInAdmin bool
// 		var adm, phone, role string  // Define the fee variable

// 		// Authenticate user in tbladmin
// 		queryAdmin := "SELECT ID, UserName FROM tbladmin WHERE UserName = ? AND Password = ?"
// 		err := db.QueryRow(queryAdmin, username, password).Scan(&userID, &username)
// 		if err == nil {
// 			foundInAdmin = true
// 			role = "admin"
// 		} else {
// 			// Authenticate user in registration table
// 			queryRegistration := "SELECT id, adm, username, phone, password FROM registration WHERE username = ? AND password = ?"
// 			err = db.QueryRow(queryRegistration, username, password).Scan(&userID, &adm, &username, &phone, &password)  // Pass address of fee
// 			if err != nil {
// 				http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
// 				return
// 			}
// 			role = "user"
// 		}

// 		// Create JWT token
// 		expirationTime := time.Now().Add(8 * time.Hour) // 8 hours expiration

// 		claims := &Claims{
// 			UserID:   userID,
// 			Username: username,
// 			Adm:      adm,
// 			Phone:    phone,
// 			Role:     role,
// 			password: password,
// 		 // Password included (consider security implications)
// 			RegisteredClaims: jwt.RegisteredClaims{
// 				ExpiresAt: jwt.NewNumericDate(expirationTime),
// 			},
// 		}
// 		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 		tokenString, err := token.SignedString(jwtSecretKey)
// 		if err != nil {
// 			http.Error(w, "Error creating JWT token", http.StatusInternalServerError)
// 			return
// 		}

// 		// Set JWT in cookie
// 		http.SetCookie(w, &http.Cookie{
// 			Name:     "auth_token",
// 			Value:    tokenString,
// 			Path:     "/",
// 			Expires:  expirationTime,
// 			HttpOnly: true,
// 		})

// 		// Redirect based on role
// 		if foundInAdmin {
// 			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
// 		} else {
// 			http.Redirect(w, r, "/parent", http.StatusSeeOther)
// 		}
// 		return
// 	}

// 	// Render the login page for GET requests
// 	api, _ := getAPIDetails(db)
// 	renderLoginPage(w, api, "")
// }

// ValidateJWT function for validating the token
func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
