package config

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func SupabaseRequest(method, path string, body any) (*http.Response, error) {
	var jsonBody []byte
	if body != nil {
		jsonBody, _ = json.Marshal(body)
	}

	req, _ := http.NewRequest(
		method,
		SUPABASE_URL+"/rest/v1/"+path,
		bytes.NewBuffer(jsonBody),
	)

	req.Header.Set("apikey", SUPABASE_KEY)
	req.Header.Set("Authorization", "Bearer "+SUPABASE_KEY)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation")

	return (&http.Client{}).Do(req)
}
