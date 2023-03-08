package schema

const (
	// low-level error messages

	InvalidRequestMethod = "Invalid request method."

	// JWT
	MissingJWTField = "Missing required `token` field."
	InvalidJWT      = "Invalid JWT."

	InvalidRequestBody  = "Invalid request body."
	InternalServerError = "Internal server error."

	// high-level error messages
	AccessDenied        = "Access denied."
	UserNotFound        = "User not found."
	TransactionNotFound = "Transaction not found."
)
