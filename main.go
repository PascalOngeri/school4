package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"feego/handlers"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

var db *sql.DB

type Class struct {
	ID   int
	Name string
}
type selectstudent struct {
	ID    int
	Adm   string
	Class string
	Fname string
	Mname string
	Lname string
	Fee   string
	Email string
	Phone string
}
type Student struct {
	FirstName        string
	MiddleName       string
	LastName         string
	Email            string
	Class            string
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
type STU struct {
	Adm      string
	Fname    string
	Mname    string
	Lname    string
	Gender   string
	Faname   string
	Maname   string
	Class    string
	Phone    string
	Phone1   string
	Address  string
	Email    string
	Fee      string
	T1       string
	T2       string
	T3       string
	Dob      string
	Image    string
	Username string
	Password string
}
type Notice struct {
	ID      int
	Title   string
	Message string
}
type User struct {
	ID       int
	Class    string
	T1       string
	T2       string
	T3       string
	Fee      string
	id       int
	Adm      string
	UserName string
	Phone    string
	Password string

	Address string
	Phone2  string
	Phone1  string
	MotherN string
	FatherN string
	Image   string
	Dob     string
	Gender  string
	Email   string
	Lname   string
	Mname   string
	Fname   string
}
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

var store = sessions.NewCookieStore([]byte("store"))

// Initialize the database connection
func initDB() {
	var err error
	//db, err = sql.Open("mysql", "root:@mesopotamia123@tcp(localhost:3306)/eduauth")
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/eduauth")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}
func getClasses() ([]Class, error) {
	rows, err := db.Query("SELECT id, class FROM classes") // Replace "classes" with your table name
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []Class
	for rows.Next() {
		var class Class
		if err := rows.Scan(&class.ID, &class.Name); err != nil {
			return nil, err
		}
		classes = append(classes, class)
	}

	return classes, nil
}

// Retrieve API details
func getAPIDetails() (API, error) {
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

// Send SMS
func sendSMS(phoneNumber, message string) error {
	query := "SELECT apikey FROM api ORDER BY id DESC LIMIT 1"
	var apiKey string
	err := db.QueryRow(query).Scan(&apiKey)
	if err != nil {
		return fmt.Errorf("error fetching API key: %v", err)
	}

	postData := map[string]interface{}{
		"message":      message,
		"msisdn":       phoneNumber,
		"callback_url": "https://callback.io/123/dlr",
	}
	jsonData, _ := json.Marshal(postData)

	req, err := http.NewRequest("POST", "https://sms-service.smsafrica.tech/message/send/transactional", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", apiKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		_, err := db.Exec("INSERT INTO logs(user, activities) VALUES('system', 'Sent Friday SMS to all users')")
		return err
	}
	_, err = db.Exec("INSERT INTO logs(user, activities) VALUES('system', 'Failed to send Friday SMS to all users')")
	return err
}

// Check and send SMS if it's Friday
func checkAndSendFridaySMS() {
	if time.Now().Weekday() == time.Friday {
		today := time.Now().Format("2006-01-02")
		query := "SELECT COUNT(*) FROM logs WHERE user='system' AND activities='Sent Friday SMS to all users' AND DATE(date) = ?"
		var count int
		err := db.QueryRow(query, today).Scan(&count)
		if err != nil || count > 0 {
			return
		}

		rows, err := db.Query("SELECT MobileNumber FROM tbladmin")
		if err != nil {
			log.Printf("Error retrieving phone numbers: %v", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var phoneNumber string
			if err := rows.Scan(&phoneNumber); err == nil {
				sendSMS(phoneNumber, "Happy Friday! From your system.")
			}
		}
	}
}

// Render login page
func renderLoginPage(w http.ResponseWriter, api API) {
	loginData := LoginData{
		Name:     api.Name,
		Icon:     api.Icon,
		Username: "", // Populate if using cookies
		Password: "", // Populate if using cookies
		Remember: false,
	}

	log.Printf("Rendering login page with: %+v", loginData)

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Error loading template: %v", err)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, loginData)
}

// Handle login
func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Parse the form data
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")
		remember := r.FormValue("remember") == "on"

		// Declare variables for user data
		var userID int
		var foundInAdmin bool
		var adm, phone string

		// Check in tbladmin
		queryAdmin := "SELECT ID, UserName FROM tbladmin WHERE UserName = ? AND Password = ?"
		err := db.QueryRow(queryAdmin, username, password).Scan(&userID, &username)
		if err == nil {
			// User found in tbladmin
			foundInAdmin = true
		} else {
			// If not found in tbladmin, check tblregistration
			queryRegistration := "SELECT id, adm, username, phone, password FROM registration WHERE username = ? AND password = ?"
			err = db.QueryRow(queryRegistration, username, password).Scan(&userID, &adm, &username, &phone, &password)
			if err != nil {
				// User not found in either table
				http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
				return
			}
		}

		// Create session and set cookies
		session, _ := store.Get(r, "id")

		// Store values in session based on which table the user is found in
		if foundInAdmin {
			session.Values["sturecmsaid"] = userID
			session.Values["username"] = username
		} else {
			session.Values["sturecmsaid"] = userID
			session.Values["adm"] = adm
			session.Values["username"] = username
			session.Values["phone"] = phone
			session.Values["password"] = password
		}

		// Save the session
		session.Save(r, w)
		// if err := session.Save(r, w); err != nil {
		// 	log.Printf("Error saving session: %v", err)
		// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		// 	return
		// }		


		// Set cookies
		// http.SetCookie(w, &http.Cookie{Name: "user_login", Value: username, Path: "/", MaxAge: 86400})
		// if remember {
		// 	http.SetCookie(w, &http.Cookie{Name: "userpassword", Value: password, Path: "/", MaxAge: 86400})
		// }
		isHTTPS := r.TLS != nil

		// Set cookies with the secure flag if HTTPS
		http.SetCookie(w, &http.Cookie{
			Name:     "user_login",
			Value:    username,
			Path:     "/",
			MaxAge:   86400, // 1 day
			SameSite: http.SameSiteLaxMode,
			Secure:   isHTTPS, // Set Secure flag if HTTPS
		})

		// If 'remember me' is checked, store the password as well
		if remember {
			http.SetCookie(w, &http.Cookie{
				Name:     "userpassword",
				Value:    password,
				Path:     "/",
				MaxAge:   86400, // 1 day
				SameSite: http.SameSiteLaxMode,
				Secure:   isHTTPS, // Set Secure flag if HTTPS
			})
		}

		// Redirect to the appropriate dashboard or parent page
		if foundInAdmin {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther) // Redirect to admin dashboard
		} else {
			http.Redirect(w, r, "/parent", http.StatusSeeOther) // Redirect to parent page
		}
		return
	}

	// Render login page if method is not POST
	api, _ := getAPIDetails()
	renderLoginPage(w, api)
}

func renderLoginPageWithError(w http.ResponseWriter, api API, errorMsg string) {
	tmpl, _ := template.New("login").ParseFiles("templates/index.html")
	data := struct {
		API      API
		ErrorMsg string
	}{
		API:      api,
		ErrorMsg: errorMsg,
	}
	tmpl.Execute(w, data)
}
func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Get database connection details from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Construct the connection string
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Open a connection to the database
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Create a new router
	router := mux.NewRouter()

	// Static files
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Create uploads folder if it doesn't exist
	err = os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating uploads directory: %v", err)
	}

	// Define routes and corresponding handlers
	router.HandleFunc("/pay", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlePayment(w, r, db)
	}).Methods("POST")

	router.HandleFunc("/reset-password", handlers.ResetPasswordHandler(db)).Methods("GET", "POST")
	router.HandleFunc("/editB", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateBusPaymentHandler(w, r, db)
	})
	router.HandleFunc("/parent", func(w http.ResponseWriter, r *http.Request) {
		handlers.HomeHandler(w, r, db)
	})
	router.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("generate") != "" {
			handlers.GenerateFeeStatement(w, r, db)
		} else if r.FormValue("generatefee") != "" {
			handlers.GenerateFeeStructure(w, r, db)
		}
	}).Methods("POST")

	router.HandleFunc("/generete", func(w http.ResponseWriter, r *http.Request) {
		handlers.GenerateFeeHandler(w, r, db)
	}).Methods(http.MethodPost)

	router.HandleFunc("/edit-compulsory-payment", handlers.EditCompulsoryPaymentHandler(db))
	router.HandleFunc("/edit-other-payment", handlers.EditOtherPaymentHandler(db)).Methods("GET", "POST")
	router.HandleFunc("/logout", handlers.LogoutHandler()).Methods("GET")

	router.HandleFunc("/payfee", func(w http.ResponseWriter, r *http.Request) {
		handlers.PayFeeHandler(w, r, db)
	}).Methods("GET", "POST")

	router.HandleFunc("/managestudent", handlers.ManageStudent(db)).Methods("GET", "POST")
	router.HandleFunc("/deletestudent", handlers.DeleteStudent(db)).Methods("GET", "POST")
	router.HandleFunc("/updatestudent", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateUserFormHandler(w, r, db)
	}).Methods("GET", "POST")

	router.HandleFunc("/setting", func(w http.ResponseWriter, r *http.Request) {
		handlers.SettingsHandler(w, r, db)
	}).Methods("GET", "POST")

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleLogin(w, r, db)
	}).Methods("GET", "POST")

	router.HandleFunc("/dashboard", handlers.Dashboard).Methods("GET")
	router.HandleFunc("/manage", func(w http.ResponseWriter, r *http.Request) {
		handlers.Manageclass(w, r, db)
	}).Methods("GET")

	router.HandleFunc("/addclass", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddClass(w, r, db)
	}).Methods("GET", "POST")

	router.HandleFunc("/addstudent", func(w http.ResponseWriter, r *http.Request) {
		handlers.Addstudent(w, r, db)
	}).Methods("GET", "POST")

	router.HandleFunc("/optionalpay", func(w http.ResponseWriter, r *http.Request) {
		handlers.OptionalPaymentHandler(w, r, db)
	}).Methods("GET", "POST")

	router.HandleFunc("/addpubnot", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddPubNot(w, r, db)
	}).Methods("GET", "POST")

	router.HandleFunc("/managepubnot", handlers.ManagePubNot(db)).Methods("GET")

	router.HandleFunc("/search", searchStudentHandler).Methods("GET", "POST")
	router.HandleFunc("/send", send).Methods("POST")

	router.HandleFunc("/adduser", handlers.ManageUser(db)).Methods("GET", "POST")
	router.HandleFunc("/logs", handlers.Logs(db)).Methods("GET")

	router.HandleFunc("/otherpayinsert", func(w http.ResponseWriter, r *http.Request) {
		handlers.Insert(w, r, db)
	}).Methods("POST")

	router.HandleFunc("/paymentinsert", func(w http.ResponseWriter, r *http.Request) {
		handlers.OptionalPaymentHandler(w, r, db)
	}).Methods("POST")

	router.HandleFunc("/transportinsert", func(w http.ResponseWriter, r *http.Request) {
		handlers.TransportPaymentHandler(w, r, db)
	}).Methods("POST")

	// Background task OptionalPaymentHandler
	router.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {
		handlers.GenerateFeeHandler(w, r, db)
	})

	go checkAndSendFridaySMS()
	router.HandleFunc("/manage-public-notice", handlers.ManagePubNot(db)).Methods("GET")
	router.HandleFunc("/delete-public-notice", handlers.DeleteNotice(db)).Methods("GET")

	router.HandleFunc("/delete-class", handlers.DeleteClass(db)).Methods("GET")
	router.HandleFunc("/edit-class", handlers.EditClass(db)).Methods("GET")
	router.HandleFunc("/update-class", handlers.UpdateClass(db)).Methods("POST")

	router.HandleFunc("/setfee", func(w http.ResponseWriter, r *http.Request) {
		handlers.SetFeeHandler(w, r, db)
	}).Methods("GET", "POST")

	router.HandleFunc("/transport", func(w http.ResponseWriter, r *http.Request) {
		handlers.FormHandler(w, r, db)
	}).Methods("GET", "POST")

	router.HandleFunc("/updatepayment", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdatePaymentHandler(w, r, db)
	}).Methods("GET")

	router.HandleFunc("/deleteother", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteOtherHandler(w, r, db)
	}).Methods("GET")

	router.HandleFunc("/deletebus", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteBusHandler(w, r, db)
	}).Methods("GET")

	router.HandleFunc("/deletecompulsory", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteCompulsoryHandler(w, r, db)
	}).Methods("GET")
	// log.Fatal(http.ListenAndServe("127.0.0.1:8080", router))
	// log.Println("Server is running on :8080")

	err = http.ListenAndServe(":8080", router)
if err != nil {
    log.Fatalf("Server failed to start: %v", err)
} 


}
func add1(i int) int {
	return i + 1
}
func searchStudentHandler(w http.ResponseWriter, r *http.Request) {
	funcMap := template.FuncMap{
		"add1": add1, // Register the add1 function
	}
	// Handle POST request (form submission)
	if r.Method == http.MethodPost {
		searchData := r.FormValue("searchdata")
		if searchData == "" {
			http.Error(w, "Please enter a search term", http.StatusBadRequest)
			return
		}

		// Query the database to search for students by their admission number (Adm)
		rows, err := db.Query("SELECT adm, fname, mname, lname, gender, faname, maname, class, phone, phone1, address, email, fee, t1, t2, t3, dob, image, username, password FROM registration WHERE adm LIKE ?", "%"+searchData+"%")
		if err != nil {
			http.Error(w, "Error querying database: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Create a slice to hold the students found
		var students []STU
		for rows.Next() {
			var student STU
			if err := rows.Scan(&student.Adm, &student.Fname, &student.Mname, &student.Lname, &student.Gender, &student.Faname, &student.Maname, &student.Class, &student.Phone, &student.Phone1, &student.Address, &student.Email, &student.Fee, &student.T1, &student.T2, &student.T3, &student.Dob, &student.Image, &student.Username, &student.Password); err != nil {
				log.Println(err)
				continue
			}
			students = append(students, student)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "Error reading from the database: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Register the custom template function

		// Render the result template with the students data
		tmpl, err := template.New("search").Funcs(funcMap).ParseFiles(
			"templates/search.html", // Update this path as needed
			"includes/header.html",
			"includes/sidebar.html",
			"includes/footer.html",
		)
		if err != nil {
			http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Error parsing template files: %v", err)
			return
		}

		// Pass the students data to the template
		// Pass the students data to the template
		err = tmpl.Execute(w, students) // 'students' is the slice of STU passed as context
		if err != nil {
			http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Error executing template: %v", err)
			return
		}

	} else {
		// Handle GET request (render search form)
		tmpl, err := template.ParseFiles(
			"templates/search.html", // Update this path as needed
			"includes/header.html",
			"includes/sidebar.html",
			"includes/footer.html",
		)
		if err != nil {
			http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Error parsing template files: %v", err)
			return
		}

		// Execute the template (empty data for initial search page)
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Error executing template: %v", err)
			return
		}
	}
}

func send(w http.ResponseWriter, r *http.Request) {
	// Parse the template files
	tmpl, err := template.ParseFiles("templates/send.html", "includes/footer.html", "includes/header.html", "includes/sidebar.html")
	if err != nil {
		// Handle the error properly, e.g., by returning a 500 status
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Data to pass to the template
	data := map[string]interface{}{
		"Title": "Manage Class", // Example dynamic data
	}

	// Execute the template and write to the response
	err = tmpl.Execute(w, data)
	if err != nil {
		// Handle the error properly
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func adduser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse the form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Unable to parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		AName := r.FormValue("adminname")
		mobno := r.FormValue("mobilenumber")
		email := r.FormValue("email")

		pass := r.FormValue("password")
		username := r.FormValue("username")
		// Log the received form data
		log.Printf("Notice Title: %s, Notice Message: %s", AName, mobno)

		// Check if form data is valid
		if mobno == "" || username == "" {
			http.Error(w, "All fields are required fields.", http.StatusBadRequest)
			return
		}

		// Insert data into the database
		_, err := db.Exec("INSERT INTO tblAdmin (AdminName,Email,UserName,password,MobileNumber) VALUES (?, ?,?,?,?)", AName, email, username, pass, mobno)
		if err != nil {
			log.Printf("Failed to insert notice: %v", err) // Log the error
			http.Error(w, "Failed to insert notice: "+err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("Notice successfully added")

		// Redirect to the form page (or any other success page)
		http.Redirect(w, r, "/adduser", http.StatusSeeOther)
		return
	}

	// Render the template for GET requests
	tmpl, err := template.ParseFiles(
		"templates/adduser.html",
		"includes/header.html",
		"includes/sidebar.html",
		"includes/footer.html",
	)
	if err != nil {
		http.Error(w, "Template parsing failed: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error parsing template files: %v", err)
		return
	}

	// Execute the template
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Template execution failed: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
		return
	}
}
