package schema

const (
	TransactionNotFoundErrorDetail               = "transaction not found"
	ThisTransactionDoesNotBelongToYouErrorDetail = "this transaction does not belong to you"
	UnknownCurrencyErrorDetail                   = "unknown currency"
)

// todo: formalize
// - missing required fields: `a`, `b`, `c`
// - missing required field `a`
// - at least one of the following fields is required: `a`, `b`, `c`
