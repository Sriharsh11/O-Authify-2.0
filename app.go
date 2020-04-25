package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Users struct {
	gorm.Model
	id       string
	name     string
	email    string
	password string
}

func main() {

	db, err := gorm.Open("postgres", "host=localhost port=5432 user=sriharsh dbname=authDB password=linuxissexy")
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()

	db.AutoMigrate(&Users{})
	//db.CreateTable(&Users{})

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
