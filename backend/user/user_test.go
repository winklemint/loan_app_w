package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func TestGetUserById(t *testing.T) {
	req, err := http.NewRequest("GET", "/user/21", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	r.HandleFunc("/user/{id}", GetUserById).Methods("GET")
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"id":21,"name":"Jane","email":"jane@example.com","contact":9876543210,"created at":"0001-01-01T00:00:00Z","last modified":"0001-01-01T00:00:00Z"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestAddUser(t *testing.T) {
	newUser := User_info{Name: "Jane", Email: "jane@example.com", Contact_Num: 9876543210, Password: "password"}
	body, err := json.Marshal(newUser)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/user", strings.NewReader(string(body)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	r.HandleFunc("/user", AddUser).Methods("POST")
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := "New student has been created with ID:"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestUpdateUser(t *testing.T) {
	updatedUser := User_info{Name: "Jane Doe", Email: "jane.doe@example.com", Contact_Num: 9876543210, Password: "newpassword", Created_At: time.Now(), Last_Modified: time.Now()}
	body, err := json.Marshal(updatedUser)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("PUT", "/user/13", strings.NewReader(string(body)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	r.HandleFunc("/user/{id}", UpdateUser).Methods("PUT")
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	fmt.Println(string(body))

	expected := "User updated successfully"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestDeleteUser(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/users/36", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	r.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
