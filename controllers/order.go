package controllers

import (
	"net/http"
	"time"

	"go-ecommerce-api/database"
	"go-ecommerce-api/models"
	"go-ecommerce-api/utils"

	"github.com/gin-gonic/gin"
)

// PlaceOrder handles the creation of a new order

// @Summary Place Order
// @Description Create a new order with the provided order items and calculate the total amount based on the product prices.
// @Tags Orders
// @Accept json
// @Produce json
// @Param order_items body []models.OrderItem true "OrderItem data"
// @Success 201 {object} utils.Response{data=models.Order} "Order placed successfully"
// @Failure 400 {object} utils.Response "Invalid input"
// @Failure 500 {object} utils.Response "Failed to place order"
// @Security ApiKeyAuth
// @Router /orders [post]
func PlaceOrder(c *gin.Context) {
	var orderInput struct {
		OrderItems []struct {
			ProductID uint `json:"product_id"`
			Quantity  int  `json:"quantity"`
		} `json:"order_items"`
	}

	if err := c.ShouldBindJSON(&orderInput); err != nil {
		c.JSON(http.StatusBadRequest, utils.GenerateResponse("failed", "Invalid input", nil, err.Error()))
		return
	}

	userID := c.MustGet("userID").(uint)
	var totalAmount float64
	orderItems := []models.OrderItem{}

	for _, item := range orderInput.OrderItems {
		var product models.Product
		if err := database.DB.First(&product, item.ProductID).Error; err != nil {
			c.JSON(http.StatusNotFound, utils.GenerateResponse("failed", "Product not found", nil, err.Error()))
			return
		}

		orderItem := models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}
		totalAmount += float64(item.Quantity) * product.Price
		orderItems = append(orderItems, orderItem)
	}

	order := models.Order{
		UserID:      userID,
		Status:      "Pending",
		TotalAmount: totalAmount,
		CreatedAt:   time.Now(),
		OrderItems:  orderItems,
	}

	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.GenerateResponse("failed", "Failed to place order", nil, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.GenerateResponse("success", "Order placed successfully", order, ""))
}

// ListUserOrders retrieves all orders for the authenticated user

// @Summary List User Orders
// @Description Fetch all orders for the currently authenticated user with order details.
// @Tags Orders
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.Order} "Orders retrieved successfully"
// @Failure 500 {object} utils.Response "Failed to retrieve orders"
// @Security ApiKeyAuth
// @Router /orders [get]
func ListUserOrders(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var orders []models.Order

	if err := database.DB.Preload("OrderItems").Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.GenerateResponse("failed", "Failed to retrieve orders", nil, err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.GenerateResponse("success", "Orders retrieved successfully", orders, ""))
}

// CancelOrder allows a user to cancel an order in the Pending status

// @Summary Cancel Order
// @Description Cancel an order if it is in the pending status. Only the owner of the order can cancel it.
// @Tags Orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} utils.Response{data=models.Order} "Order cancelled successfully"
// @Failure 400 {object} utils.Response "Only pending orders can be cancelled"
// @Failure 404 {object} utils.Response "Order not found"
// @Failure 403 {object} utils.Response "Unauthorized action"
// @Failure 500 {object} utils.Response "Failed to cancel order"
// @Security ApiKeyAuth
// @Router /orders/{id} [delete]
func CancelOrder(c *gin.Context) {
	orderID := c.Param("id")
	var order models.Order

	if err := database.DB.First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.GenerateResponse("failed", "Order not found", nil, err.Error()))
		return
	}

	if order.UserID != c.MustGet("userID").(uint) {
		c.JSON(http.StatusForbidden, utils.GenerateResponse("failed", "You are not authorized to cancel this order", nil, ""))
		return
	}

	if order.Status != "Pending" {
		c.JSON(http.StatusBadRequest, utils.GenerateResponse("failed", "Only pending orders can be cancelled", nil, ""))
		return
	}

	order.Status = "Cancelled"
	if err := database.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.GenerateResponse("failed", "Failed to cancel order", nil, err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.GenerateResponse("success", "Order cancelled successfully", order, ""))
}

// UpdateOrderStatus allows an admin to update the status of an order

// @Summary Update Order Status
// @Description Admin can update the status of an order (e.g., to shipped, delivered, etc.).
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param status body string true "New status of the order" Enums(pending, shipped, delivered, cancelled)
// @Success 200 {object} utils.Response{data=models.Order} "Order status updated successfully"
// @Failure 400 {object} utils.Response "Invalid input"
// @Failure 404 {object} utils.Response "Order not found"
// @Failure 403 {object} utils.Response "Unauthorized action"
// @Failure 500 {object} utils.Response "Failed to update order status"
// @Security ApiKeyAuth
// @Router /orders/{id}/status [put]
func UpdateOrderStatus(c *gin.Context) {
	if !utils.IsAdmin(c) {
		c.JSON(http.StatusForbidden, utils.GenerateResponse("failed", "You are not authorized to perform this action", nil, ""))
		return
	}

	orderID := c.Param("id")
	var input struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.GenerateResponse("failed", "Invalid input", nil, err.Error()))
		return
	}

	var order models.Order
	if err := database.DB.First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.GenerateResponse("failed", "Order not found", nil, err.Error()))
		return
	}

	order.Status = input.Status
	if err := database.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.GenerateResponse("failed", "Failed to update order status", nil, err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.GenerateResponse("success", "Order status updated successfully", order, ""))
}
