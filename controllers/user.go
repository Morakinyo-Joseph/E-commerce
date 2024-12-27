package controllers

import (
	"go-ecommerce-api/models"
	"go-ecommerce-api/utils"
	"net/http"

	"go-ecommerce-api/database"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser handles user registration
// @Summary Register User
// @Description Register a new user by creating a hashed password and saving the user information in the database.
// @Tags Users
// @Accept json
// @Produce json
// @Param input body struct { Email string `json:"email" binding:"required,email"`; Password string `json:"password" binding:"required,min=6"`; Role string `json:"role"` } true "User registration data"
// @Success 201 {object} utils.Response{data=gin.H} "User created successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 500 {object} utils.Response "Failed to create user"
// @Router /register [post]
func RegisterUser(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.GenerateResponse("failed", "Invalid request data", nil, err.Error()))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.GenerateResponse("failed", "Failed to hash password", nil, err.Error()))
		return
	}

	user := models.User{
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		Role:         input.Role,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.GenerateResponse("failed", "Failed to create user", nil, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.GenerateResponse("success", "User created successfully", gin.H{
		"id":        user.ID,
		"email":     user.Email,
		"role":      user.Role,
		"createdAt": user.CreatedAt,
	}, ""))
}

// LoginUser handles user login and generates a JWT token

// @Summary Login User
// @Description Authenticate a user and generate a JWT token upon successful login.
// @Tags Users
// @Accept json
// @Produce json
// @Param input body models.UserLoginInput true "User login credentials"
// @Success 200 {object} utils.Response{data=gin.H} "Login successful"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 401 {object} utils.Response "Invalid email or password"
// @Failure 500 {object} utils.Response "Failed to generate token"
// @Router /login [post]
func LoginUser(c *gin.Context) {
	var input models.UserLoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.GenerateResponse("failed", "Invalid request data", nil, err.Error()))
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, utils.GenerateResponse("failed", "Invalid email or password", nil, ""))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, utils.GenerateResponse("failed", "Invalid email or password", nil, ""))
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.GenerateResponse("failed", "Failed to generate token", nil, err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.GenerateResponse("success", "Login successful", gin.H{"token": token}, ""))
}
