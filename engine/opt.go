package engine

import (
	"go.uber.org/zap"

	"github.com/zoulux/crawler-go/collect"
)

var defaultOptions = options{
	Logger: zap.NewNop(),
}

type options struct {
	WorkCount int
	Logger    *zap.Logger
	Fetcher   collect.Fetcher
	Seeds     []*collect.Task
	scheduler Scheduler
}

type Option func(opts *options)

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.Logger = logger
	}
}
func WithFetcher(fetcher collect.Fetcher) Option {
	return func(opts *options) {
		opts.Fetcher = fetcher
	}
}

func WithWorkCount(workCount int) Option {
	return func(opts *options) {
		opts.WorkCount = workCount
	}
}

func WithSeeds(seeds []*collect.Task) Option {
	return func(opts *options) {
		opts.Seeds = seeds
	}
}

func WithScheduler(schedule Scheduler) Option {
	return func(opts *options) {
		opts.scheduler = schedule
	}
}
