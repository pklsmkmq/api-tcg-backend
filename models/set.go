package models

// PokemonSet untuk membaca/menyimpan ke database Supabase Anda
type PokemonSet struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Logo        string `json:"logo"`
	Price       int    `json:"price"`
	Description string `json:"description"`
}

// TCGdexSet untuk membaca respons dari API TCGdex saat proses sinkronisasi
type TCGdexSet struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}
