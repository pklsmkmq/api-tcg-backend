package controllers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"go-supabase-api/config"
	"go-supabase-api/models"

	"github.com/gin-gonic/gin"
)

func BuyPack(c *gin.Context) {
	// --- Langkah 1: Terima Request ---
	var req struct {
		SetID string `json:"set_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.SetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Format JSON salah, pastikan mengirim set_id"})
		return
	}

	// Ambil user_id dari JWT token
	userIDContext, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	idStr, ok := userIDContext.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Format User ID di token tidak valid (harus string/UUID)"})
		return
	}

	// --- Langkah 2: Validasi Saldo & Harga Pack ---
	// A. Cek harga pack dari DB Supabase
	respSet, err := config.SupabaseRequest("GET", fmt.Sprintf("pokemon_sets?id=eq.%s&select=price,name", req.SetID), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil data pack dari database"})
		return
	}
	defer respSet.Body.Close()

	var sets []models.PokemonSet
	json.NewDecoder(respSet.Body).Decode(&sets)
	if len(sets) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Set/Pack tidak ditemukan di database"})
		return
	}
	packPrice := float64(sets[0].Price)
	packName := sets[0].Name

	// B. Cek saldo user dari DB Supabase
	respUser, err := config.SupabaseRequest("GET", fmt.Sprintf("users?id=eq.%s&select=saldo_uang", idStr), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menghubungi database user"})
		return
	}
	defer respUser.Body.Close()

	var users []models.User
	json.NewDecoder(respUser.Body).Decode(&users)
	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User tidak ditemukan"})
		return
	}

	currentSaldo := users[0].SaldoUang
	if currentSaldo < packPrice {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":    "Saldo tidak mencukupi untuk membeli pack ini",
			"harga_pack": packPrice,
			"saldo_anda": currentSaldo,
		})
		return
	}

	// --- Langkah 3: Fetch API TCGdex ---
	apiURL := fmt.Sprintf("https://api.tcgdex.net/v2/en/sets/%s", req.SetID)
	respTCG, err := http.Get(apiURL)
	if err != nil || respTCG.StatusCode != 200 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil data kartu dari server TCGdex"})
		return
	}
	defer respTCG.Body.Close()

	var setDetail models.TCGdexSetDetail
	if err := json.NewDecoder(respTCG.Body).Decode(&setDetail); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal membaca respons dari TCGdex"})
		return
	}
	if len(setDetail.Cards) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Pack ini kosong (tidak memiliki kartu)"})
		return
	}

	// --- Langkah 4: Pilih 5 Kartu Secara Random ---
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var obtainedCards []models.CardCollection

	for i := 0; i < 5; i++ {
		randomIndex := r.Intn(len(setDetail.Cards))
		selectedCard := setDetail.Cards[randomIndex]

		imgUrl := selectedCard.Image
		if imgUrl != "" {
			imgUrl = fmt.Sprintf("%s/low.png", imgUrl) // Menggunakan gambar resolusi tinggi
		}

		obtainedCards = append(obtainedCards, models.CardCollection{
			UserID:   idStr,
			CardID:   selectedCard.ID,
			CardName: selectedCard.Name,
			ImageURL: imgUrl,
			SetID:    req.SetID,
		})
	}

	// --- Langkah 5 & 6: Hitung & Update Saldo Baru (DB 1) ---
	newSaldo := currentSaldo - packPrice
	updateSaldoResp, _ := config.SupabaseRequest("PATCH", fmt.Sprintf("users?id=eq.%s", idStr), gin.H{
		"saldo_uang": newSaldo,
	})
	defer updateSaldoResp.Body.Close()

	// --- Langkah 7: Insert ke Tabel Transactions (DB 2) ---
	txData := map[string]interface{}{
		"user_id":        idStr,
		"total_amount":   packPrice,
		"status":         "success",
		"recipient_name": "Gacha Pack: " + packName, // Keterangan transaksi
	}

	respTx, err := config.SupabaseRequest("POST", "transactions", txData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mencatat transaksi"})
		return
	}
	defer respTx.Body.Close()

	// Ambil ID transaksi yang baru saja terbuat
	var createdTxs []models.Transaction
	json.NewDecoder(respTx.Body).Decode(&createdTxs)

	var newTxID string
	if len(createdTxs) > 0 {
		newTxID = createdTxs[0].ID
	}

	// --- Langkah 8: Insert 5 Kartu ke Tabel card_collection (DB 3) ---
	// Masukkan ID Transaksi ke setiap kartu agar riwayatnya terhubung
	for i := range obtainedCards {
		obtainedCards[i].TransactionID = newTxID
	}

	respSaveCards, err := config.SupabaseRequest("POST", "card_collection", obtainedCards)
	if err == nil {
		defer respSaveCards.Body.Close()
	}

	// --- Langkah 9: Kirim Respons ke Frontend ---
	c.JSON(http.StatusOK, gin.H{
		"message":        "Berhasil membeli pack dan mendapatkan kartu!",
		"transaction_id": newTxID,
		"pack_dibeli":    packName,
		"sisa_saldo":     newSaldo,
		"kartu_didapat":  obtainedCards,
	})
}
