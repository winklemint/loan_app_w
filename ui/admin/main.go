// Amitabh  sir modified response json struct
package main

import (
	//control "admin/controller"

	con "admin/Config"

	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var db *sql.DB

var tmpl = template.Must(template.ParseGlob("form/*.html"))

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the URL from the request
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	// Create a new request to the external API
	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	// Set headers from the original request to the new request
	req.Header = r.Header

	// Make a request to the external API
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the appropriate headers to enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	// Allow the following headers
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Referer, Sec-Ch-Ua, Sec-Ch-Ua-Mobile, Sec-Ch-Ua-Platform, User-Agent")

	// Handle preflight OPTIONS requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Copy the response headers to the actual response writer
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Write the response body to the actual response writer
	w.WriteHeader(resp.StatusCode)
	_, err = w.Write(body)
	if err != nil {
		log.Println(err)
	}
}

func TemplatePage(w http.ResponseWriter, r *http.Request) {

	tmpl.ExecuteTemplate(w, "admin.html", r)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "admin.html", r)
	//http.FileServer(http.Dir("form"))
}

func LandingPage(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index.html", r)
}

func main() {

	config, err := con.LoadConfig("Config/config.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	db, err = con.ConnectDB(config)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	defer db.Close()

	file, err := os.OpenFile("logrus.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logrus.log")
		panic(err)
	}
	logrus.SetOutput(file)
	logrus.SetLevel(logrus.TraceLevel)
	r := mux.NewRouter()

	protectedRoutes := r.PathPrefix("/").Subrouter()
	protectedRoutes.Use(IsAuthenticated)

	// Define routes
	protectedRoutes.PathPrefix("/form").Handler(http.StripPrefix("/form/", http.FileServer(http.Dir("form")))) //localhost:9000/form/ endpoint for accessing forms
	// r.HandleFunc("/index", TemplatePage).Methods("GET")

	r.HandleFunc("/proxy", ProxyHandler).Methods("GET", "POST", "PATCH", "DELETE", "OPTIONS")
	r.HandleFunc("/admin/login", LoginPage).Methods("GET")
	//r.HandleFunc("/landing", LandingPage).Methods("GET")
	r.HandleFunc("/login", AdminLoginHandler).Methods("POST")
	r.HandleFunc("/subadmin/get", GetAdmin).Methods("GET")
	r.HandleFunc("/get/admin", GetAdminByID).Methods("GET")
	r.HandleFunc("/admin/logout", LogOut).Methods("DELETE")
	r.HandleFunc("/admin/add", AdminInsert).Methods("POST")
	r.HandleFunc("/admin/update/{id}", UpadteAdmin).Methods("POST")
	r.HandleFunc("/admin/delete/{id}", AdminDelete).Methods("DELETE")
	r.HandleFunc("/get/{id}", AdminById).Methods("GET")

	// Sleep for a while to allow time for the API requests to complete
	//time.Sleep(5 * time.Second)
	// Create a reverse proxy for the React application
	// Register the /react route with the POST method and ProxyHandlerReact as the handler
	//r.HandleFunc("/", ProxyHandlerReact).Methods(http.MethodPost)

	// Create a reverse proxy for the React application
	reactURL, _ := url.Parse("http://localhost:3000")
	proxy := httputil.NewSingleHostReverseProxy(reactURL)

	// Register the reverse proxy as the default handler
	r.PathPrefix("/react").Handler(proxy)

	fmt.Println("Data fetched from the backend successfully.")

	http.ListenAndServe(":9000", r) // Start the server and listen for incoming requests
}

// Define a custom response writer to capture the output of template execution
type responseWriter struct {
	body []byte
}

func (rw *responseWriter) Header() http.Header {
	return make(http.Header)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	rw.body = data

	return len(data), nil
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	// Do nothing
}
func (rw *responseWriter) WriteTo(w http.ResponseWriter) {
	_, err := w.Write(rw.body)
	if err != nil {
		log.Fatal(err)
	}
}
func ProxyHandlerReact(w http.ResponseWriter, r *http.Request) {
	// Extract the URL from the request
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	// Create a new request to the external API
	req, err := http.NewRequest(http.MethodPost, url, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	// Set headers from the original request to the new request
	req.Header = r.Header

	// Make a request to the external API
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the appropriate headers to enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	// Allow the following headers
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Referer, Sec-Ch-Ua, Sec-Ch-Ua-Mobile, Sec-Ch-Ua-Platform, User-Agent")

	// Handle preflight OPTIONS requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Copy the response headers to the actual response writer
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Write the response body to the actual response writer
	w.WriteHeader(resp.StatusCode)
	_, err = w.Write(body)
	if err != nil {
		log.Println(err)
	}
}
