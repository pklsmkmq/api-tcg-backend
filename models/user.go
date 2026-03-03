package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	SaldoUang    float64   `json:"saldo_uang"`
	CreatedAt    time.Time `json:"created_at"`
}
