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

type PostBuilderButton struct {
	Text string `json:"text"`
	URL  string `json:"url"`
}

type PostBuilderState struct {
	MediaType   string              `json:"media_type"`
	MediaFileID string              `json:"media_file_id"`
	Title       string              `json:"title"`
	Body        string              `json:"body"`
	Footer      string              `json:"footer"`
	Buttons     []PostBuilderButton `json:"buttons"`
	Step        string              `json:"step"`
}
