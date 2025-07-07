package repositories

import "golang.org/x/net/context"

type Request struct {
	Caption string `json:"caption"`
}

func (r *ChannelRepository) UpdateDefaultCaption(ctx context.Context, channelID int64) {

}
