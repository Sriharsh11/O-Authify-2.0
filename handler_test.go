package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

//load environment variables
func Init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

var test_name, test_email, test_password string

func TestEnterIntoDB(t *testing.T) {
	test_name, test_name_exists := os.LookupEnv("DUMMY_NAME")
	test_email, test_email_exists := os.LookupEnv("DUMMY_EMAIL")
	test_password, test_password_exists := os.LookupEnv("DUMMY_PASSWORD")
	if test_name_exists && test_email_exists && test_password_exists {
		db, mock, err := sqlmock.New()
		hashed_test_password, _ := HashPassword(test_password)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO USERS3 (NAME, EMAIL, PASSWORD) VALUES (` + test_name + `,` + test_email + `,` + hashed_expected_password + `)`).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		if err := EnterIntoDB(test_name, test_email, test_password); err != true {
			t.Errorf("error was not expected while updating stats")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}

func TestCheckForExistingUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectBegin()
	mock.ExpectExec(`SELECT EMAIL FROM USERS3 WHERE EMAIL=` + test_email).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	if err := CheckForExistingUser(test_email, test_password); err != true {
		t.Errorf("error was not expected while updating stats")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAddUsers(t *testing.T) {
	w := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/addUser", AddUsers)
	params := url.Values{}
	params.Add("NAME", test_name)
	params.Add("EMAIL", test_email)
	params.Add("PASSWORD", test_password)
	registrationPayload := params.Encode()
	req, err := http.NewRequest("POST", "/addUser", strings.NewReader(registrationPayload))
	if err != nil {
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(registrationPayload)))
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fail()
	}
}

func TestAuthenticateUsers(t *testing.T) {
	w := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/oauth", AuthenticateUsers)
	params := url.Values{}
	params.Add("EMAIL", test_email)
	params.Add("PASSWORD", test_password)
	registrationPayload := params.Encode()
	req, err := http.NewRequest("POST", "/oauth", strings.NewReader(registrationPayload))
	if err != nil {
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(registrationPayload)))
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fail()
	}
}

func TestHomeAccess(t *testing.T) {
	w := httptest.NewRecorder()
	router := gin.Default()
	router.GET("/oauth", HomeAccess)
	token := GenerateAccessToken(test_email, test_password)
	req, err := http.NewRequest("GET", "/home", nil)
	if err != nil {
		t.Fail()
	}
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fail()
	}
}
