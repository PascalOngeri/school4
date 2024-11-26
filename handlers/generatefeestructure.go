package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
)

// GenerateFeeHandler handles PDF generation for fee statements
func GenerateFeeHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Get the admission number from the form data
	adm := r.FormValue("genclass")
	if adm == "" {
		http.Error(w, "Admission number is required", http.StatusBadRequest)
		return
	}

	// Fetch payment details from the database
	rows, err := db.Query("SELECT id, date, amount, bal FROM payment WHERE adm = ? ORDER BY id ASC", adm)
	if err != nil {
		log.Printf("Error querying payment records: %v\n", err)
		http.Error(w, "Error fetching payment records", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Add logo
	logoPath := filepath.Join("assets", "images", "logo.png")
	pdf.ImageOptions(logoPath, 80, 10, 50, 0, false, gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")

	// Add title and admission number
	pdf.SetFont("Arial", "B", 16)
	pdf.Ln(20) // Add vertical space
	pdf.CellFormat(0, 10, "School Name", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 10, fmt.Sprintf("Admission Number: %s Fee Statement", adm), "", 1, "C", false, 0, "")

	// Table headers
	pdf.Ln(10) // Add space before the table
	pdf.SetFont("Arial", "B", 12)
	headers := []string{"Payment No.", "Date", "Amount", "Balance", "Status"}
	widths := []float64{40, 38, 38, 38, 38}
	for i, header := range headers {
		pdf.CellFormat(widths[i], 10, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	// Table data
	pdf.SetFont("Arial", "", 10)
	for rows.Next() {
		var id int
		var date string
		var amount, balance float64
		if err := rows.Scan(&id, &date, &amount, &balance); err != nil {
			log.Printf("Error scanning row: %v\n", err)
			continue
		}

		status := "Received"
		pdf.CellFormat(40, 10, fmt.Sprintf("%d", id), "1", 0, "C", false, 0, "")
		pdf.CellFormat(38, 10, date, "1", 0, "C", false, 0, "")
		pdf.CellFormat(38, 10, fmt.Sprintf("%.2f", amount), "1", 0, "C", false, 0, "")
		pdf.CellFormat(38, 10, fmt.Sprintf("%.2f", balance), "1", 0, "C", false, 0, "")
		pdf.CellFormat(38, 10, status, "1", 1, "C", false, 0, "")
	}

	if rows.Err() != nil {
		log.Printf("Error iterating rows: %v\n", rows.Err())
		http.Error(w, "Error processing payment records", http.StatusInternalServerError)
		return
	}

	// Set the response headers to trigger download
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="fee_statement.pdf"`)

	// Output the PDF to the response
	if err := pdf.Output(w); err != nil {
		log.Printf("Error generating PDF: %v\n", err)
		http.Error(w, "Error generating PDF", http.StatusInternalServerError)
	}
}
