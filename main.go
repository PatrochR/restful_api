package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
}

func main() {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User{})
	router := gin.Default()

	// register
	router.POST("/auth/register", func(c *gin.Context) {
		var body struct {
			Email    string
			Password string
		}
		if c.Bind(&body) != nil {
			log.Println("error :register binding a body")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "something wrong!",
			})
			return
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
		if err != nil {
			log.Println("error :user hash password")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "something wrong!",
			})
			return
		}

		user := User{Email: body.Email, Password: string(passwordHash)}
		result := db.Create(&user)
		if result.Error != nil {
			log.Println("error :user creat")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "something wrong!",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	//login
	router.POST("/auth/login", func(c *gin.Context) {
		var body struct {
			Email    string
			Password string
		}
		if c.Bind(&body) != nil {
			log.Println("error :login binding a body")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "something wrong!",
			})
			return
		}
		var user User
		db.First(&user, "Email = ?", body.Email)

		if user.ID == 0 {
			log.Println("error :login email is not find")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid email or password!",
			})
			return
		}

		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
		if err != nil {
			log.Println("error :login password is wrong")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid email or password!",
			})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.ID,
			"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
		if err != nil {
			log.Println("massage : " + err.Error())
			log.Println("error :login generate token")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "something wrong!",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"token": tokenString,
		})
	})
	router.Run()
}
