package controllers

import (
	"net/http"

	"go-ecommerce-api/database"
	"go-ecommerce-api/models"
	"go-ecommerce-api/utils"

	"github.com/gin-gonic/gin"
)

// CreateProduct handles creating a new product (admin only)
// @Summary Create a new product
// @Description Admins can create a new product by providing necessary details
// @Tags Products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product data"
// @Success 201 {object} utils.Response "Product created successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 500 {object} utils.Response "Failed to create product"
// @Router /products [post]
func CreateProduct(c *gin.Context) {
	var input models.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.GenerateResponse("failed", "Invalid request data", nil, err.Error()))
		return
	}

	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.GenerateResponse("failed", "Failed to create product", nil, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.GenerateResponse("success", "Product created successfully", input, ""))
}

// GetProducts handles retrieving all products
// @Summary Retrieve all products
// @Description Get a list of all products
// @Tags Products
// @Produce json
// @Success 200 {object} utils.Response "Products retrieved successfully"
// @Failure 500 {object} utils.Response "Failed to fetch products"
// @Router /products [get]
func GetProducts(c *gin.Context) {
	var products []models.Product
	if err := database.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.GenerateResponse("failed", "Failed to fetch products", nil, err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.GenerateResponse("success", "Products retrieved successfully", products, ""))
}

// UpdateProduct handles updating a product (admin only)
// @Summary Update an existing product
// @Description Admins can update product details by providing the product ID and new data
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body models.Product true "Updated product data"
// @Success 200 {object} utils.Response "Product updated successfully"
// @Failure 400 {object} utils.Response "Invalid request data"
// @Failure 403 {object} utils.Response "Unauthorized action"
// @Failure 404 {object} utils.Response "Product not found"
// @Failure 500 {object} utils.Response "Failed to update product"
// @Router /products/{id} [put]
func UpdateProduct(c *gin.Context) {
	if !utils.IsAdmin(c) {
		c.JSON(http.StatusForbidden, utils.GenerateResponse("failed", "You are not authorized to perform this action", nil, ""))
		return
	}

	productID := c.Param("id")
	var input models.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.GenerateResponse("failed", "Invalid request data", nil, err.Error()))
		return
	}

	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.GenerateResponse("failed", "Product not found", nil, err.Error()))
		return
	}

	product.Name = input.Name
	product.Description = input.Description
	product.Price = input.Price
	product.Stock = input.Stock
	if err := database.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.GenerateResponse("failed", "Failed to update product", nil, err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.GenerateResponse("success", "Product updated successfully", product, ""))
}

// DeleteProduct handles deleting a product (admin only)
// @Summary Delete a product
// @Description Admins can delete a product by providing the product ID
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} utils.Response "Product deleted successfully"
// @Failure 403 {object} utils.Response "Unauthorized action"
// @Failure 404 {object} utils.Response "Product not found"
// @Failure 500 {object} utils.Response "Failed to delete product"
// @Router /products/{id} [delete]
func DeleteProduct(c *gin.Context) {
	if !utils.IsAdmin(c) {
		c.JSON(http.StatusForbidden, utils.GenerateResponse("failed", "You are not authorized to perform this action", nil, ""))
		return
	}

	productID := c.Param("id")
	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.GenerateResponse("failed", "Product not found", nil, err.Error()))
		return
	}

	if err := database.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.GenerateResponse("failed", "Failed to delete product", nil, err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.GenerateResponse("success", "Product deleted successfully", nil, ""))
}
