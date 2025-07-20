package entities

type (
	AppCount struct {
		AppName string `json:"app_name"`
		Total   int    `json:"total"`
	}
	// For Response
	ResponseGetAppCount struct {
		Message string     `json:"message" example:"Track fetched"`
		Status  string     `json:"status" example:"success"`
		Data    []AppCount `json:"data"`
	}
)
