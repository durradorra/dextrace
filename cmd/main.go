package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type SibionicUser struct {
	Email string `json:"email"`
	Password string `json:"password"`
}


func main() {
	
	var user SibionicUser
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	user.Email = os.Getenv("USER_EMAIL")
	user.Password = os.Getenv("USER_PASSWORD")

	fmt.Println(user)

	r := gin.Default()

	
	r.GET("/data", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	
	r.Run(":8080")
}