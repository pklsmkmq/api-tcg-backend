package controllers

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"go-supabase-api/config"
	"go-supabase-api/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var req struct {
		Name         string `json:"username"`
		Email        string `json:"email"`
		PasswordHash string `json:"password_hash"`
		Confirm      string `json:"confirm_password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": "Format JSON salah"})
		return
	}

	if req.PasswordHash != req.Confirm {
		c.JSON(400, gin.H{"message": "Password tidak sama"})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.PasswordHash), 10)

	resp, err := config.SupabaseRequest("POST", "users", gin.H{
		"username":      req.Name,
		"email":         req.Email,
		"password_hash": string(hash),
		"saldo_uang":    10000,
	})

	if err != nil {
		c.JSON(500, gin.H{"message": "Gagal menghubungi database", "error": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		c.JSON(resp.StatusCode, gin.H{"message": "Register gagal", "detail": errorResponse})
		return
	}

	c.JSON(201, gin.H{"message": "Register berhasil"})
}

func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": "Format JSON salah"})
		return
	}

	resp, _ := config.SupabaseRequest(
		"GET",
		fmt.Sprintf("users?username=eq.%s&select=*", req.Username),
		nil,
	)

	var users []models.User
	json.NewDecoder(resp.Body).Decode(&users)

	if len(users) == 0 {
		c.JSON(401, gin.H{"message": "Username tidak ditemukan"})
		return
	}

	user := users[0]

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		c.JSON(401, gin.H{"message": "Password salah"})
		return
	}

	token, _ := generateJWT(user)

	// ⬇️ RESPONSE BARU
	c.JSON(200, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"saldo":    user.SaldoUang,
		},
	})
}

func generateJWT(user models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key" // fallback, should use env variable
	}

	claims := jwt.MapClaims{
		"id":       user.ID,
		"email":    user.Email,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24 hours expiration
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
