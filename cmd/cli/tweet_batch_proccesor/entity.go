package tweet_batch_proccesor

import (
	"github.com/uala-challenge/event-processor/internal/platfrom/consume_tweet_event_sqs"
	"github.com/uala-challenge/event-processor/internal/task_manager"
)

type Config struct {
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

type Dependencies struct {
	Processor   consume_tweet_event_sqs.Service
	TaskManager task_manager.Manager
	Config      Config
}
