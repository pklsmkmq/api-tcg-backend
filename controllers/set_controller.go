package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go-supabase-api/config"
	"go-supabase-api/models"

	"github.com/gin-gonic/gin"
)

// 1. ENDPOINT UNTUK FRONTEND: Mengambil data langsung dari DB Anda
func GetPokemonSets(c *gin.Context) {
	resp, err := config.SupabaseRequest("GET", "pokemon_sets?select=*", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menghubungi database"})
		return
	}
	defer resp.Body.Close()

	var sets []models.PokemonSet
	json.NewDecoder(resp.Body).Decode(&sets)

	c.JSON(http.StatusOK, sets)
}

// 2. ENDPOINT UNTUK ADMIN: Tarik data dari TCGdex dan simpan ke DB Supabase
func SyncSets(c *gin.Context) {
	// Ambil dari API TCGdex
	resp, err := http.Get("https://api.tcgdex.net/v2/en/sets")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil data dari TCGdex"})
		return
	}
	defer resp.Body.Close()

	var tcgdexSets []models.TCGdexSet
	if err := json.NewDecoder(resp.Body).Decode(&tcgdexSets); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memproses respons"})
		return
	}

	var formattedSets []models.PokemonSet

	// LOOPING UNTUK SELURUH DATA TANPA LIMIT
	for _, set := range tcgdexSets {
		logoUrl := set.Logo
		if logoUrl != "" {
			logoUrl = fmt.Sprintf("%s.png", logoUrl)
		}

		formattedSets = append(formattedSets, models.PokemonSet{
			ID:          set.ID,
			Name:        set.Name,
			Logo:        logoUrl,
			Price:       50000, // Harga default awal
			Description: fmt.Sprintf("Pack booster %s resmi.", set.Name),
		})
	}

	// Simpan seluruh array formattedSets ke tabel pokemon_sets di Supabase
	insertResp, err := config.SupabaseRequest("POST", "pokemon_sets", formattedSets)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menyimpan ke database"})
		return
	}
	defer insertResp.Body.Close()

	if insertResp.StatusCode >= 400 {
		var errResp map[string]interface{}
		json.NewDecoder(insertResp.Body).Decode(&errResp)
		c.JSON(insertResp.StatusCode, gin.H{"message": "Gagal insert ke Supabase", "error": errResp})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Berhasil sinkronisasi SEMUA data set Pokémon ke database",
		"total_inserted": len(formattedSets), // Menampilkan total data yang berhasil dimasukkan
	})
}
