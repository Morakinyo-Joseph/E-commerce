package utils

import (
	"go-ecommerce-api/models"
	"time"

	"net/http"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

// GenerateJWT generates a JWT token
func GenerateJWT(userID uint, email string, role string) (string, error) {
	claims := models.JwtClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "go-ecommerce-api", // Issuer for identification
		},
	}

	// Create a new token with the claims and signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	return token.SignedString([]byte("your_secret_key"))
}

// IsAdmin checks if the user is an admin based on their role in the context
func IsAdmin(c *gin.Context) bool {
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusForbidden, GenerateResponse("failed", "Role not found in context", nil, ""))
		c.Abort()
		return false
	}

	return role.(string) == "admin"
}
