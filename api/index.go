package handler

import (
	"go-supabase-api/config"
	"go-supabase-api/routes"
	"net/http"

	"github.com/gin-gonic/gin"
)

var app *gin.Engine

func init() {
	// Load environment variables (from .env or system env)
	config.LoadEnv()

	// Initialize Gin
	app = gin.New()
	app.Use(gin.Recovery())

	// CORS configuration
	app.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Setup Routes
	routes.SetupRoutes(app)
}

// Handler is the entry point for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
