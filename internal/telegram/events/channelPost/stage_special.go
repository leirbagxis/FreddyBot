package channelpost

import (
	"github.com/leirbagxis/FreddyBot/internal/container"
)

func StageSpecialFlows(c *container.AppContainer) Stage {
	return func(pCtx *ProcessingContext) error {
		post := pCtx.Update.ChannelPost
		if post == nil {
			return nil
		}

		handled, err := TryHandleNewPack(pCtx.Ctx, pCtx.Bot, *pCtx.Channel, *post)
		if err != nil {
			return err
		}

		if handled {
			pCtx.StopPipeline = true
		}

		return nil
	}
}
