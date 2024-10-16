package controllers

import (
	"denis/first/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func (ctrl *UserController) GetUsers(req *gin.Context) {
	var users []models.User
	if err := ctrl.DB.Find(&users).Error; err != nil {
		req.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.JSON(http.StatusOK, users)
}

func (ctrl *UserController) CreateUser(req *gin.Context) {
	var user models.CreateUser
	if err := req.ShouldBindJSON(&user); err != nil {
		req.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := models.ValidateUser(&user.Data); err != nil {
		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors[err.Field()] = err.Tag()
		}
		req.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": validationErrors,
		})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Data.Password), bcrypt.DefaultCost)
	if err != nil {
		req.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Data.Password = string(hashedPassword)

	if err := ctrl.DB.Create(&user.Data).Error; err != nil {
		req.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}

func (ctrl *UserController) UpdateUser(req *gin.Context) {
	var user models.CreateUser
	if err := ctrl.DB.First(&user.Data, req.Param("id")).Error; err != nil {
		req.JSON(http.StatusNotFound, gin.H{"message": "User Not Found"})
		return
	}
	if err := req.ShouldBindJSON(&user); err != nil {
		req.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := models.ValidateUser(&user.Data); err != nil {
		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors[err.Field()] = err.Tag()
		}
		req.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": validationErrors,
		})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Data.Password), bcrypt.DefaultCost)
	if err != nil {
		req.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Data.Password = string(hashedPassword)
	if err := ctrl.DB.Save(&user.Data).Error; err != nil {
		req.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}

func (ctrl *UserController) DeleteUser(req *gin.Context) {
	var user models.User
	if err := ctrl.DB.First(&user, req.Param("id")).Error; err != nil {
		req.JSON(http.StatusNotFound, gin.H{"message": "User Not Found"})
		return
	}

	ctrl.DB.Delete(&user)
	req.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})

}
