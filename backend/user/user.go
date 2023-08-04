package user

import (
	con "backend/Config"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"

	"github.com/go-playground/validator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type User_info struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Contact_Num int    `json:"contact"`
	Password    string `json:"password,omitempty"`
	Created_At    time.Time `json:"created_at,omitempty"`
	Last_Modified time.Time `json:"last_modified,omitempty"`
}

func init() {
	var err error
	_, err = con.GetDB()
	if err != nil {
		panic(fmt.Errorf("failed to initialize database: %v", err))
	}
}

type Error struct {
	Message string `json:"message"`
}
type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error *Error      `json:"error,omitempty"`
}

func isEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
	return emailRegex.MatchString(email)
}

func ValidateUser(user User_info, validate *validator.Validate) error {
	// Check that user name is not empty
	if user.Name == "" {
		return errors.New("user name is required.")
	}
	// Check that phone no. is exactly 10 digits
	if match, _ := regexp.MatchString(`^\d{10}$`, fmt.Sprint(user.Contact_Num)); !match {
		return errors.New("Phone number must be exactly 10 digits.")
	}

	// Check that password is not empty
	if user.Password == "" {
		return errors.New("Password is required.")
	}
	// Check that email is valid
	if !isEmailValid(user.Email) {
		return errors.New("Invalid email format.")
	}
	// If all checks pass, return nill
	return nil
}

// func GetAllUser(w http.ResponseWriter, r *http.Request) {
// 	config, err := con.LoadConfig("Config/config.yaml")
// 	if err != nil {
// 		panic(err)
// 	}

// 	con.ConnectDB(config)
// 	defer con.CloseDB()

// 	rows, err := db.Query("SELECT id, user_name, user_email, user_contact_num,created_at,last_modified FROM user_info")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	var users []User_info

// 	for rows.Next() {
// 		var user User_info
// 		var createdAt string
// 		var lastAt string

// 		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Contact_Num, &createdAt, &lastAt)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		Createdstr, err := time.Parse("2006-01-02 15:04:05", createdAt)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 		user.Created_At = Createdstr
// 		laststr, err := time.Parse("2006-01-02 15:04:05", lastAt)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 		user.Created_At = laststr

// 		users = append(users, user)
// 	}
// 	response := Response{
// 		Data: users,
// 	}

// 	jsonResponse, err := json.Marshal(response)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	// w.WriteHeader(http.StatusOK)
// 	w.Write(jsonResponse)

// 	// if err := rows.Err(); err != nil {
// 	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
// 	// 	return
// 	// }
// 	// Process the retrieved users as needed
// 	// ...
// }

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	db, err := con.GetDB()
	if err != nil {
		panic(fmt.Errorf("failed to initialize database: %v", err))
	}
	rows, err := db.Query("SELECT id, user_name, user_email, user_contact_num,created_at,last_modified FROM user_info")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User_info

	for rows.Next() {
		var user User_info
		var CreatedAtstr, LastAtStr string
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Contact_Num, &CreatedAtstr, &LastAtStr)

		CreatedAt, err := time.Parse("2006-01-02 15:04:05", CreatedAtstr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		user.Created_At = CreatedAt
		LastAt, err := time.Parse("2006-01-02 15:04:05", LastAtStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		user.Last_Modified = LastAt

		users = append(users, user)
	}

	response := Response{
		Data: users,
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
func GetUserById(w http.ResponseWriter, r *http.Request) {
	db, err := con.GetDB()
	if err != nil {
		panic(fmt.Errorf("failed to initialize database: %v", err))
	}

	params := mux.Vars(r)
	userID := params["id"]
	var user User_info

	error := db.QueryRow("SELECT id, user_name, user_email, user_contact_num ,created_at,last_modified  FROM user_info WHERE id=?", userID).Scan(&user.Id, &user.Name, &user.Email, &user.Contact_Num, &user.Created_At, &user.Last_Modified)
	if err != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return

	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userJSON)
}

// func Admin(w http.ResponseWriter, r *http.Request) {

// 	config, err := con.LoadConfig("Config/config.yaml")
// 	if err != nil {
// 		panic(err)
// 	}

// 	con.ConnectDB(config)
// 	defer con.CloseDB()
// 	tmpl, err := template.ParseGlob("../Views/admin.html/")

// 	w.Write([]byte(err.Error()))
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	err = tmpl.Execute(w, nil)
// 	w.Write([]byte(err.Error()))
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

func AddUser(w http.ResponseWriter, r *http.Request) {
	db, err := con.GetDB()
	if err != nil {
		panic(fmt.Errorf("failed to initialize database: %v", err))
	}
	var user User_info
	errors := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, errors.Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()

	// validate loan details
	if err := ValidateUser(user, validate); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&Response{Error: &Error{Message: err.Error()}})
		return
	}
	maxPasswordLength := 120 // Modify this value based on your database column's maximum length

	if len(user.Password) > maxPasswordLength {
		user.Password = user.Password[:maxPasswordLength]
	}

	// Encrypt the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hashedPasswordStr := hex.EncodeToString(hashedPassword)

	// Store the fixed-length representation of the hashed password
	result, err := db.Exec("INSERT INTO user_info (user_name, user_email, user_contact_num, user_password, created_at, last_modified) VALUE (?, ?, ?, ?, NOW(), NOW())", user.Name, user.Email, user.Contact_Num, hashedPasswordStr)
	if err != nil {
		log.Fatal(err)
	}
	otp := "http://127.0.0.1:3000/login"

	storeOTP(user.Email, otp)

	err = sendOTPEmail1(user.Email, otp)
	if err != nil {
		http.Error(w, "Failed to send OTP via email", http.StatusInternalServerError)
		return

	}
	fmt.Fprint(w, "OTP sent successfully!")

	id, _ := result.LastInsertId()
	fmt.Fprintf(w, "New User has been created with ID: %d", id)
	logrus.WithFields(logrus.Fields{
		"User_name":     user.Name,
		"User_Email":    user.Email,
		"User_Contact":  user.Contact_Num,
		"User_Password": "Encrypted",
	}).Infoln("Added New User Successfully")
}
func sendOTPEmail1(email, otp string) error {

	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	smtpUsername := "raksha.rajput98@gmail.com"
	smtpPassword := "gmrzimywbfdtxslt"

	m := gomail.NewMessage()
	m.SetHeader("From", "raksha.rajput98@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Verification Successful - Welcome to Our Platform!")
	m.SetBody("text/plain", "Hi, \n We are delighted to inform you that your account verification was successful! You are now officially a member of our platform. \n Get started, simply click on the link below to log in to your account:"+otp)

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUsername, smtpPassword)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)

	}

	fmt.Println("Email sent successfully")

	return nil
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	db, err := con.GetDB()
	if err != nil {
		panic(fmt.Errorf("failed to initialize database: %v", err))
	}
	params := mux.Vars(r)
	userID := params["id"]

	result, err := db.Exec("DELETE FROM user_info WHERE id=?", userID)
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
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User deleted successfully")
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db, err := con.GetDB()
	if err != nil {
		panic(fmt.Errorf("failed to initialize database: %v", err))
	}
	params := mux.Vars(r)
	userID := params["id"]

	var user User_info
	errr := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, errr.Error(), http.StatusBadRequest)
		return
	}
	result, err := db.Exec("UPDATE user_info SET user_name=?, user_email=?, user_contact_num=?, user_password=?, last_modified=NOW() WHERE id=?", user.Name, user.Email, user.Contact_Num, user.Password, userID)
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
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)

	logrus.WithFields(logrus.Fields{
		"User_Id":       user.Id,
		"User_Name":     user.Name,
		"User_Email":    user.Email,
		"User_Contact":  user.Contact_Num,
		"User_Password": user.Password,
	}).Info("Upadte User Details")

}

func UserCount(w http.ResponseWriter, r *http.Request) {
	db, err := con.GetDB()
	if err != nil {
		panic(fmt.Errorf("failed to initialize database: %v", err))
	}

	//db := con.GetDB()

	rows, err := db.Query("SELECT COUNT(*) as counting FROM user_info")
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
func AdminUserSoftDelete(w http.ResponseWriter, r *http.Request) {
	db, err := con.GetDB()
	if err != nil {
		panic(fmt.Errorf("failed to initialize database: %v", err))
	}
	params := mux.Vars(r)
	UserID := params["id"]

	// Check if the request method is allowed
	if r.Method != http.MethodPatch {
		w.Header().Set("Allow", http.MethodPatch)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Update is_delete field
	result, err := db.Exec("UPDATE user_info SET is_delete=1 WHERE id=?", UserID)
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
		http.Error(w, "user ID not found", http.StatusNotFound)
		return
	}

}
