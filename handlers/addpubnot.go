package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

// AddPubNot handles adding a new public notice to the database
func AddPubNot(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			log.Printf("ERROR: Unable to parse form data: %v", err)
			http.Error(w, "Bad Request: Unable to parse form data", http.StatusBadRequest)
			return
		}

		// Retrieve form values
		nottitle := r.FormValue("nottitle")
		notmsg := r.FormValue("notmsg")

		// Log the received form data for debugging
		log.Printf("INFO: Received form data - Notice Title: %s, Notice Message: %s", nottitle, notmsg)

		// Validate form data
		if nottitle == "" || notmsg == "" {
			log.Printf("ERROR: Missing required fields - Notice Title or Notice Message is empty")
			http.Error(w, "Notice Title and Message are required fields.", http.StatusBadRequest)
			return
		}

		// Insert data into the database
		_, err = db.Exec("INSERT INTO tblpublicnotice (NoticeTitle, NoticeMessage) VALUES (?, ?)", nottitle, notmsg)
		if err != nil {
			log.Printf("ERROR: Failed to insert notice into database: %v", err)
			http.Error(w, "Internal Server Error: Failed to add notice", http.StatusInternalServerError)
			return
		}

		log.Printf("INFO: Notice successfully added - Title: %s", nottitle)

		// Redirect to the form page (or a confirmation page)
		http.Redirect(w, r, "/addpubnot", http.StatusSeeOther)
		return
	}

	// Render the form for GET requests
	tmpl, err := template.ParseFiles(
		"templates/addpubnotice.html",
		"includes/header.html",
		"includes/sidebar.html",
		"includes/footer.html",
	)
	if err != nil {
		log.Printf("ERROR: Failed to parse template files for AddPubNot page: %v", err)
		http.Error(w, "Internal Server Error: Unable to load page", http.StatusInternalServerError)
		return
	}

	// Execute the template
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("ERROR: Failed to render AddPubNot page: %v", err)
		http.Error(w, "Internal Server Error: Unable to render page", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: Successfully rendered AddPubNot page")
}
