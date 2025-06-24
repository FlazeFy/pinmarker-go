package entities

type (
	Metadata struct {
		Limit      int `json:"limit"`
		Page       int `json:"page"`
		Total      int `json:"total"`
		TotalPages int `json:"total_pages"`
	}
	// For Response
	ResponseBadRequest struct {
		Message string `json:"message" example:"app_source is not valid"`
		Status  string `json:"status" example:"failed"`
	}
	ResponseNotFound struct {
		Message string `json:"message" example:"track not found"`
		Status  string `json:"status" example:"failed"`
	}
)
