package main

import (
	con "backend/Config"
	"backend/lead"
	"backend/loan"
	"database/sql"
	"fmt"
	"os"

	//"html/template"

	user "backend/user"

	"log"
	"net/http"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

//  var tmpl = template.Must(template.ParseGlob("form/*.html"))

//	func TemplatePage(w http.ResponseWriter, r *http.Request) {
//		tmpl.ExecuteTemplate(w, "index.html", nil)
//	}
var Db *sql.DB

func main() {
	file, err := os.OpenFile("logrus.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logrus.log")
		panic(err)
	}

	logrus.SetOutput(file)
	logrus.SetLevel(logrus.TraceLevel)

	// Load the config file

	config, err := con.LoadConfig("Config/config.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	Db, err = con.ConnectDB(config)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	defer Db.Close()
	r := mux.NewRouter()

	// logrus.Info("info")

	//r.PathPrefix("/form/").Handler(http.StripPrefix("/form/", http.FileServer(http.Dir("form")))) //localhost:9000/form/ endpoint for accessing forms
	//r.HandleFunc("/", TemplatePage).Methods("GET")
	r.HandleFunc("/lead/get/{id}", lead.LeadIndexAll).Methods("GET")
	r.HandleFunc("/lead/add", lead.InsertLead).Methods("POST")
	r.HandleFunc("/lead/update/{id}", lead.UpdateLead).Methods("PATCH")
	r.HandleFunc("/lead/delete/{id}", lead.DeleteLead).Methods("DELETE")

	//Admin Panel's endpoint
	r.HandleFunc("/lead/admin", lead.LeadIndex).Methods("GET")
	r.HandleFunc("/lead/top6", lead.LeadIndexTop6).Methods("GET")
	r.HandleFunc("/lead/count", lead.LeadCount).Methods("GET")
	r.HandleFunc("/admin/delete/{id}", lead.AdminSoftDelete).Methods("PATCH")
	r.HandleFunc("/lead/graph", lead.LeadGraph).Methods("GET")

	// Register the route for the GET request
	//r.HandleFunc("/lead_data", lead.AdminSoftDelete).Methods(http.MethodGet)

	//r.HandleFunc("/lead/upd/{id}", lead.UpdateLeadTest).Methods("PATCH")
	//r.HandleFunc("/lead/admin/update/{id}", lead.UpdateLeadDashboard).Methods("PATCH")
	//r.HandleFunc("/lead/handle/{id}", lead.HandleUpdateLead).Methods("PATCH")

	/////////////////////////////////////////////////////////////////////

	r.HandleFunc("/loan/index/{id}", loan.GetLoanDetails).Methods("GET")
	r.HandleFunc("/loan/insert", loan.InsertLoanDetails).Methods("POST")
	r.HandleFunc("/loan/update/{id}", loan.UpdateLoanDetails).Methods("PATCH")
	r.HandleFunc("/loan/delete/{id}", loan.DeleteLoanDetails).Methods("DELETE")
	r.HandleFunc("/loan/all", loan.LoanIndex1).Methods("GET")

	//all loan details for admin panel
	r.HandleFunc("/loan/get", loan.LoanIndex).Methods("GET")
	r.HandleFunc("/loan/count", loan.LoanCount).Methods("GET")
	r.HandleFunc("/admin/soft/delete/{id}", loan.AdminSoftDeleteLoan).Methods("PATCH")

	r.HandleFunc("/user/add", user.AddUser).Methods("POST")
	r.HandleFunc("/user/get/{id}", user.GetUserById).Methods("GET")
	r.HandleFunc("/user/update/{id}", user.UpdateUser).Methods("PATCH")
	r.HandleFunc("/user/delete/{id}", user.DeleteUser).Methods("DELETE")
	r.HandleFunc("/user/count", user.UserCount).Methods("GET")
	r.HandleFunc("/user/soft/delete/{id}", user.AdminUserSoftDelete).Methods("PATCH")
	r.HandleFunc("/user/all", user.GetAllUser).Methods("GET")
	// r.HandleFunc("/user/home", user.Home).Methods("GET")
	r.HandleFunc("/user/login", user.Login).Methods("POST")

	r.HandleFunc("/forgotpassword", user.ForgotPasswordHandler).Methods("POST")
	r.HandleFunc("/otpverify", user.VerifyOTPhandler).Methods("POST")
	r.HandleFunc("/resetpassword", user.ResetPasswordHandler).Methods("POST")
	r.HandleFunc("/user/logout", user.LogOut).Methods("DELETE")
	r.HandleFunc("/device/login", user.LoginWithDevice).Methods("POST")
	r.HandleFunc("/captcha", user.CaptchaHandler).Methods("GET")

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", r))

}
