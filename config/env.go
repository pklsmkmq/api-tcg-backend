package config

import (
	"github.com/joho/godotenv"
)

var (
	SUPABASE_URL string
	SUPABASE_KEY string
	JWT_SECRET   []byte
)

func LoadEnv() {
	env, err := godotenv.Read()
	if err != nil {
		panic("Error loading .env file")
	}

	SUPABASE_URL = env["SUPABASE_URL"]
	SUPABASE_KEY = env["SUPABASE_ANON_KEY"]
	JWT_SECRET = []byte(env["JWT_SECRET"])

	if SUPABASE_URL == "" || SUPABASE_KEY == "" {
		panic("ENV SUPABASE_URL dan SUPABASE_ANON_KEY wajib di set")
	}
}
