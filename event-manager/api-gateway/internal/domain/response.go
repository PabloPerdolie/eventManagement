package domain

// ErrorResponse представляет структуру ответа при ошибке
// @Description Ответ сервера при возникновении ошибки
type ErrorResponse struct {
	Error   string `json:"error" example:"Unauthorized"`
	Message string `json:"message" example:"Invalid or expired token"`
}

// SuccessResponse представляет структуру ответа при успешном выполнении
// @Description Ответ сервера при успешном выполнении операции
type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message,omitempty" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
}

// TokenResponse представляет структуру ответа с токеном
// @Description Ответ с токеном доступа
type TokenResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn   int64  `json:"expires_in" example:"3600"`
}

func NewSuccessResponse(message string, data interface{}) SuccessResponse {
	return SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func NewErrorResponse(errorType string, message string) ErrorResponse {
	return ErrorResponse{
		Error:   errorType,
		Message: message,
	}
}
