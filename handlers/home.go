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
	log.Println("Received request to load the home page")

	// Retrieve values from query parameters or headers
	adm := r.URL.Query().Get("adm")
	username := r.URL.Query().Get("username")
	phone := r.URL.Query().Get("phone")
	password := r.URL.Query().Get("password")

	// Log the user data (e.g., admission number, username) for debugging purposes
	log.Printf("Loading home page for user: %s, Admission Number: %s", username, adm)

	// Validate that the required parameters are present
	if adm == "" || username == "" || phone == "" || password == "" {
		log.Printf("Missing required parameters: adm=%s, username=%s, phone=%s, password=%s", adm, username, phone, password)
		http.Error(w, "Missing user data", http.StatusBadRequest)
		return
	}

	// Fetch payment history
	log.Println("Fetching payment history for admission number:", adm)
	paymentRows, err := db.Query("SELECT id, adm, date, amount, bal FROM payment WHERE adm = ?", adm)
	if err != nil {
		log.Printf("Failed to fetch payments: %v", err)
		http.Error(w, "Internal server error while fetching payments", http.StatusInternalServerError)
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
	log.Printf("Fetched %d payments for admission number: %s", len(payments), adm)

	// Fetch notices
	log.Println("Fetching public notices")
	noticeRows, err := db.Query("SELECT NoticeTitle, NoticeMessage FROM tblpublicnotice")
	if err != nil {
		log.Printf("Failed to fetch notices: %v", err)
		http.Error(w, "Internal server error while fetching notices", http.StatusInternalServerError)
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
	log.Printf("Fetched %d notices", len(notices))

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
	log.Println("Parsing home page template")
	tmpl, err := template.ParseFiles("templates/parent.html", "includes/footer.html")
	if err != nil {
		log.Printf("Error loading template: %v", err)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	// Execute the template
	log.Println("Rendering home page template with user and payment data")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}

	log.Println("Home page rendered successfully")
}
