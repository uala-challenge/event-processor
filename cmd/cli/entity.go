package main

import (
	"github.com/uala-challenge/event-processor/cmd/cli/tweet_batch_proccesor"
	"github.com/uala-challenge/event-processor/internal/platfrom/consume_tweet_event_sqs"
	"github.com/uala-challenge/event-processor/internal/task_manager"
)

type Service interface {
	run()
}
type repositories struct {
	ConsumeTweet consume_tweet_event_sqs.Service
}

type useCases struct {
	Tasker task_manager.Manager
}

type batchProcessor struct {
	TweetRunner *tweet_batch_proccesor.Runner
}
