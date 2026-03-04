package models

type Transaction struct {
	ID          string `json:"id,omitempty"`
	UserID      string `json:"user_id"`
	Description string `json:"description"` // Sebelumnya recipient_name
	TotalAmount int    `json:"total_amount"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at,omitempty"`
}