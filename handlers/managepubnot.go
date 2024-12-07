package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

// Notice represents a public notice with a title and message
type Notice struct {
	Title   string
	Message string
	Date    string
	ID      interface{}
}

// ManagePubNot handles public notices by fetching them from the database
func ManagePubNot(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log incoming request
		log.Println("Received request for managing public notices")

		// Query to fetch public notices
		query := "SELECT ID, NoticeTitle, NoticeMessage FROM tblpublicnotice"
		rows, err := db.Query(query)
		if err != nil {
			// Log error during database query
			log.Printf("Database query failed: %v", err)
			http.Error(w, "Failed to fetch notices from the database.", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var notices []Notice
		// Process rows
		for rows.Next() {
			var notice Notice
			if err := rows.Scan(&notice.ID, &notice.Title, &notice.Message); err != nil {
				// Log error during row scanning
				log.Printf("Error scanning row: %v", err)
				http.Error(w, "Error processing data from the database.", http.StatusInternalServerError)
				return
			}
			notices = append(notices, notice)
		}

		// Check for errors during row iteration
		if err := rows.Err(); err != nil {
			// Log error during row iteration
			log.Printf("Error iterating rows: %v", err)
			http.Error(w, "Error reading notices from the database.", http.StatusInternalServerError)
			return
		}

		// Log number of notices retrieved
		log.Printf("Successfully fetched %d notices from the database", len(notices))

		// Parse template files
		tmpl, err := template.ParseFiles(
			"templates/managepubnot.html",
			"includes/header.html",
			"includes/sidebar.html",
			"includes/footer.html",
		)
		if err != nil {
			// Log error during template parsing
			log.Printf("Template parsing failed: %v", err)
			http.Error(w, "Failed to load page templates.", http.StatusInternalServerError)
			return
		}

		// Pass data to the template
		data := map[string]interface{}{
			"Title":   "Manage Public Notice",
			"Notices": notices,
		}

		// Render the template
		if err := tmpl.Execute(w, data); err != nil {
			// Log error during template execution
			log.Printf("Template execution failed: %v", err)
			http.Error(w, "Failed to render the page.", http.StatusInternalServerError)
			return
		}

		// Log successful rendering
		log.Println("Successfully rendered manage public notice page")
	}
}
