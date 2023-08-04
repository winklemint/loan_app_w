package loan

import (
	con "backend/Config"
	"backend/user"
	"database/sql"
	"encoding/json"
	"strings"

	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"

	"net/http"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"

	_ "backend/user"

	"github.com/gorilla/mux"
)

type Loan_details struct {
	ID                   int       `json:"id,omitempty"`
	Loan_type            string    `json:"loan_type,omitempty" validate:"required"`
	Employment_type      string    `json:"employment_type,omitempty" `
	Loan_amount          float64   `json:"loan_amount,omitempty" validate:"gt=0"`
	Gross_monthly_income float64   `json:"gross_monthly_income"`
	Pincode              int       `json:"pincode,omitempty" validate:"len=6"`
	Tenure               int       `json:"tenure,omitempty" validate:"gt=0"`
	Created_at           time.Time `json:"created_at,omitempty"`
	Last_modified        time.Time `json:"last_modified,omitempty"`
	Status               string    `json:"status"`
	Remark               string    `json:"remark"`
	Admin_Name           string    `json:"admin_name"`
}

var db *sql.DB
var loan Loan_details

func init() {
	var err error
	db, err = con.GetDB()
	if err != nil {
		panic(fmt.Errorf("failed to initialize database: %v", err))
	}
}

func ValidateLoan(loan Loan_details, validate *validator.Validate) error {
	// Check that loan_type is not empty
	if loan.Loan_type == "" {
		return errors.New("Loan type is required.")
	}

	// Check that loan_amount is greater than zero
	if loan.Loan_amount <= 0 {
		return errors.New("Loan amount must be greater than zero.")
	}

	// Check that tenure is greater than zero
	if loan.Tenure <= 0 {
		return errors.New("Tenure must be greater than zero.")
	}

	// Check that pincode is exactly 6 digits
	if match, _ := regexp.MatchString(`^\d{6}$`, fmt.Sprint(loan.Pincode)); !match {
		return errors.New("Pincode must be exactly 6 digits.")
	}

	// Check that employment_type is not empty
	if loan.Employment_type == "" {
		return errors.New("Employment type is required.")
	}

	if loan.Gross_monthly_income <= 0 {
		return errors.New("Gross monthly income must be greater than zero.")
	}

	// If all checks pass, return nil
	return nil
}

type Error struct {
	Message string `json:"message"`
}
type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error *Error      `json:"error,omitempty"`
}

// func InsertLoanDetails(w http.ResponseWriter, r *http.Request) {

// 	config, err := con.LoadConfig("Config/config.yaml")
// 	if err != nil {
// 		panic(err)
// 	}

// 	con.ConnectDB(config)
// 	defer con.CloseDB()

// 	//r.Header.Get("User-Agent")
// 	//r.Header("User-Agent")

// 	var loan Loan_details
// 	if err := json.NewDecoder(r.Body).Decode(&loan); err != nil {
// 		logrus.Errorf("Error decoding request body: %v", err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(&Response{Error: &Error{Message: err.Error()}})
// 		return
// 	}

// 	//create validator instance
// 	validate := validator.New()

// 	// validate loan details
// 	if err := ValidateLoan(loan, validate); err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(&Response{Error: &Error{Message: err.Error()}})
// 		return
// 	}

// 	result, err := db.Exec("INSERT INTO loan_details_table(loan_type,loan_amount,pincode,tenure,employment_type,gross_monthly_income,created_at,last_modified) VALUES(?,?,?,?,?,?,NOW(),NOW())", loan.Loan_type, loan.Loan_amount, loan.Pincode, loan.Tenure, loan.Employment_type, loan.Gross_monthly_income)
// 	if err != nil {
// 		logrus.Errorf("Error executing SQL query: %v", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode(&Response{Error: &Error{Message: err.Error()}})
// 		return
// 	}

// 	id, _ := result.LastInsertId()
// 	loan.ID = int(id)
// 	logrus.WithFields(logrus.Fields{
// 		"loan_id":              loan.ID,
// 		"Loan_Type":            loan.Loan_type,
// 		"Loan_Amount":          loan.Loan_amount,
// 		"Tenure":               loan.Tenure,
// 		"Pincode":              loan.Pincode,
// 		"Employment_Type":      loan.Employment_type,
// 		"Gross_Monthly_Income": loan.Gross_monthly_income,
// 	}).Info("Loan Details Inserted")

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(&Response{Data: loan})
// }

func InsertLoanDetails(w http.ResponseWriter, r *http.Request) {
	config, err := con.LoadConfig("Config/config.yaml")
	if err != nil {
		logrus.Error("Failed to load config:", err)
		//http.Error(w, "Failed to load config", http.StatusInternalServerError)
		return
	}

	con.ConnectDB(config)
	defer con.CloseDB()

	// Set the appropriate headers to enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Referer, Sec-Ch-Ua, Sec-Ch-Ua-Mobile, Sec-Ch-Ua-Platform, User-Agent, Authorization")

	// Handle preflight OPTIONS requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Check request method
	if r.Method != http.MethodPost {
		//http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve the JWT token from the request header
	authHeader := r.Header.Get("Cookie")
	if authHeader == "" {
		logrus.Warnln("Authorization header missing")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&Response{Error: &Error{Message: "Unauthorized"}})
		return
	}

	// Extract the token from the Authorization header
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	session := &user.Session{
		Token: tokenString,
	}

	// Use the GetUserIDFromSession function from the login package to retrieve the userID
	userID, err := user.GetUserIDFromSession(session)
	if err != nil {
		logrus.Errorf("Error retrieving userID from session: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&Response{Error: &Error{Message: "Unauthorized"}})
		return
	}

	var loan Loan_details
	if err := json.NewDecoder(r.Body).Decode(&loan); err != nil {
		logrus.Errorf("Error decoding request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&Response{Error: &Error{Message: err.Error()}})
		return
	}

	// Create validator instance
	validate := validator.New()

	// Validate loan details
	if err := ValidateLoan(loan, validate); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&Response{Error: &Error{Message: err.Error()}})
		return
	}

	// Connect to the database
	db, err := con.GetDB()
	if err != nil {
		logrus.Errorf("Error connecting to the database: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&Response{Error: &Error{Message: "Internal Server Error"}})
		return
	}
	defer db.Close()

	result, err := db.Exec("INSERT INTO loan_details_table (loan_type, employment_type, loan_amount, gross_monthly_income, pincode, tenure,   created_at, last_modified, user_id) VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW(), ?)",
		loan.Loan_type, loan.Employment_type, loan.Loan_amount, loan.Gross_monthly_income, loan.Pincode, loan.Tenure, userID)
	if err != nil {
		logrus.Errorf("Error executing SQL query: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&Response{Error: &Error{Message: "Internal Server Error"}})
		return
	}

	id, _ := result.LastInsertId()
	loan.ID = int(id)

	// var lead lead.LeadInfo

	// results, err := db.Exec("INSERT INTO lead_table (loan_type, loan_amount, pincode, tenure, employment_type, gross_monthly_income, created_at, last_modified, user_id) VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW(), ?)",
	// 	lead.Loan_type, lead.Loan_amount, lead.Pincode, lead.Tenure, lead.Employment_type, lead.Gross_monthly_income, userID)
	// if err != nil {
	// 	logrus.Errorf("Error executing SQL query: %v", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	json.NewEncoder(w).Encode(&Response{Error: &Error{Message: "Internal Server Error"}})
	// 	return
	// }

	// ids, _ := results.LastInsertId()
	// loan.ID = int(ids)

	logrus.WithFields(logrus.Fields{
		"loan_id":              loan.ID,
		"Loan_Type":            loan.Loan_type,
		"Loan_Amount":          loan.Loan_amount,
		"Tenure":               loan.Tenure,
		"Pincode":              loan.Pincode,
		"Employment_Type":      loan.Employment_type,
		"Gross_Monthly_Income": loan.Gross_monthly_income,
	}).Info("Loan Details Inserted")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&Response{Data: loan})
}

func UpdateLoanDetails(w http.ResponseWriter, r *http.Request) {

	config, err := con.LoadConfig("Config/config.yaml")
	if err != nil {
		panic(err)
	}

	con.ConnectDB(config)
	defer con.CloseDB()

	params := mux.Vars(r)
	LoanID := params["id"]

	var loan Loan_details
	//fmt.Println("r.body", r.Body)
	error := json.NewDecoder(r.Body).Decode(&loan)
	if error != nil {
		http.Error(w, error.Error(), http.StatusBadRequest)
		return
	}
	result, err := db.Exec("UPDATE loan_details_table SET loan_type=?, loan_amount=?, tenure=?, Pincode=?, employment_type=?, gross_monthly_income=?,remark=?,admin_name=? WHERE id=?", loan.Loan_type, loan.Loan_amount, loan.Tenure, loan.Pincode, loan.Employment_type, loan.Gross_monthly_income, loan.Remark, loan.Admin_Name, LoanID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Loan ID not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Loan ID updated successfully")
	logrus.WithFields(logrus.Fields{
		"loan_id":              LoanID,
		"Loan_Type":            loan.Loan_type,
		"Loan_Amount":          loan.Loan_amount,
		"Tenure":               loan.Tenure,
		"Pincode":              loan.Pincode,
		"Employment_Type":      loan.Employment_type,
		"Gross_Monthly_Income": loan.Gross_monthly_income,
	}).Info("Update Loan Table Successfully")

	// err = templates.ExecuteTemplate(w, "http://localhost:9000/form/update_loan.html", loan)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

}

func DeleteLoanDetails(w http.ResponseWriter, r *http.Request) {

	config, err := con.LoadConfig("Config/config.yaml")
	if err != nil {
		panic(err)
	}

	con.ConnectDB(config)
	defer con.CloseDB()

	params := mux.Vars(r)
	loanID := params["id"]

	result, err := db.Exec("DELETE FROM loan_details_table WHERE id=?", loanID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Loan not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "loan deleted successfully")
	logrus.WithFields(logrus.Fields{
		"loan_id": loan.ID,
	}).Warnln("Delete Successfully")
}

//var templates = template.Must(template.ParseFiles("form/loan.html"))

func GetLoanDetails(w http.ResponseWriter, r *http.Request) {
	config, err := con.LoadConfig("Config/config.yaml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	con.ConnectDB(config)
	logrus.Infoln("Connection Start")
	defer con.CloseDB()

	params := mux.Vars(r)
	loanID := params["id"]

	rows, err := db.Query("SELECT id, loan_type, loan_amount, tenure, pincode, employment_type, gross_monthly_income, status, created_at, last_modified,remark,admin_name FROM loan_details_table WHERE id=?", loanID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	var createdAtStr string
	var lastAtStr string
	var loans []Loan_details

	for rows.Next() {
		var loan Loan_details
		//var statusText string
		err = rows.Scan(&loan.ID, &loan.Loan_type, &loan.Loan_amount, &loan.Tenure, &loan.Pincode, &loan.Employment_type, &loan.Gross_monthly_income, &loan.Status, &createdAtStr, &lastAtStr, &loan.Remark, &loan.Admin_Name)
		if err != nil {
			logrus.Debugln("Data Not Fatch")
		}
		createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		lastAt, err := time.Parse("2006-01-02 15:04:05", lastAtStr)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		loan.Created_at = createdAt
		loan.Last_modified = lastAt

		loans = append(loans, loan)
	}

	response := Response{
		Data: loans,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	//fmt.Fprint(w, "Lead ID updated successfully")
	logrus.WithFields(logrus.Fields{
		"loan_id":              loanID,
		"Loan_Type":            loan.Loan_type,
		"Loan_Amount":          loan.Loan_amount,
		"Tenure":               loan.Tenure,
		"Pincode":              loan.Pincode,
		"Employment_Type":      loan.Employment_type,
		"Gross_Monthly_Income": loan.Gross_monthly_income,
	}).Infoln("Loan Details Fetach Successfully")

}

//retrieve all data from loan

func LoanIndex(w http.ResponseWriter, r *http.Request) {
	config, err := con.LoadConfig("Config/config.yaml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	con.ConnectDB(config)
	defer con.CloseDB()

	rows, err := db.Query("SELECT loan_details_table.id, loan_details_table.loan_type, loan_details_table.loan_amount, loan_details_table.tenure, loan_details_table.pincode, loan_details_table.employment_type, loan_details_table.gross_monthly_income, loan_details_table.status, loan_details_table.created_at, loan_details_table.last_modified,loan_details_table.remark, admin.username FROM loan_details_table JOIN admin ON loan_details_table.admin_name = admin.id WHERE is_delete = 0 ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var loans []Loan_details

	for rows.Next() {

		var loan Loan_details
		var CreatedAtstr, LastAtStr string

		err = rows.Scan(&loan.ID, &loan.Loan_type, &loan.Loan_amount, &loan.Tenure, &loan.Pincode, &loan.Employment_type, &loan.Gross_monthly_income, &loan.Status, &CreatedAtstr, &LastAtStr, &loan.Remark, &loan.Admin_Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// loan.Status = getStatusValue(statusText)
		CreatedAt, err := time.Parse("2006-01-02 15:04:05", CreatedAtstr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		loan.Created_at = CreatedAt
		LastAt, err := time.Parse("2006-01-02 15:04:05", LastAtStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		loan.Last_modified = LastAt

		loans = append(loans, loan)
	}

	response := Response{
		Data: loans,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}

func LoanCount(w http.ResponseWriter, r *http.Request) {
	config, err := con.LoadConfig("./Config/config.yaml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	con.ConnectDB(config)
	defer con.CloseDB()

	//db := con.GetDB()

	rows, err := db.Query("SELECT COUNT(*) as counting FROM loan_details_table")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var counting int

	for rows.Next() {
		err = rows.Scan(&counting)
		if err != nil {
			panic(err)
		}
	}

	// loans := make([]loan_info, 0)
	// loan := loan_info{} // Assuming loan_info is a struct type

	// loans = append(loans, loan)

	// response := Response{
	// 	Data: loans,
	// }

	jsonResponse, err := json.Marshal(counting)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

//admin soft delete

func AdminSoftDeleteLoan(w http.ResponseWriter, r *http.Request) {
	config, err := con.LoadConfig("Config/config.yaml")
	if err != nil {
		panic(err)
	}

	con.ConnectDB(config)
	defer con.CloseDB()

	params := mux.Vars(r)
	LeadID := params["id"]

	if r.Method != http.MethodPatch {
		w.Header().Set("Allow", http.MethodPatch)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Update is_delete field
	result, err := db.Exec("UPDATE loan_details_table SET is_delete=1 WHERE id=?", LeadID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Lead ID not found", http.StatusNotFound)
		return
	}
	// w.WriteHeader(http.StatusOK)
	// fmt.Fprintf(w, "Loan ID updated successfully")

}
func LoanIndex1(w http.ResponseWriter, r *http.Request) {
	config, err := con.LoadConfig("Config/config.yaml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	con.ConnectDB(config)
	defer con.CloseDB()

	rows, err := db.Query("SELECT id, loan_type, loan_amount, tenure, pincode, employment_type, gross_monthly_income, status, created_at, last_modified FROM loan_details_table")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var loans []Loan_details

	for rows.Next() {
		var loan Loan_details
		var CreatedAtstr, LastAtStr string

		err = rows.Scan(&loan.ID, &loan.Loan_type, &loan.Loan_amount, &loan.Tenure, &loan.Pincode, &loan.Employment_type, &loan.Gross_monthly_income, &loan.Status, &CreatedAtstr, &LastAtStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// loan.Status = getStatusValue(statusText)
		CreatedAt, err := time.Parse("2006-01-02 15:04:05", CreatedAtstr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		loan.Created_at = CreatedAt
		LastAt, err := time.Parse("2006-01-02 15:04:05", LastAtStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		loan.Last_modified = LastAt

		loans = append(loans, loan)
	}

	response := Response{
		Data: loans,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}
