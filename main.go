package main

import (
	database "basic/config"
	"basic/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDB()
	router := gin.Default()
	routes.SetupRoutes(router, database.DB)
	log.Println("Server running at http://localhost:8000")
	router.Run(":8000")

}
