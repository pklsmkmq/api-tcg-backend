package models

import "time"

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	SaldoUang    float64   `json:"saldo_uang"`
	Role         string    `json:"role"` // <-- TAMBAHAN BARU
	CreatedAt    time.Time `json:"created_at"`
}
