package types

type TransferChannelRequest struct {
	OldOwnerID int64 `json:"oldOwnerId"`
	NewOwnerID int64 `json:"newOwnerId"`
	ChannelID  int64 `json:"channelId"`
}
