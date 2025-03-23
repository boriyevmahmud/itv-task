package models

type ErrorResponse struct {
	Code    int    `json:"code"`    // HTTP status code
	Message string `json:"message"` // Error message
	Detail  string `json:"detail"`  // Optional detailed error message
}

func NewErrorResponse(code int, message string, detail string) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}
