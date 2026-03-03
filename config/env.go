package config

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	SUPABASE_URL string
	SUPABASE_KEY string
	JWT_SECRET   []byte
)

func LoadEnv() {
	// Membaca file .env menjadi bentuk map (tidak memengaruhi/membaca OS environment Windows)
	envMap, err := godotenv.Read()

	if err == nil {
		// KONDISI 1: File .env DITEMUKAN (Berarti jalan di Lokal)
		// Kita ambil nilainya murni dari dalam file .env, mengabaikan variabel OS Windows
		SUPABASE_URL = envMap["SUPABASE_URL"]
		SUPABASE_KEY = envMap["SUPABASE_ANON_KEY"]

		if jwtSecret, exists := envMap["JWT_SECRET"]; exists && jwtSecret != "" {
			JWT_SECRET = []byte(jwtSecret)
		}
	} else {
		// KONDISI 2: File .env TIDAK DITEMUKAN (Berarti jalan di Vercel/Production)
		// Vercel tidak memakai file .env, melainkan menyuntikkannya langsung ke OS
		SUPABASE_URL = os.Getenv("SUPABASE_URL")
		SUPABASE_KEY = os.Getenv("SUPABASE_ANON_KEY")

		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret != "" {
			JWT_SECRET = []byte(jwtSecret)
		}
	}

	// Validasi untuk memastikan variabel wajib tidak kosong
	if SUPABASE_URL == "" || SUPABASE_KEY == "" {
		panic("ENV SUPABASE_URL dan SUPABASE_ANON_KEY wajib di set")
	}
}
