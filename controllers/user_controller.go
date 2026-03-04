package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go-supabase-api/config"
	"go-supabase-api/models"

	"github.com/gin-gonic/gin"
)

func TopupSaldo(c *gin.Context) {
	// 1. Tangkap input jumlah topup dari body request
	var req struct {
		Amount float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Format JSON salah"})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Jumlah topup harus lebih dari 0"})
		return
	}

	// 2. Ambil user_id dari JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User ID tidak ditemukan di token"})
		return
	}

	// 3. Konversi dan Validasi Tipe Data ID
	var idStr string
	switch v := userID.(type) {
	case string:
		// Jika tipenya string, berarti ini adalah Token Baru yang benar (UUID)
		idStr = v
	case float64:
		// Jika tipenya float64 (angka), berarti ini adalah Token Lama yang merekam angka 0
		c.JSON(http.StatusUnauthorized, gin.H{
			"message":          "Token lama terdeteksi. Silakan Login ulang di Postman untuk mendapatkan Token baru.",
			"tipe_terdeteksi":  "float64/angka",
			"nilai_terdeteksi": v,
		})
		return
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Format User ID tidak dikenali"})
		return
	}

	// 4. Tarik data user saat ini dari DB (menggunakan %s karena ID sekarang string)
	resp, err := config.SupabaseRequest("GET", fmt.Sprintf("users?id=eq.%s&select=saldo_uang", idStr), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menghubungi database"})
		return
	}
	defer resp.Body.Close()

	var users []models.User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil || len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User tidak ditemukan di database"})
		return
	}

	currentSaldo := users[0].SaldoUang
	newSaldo := currentSaldo + req.Amount

	// 5. Update saldo baru ke Supabase (menggunakan %s)
	updateResp, err := config.SupabaseRequest("PATCH", fmt.Sprintf("users?id=eq.%s", idStr), gin.H{
		"saldo_uang": newSaldo,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengupdate saldo"})
		return
	}
	defer updateResp.Body.Close()

	if updateResp.StatusCode >= 400 {
		var errResp map[string]interface{}
		json.NewDecoder(updateResp.Body).Decode(&errResp)
		c.JSON(updateResp.StatusCode, gin.H{"message": "Gagal menyimpan saldo", "error": errResp})
		return
	}

	// 6. Kembalikan respons sukses
	c.JSON(http.StatusOK, gin.H{
		"message":          "Topup berhasil",
		"topup_amount":     req.Amount,
		"saldo_sebelumnya": currentSaldo,
		"saldo_baru":       newSaldo,
	})
}
