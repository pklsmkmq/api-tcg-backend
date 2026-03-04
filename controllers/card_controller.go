package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go-supabase-api/config"
	"go-supabase-api/models"

	"github.com/gin-gonic/gin"
)

// GetMyCards mengambil semua kartu milik user yang sedang login
func GetMyCards(c *gin.Context) {
	// 1. Ambil user_id dari Token
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

	// 2. Tarik data dari Supabase tabel card_collection khusus user_id ini
	// Kita urutkan dari yang paling baru didapatkan (acquired_at.desc)
	endpoint := fmt.Sprintf("card_collection?user_id=eq.%s&order=acquired_at.desc", idStr)
	resp, err := config.SupabaseRequest("GET", endpoint, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menghubungi database"})
		return
	}
	defer resp.Body.Close()

	var cards []models.CardCollection
	if err := json.NewDecoder(resp.Body).Decode(&cards); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memproses data kartu"})
		return
	}

	// 3. Jika user belum punya kartu, pastikan mengembalikan array kosong [], bukan null
	if cards == nil {
		cards = []models.CardCollection{}
	}

	// 4. Kirim respons ke Frontend
	c.JSON(http.StatusOK, gin.H{
		"message":     "Berhasil mengambil koleksi kartu",
		"total_kartu": len(cards),
		"data":        cards,
	})
}