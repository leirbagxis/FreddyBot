package cache

type ChannelPayload struct {
	ChannelID int64  `json:"channel_id"`
	Title     string `json:"title"`
	OwnerID   int64  `json:"owner_id"`
}

type Session struct {
	Key     string         `json:"key"`
	Payload ChannelPayload `json:"payload"`
}
