package models

type Product struct {
	ID          int     `json:"id,omitempty"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Price       int     `json:"price"`
	ImageURL    string  `json:"image_url"`
	Description string  `json:"description"`
	Rating      float64 `json:"rating"`
	CreatedAt   *string `json:"created_at,omitempty"`
	UpdatedAt   *string `json:"updated_at,omitempty"`
}
