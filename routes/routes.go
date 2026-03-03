package routes

import (
	"go-supabase-api/controllers"

	"github.com/gin-gonic/gin"
)

// "go-supabase-api/middleware"

func SetupRoutes(r *gin.Engine) {
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}
	// r.GET("/products", controllers.GetProducts)

	// // User routes
	// r.POST("/transactions", middleware.AuthMiddleware(), controllers.CreateTransaction)

	// // Admin-only routes - hanya admin yang bisa create/update/delete
	// product := r.Group("/products")
	// product.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
	// {
	// 	// product.GET("", controllers.GetProducts)
	// 	product.GET("/:id", controllers.GetProduct)
	// 	product.POST("", controllers.CreateProduct)
	// 	product.PUT("/:id", controllers.UpdateProduct)
	// 	product.DELETE("/:id", controllers.DeleteProduct)
	// }
}
