package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go-supabase-api/config"
	"go-supabase-api/models"

	"github.com/gin-gonic/gin"
)

// GetTransactions mengambil riwayat transaksi digital milik user yang sedang login
func GetTransactions(c *gin.Context) {
	// Ambil user_id (string/UUID) dari Token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	idStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Format User ID di token tidak valid"})
		return
	}

	// Tarik data transaksi dari Supabase khusus untuk user ini, diurutkan dari yang terbaru
	endpoint := fmt.Sprintf("transactions?user_id=eq.%s&order=created_at.desc", idStr)
	resp, err := config.SupabaseRequest("GET", endpoint, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menghubungi database"})
		return
	}
	defer resp.Body.Close()

	var transactions []models.Transaction
	if err := json.NewDecoder(resp.Body).Decode(&transactions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memproses data transaksi"})
		return
	}

	// Jika kosong, kembalikan array kosong agar Frontend tidak error (bukan null)
	if transactions == nil {
		transactions = []models.Transaction{}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil riwayat transaksi",
		"data":    transactions,
	})
}
