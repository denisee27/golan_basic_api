package controllers

import (
	"denis/first/jwtKey"
	"denis/first/models"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

type LoginRequest struct {
	Email     string `json:"email"`
	Passsword string `json:"password"`
}

type payloadJwt struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func (ctrl *AuthController) Login(req *gin.Context) {
	var loginReq LoginRequest

	//Validasi & Object Parsing
	if err := req.ShouldBindJSON(&loginReq); err != nil {
		req.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload"})
		return
	}

	//Validasi Email
	var user models.User
	if err := ctrl.DB.Where("email = ?", loginReq.Email).First(&user).Error; err != nil {
		req.JSON(http.StatusUnauthorized, gin.H{"message": "User Not Found"})
		return
	}

	//Validasi Password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Passsword))
	if err != nil {
		req.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		return
	}

	//JWT Secret Key
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		generateKey, err := jwtKey.JwtSecretKey(req)
		if err != nil {
			req.JSON(http.StatusUnauthorized, gin.H{"message": "Could not generate or write to .env file"})
			return
		}
		jwtSecret = generateKey
		fmt.Println("Generated and saved new JWT Secret Key:", jwtSecret)
	}

	// Create Token
	expirationTime := time.Now().Add(24 * time.Hour)
	payloadJwt := &payloadJwt{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payloadJwt) // Ubah ke SigningMethodHS256 jika tidak menggunakan ES256
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		req.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token"})
		return
	}
	req.JSON(http.StatusOK, gin.H{
		"email":       user.Email,
		"name":        user.Name,
		"acces_token": tokenString,
	})

}
