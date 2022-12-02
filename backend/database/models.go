package database

type User struct {
	Username     string   `json:"username"`
	Password     string   `json:"password"`
	BaseCurrency Currency `json:"base_currency"`
}

type Transaction struct {
	Amount      int      `json:"amount"`
	Currency    Currency `json:"currency"`
	Description string   `json:"description"`
	Timestamp   int8     `json:"timestamp"`
}

type Currency string
