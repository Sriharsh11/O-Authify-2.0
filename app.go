package main

import (
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

//hash password
func HashPassword(password string) (string, error) {
	hashed_password_in_bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hashed_password_in_bytes), err
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
	router.Run(":3000")

	// router := gin.Default()
	// router.POST("/test", func(c *gin.Context) {
	// 	email := c.PostForm("email")
	// 	password := c.PostForm("password")
	// 	c.JSON(200, gin.H{
	// 		"message": email + " " + password,
	// 	})
	// })
	// router.Run(":3000")
}
