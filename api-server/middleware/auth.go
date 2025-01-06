package middleware

import (
	"api-server/models"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		jwtSecret := os.Getenv("JWT_SECRET")

		token, err := jwt.ParseWithClaims(tokenString, &models.SupabaseJWTPayload{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		claims := token.Claims.(*models.SupabaseJWTPayload)
		if claims.Exp < time.Now().Unix() {
			c.AbortWithStatusJSON(401, gin.H{"error": "Token expired"})
			return
		}

		c.Set("user", claims)
		c.Next()
	}
}
