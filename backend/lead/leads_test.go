// package lead

// import (
// 	//con "backend/Config"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	"github.com/gorilla/mux"
// )

// func TestLeadIndex(t *testing.T) {

// 	// config, err := con.LoadConfig("Config/config.yaml")
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// con.ConnectDB(config)
// 	// defer con.CloseDB()

// 	req, err := http.NewRequest("GET", "/lead/191", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	rr := httptest.NewRecorder()
// 	r := mux.NewRouter()
// 	r.HandleFunc("/lead/{id}", LeadIndex).Methods("GET")
// 	r.ServeHTTP(rr, req)
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusOK)
// 	}
// 	expected := `{"id":1,"created_at":"0001-01-01T00:00:00Z","last_modified":"0001-01-01T00:00:00Z","status":"testtestetst"}`
// 	if rr.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %v want %v",
// 			rr.Body.String(), expected)
// 	}
// }
// func TestInsertLead(t *testing.T) {

// 	// config, err := con.LoadConfig("Config/config.yaml")
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// con.ConnectDB(config)
// 	// defer con.CloseDB()
// 	newLead := LeadInfo{Status: 1}
// 	body, err := json.Marshal(newLead)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req, err := http.NewRequest("POST", "/leads", strings.NewReader(string(body)))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req.Header.Set("Content-type", "application/json")
// 	rr := httptest.NewRecorder()
// 	r := mux.NewRouter()
// 	r.HandleFunc("/leads", InsertLead).Methods("POST")
// 	r.ServeHTTP(rr, req)
// 	if status := rr.Code; status != http.StatusCreated {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusCreated)
// 	}
// }

// func TestUpdateLead(t *testing.T) {
// 	leadID := "92"
// 	updateLead := LeadInfo{Status: 2}

// 	// Create request body
// 	body, err := json.Marshal(updateLead)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Create HTTP PUT request with lead ID parameter
// 	req, err := http.NewRequest("PATCH", "/update/"+leadID, strings.NewReader(string(body)))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Set request headers
// 	req.Header.Set("Content-type", "application/json")

// 	// Create router and HTTP response recorder
// 	rr := httptest.NewRecorder()
// 	r := mux.NewRouter()
// 	r.HandleFunc("/update/{id}", UpdateLead).Methods("PATCH")

// 	// Make HTTP PUT request and record response
// 	r.ServeHTTP(rr, req)

// 	// Check response status code
// 	if status := rr.Code; status < 200 || status > 299 {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusOK)
// 	}

// 	// Check response body
// 	expected := "Lead updated successfully"
// 	if rr.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %v want %v",
// 			rr.Body.String(), expected)
// 	}
// }

//	func TestDeleteLead(t *testing.T) {
//		req, err := http.NewRequest("DELETE", "/u/191", nil)
//		if err != nil {
//			t.Fatal(err)
//		}
//		rr := httptest.NewRecorder()
//		r := mux.NewRouter()
//		r.HandleFunc("/u/{id}", DeleteLead).Methods("DELETE")
//		r.ServeHTTP(rr, req)
//		if status := rr.Code; status != http.StatusOK {
//			t.Errorf("handler returned wrong status code: got %v want %v",
//				status, http.StatusOK)
//		}
//	}
package lead_test

import (
	con "backend/Config"
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/lead"

	"github.com/stretchr/testify/assert"
)

var db *sql.DB

func init() {
	var err error
	db, err = con.GetDB()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
}

func TestInsertLead(t *testing.T) {

	// Load configuration from config.yaml
	config, err := con.LoadConfig("./Config/config.yaml")
	assert.NoError(t, err)

	// Establish database connection
	db, err := con.ConnectDB(config)
	assert.NoError(t, err)
	defer db.Close()

	// Set the database connection for the lead package
	lead.SetDB(db)

	// Create a test HTTP server
	ts := httptest.NewServer(http.HandlerFunc(lead.InsertLead))
	defer ts.Close()

	// Prepare a sample lead
	leadData := lead.LeadInfo{
		Loan_type:            "Home Loan",
		Loan_amount:          50000,
		Tenure:               12,
		Pincode:              123456,
		Employment_type:      "Full-time",
		Gross_monthly_income: 30000,
	}

	// Convert leadData to JSON
	jsonStr, _ := json.Marshal(leadData)

	// Send POST request to the InsertLead endpoint
	resp, err := http.Post(ts.URL+"/lead", "application/json", bytes.NewBuffer(jsonStr))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// TODO: Assert the response body and check the inserted lead data
}

func TestUpdateLead(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(lead.UpdateLead))
	defer ts.Close()

	// Prepare a sample lead
	leadData := lead.LeadInfo{
		ID:                   1, // Assuming lead with ID 1 exists
		Loan_type:            "Home Loan",
		Loan_amount:          50000,
		Tenure:               12,
		Pincode:              123456,
		Employment_type:      "Full-time",
		Gross_monthly_income: 30000,
	}

	// Convert leadData to JSON
	jsonStr, _ := json.Marshal(leadData)

	// Create a PUT request with the lead ID in the URL
	req, _ := http.NewRequest("PUT", ts.URL+"/lead/1", bytes.NewBuffer(jsonStr))
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// TODO: Assert the response body and check the updated lead data
}

func TestDeleteLead(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(lead.DeleteLead))
	defer ts.Close()

	// Create a DELETE request with the lead ID in the URL
	req, _ := http.NewRequest("DELETE", ts.URL+"/lead/116", nil)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// TODO: Assert the response body and check the deletion confirmation
}

func TestLeadIndex(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(lead.LeadIndexAll))
	defer ts.Close()

	// Create a GET request with the lead ID in the URL
	req, _ := http.NewRequest("GET", ts.URL+"/lead/116", nil)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// TODO: Assert the response body and check the lead details
}
