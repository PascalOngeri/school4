package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// PaymentRecord defines the structure for each payment record
type PaymentRecord struct {
	ID      int
	Date    string
	Amount  float64
	Balance float64
}

// GenerateFeeHandler generates the fee statement for a given admission number
func GenerateFeeHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Check if request method is POST
	if r.Method != http.MethodPost {
		log.Println("Invalid request method: ", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve admission number from form data
	adm := r.FormValue("adm")
	if adm == "" {
		log.Println("Admission number is required")
		http.Error(w, "Admission number is required", http.StatusBadRequest)
		return
	}

	// Query the database to get payment records for the given admission number
	rows, err := db.Query("SELECT id, date, amount, bal FROM payment WHERE adm = ? ORDER BY id ASC", adm)
	if err != nil {
		log.Printf("Database query error for admission number %s: %v\n", adm, err)
		http.Error(w, "Error fetching payment records", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a slice to hold the payment records
	var payments []PaymentRecord

	// Populate the slice with the query results
	for rows.Next() {
		var payment PaymentRecord
		if err := rows.Scan(&payment.ID, &payment.Date, &payment.Amount, &payment.Balance); err != nil {
			log.Printf("Failed to scan row for admission number %s: %v\n", adm, err)
			http.Error(w, "Error processing payment records", http.StatusInternalServerError)
			return
		}
		payments = append(payments, payment)
	}

	// Check for any row iteration errors
	if rows.Err() != nil {
		log.Printf("Error iterating over rows for admission number %s: %v\n", adm, rows.Err())
		http.Error(w, "Error processing payment records", http.StatusInternalServerError)
		return
	}

	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set logo path and add the logo to the document
	logoPath := filepath.Join("assets", "images", "logo.png")
	if _, err := os.Stat(logoPath); os.IsNotExist(err) {
		log.Printf("Logo file not found: %v\n", err)
		http.Error(w, "Logo file not found", http.StatusInternalServerError)
		return
	}
	pdf.ImageOptions(logoPath, 80, 10, 50, 0, false, gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")

	// Add school name below the logo
	pdf.SetFont("Arial", "B", 16)
	pdf.Ln(50) // Adjust the vertical space as needed
	pdf.CellFormat(0, 10, "INFINITY SCHOOLS", "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Add table headers
	pdf.SetFont("Arial", "B", 12)
	headers := []string{"Payment No.", "Date", "Amount", "Balance", "Status"}
	widths := []float64{40, 38, 38, 38, 38}
	for i, header := range headers {
		pdf.CellFormat(widths[i], 10, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	// Set font for table rows
	pdf.SetFont("Arial", "", 10)

	// Iterate over the payment records and add data to the PDF
	for _, payment := range payments {
		status := "Received"
		pdf.CellFormat(40, 10, fmt.Sprintf("%d", payment.ID), "1", 0, "C", false, 0, "")
		pdf.CellFormat(38, 10, payment.Date, "1", 0, "C", false, 0, "")
		pdf.CellFormat(38, 10, fmt.Sprintf("%.2f", payment.Amount), "1", 0, "C", false, 0, "") // Display amount
		pdf.CellFormat(38, 10, fmt.Sprintf("%.2f", payment.Balance), "1", 0, "C", false, 0, "")
		pdf.CellFormat(38, 10, status, "1", 1, "C", false, 0, "")
	}

	// // Check for any PDF generation errors
	// if pdf.Err() != nil {
	// 	log.Println("Error generating PDF: ", pdf.Err())
	// 	http.Error(w, "Error generating PDF", http.StatusInternalServerError)
	// 	return
	// }

	// Create output directory
	outputDir := "generated_pdfs"
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Printf("Error creating directory %s: %v", outputDir, err)
		http.Error(w, "Error creating directory", http.StatusInternalServerError)
		return
	}

	// Generate unique file name for the PDF
	timestamp := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("fee_statement_%s_%s.pdf", adm, timestamp)
	filePath := filepath.Join(outputDir, fileName)

	// Save the PDF to the file
	err = pdf.OutputFileAndClose(filePath)
	if err != nil {
		log.Printf("Error saving PDF file %s: %v", filePath, err)
		http.Error(w, "Error saving PDF", http.StatusInternalServerError)
		return
	}

	// Log the success of the PDF generation
	log.Printf("Fee statement generated successfully for admission number %s. Saved to %s\n", adm, filePath)

	// Redirect to the parent page (or wherever you want)
	http.Redirect(w, r, "/parent", http.StatusSeeOther)
}
