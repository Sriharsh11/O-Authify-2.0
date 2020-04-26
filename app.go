package main

import (
	"fmt"
	"net/http"

	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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

//hash password
func HashPassword(password string) (string, error) {
	hashed_password_in_bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hashed_password_in_bytes), err
}

//check passowrd against the hash stored in the database
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func main() {

	//connect to database
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=authdb password=linuxissexy")
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()

	db.AutoMigrate(&Users3{})

	router := gin.Default()

	//add new user to table Users3
	router.POST("/addUser", func(c *gin.Context) {
		name := c.PostForm("NAME")
		email := c.PostForm("EMAIL")
		password := c.PostForm("PASSWORD")
		hashed_password, err := HashPassword(password) //store hashed password in database
		if err != nil {
			panic(err)
		}
		if name != "" && email != "" && hashed_password != "" {
			newUser := Users3{NAME: name, EMAIL: email, PASSWORD: hashed_password}
			db.NewRecord(newUser)
			db.Create(&newUser)
		} else {
			panic("failed to add new user")
		}
	})

	//return access tokens to verified users
	router.POST("/oauth", func(c *gin.Context) {
		email := c.PostForm("EMAIL")
		password := c.PostForm("PASSWORD")
		var password_db_list login
		var email_db_list login
		if email != "" && password != "" {
			fmt.Println(db.Table("users3").Select("EMAIL").Where("EMAIL = ?", email).Scan(&email_db_list))
			db.Table("users3").Select("EMAIL").Where("EMAIL = ?", email).Scan(&email_db_list)
			fmt.Println(email_db_list)
			email_db := email_db_list.EMAIL
			if email_db != "" {
				db.Table("users3").Select("PASSWORD").Where("EMAIL = ?", email).Scan(&password_db_list)
				fmt.Println(password_db_list)
				password_db := password_db_list.PASSWORD
				if CheckPasswordHash(password, password_db) {
					payload := `{"security":"OAuth 2.0"}`
					key := []byte{97, 48, 97, 50, 97, 98, 100, 56, 45, 54, 49, 54, 50, 45, 52, 49, 99, 51, 45, 56, 51, 100, 54, 45, 49, 99, 102, 53, 53, 57, 98, 52, 54, 97, 102, 99}
					token, err := jose.Sign(payload, jose.HS256, key) //using HS256 algorithm for creating JWT
					if err == nil {
						c.JSON(http.StatusOK, token)
					} else {
						panic("failed to generate token")
					}
				} else {
					panic("invalid credentials")
				}
			} else {
				panic("user does not exist")
			}
		} else {
			panic("fields are empty")
		}
	})

	router.GET("/home", func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		sharedKey := []byte{97, 48, 97, 50, 97, 98, 100, 56, 45, 54, 49, 54, 50, 45, 52, 49, 99, 51, 45, 56, 51, 100, 54, 45, 49, 99, 102, 53, 53, 57, 98, 52, 54, 97, 102, 99}
		payload, _, err := jose.Decode(token, sharedKey)
		if err == nil {
			c.JSON(http.StatusOK, "Access Granted for "+payload) //Granting access only to authorised users
		} else {
			c.AbortWithStatus(401)
		}
	})

	router.Run(":3000")
}
