package cache

type ChannelPayload struct {
	ChannelID  int64  `json:"channel_id"`
	Title      string `json:"title,omitempty"`
	OwnerID    int64  `json:"owner_id"`
	NewOwnerID int64  `json:"new_owner_id,omitempty"`
}

type Session struct {
	Key     string         `json:"key"`
	Payload ChannelPayload `json:"payload"`
}
