package database

type Currency string

type User struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	IsSuperuser bool   `json:"is_superuser"`

	BaseCurrency Currency `json:"base_currency"`
}

type Transaction struct {
	Amount      int      `json:"amount"`
	Currency    Currency `json:"currency"`
	Description string   `json:"description"`
	Timestamp   int8     `json:"timestamp"`
}
