package tweet_proccesor

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/uala-challenge/event-processor/kit"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
)

type Service interface {
	Accept(ctx context.Context, tweet kit.Tweet) error
}

type Dependencies struct {
	Client *redis.Client
	Log    log.Service
}
