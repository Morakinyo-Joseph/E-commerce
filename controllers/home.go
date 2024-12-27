package controllers

import (
	"go-ecommerce-api/utils"

	"github.com/gin-gonic/gin"

	"net/http"
)

func HomeView(c *gin.Context) {
	response := utils.GenerateResponse("success", "Ecommerce API. Visit the documentation for insight.", nil, "")
	c.JSON(http.StatusOK, response)
}
