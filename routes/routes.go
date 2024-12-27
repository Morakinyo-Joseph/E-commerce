package routes

import (
	"go-ecommerce-api/controllers"
	"go-ecommerce-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	// Public routes (No authentication needed)
	router.POST("/register", controllers.RegisterUser)
	router.POST("/login", controllers.LoginUser)

	// Protected routes (Requires JWT)
	protected := router.Group("/api")
	protected.Use(middleware.JWTMiddleware())

	// Product routes (Admin-only access)
	protected.POST("/products", controllers.CreateProduct)
	protected.PUT("/products/:id", controllers.UpdateProduct)
	protected.DELETE("/products/:id", controllers.DeleteProduct)
	protected.GET("/products", controllers.GetProducts)

	// Order routes
	protected.POST("/orders", controllers.PlaceOrder)
	protected.GET("/orders", controllers.ListUserOrders)
	protected.PUT("/orders/:id/cancel", controllers.CancelOrder)
	protected.PUT("/orders/:id", controllers.UpdateOrderStatus)

	return router
}
