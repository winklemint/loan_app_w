package lead

import (
	con "backend/Config"
	user "backend/user"
	"log"
	"strings"

	//"backend/lead"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	//"strings"
	//"text/template"

	//"log"
	"regexp"

	"github.com/go-playground/validator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"

	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type LeadInfo struct {
	ID                   int       `json:"id,omitempty"`
	Loan_type            string    `json:"loan_type,omitempty"`
	Employment_type      string    `json:"employment_type,omitempty" `
	Loan_amount          float64   `json:"loan_amount,omitempty"`
	Gross_monthly_income float64   `json:"gross_monthly_income"`
	Pincode              int       `json:"pincode,omitempty"`
	Tenure               int       `json:"tenure,omitempty"`
	Status               string    `json:"status,omitempty"`
	Created_at           time.Time `json:"created_at,omitempty"`
	Last_modified        time.Time `json:"last_modified,omitempty"`
	Remark               string    `json:"remark"`
	Admin_Name           string    `json:"admin_name"`
}

var db *sql.DB
var lead LeadInfo

// var tmpl = template.Must(template.ParseGlob("form/*.html"))

func init() {
	var err error
	db, err = con.GetDB()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
}

func SetDB(database *sql.DB) {
	db = database

}

func ValidateLead(lead LeadInfo, validate *validator.Validate) error {
	// Check that loan_type is not empty
	if lead.Loan_type == "" {
		return errors.New("Loan type is required.")
	}

	// Check that loan_amount is greater than zero
	if lead.Loan_amount <= 0 {
		return errors.New("Loan amount must be greater than zero.")
	}

	// Check that tenure is greater than zero
	if lead.Tenure <= 0 {
		return errors.New("Tenure must be greater than zero.")
	}

	// Check that pincode is exactly 6 digits
	if match, _ := regexp.MatchString(`^\d{6}$`, fmt.Sprint(lead.Pincode)); !match {
		return errors.New("Pincode must be exactly 6 digits.")
	}

	// Check that employment_type is not empty
	if lead.Employment_type == "" {
		return errors.New("Employment type is required.")
	}

	if lead.Gross_monthly_income <= 0 {
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

func InsertLead(w http.ResponseWriter, r *http.Request) {
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

	var lead LeadInfo
	if err := json.NewDecoder(r.Body).Decode(&lead); err != nil {
		logrus.Errorf("Error decoding request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&Response{Error: &Error{Message: err.Error()}})
		return
	}

	// Create validator instance
	validate := validator.New()

	// Validate loan details
	if err := ValidateLead(lead, validate); err != nil {
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

	result, err := db.Exec("INSERT INTO lead_table (loan_type, employment_type, loan_amount, gross_monthly_income, pincode, tenure, created_at, last_modified, user_id) VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW(), ?)",
		lead.Loan_type, lead.Employment_type, lead.Loan_amount, lead.Gross_monthly_income, lead.Pincode, lead.Tenure, userID)
	if err != nil {
		logrus.Errorf("Error executing SQL query: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&Response{Error: &Error{Message: "Internal Server Error"}})
		return
	}

	id, _ := result.LastInsertId()
	lead.ID = int(id)

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
		"loan_id":              lead.ID,
		"Loan_Type":            lead.Loan_type,
		"Loan_Amount":          lead.Loan_amount,
		"Tenure":               lead.Tenure,
		"Pincode":              lead.Pincode,
		"Employment_Type":      lead.Employment_type,
		"Gross_Monthly_Income": lead.Gross_monthly_income,
	}).Info("Loan Details Inserted")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&Response{Data: lead})
}

func DeleteLead(w http.ResponseWriter, r *http.Request) {

	config, err := con.LoadConfig("./Config/config.yaml")
	if err != nil {
		panic(err)
	}

	con.ConnectDB(config)
	defer con.CloseDB()

	params := mux.Vars(r)
	leadID := params["id"]

	result, err := db.Exec("DELETE FROM lead_table WHERE id=?", leadID)
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
		http.Error(w, "Lead not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Lead ID updated successfully")
	logrus.WithFields(logrus.Fields{
		"lead_id": leadID,
	}).Warnln("Delete Successfullly")
}

// Lead Dashboard code
func LeadIndexAll(w http.ResponseWriter, r *http.Request) {
	config, err := con.LoadConfig("./Config/config.yaml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	con.ConnectDB(config)
	defer con.CloseDB()

	params := mux.Vars(r)
	leadID := params["id"]

	rows, err := db.Query("SELECT user_info.id, user_info.user_name, user_info.user_email, user_info.user_contact_num , lead_table.id,lead_table.loan_type, lead_table.loan_amount, lead_table.tenure, lead_table.pincode, lead_table.employment_type, lead_table.gross_monthly_income, lead_table.status, lead_table.created_at, lead_table.last_modified,lead_table.remark,admin.username FROM lead_table INNER JOIN user_info ON lead_table.user_id = user_info.id  JOIN admin ON lead_table.admin_id = admin.id WHERE lead_table.id=?", leadID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var data []map[string]interface{}
	for rows.Next() {
		var user user.User_info
		var lead LeadInfo
		var CreatedAtstr, lastAtStr string
		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Contact_Num, &lead.ID, &lead.Loan_type, &lead.Loan_amount, &lead.Tenure, &lead.Pincode, &lead.Employment_type, &lead.Gross_monthly_income, &lead.Status, &CreatedAtstr, &lastAtStr, &lead.Remark, &lead.Admin_Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		CreatedAt, err := time.Parse("2006-01-02 15:04:05", CreatedAtstr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		lead.Created_at = CreatedAt

		lastAt, err := time.Parse("2006-01-02 15:04:05", lastAtStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		lead.Last_modified = lastAt

		leadMap := map[string]interface{}{
			"id":                   lead.ID,
			"loan_type":            lead.Loan_type,
			"employment_type":      lead.Employment_type,
			"loan_amount":          lead.Loan_amount,
			"gross_monthly_income": lead.Gross_monthly_income,
			"pincode":              lead.Pincode,
			"tenure":               lead.Tenure,
			"status":               lead.Status,
			"created_at":           lead.Created_at,
			"user_id":              user.Id,
			"user_name":            user.Name,
			"user_email":           user.Email,
			"user_contact":         user.Contact_Num,
			"remark":               lead.Remark,
			"admin_name":           lead.Admin_Name,
		}

		data = append(data, leadMap)
	}

	jsonResponse, err := json.Marshal(Response{
		Data: data,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// func getStatusValue(statusText string) int {
// 	switch statusText {
// 	case "Pending":
// 		return 1
// 	case "Approved":
// 		return 2
// 	case "Declined":
// 		return 3
// 	default:
// 		return 0
// 	}
// }

// var templates = template.Must(template.ParseFiles("form/update.html"))

func UpdateLead(w http.ResponseWriter, r *http.Request) {

	db, err := con.GetDB()

	params := mux.Vars(r)
	LeadID := params["id"]

	// Check if the request method is allowed
	if r.Method != http.MethodPatch {
		w.Header().Set("Allow", http.MethodPatch)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var user user.User_info

	var lead LeadInfo
	//fmt.Println("r.Body ================== >", r.Body)
	err = json.NewDecoder(r.Body).Decode(&lead)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//fmt.Printf("%#v", lead)
	result, err := db.Exec("UPDATE lead_table INNER JOIN  user_info ON lead_table.user_id = user_info.id SET  lead_table.loan_type=?, lead_table.loan_amount=?, lead_table.tenure=?, lead_table.pincode=?, lead_table.employment_type=?,lead_table.gross_monthly_income=?,lead_table.status=?,lead_table.remark=?, user_info.user_name=?,user_info.user_email=?,user_info.user_contact_num=? WHERE lead_table.id=? ", lead.Loan_type, lead.Loan_amount, lead.Tenure, lead.Pincode, lead.Employment_type, lead.Gross_monthly_income, lead.Status, lead.Remark, user.Name, user.Email, user.Contact_Num, LeadID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//fmt.Println(rowsAffected)

	if rowsAffected == 0 {
		http.Error(w, "Lead ID not found", http.StatusNotFound)
		return
	} else {
		fmt.Println("lead update succesfully")
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Lead ID updated successfully")
	logrus.WithFields(logrus.Fields{
		"Lead_Id":              LeadID,
		"Lead_Type":            lead.Loan_type,
		"Loan_Amount":          lead.Loan_amount,
		"Tenure":               lead.Tenure,
		"Pincode":              lead.Pincode,
		"Employment_Type":      lead.Employment_type,
		"Gross_Monthly_Income": lead.Gross_monthly_income,
	}).Info("Update Lead Table Successfully")
}

// retrieve all leads
func LeadIndex(w http.ResponseWriter, r *http.Request) {
	config, err := con.LoadConfig("./Config/config.yaml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	con.ConnectDB(config)
	defer con.CloseDB()
	query := `SELECT user_info.id, user_info.user_name, user_info.user_email, user_info.user_contact_num , lead_table.id,lead_table.loan_type, lead_table.loan_amount, lead_table.tenure, lead_table.pincode, lead_table.employment_type, lead_table.gross_monthly_income, lead_table.status, lead_table.created_at,lead_table.remark,admin.username FROM lead_table INNER JOIN user_info ON lead_table.user_id = user_info.id  JOIN admin ON lead_table.admin_id = admin.id ; `
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var data []map[string]interface{}
	for rows.Next() {
		var user user.User_info
		var lead LeadInfo
		var CreatedAtstr string
		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Contact_Num, &lead.ID, &lead.Loan_type, &lead.Loan_amount, &lead.Tenure, &lead.Pincode, &lead.Employment_type, &lead.Gross_monthly_income, &lead.Status, &CreatedAtstr, &lead.Remark, &lead.Admin_Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		CreatedAt, err := time.Parse("2006-01-02 15:04:05", CreatedAtstr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		lead.Created_at = CreatedAt

		// Create a map for each lead and user, then append them to the data slice
		leadMap := map[string]interface{}{
			"id":                   lead.ID,
			"loan_type":            lead.Loan_type,
			"employment_type":      lead.Employment_type,
			"loan_amount":          lead.Loan_amount,
			"gross_monthly_income": lead.Gross_monthly_income,
			"pincode":              lead.Pincode,
			"tenure":               lead.Tenure,
			"status":               lead.Status,
			"created_at":           lead.Created_at,
			"user_id":              user.Id,
			"user_name":            user.Name,
			"user_email":           user.Email,
			"user_contact":         user.Contact_Num,
			"remark":               lead.Remark,
			"admin_name":           lead.Admin_Name,
		}

		data = append(data, leadMap)
	}

	jsonResponse, err := json.Marshal(Response{
		Data: data,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// Retrieving only top 5 ids

func LeadIndexTop6(w http.ResponseWriter, r *http.Request) {
	db, err := con.GetDB()
	rows, err := db.Query("SELECT lead_table.id, lead_table.loan_type, lead_table.loan_amount, lead_table.tenure, lead_table.pincode, lead_table.employment_type, lead_table.gross_monthly_income, lead_table.status, lead_table.created_at, lead_table.last_modified,user_info.id, user_info.user_name, user_info.user_email, user_info.user_contact_num FROM lead_table INNER JOIN user_info ON lead_table.user_id = user_info.id ORDER BY lead_table.id DESC LIMIT 6")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var data []map[string]interface{}

	for rows.Next() {
		var lead LeadInfo
		var user user.User_info
		var statusText, createdAtStr, lastAtStr string
		err = rows.Scan(&lead.ID, &lead.Loan_type, &lead.Loan_amount, &lead.Tenure, &lead.Pincode, &lead.Employment_type, &lead.Gross_monthly_income, &statusText, &createdAtStr, &lastAtStr, &user.Id, &user.Name, &user.Email, &user.Contact_Num)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		lead.Created_at = createdAt

		lastAt, err := time.Parse("2006-01-02 15:04:05", lastAtStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		lead.Last_modified = lastAt

		// Create a map for each lead and user, then append them to the data slice
		item := map[string]interface{}{

			"id":                   lead.ID,
			"loan_type":            lead.Loan_type,
			"loan_amount":          lead.Loan_amount,
			"tenure":               lead.Tenure,
			"pincode":              lead.Pincode,
			"employment_type":      lead.Employment_type,
			"gross_monthly_income": lead.Gross_monthly_income,
			"status":               statusText,
			"created_at":           createdAt,
			"last_modified":        lastAt,
			"user_id":              user.Id,
			"user_name":            user.Name,
			"user_email":           user.Email,
			"user_contact":         user.Contact_Num,
		}

		data = append(data, item)
	}

	jsonResponse, err := json.Marshal(Response{
		Data: data,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func LeadCount(w http.ResponseWriter, r *http.Request) {
	config, err := con.LoadConfig("./Config/config.yaml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	con.ConnectDB(config)
	defer con.CloseDB()

	//db := con.GetDB()

	rows, err := db.Query("SELECT COUNT(*) as counting FROM lead_table")
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

	// leads := make([]LeadInfo, 0)
	// lead := LeadInfo{} // Assuming LeadInfo is a struct type

	// leads = append(leads, lead)

	// response := Response{
	// 	Data: leads,
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

func AdminSoftDelete(w http.ResponseWriter, r *http.Request) {
	config, err := con.LoadConfig("Config/config.yaml")
	if err != nil {
		panic(err)
	}

	con.ConnectDB(config)
	defer con.CloseDB()

	params := mux.Vars(r)
	LeadID := params["id"]

	// Check if the request method is allowed
	if r.Method != http.MethodPatch {
		w.Header().Set("Allow", http.MethodPatch)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Update is_delete field
	result, err := db.Exec("UPDATE lead_table SET is_delete=1 WHERE id=?", LeadID)
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

}

func LeadGraph(w http.ResponseWriter, r *http.Request) {
	config, err := con.LoadConfig("./Config/config.yaml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	con.ConnectDB(config)
	defer con.CloseDB()

	rows, err := db.Query("SELECT status FROM lead_table")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var statuses []string

	for rows.Next() {
		var status string
		err = rows.Scan(&status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		statuses = append(statuses, status)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(statuses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
