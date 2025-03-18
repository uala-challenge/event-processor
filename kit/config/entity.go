package config

import (
	"github.com/uala-challenge/event-processor/cmd/cli/tweet_batch_proccesor"
)

type BatchConfig struct {
	Tweets tweet_batch_proccesor.Config `json:"tweets"`
}
