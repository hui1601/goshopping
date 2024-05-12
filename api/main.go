package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"main/routers"
)

func registerRoutes(router *gin.Engine) {
	routers.RegisterRoutes(router)
}
func main() {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8000", "http://localhost:3000"}
	config.AllowCredentials = true
	router.Use(cors.New(config))
	registerRoutes(router)
	router.Run(":8080")
}
