package jwtKey

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type PayloadJwt struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func JwtSecretKey(req *gin.Context) (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	secretKey := base64.URLEncoding.EncodeToString(bytes)
	fmt.Println("Generated Secret Key:", secretKey)
	if err := writeSecretKey(secretKey); err != nil {
		req.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write secret key to .env"})
		return "", err
	}
	return secretKey, nil
}

func writeSecretKey(secret string) error {
	envFile, err := os.OpenFile(".env", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer envFile.Close()
	_, err = io.WriteString(envFile, fmt.Sprintf("JWT_SECRET=%s", secret))
	if err != nil {
		return err
	}
	return nil
}
