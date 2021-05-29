package models

// ErrorResponse : This is error model.
type ErrorResponse struct {
	StatusCode int    `json:"status"`
	Error      string `json:"error"`
}
