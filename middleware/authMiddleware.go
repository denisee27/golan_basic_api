package middleware

import (
	"basic/jwtKey"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TokenValidation() gin.HandlerFunc {
	return func(req *gin.Context) {
		if req.Request.URL.Path == "/auth/login" {
			req.Next()
			return
		}
		authHeader := req.GetHeader("Authorization")
		if authHeader == "" {
			req.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Authorization header missing"})
			req.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &jwtKey.PayloadJwt{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			req.JSON(http.StatusUnauthorized, gin.H{
				"status": os.Getenv("JWT_SECRET"),
				// "status": http.StatusUnauthorized,
				"error": "Invalid token"})
			req.Abort()
			return
		}
		if claims, ok := token.Claims.(*jwtKey.PayloadJwt); ok && token.Valid {
			// Periksa apakah token sudah kadaluarsa
			if claims.ExpiresAt.Time.Before(time.Now()) {
				req.JSON(http.StatusUnauthorized, gin.H{"error": "Token is expired"})
				req.Abort()
				return
			}

			// Token valid, bisa melanjutkan request
			// Anda dapat menggunakan klaim seperti email di dalam request
			req.Set("email", claims.Email)

			// Lanjutkan ke handler berikutnya
			req.Next()
		} else {
			// Jika klaim tidak valid
			req.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			req.Abort()
		}
	}
}
