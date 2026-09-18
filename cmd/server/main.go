package main

import (
	"log"
	"net/http"
	"team-access-control/internal/config"
	"team-access-control/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main(){


	if err := godotenv.Load() ; err !=nil{
			log.Println("No .env file found")
	}
	cfg := config.Load()

	db , err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
		defer db.Close()
	router := gin.Default()
	router.GET("/health" , func(c *gin.Context){
		c.JSON(http.StatusOK , gin.H{
			"status":"ok",
		})
	})
	router.Run(":"+ cfg.Port)
}