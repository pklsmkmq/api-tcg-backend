package controllers

import (
	"encoding/json"
	"net/http"

	"go-supabase-api/config"
	"go-supabase-api/models"

	"time"

	"github.com/gin-gonic/gin"
)

func CreateTransaction(c *gin.Context) {
	var input struct {
		models.Transaction
		Details []models.TransactionDetail `json:"details"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Format JSON salah"})
		return
	}

	// Ambil UserID dari context (set via AuthMiddleware).
	// Jika token ada, kita prioritaskan ID dari token.
	userID, exists := c.Get("user_id")
	if exists {
		// JWT claims unmarshal numbers as float64
		if floatID, ok := userID.(float64); ok {
			input.Transaction.UserID = int(floatID)
		} else {
			// Fallback: Jika ternyata tersimpan sebagai int (misal dari middleware custom lain)
			if intID, ok := userID.(int); ok {
				input.Transaction.UserID = intID
			}
		}
	}

	// Validasi sederhana
	if input.Transaction.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User ID diperlukan"})
		return
	}
	if len(input.Details) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Detail transaksi diperlukan"})
		return
	}

	// 1. Simpan Transaksi Header
	// Gunakan struct terpisah atau map agar field Details tidak terkirim dua kali (sekali di header, sekali di detail)
	// Namun models.Transaction punya field Details `json:"details,omitempty"`
	// Kita perlu set Details di struct transaction ke nil atau kosong sebelum kirim header

	headerTx := input.Transaction

	// Set waktu ke Asia/Jakarta (WIB)
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Fallback jika timezone db tidak ada
		loc = time.FixedZone("Asia/Jakarta", 7*60*60)
	}
	headerTx.CreatedAt = time.Now().In(loc).Format(time.RFC3339)

	headerTx.Details = nil // Pastikan tidak dikirim ke endpoint transactions

	resp, err := config.SupabaseRequest("POST", "transactions", headerTx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menghubungi database"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		c.JSON(resp.StatusCode, gin.H{"message": "Gagal menyimpan transaksi", "error": errResp})
		return
	}

	// Parse response untuk dapat ID
	var createdTransactions []models.Transaction
	if err := json.NewDecoder(resp.Body).Decode(&createdTransactions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memproses respons transaksi"})
		return
	}
	if len(createdTransactions) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Transaksi berhasil disimpan tetapi tidak ada data kembali"})
		return
	}
	newTransactionID := createdTransactions[0].ID

	// 2. Simpan Transaksi Detail
	for i := range input.Details {
		input.Details[i].TransactionID = newTransactionID
		// Pastikan ID detail kosong agar digenerate DB, atau set jika perlu
		input.Details[i].ID = ""
	}

	respDetail, err := config.SupabaseRequest("POST", "transaction_details", input.Details)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menghubungi database untuk detail"})
		return
	}
	defer respDetail.Body.Close()

	if respDetail.StatusCode >= 400 {
		var errResp map[string]interface{}
		json.NewDecoder(respDetail.Body).Decode(&errResp)
		// Warning: Header sudah tersimpan, detail gagal. Inconsistent state.
		c.JSON(respDetail.StatusCode, gin.H{
			"message":        "Gagal menyimpan detail transaksi. Transaksi header tersimpan.",
			"transaction_id": newTransactionID,
			"error":          errResp,
		})
		return
	}

	// Kembalikan hasil lengkap
	result := createdTransactions[0]
	result.Details = input.Details
	c.JSON(http.StatusCreated, result)
}
