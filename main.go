package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
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

//global variables used only because almost all functions need them
var Db *gorm.DB
var server_port string
var router *gin.Engine
var AtJwtKey []byte

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
		return false
	} else {
		newUser := Users3{NAME: name, EMAIL: email, PASSWORD: hashed_password}
		Db.NewRecord(newUser)
		Db.Create(&newUser)
		return true
	}
}

//generates access token
func GenerateAccessToken(email, password string) (string, error) {
	//expiration time is 60 minutes
	expirationTimeAccessToken := time.Now().Add(60 * time.Minute).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["exp"] = expirationTimeAccessToken
	claims["sub"] = password
	tokenString, err := token.SignedString(AtJwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

//checks the entered user against existing users in database
func CheckForExistingUser(email, password string) bool {
	var password_DB_list login
	var email_DB_list login
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
		c.AbortWithStatusJSON(406, gin.H{"status": 406, "message": "Failed to add user"})
	}
}

//authenticate existing users and return access tokens to authenticated users
func AuthenticateUsers(c *gin.Context) {
	email := c.PostForm("EMAIL")
	password := c.PostForm("PASSWORD")
	if email != "" && password != "" {
		if CheckForExistingUser(email, password) {
			token, err := GenerateAccessToken(email, password)
			if err == nil {
				c.JSON(http.StatusOK, token)
			} else {
				c.AbortWithStatusJSON(406, gin.H{"status": 406, "message": err})
			}
		} else {
			c.AbortWithStatusJSON(406, gin.H{"status": 406, "message": "user does not exist"})
		}
	} else {
		c.AbortWithStatusJSON(406, gin.H{"status": 406, "message": "fields are empty"})
	}
}

//give access only to authorised users which have an access token
func HomeAccess(c *gin.Context) {
	clientToken := c.GetHeader("Authorization")
	if clientToken == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Authorization Token is required"})
	}
	claims := jwt.MapClaims{}
	extractedToken := strings.Split(clientToken, "Bearer ")
	if len(extractedToken) == 2 {
		clientToken = strings.TrimSpace(extractedToken[1])
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Incorrect Format of Authorization Token "})
	}
	parsedToken, err := jwt.ParseWithClaims(clientToken, claims, func(token *jwt.Token) (interface{}, error) {
		return AtJwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Invalid Token Signature"})
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Bad Request"})
	}
	if !parsedToken.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Invalid Token"})
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Access Granted"})
}

func main() {
	//get all the environment variables
	db_host, host_exists := os.LookupEnv("DB_HOST")
	db_port, port_exists := os.LookupEnv("DB_PORT")
	db_user, user_exists := os.LookupEnv("DB_USER")
	db_name, name_exists := os.LookupEnv("DB_NAME")
	db_password, password_exists := os.LookupEnv("DB_PASSWORD")
	shared_key, shared_key_exists := os.LookupEnv("SHARED_KEY")
	var server_port_exists bool
	server_port, server_port_exists = os.LookupEnv("SERVER_PORT")
	if shared_key_exists {
		AtJwtKey = []byte(shared_key)
	}

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
