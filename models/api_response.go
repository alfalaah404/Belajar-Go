package models

// ApiResponse represents the standard API response format
type ApiResponse struct {
	Status  string      `json:"status" example:"success"`
	Message string      `json:"message" example:"Success"`
	Data    interface{} `json:"data"`
}
