package channelpost

import (
	"fmt"

	"github.com/leirbagxis/FreddyBot/pkg/logger"
)

// Stage defines a single step in the processing pipeline.
type Stage func(ctx *ProcessingContext) error

// Pipeline orchestrates the execution of multiple stages.
type Pipeline struct {
	Name   string
	stages []Stage
}

// NewPipeline initializes a pipeline with the provided stages.
func NewPipeline(name string, stages ...Stage) *Pipeline {
	return &Pipeline{
		Name:   name,
		stages: stages,
	}
}

// Execute runs the pipeline for a given context starting from the first stage.
func (p *Pipeline) Execute(ctx *ProcessingContext) error {
	return p.ExecuteFrom(ctx, 0)
}

// ExecuteFrom runs the pipeline starting from a specific stage index.
func (p *Pipeline) ExecuteFrom(ctx *ProcessingContext, startIndex int) error {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("PIPELINE", "[%s] Recovered from panic: %v", p.Name, r)
			ctx.Error = fmt.Errorf("pipeline panic: %v", r)
		}
	}()

	for i := startIndex; i < len(p.stages); i++ {
		if ctx.StopPipeline {
			logger.Bot("⏹️ Pipeline [%s] stopped at stage %d", p.Name, i+1)
			break
		}

		if err := p.stages[i](ctx); err != nil {
			ctx.Error = err
			logger.Error("PIPELINE", "[%s] Error at stage %d: %v", p.Name, i+1, err)
			return err
		}
	}

	return nil
}
