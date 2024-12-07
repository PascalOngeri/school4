package handlers

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Student struct {
	FirstName  string
	MiddleName string
	LastName   string
	Email      string
	Class      string

	Gender           string
	DOB              string
	AdmissionNumber  string
	Image            string
	FatherName       string
	MotherName       string
	ContactNumber    string
	AltContactNumber string
	Address          string
	UserName         string
	Password         string
}

type Class struct {
	ID   int
	Name string

	Class string
	T1    float64
	T2    float64
	T3    float64
	Fee   float64
}

func GetClassDetails(db *sql.DB, class string) (float64, float64, float64, float64, error) {
	query := `SELECT t1, t2, t3, fee FROM classes WHERE class = ?`
	var t1, t2, t3, fee float64
	err := db.QueryRow(query, class).Scan(&t1, &t2, &t3, &fee)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("[INFO] No class details found for class: %s", class)
			return 0, 0, 0, 0, nil
		}
		log.Printf("[ERROR] Failed to fetch class details: %v", err)
		return 0, 0, 0, 0, err
	}
	log.Printf("[DEBUG] Retrieved class details for %s: T1=%.2f, T2=%.2f, T3=%.2f, Fee=%.2f", class, t1, t2, t3, fee)
	return t1, t2, t3, fee, nil
}

func Addstudent(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Printf("[INFO] Addstudent handler invoked")
	tmpl, err := template.ParseFiles(
		"templates/addstudent.html",
		"includes/header.html",
		"includes/sidebar.html",
		"includes/footer.html",
	)
	if err != nil {
		log.Printf("[ERROR] Failed to parse templates: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var classes []Class
	rows, err := db.Query("SELECT id, class FROM classes")
	if err != nil {
		log.Printf("[ERROR] Database query failed: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var class Class
		if err := rows.Scan(&class.ID, &class.Name); err != nil {
			log.Printf("[ERROR] Error scanning row: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		classes = append(classes, class)
	}
	if err := rows.Err(); err != nil {
		log.Printf("[ERROR] Error iterating rows: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		log.Printf("[INFO] Handling POST request to add student")
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			log.Printf("[ERROR] Failed to parse form: %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		student := Student{
			FirstName:        r.FormValue("fname"),
			MiddleName:       r.FormValue("mname"),
			LastName:         r.FormValue("lname"),
			Email:            r.FormValue("stuemail"),
			Class:            r.FormValue("stuclass"),
			Gender:           r.FormValue("gender"),
			DOB:              r.FormValue("dob"),
			AdmissionNumber:  r.FormValue("stuid"),
			FatherName:       r.FormValue("faname"),
			MotherName:       r.FormValue("maname"),
			ContactNumber:    r.FormValue("connum"),
			AltContactNumber: r.FormValue("altconnum"),
			Address:          r.FormValue("address"),
			UserName:         r.FormValue("uname"),
			Password:         hashPassword(r.FormValue("password")),
		}

		file, handler, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			uploadDir := "uploads/"
			os.MkdirAll(uploadDir, os.ModePerm)

			fileName := filepath.Join(uploadDir, randomFileName(handler.Filename))
			destFile, err := os.Create(fileName)
			if err == nil {
				defer destFile.Close()
				_, _ = destFile.ReadFrom(file)
				student.Image = fileName
				log.Printf("[DEBUG] Uploaded file saved as: %s", fileName)
			} else {
				log.Printf("[ERROR] Failed to save uploaded file: %v", err)
			}
		} else {
			log.Printf("[INFO] No image uploaded for student")
		}

		t1, t2, t3, fee, err := GetClassDetails(db, student.Class)
		if err != nil {
			log.Printf("[ERROR] Failed to fetch class details: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		query := `
			INSERT INTO registration (
				adm, fname, mname, lname, gender, faname, maname, class, phone, phone1,
				address, email, fee, t1, t2, t3, dob, image, username, password
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		_, err = db.Exec(query,
			student.AdmissionNumber, student.FirstName, student.MiddleName, student.LastName,
			student.Gender, student.FatherName, student.MotherName, student.Class,
			student.ContactNumber, student.AltContactNumber, student.Address,
			student.Email, fee, t1, t2, t3, student.DOB, student.Image,
			student.UserName, student.Password,
		)
		if err != nil {
			log.Printf("[ERROR] Failed to insert student into database: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		log.Printf("[INFO] Successfully added student: %s %s", student.FirstName, student.LastName)
		http.Redirect(w, r, "/addstudent", http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{
		"Title":   "Add Students",
		"Classes": classes,
	}
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("[ERROR] Failed to render template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func hashPassword(password string) string {
	return password
}

func randomFileName(original string) string {
	buf := make([]byte, 8)
	_, _ = rand.Read(buf)
	return fmt.Sprintf("%x_%s", buf, original)
}
