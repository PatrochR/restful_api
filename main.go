package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type user struct {
	Email    string `gorm:"unique"`
	Password string
}

func main() {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&user{})
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		var body struct {
			Email    string
			Password string
		}
		if c.Bind(&body) != nil {
			log.Fatal("error : binding a body!")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "something wrong!",
			})
			return
		}
	})
	router.Run()
}
