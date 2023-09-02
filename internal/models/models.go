// Package models package models contains response models that are returned.
package models

// User represents response containing information about user.
type User struct {
	Username string `json:"username" example:"pipka5000"`
}

// Transaction represents response containing transaction information.
type Transaction struct {
	UUID string `json:"uuid" example:"c81ab774-3f96-40e8-9ebd-170e303a682e"`

	Amount      int    `json:"amount" example:"-999"`
	Currency    string `json:"currency" example:"USD"`
	Description string `json:"description" example:"Bought some donuts for $9.99..."`
	Date        string `json:"date" example:"2023-09-02T12:38:10+03:00"`

	CreatedAt string `json:"created_at" example:"2023-09-02T12:38:10+03:00"`
	UpdatedAt string `json:"updated_at" example:"2023-09-02T12:38:10+03:00"`
}

// Summary represents summary of transactions, returned by transactionsReadSummary handler.
type Summary struct {
	Currency string `json:"currency" example:"USD"`

	Income  int `json:"income" example:"5000"`
	Outcome int `json:"outcome" example:"1000"`
	Total   int `json:"total" example:"4000"`

	TransactionsCount int `json:"transactions_count" example:"6"`
}

// Error represents response containing information about API error.
type Error struct {
	ErrorMessage string   `json:"error_message" example:"example error message"`
	ErrorDetails []string `json:"error_details"`
}
