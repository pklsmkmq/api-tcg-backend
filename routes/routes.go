package routes

import (
	"net/http" // <-- Tambahkan import http

	"go-supabase-api/controllers"
	"go-supabase-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// 1. ENDPOINT ROOT (Dokumentasi Mini)
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"app":     "Pokémon TCG API",
			"status":  "running",
			"version": "1.0.0",
			"endpoints": gin.H{
				"Authentication": gin.H{
					"POST /api/auth/register": "Mendaftarkan user baru (Public)",
					"POST /api/auth/login":    "Login dan mendapatkan token JWT (Public)",
				},
				"Users": gin.H{
					"POST /api/users/topup": "Topup saldo user (Butuh Token)", // <-- Dokumentasi baru
				},
				"Pokémon Sets": gin.H{
					"GET /api/pokemon/sets":       "Melihat daftar pack/set Pokémon (Butuh Token)",
					"POST /api/pokemon/sets/sync": "Sinkronisasi data dari TCGdex ke DB (Hanya Admin)",
				},
				// Catatan: Anda bisa menambahkan Products dan Transactions ke sini nanti
				// saat route-nya sudah Anda buka (uncomment).
			},
		})
	})

	// 2. ENDPOINT AUTHENTICATION
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	// 3. ENDPOINT POKEMON SETS
	pokemon := r.Group("/api/pokemon")
	pokemon.Use(middleware.AuthMiddleware())
	{
		// User biasa & Admin
		pokemon.GET("/sets", controllers.GetPokemonSets)

		// Hanya Admin
		pokemon.POST("/sets/sync", middleware.AdminOnly(), controllers.SyncSets)

		pokemon.GET("/my-cards", controllers.GetMyCards)
	}

	// 4. ENDPOINT USERS (Baru)
	users := r.Group("/api/users")
	users.Use(middleware.AuthMiddleware()) // Wajib pakai token
	{
		users.POST("/topup", controllers.TopupSaldo)
		users.GET("/transactions", controllers.GetTransactions)
	}

	store := r.Group("/api/store")
	store.Use(middleware.AuthMiddleware()) // Wajib pakai token JWT
	{
		// Endpoint untuk membeli pack
		store.POST("/buy-pack", controllers.BuyPack)
	}
}
