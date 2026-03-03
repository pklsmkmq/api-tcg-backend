package models

type Transaction struct {
	ID               string              `json:"id,omitempty"`
	UserID           int                 `json:"user_id"`
	RecipientName    string              `json:"recipient_name"`
	RecipientAddress string              `json:"recipient_address"`
	RecipientPhone   string              `json:"recipient_phone"`
	Subtotal         int                 `json:"subtotal"`
	TotalAmount      int                 `json:"total_amount"`
	Status           string              `json:"status"`
	CreatedAt        string              `json:"created_at,omitempty"`
	Details          []TransactionDetail `json:"details,omitempty"`
}

type TransactionDetail struct {
	ID            string `json:"id,omitempty"`
	TransactionID string `json:"transaction_id"`
	ProductID     int    `json:"product_id"`
	Quantity      int    `json:"quantity"`
	Price         int    `json:"price"`
	Subtotal      int    `json:"subtotal"`
}
