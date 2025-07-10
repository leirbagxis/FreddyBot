package types

type CreateCustomCaptionRequest struct {
	Code        string `json:"code" binding:"required"`
	Caption     string `json:"caption" binding:"required"`
	LinkPreview bool   `json:"linkPreview"`
}

type CreateCustomCaptionResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}
