package main

import (
	database "denis/first/config"
	"denis/first/routes"
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
