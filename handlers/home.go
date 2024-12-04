package handlers

import (
	"database/sql"

	"html/template"
	"log"
	"net/http"
)

// HomePageData is the structure used to hold data for rendering the home page template
type HomePageData struct {
	Title           string
	Username        string
	AdmissionNumber string
	Password        string
	Phone           string
	Payments        []Payment
	Notices         []Notice
}

// HomeHandler handles requests to the home page
func HomeHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
	if claims.Role != "user" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// Log authenticated user info for debugging
	log.Printf("Authenticated user: %s, Role: %s, Admission Number: %s, Password: %s, Phone: %s",
		claims.Username, claims.Role, claims.Adm, claims.password, claims.Phone)

	adm := claims.Adm
	username := claims.Username
	phone := claims.Phone
	password := claims.password

	// Retrieve values from session

	// Fetch payment history
	paymentRows, err := db.Query("SELECT id, adm, date, amount, bal FROM payment WHERE adm = ? ORDER BY id DESC", adm)
	if err != nil {
		log.Printf("Failed to fetch payments: %v", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	defer paymentRows.Close()

	var payments []Payment
	for paymentRows.Next() {
		var p Payment
		err := paymentRows.Scan(&p.SNo, &p.RegNo, &p.Date, &p.Amount, &p.Balance)
		if err != nil {
			log.Printf("Failed to scan payment: %v", err)
			continue
		}
		payments = append(payments, p)
	}

	// Fetch notices
	noticeRows, err := db.Query("SELECT NoticeTitle, NoticeMessage,CreationDate FROM tblpublicnotice")
	if err != nil {
		log.Printf("Failed to fetch notices: %v", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	defer noticeRows.Close()

	var notices []Notice
	for noticeRows.Next() {
		var n Notice
		err := noticeRows.Scan(&n.Title, &n.Message, &n.Date)
		if err != nil {
			log.Printf("Failed to scan notice: %v", err)
			continue
		}
		notices = append(notices, n)
	}

	// Prepare data for the template
	data := HomePageData{
		Title:           "Infinityschools Analytics",
		Username:        username,
		AdmissionNumber: adm,
		Password:        password,
		Phone:           phone,
		Payments:        payments,
		Notices:         notices,
	}

	// Render the template
	tmpl, err := template.ParseFiles("templates/parent.html", "includes/footer.html")
	if err != nil {
		log.Printf("Error loading template: %v", err)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
