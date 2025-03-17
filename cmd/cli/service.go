package main

import (
	"context"
	"github.com/uala-challenge/event-processor/cmd/cli/tweet_batch_proccesor"
	"github.com/uala-challenge/event-processor/internal/platfrom/consume_tweet_event_sqs"
	"github.com/uala-challenge/event-processor/internal/task_manager"
	"github.com/uala-challenge/simple-toolkit/pkg/simplify/app_engine"
)

type app struct {
	batchApplication app_engine.Engine
	repositories     repositories
	useCases         useCases
	batchProcessors  batchProcessor
}

var _ Service = (*app)(nil)

func NewService() *app {
	a := *app_engine.NewApp()
	return &app{batchApplication: a}
}

func (a *app) initRepositories() {
	a.repositories.ConsumeTweet = consume_tweet_event_sqs.NewService(consume_tweet_event_sqs.Dependencies{
		Client: a.batchApplication.SQSClient,
		Log:    a.batchApplication.Log,
	})

}

func (a *app) initUseCases() {
	a.useCases.Tasker = task_manager.NewManager(5, a.repositories.ConsumeTweet.ProcessMessage)
}

func (a *app) initBatchProcessors() {
	a.batchProcessors.TweetRunner = tweet_batch_proccesor.
		NewRunner(a.repositories.ConsumeTweet, a.useCases.Tasker, "http://sqs.us-east-1.localhost.localstack.cloud:4566/000000000000/tweets")
}

func (a *app) run() {
	a.initRepositories()
	a.initUseCases()
	a.initBatchProcessors()
	ctx := context.Background()
	a.batchProcessors.TweetRunner.Run(ctx)

}
