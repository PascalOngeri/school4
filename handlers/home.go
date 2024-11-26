package handlers

import (
	"database/sql"
	"fmt"
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
	session, err := store.Get(r, "store")
	if err != nil {
		log.Printf("Failed to retrieve session: %v", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	// Check if user is logged in
	if session.Values["adm"] == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Retrieve values from session
	adm := fmt.Sprintf("%v", session.Values["adm"])
	username := fmt.Sprintf("%v", session.Values["username"])
	phone := fmt.Sprintf("%v", session.Values["phone"])
	password := fmt.Sprintf("%v", session.Values["password"])

	// Fetch payment history
	paymentRows, err := db.Query("SELECT id, adm, date, amount, bal FROM payment WHERE adm = ?", adm)
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
	noticeRows, err := db.Query("SELECT NoticeTitle, NoticeMessage FROM tblpublicnotice")
	if err != nil {
		log.Printf("Failed to fetch notices: %v", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	defer noticeRows.Close()

	var notices []Notice
	for noticeRows.Next() {
		var n Notice
		err := noticeRows.Scan(&n.Title, &n.Message)
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
