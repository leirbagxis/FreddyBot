package channelpost

import (
	"context"

	"github.com/leirbagxis/FreddyBot/internal/container"
)

func StageQueue(c *container.AppContainer, executionPipeline *Pipeline) Stage {
	return func(pCtx *ProcessingContext) error {
		// We need a fresh context for the worker that won't be cancelled when the handler returns
		workerCtx := *pCtx // Shallow copy
		workerCtx.Ctx = context.Background()
		workerCtx.StopPipeline = false
		
		// Use the existing messageQueue for rate limiting and concurrency control
		messageQueue.AddV2ToQueue(&workerCtx, executionPipeline)
		
		// Once added to the queue, this thread's work is done
		pCtx.StopPipeline = true
		return nil
	}
}
