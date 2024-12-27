package middleware

import (
	"fmt"
	"go-ecommerce-api/models"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWTMiddleware checks if the user is authenticated and authorized
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println(c)
		// Get the token from the Authorization header
		tokenString := c.GetHeader("Authorization")
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		// Parse the token
		token, err := jwt.ParseWithClaims(tokenString, &models.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("your_secret_key"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract claims from the token
		claims, ok := token.Claims.(*models.JwtClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Attach the user info to the context
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		// Check if user is admin
		fmt.Println("++++++++++++++++++++", c)
		if !IsAdmin(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to perform this action"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// IsAdmin checks if the user is an admin based on their role
func IsAdmin(c *gin.Context) bool {
	role := c.MustGet("role").(string)
	return role == "admin"
}
