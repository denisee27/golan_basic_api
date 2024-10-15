package main

import (
	"log";
	"github.com/gin-gonic/gin"
	"denis/first/routes"
	"denis/first/config"
);

func main(){
	database.ConnectDB()
	router := gin.Default()
	routes.SetupRoutes(router, database.DB)
	log.Println("Server running at http://localhost:8000")
	router.Run(":8000")

}	