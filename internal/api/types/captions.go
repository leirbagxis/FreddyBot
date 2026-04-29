package types

type CaptionDefaultUpdateRequest struct {
	Caption string `json:"caption" binding:"required"`
}

type ReactionsUpdateRequest struct {
	Reactions string `json:"reactions"`
}

type ReactionPositionUpdateRequest struct {
	ReactionPosition int `json:"reactionPosition"`
}

type CaptionUpdateResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}
