package models

import "time"

// CardCollection merepresentasikan tabel card_collection di database
type CardCollection struct {
	ID            string    `json:"id,omitempty"`
	UserID        string    `json:"user_id"`
	TransactionID string    `json:"transaction_id,omitempty"`
	CardID        string    `json:"card_id"`
	CardName      string    `json:"card_name"`
	ImageURL      string    `json:"image_url"`
	SetID         string    `json:"set_id"`
	AcquiredAt    time.Time `json:"acquired_at,omitempty"`
}

// TCGdexCard merepresentasikan respons data kartu 1 biji dari API TCGdex
type TCGdexCard struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

// TCGdexSetDetail merepresentasikan respons API TCGdex saat kita memanggil /sets/{id}
type TCGdexSetDetail struct {
	ID    string       `json:"id"`
	Name  string       `json:"name"`
	Cards []TCGdexCard `json:"cards"`
}
