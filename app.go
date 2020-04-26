package main

import (
	"log"
	"net/http"
	"os"

	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type Users3 struct {
	ID       uint `gorm:"primary_key"`
	NAME     string
	EMAIL    string
	PASSWORD string
}

type login struct {
	EMAIL    string
	PASSWORD string
}

var Db *gorm.DB
var server_port string
var router *gin.Engine
var password_DB_list login
var email_DB_list login

//load environment variables
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

//hash password before storing it in the database
func HashPassword(password string) (string, error) {
	hashed_password_in_bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hashed_password_in_bytes), err
}

//check passowrd against the hash stored in the database
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//enter users in the database
func EnterIntoDB(name, email, password string) bool {
	hashed_password, err := HashPassword(password) //store hashed password in database
	if err != nil {
		panic(err)
	}
	newUser := Users3{NAME: name, EMAIL: email, PASSWORD: hashed_password}
	Db.NewRecord(newUser)
	Db.Create(&newUser)
	return true
}

//generates access token
func GenerateAccessToken() (string, error) {
	payload := `{"security":"OAuth 2.0"}`
	sharedKey := []byte{99, 75, 63}
	token, err := jose.Sign(payload, jose.HS256, sharedKey) //using HS256 algorithm for creating JWT
	return token, err
}

//checks the entered user against existing users in database
func CheckForExistingUser(email, password string) bool {
	Db.Table("users3").Select("EMAIL").Where("EMAIL = ?", email).Scan(&email_DB_list)
	email_DB := email_DB_list.EMAIL
	if email_DB != "" {
		Db.Table("users3").Select("PASSWORD").Where("EMAIL = ?", email).Scan(&password_DB_list)
		password_DB := password_DB_list.PASSWORD
		return CheckPasswordHash(password, password_DB)
	} else {
		return false
	}
}

//add new users in database
func AddUsers(c *gin.Context) {
	name := c.PostForm("NAME")
	email := c.PostForm("EMAIL")
	password := c.PostForm("PASSWORD")
	if name != "" && email != "" && password != "" {
		if EnterIntoDB(name, email, password) {
			c.String(200, "Added new user successfully")
		}
	} else {
		panic("failed to add new user")
	}
}

//authenticate existing users and return access tokens to authenticated users
func AuthenticateUsers(c *gin.Context) {
	email := c.PostForm("EMAIL")
	password := c.PostForm("PASSWORD")
	if email != "" && password != "" {
		if CheckForExistingUser(email, password) {
			token, err := GenerateAccessToken()
			if err == nil {
				c.JSON(http.StatusOK, token)
			} else {
				panic("failed to generate token")
			}
		} else {
			panic("user does not exist")
		}
	} else {
		panic("fields are empty")
	}
}

//give access only to authorised users which have an access token
func HomeAccess(c *gin.Context) {
	token := c.Request.Header.Get("token")
	sharedKey := []byte{99, 75, 63}
	payload, _, err := jose.Decode(token, sharedKey)
	if err == nil {
		c.String(http.StatusOK, "Access Granted "+payload) //Granting access only to authorised users
	} else {
		c.AbortWithStatus(401)
	}
}

func main() {
	//get all the environment variables
	db_host, host_exists := os.LookupEnv("DB_HOST")
	db_port, port_exists := os.LookupEnv("DB_PORT")
	db_user, user_exists := os.LookupEnv("DB_USER")
	db_name, name_exists := os.LookupEnv("DB_NAME")
	db_password, password_exists := os.LookupEnv("DB_PASSWORD")
	var server_port_exists bool
	server_port, server_port_exists = os.LookupEnv("SERVER_PORT")

	if host_exists && port_exists && user_exists && name_exists && password_exists && server_port_exists {
		var err error
		Db, err = gorm.Open("postgres", "host="+db_host+" port="+db_port+" user="+db_user+" dbname="+db_name+" password="+db_password)
		if err != nil {
			panic("failed to connect to database")
		}
		defer Db.Close()

		Db.AutoMigrate(&Users3{})

		router := gin.Default()

		router.POST("/addUser", AddUsers)
		router.POST("/oauth", AuthenticateUsers)
		router.GET("/home", HomeAccess)

		router.Run(":" + server_port)
	}
}
