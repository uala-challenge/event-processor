package main

import (
	"context"
	"github.com/uala-challenge/event-processor/internal/platfrom/consume_tweet_event_sqs"
	"github.com/uala-challenge/event-processor/internal/platfrom/tweet_proccesor"
	"github.com/uala-challenge/event-processor/kit/config"
	"github.com/uala-challenge/simple-toolkit/pkg/simplify/app_engine"

	"github.com/uala-challenge/event-processor/cmd/cli/tweet_batch_proccesor"
	"github.com/uala-challenge/event-processor/internal/task_manager"
)

type app struct {
	batchApplication app_engine.Engine
	repositories     repositories
	useCases         useCases
	batchProcessors  batchProcessor
	batchConfig      config.BatchConfig
}

var _ Service = (*app)(nil)

func NewService() *app {
	a := *app_engine.NewApp()
	return &app{batchApplication: a}
}

func (a *app) LoadConfig() {
	a.batchConfig = app_engine.GetConfig[config.BatchConfig](a.batchApplication.BatchConfig)
}

func (a *app) initRepositories() {
	a.repositories.TweetProcessor = tweet_proccesor.NewService(tweet_proccesor.Dependencies{
		Client: a.batchApplication.RedisClient,
		Log:    a.batchApplication.Log,
	})
	a.repositories.ConsumeTweet = consume_tweet_event_sqs.NewService(consume_tweet_event_sqs.Dependencies{
		Client:    a.batchApplication.SQSClient,
		Log:       a.batchApplication.Log,
		Processor: a.repositories.TweetProcessor,
	})
}

func (a *app) initUseCases() {
	a.useCases.Tasker = task_manager.NewManager(5, a.repositories.ConsumeTweet.ProcessMessage)
}

func (a *app) initBatchProcessors() {
	a.batchProcessors.TweetRunner = tweet_batch_proccesor.
		NewRunner(tweet_batch_proccesor.Dependencies{
			Processor:   a.repositories.ConsumeTweet,
			TaskManager: a.useCases.Tasker,
			Config:      a.batchConfig.Tweets,
		})
}

func (a *app) run() {
	a.LoadConfig()
	a.initRepositories()
	a.initUseCases()
	a.initBatchProcessors()
	ctx := context.Background()
	a.batchProcessors.TweetRunner.Run(ctx)

}
