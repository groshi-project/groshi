package schemas

type ErrorResponse struct {
	//ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}
