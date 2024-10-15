package routes

import (
	"denis/first/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	userController := &controllers.UserController{DB: db}
	router.GET("/users", userController.GetUsers)
	router.POST("/users/create", userController.CreateUser)
	router.PUT("/users/update/:id", userController.UpdateUser)
	router.DELETE("/users/delete/:id", userController.DeleteUser)
}
