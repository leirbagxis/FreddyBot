package types

type CaptionDefaultUpdateRequest struct {
	Caption string `json:"caption" binding:"required"`
}

type CaptionUpdateResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}
